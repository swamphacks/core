package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WaitlistInvitationSender func(
	context.Context, string, string, time.Time,
) error

type WaitlistDeliveryResult struct {
	Sent   int
	Failed int
}

// A crash after SES accepts an email but before the transaction commits
// can cause a duplicate on retry.
func (s *ApplicationService) DeliverAdmissionWaitlistInvitations(
	ctx context.Context, send WaitlistInvitationSender,
) (WaitlistDeliveryResult, error) {
	var result WaitlistDeliveryResult
	if send == nil {
		return result, errors.New("invitation sender is required")
	}

	for i := 0; i < 50; i++ {
		processed := false
		sent, failed := false, false
		err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
			var outboxID, applicationID uuid.UUID
			err := tx.QueryRow(ctx, `
                SELECT o.id, o.application_id
                FROM waitlist_invitation_outbox o
                JOIN applications a ON a.id = o.application_id
                JOIN hackathons h ON h.id = a.hackathon_id
                JOIN waitlist_dispatch_policies p ON p.hackathon_id = h.id
                WHERE h.is_active AND p.enabled
                  AND o.sent_at IS NULL AND o.cancelled_at IS NULL
                  AND o.next_attempt_at <= now()
                ORDER BY o.next_attempt_at, o.created_at, o.id
                LIMIT 1
            `).Scan(&outboxID, &applicationID)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			if err != nil {
				return err
			}

			var status string
			if err := tx.QueryRow(ctx, `
                SELECT status::text FROM applications WHERE id=$1 FOR UPDATE
            `, applicationID).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil
				}
				return err
			}

			var recipient, name string
			var offeredAt time.Time
			var attempts int
			err = tx.QueryRow(ctx, `
                SELECT recipient, first_name, offered_at, attempts
                FROM waitlist_invitation_outbox
                WHERE id=$1 AND sent_at IS NULL AND cancelled_at IS NULL
                  AND next_attempt_at <= now()
                FOR UPDATE
            `, outboxID).Scan(&recipient, &name, &offeredAt, &attempts)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			if err != nil {
				return err
			}
			processed = true

			var currentOffer bool
			if err := tx.QueryRow(ctx, `
                SELECT EXISTS (
                    SELECT 1 FROM application_waitlist_offers
                    WHERE application_id=$1 AND offered_at=$2
                )
            `, applicationID, offeredAt).Scan(&currentOffer); err != nil {
				return err
			}
			if status != "accepted" || !currentOffer {
				_, err := tx.Exec(ctx, `
                    UPDATE waitlist_invitation_outbox
                    SET cancelled_at=now() WHERE id=$1
                `, outboxID)
				return err
			}

			// Recheck after acquiring the application and outbox locks.
			var closesAt time.Time
			if err := tx.QueryRow(ctx, `
            SELECT p.invitations_close_at
            FROM applications a
            JOIN waitlist_dispatch_policies p ON p.hackathon_id=a.hackathon_id
            WHERE a.id=$1
        `, applicationID).Scan(&closesAt); err != nil {
				return err
			}
			if !time.Now().Before(closesAt) {
				if _, err := tx.Exec(ctx, `
                UPDATE applications SET status='withdrawn', updated_at=now()
                WHERE id=$1 AND status='accepted'
            `, applicationID); err != nil {
					return err
				}
				_, err := tx.Exec(ctx, `
                UPDATE waitlist_invitation_outbox
                SET cancelled_at=now(), last_error='Invitation delivery window closed'
                WHERE id=$1
            `, outboxID)
				return err
			}

			deadline := time.Now().UTC().Add(48 * time.Hour)
			if sendErr := send(ctx, recipient, name, deadline); sendErr != nil {
				delay := time.Minute
				for n := 0; n < attempts && delay < time.Hour; n++ {
					delay *= 2
				}
				if delay > time.Hour {
					delay = time.Hour
				}
				_, err := tx.Exec(ctx, `
                    UPDATE waitlist_invitation_outbox
                    SET attempts=attempts+1, last_error=$2, next_attempt_at=$3
                    WHERE id=$1
                `, outboxID, sendErr.Error(), time.Now().UTC().Add(delay))
				failed = err == nil
				return err
			}

			if _, err := tx.Exec(ctx, `
                UPDATE application_waitlist_offers
                SET confirmation_deadline=$2
                WHERE application_id=$1
            `, applicationID, deadline); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `
                UPDATE waitlist_invitation_outbox
                SET confirmation_deadline=$2, sent_at=now(),
                    attempts=attempts+1, last_error=NULL
                WHERE id=$1
            `, outboxID, deadline); err != nil {
				return err
			}
			sent = true
			return nil
		})
		if err != nil {
			return result, err
		}
		if sent {
			result.Sent++
		}
		if failed {
			result.Failed++
		}
		if !processed {
			break
		}

		if i < 49 {
			timer := time.NewTimer(1500 * time.Millisecond)
			select {
			case <-ctx.Done():
				timer.Stop()
				return result, ctx.Err()
			case <-timer.C:
			}
		}
	}
	return result, nil
}
