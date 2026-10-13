package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WaitlistReceiptSender func(context.Context, string, string, string, string) error

// SES delivery and the database commit are separate operations.
// A crash after SES accepts a message can cause a duplicate on retry.
func (s *ApplicationService) DeliverWaitlistReceipts(
	ctx context.Context, send WaitlistReceiptSender,
) error {
	if send == nil {
		return errors.New("waitlist receipt sender is required")
	}
	for i := 0; i < 20; i++ {
		processed := false
		err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
			var id uuid.UUID
			var kind, source, recipient, name string
			var attempts int
			var eligible bool
			err := tx.QueryRow(ctx, `
				SELECT o.id,o.kind,o.signup_source,o.attempts,
					COALESCE(
						NULLIF(trim(u.preferred_email),''),
						NULLIF(trim(u.email),''),
						NULLIF(trim(a.application->>'preferredEmail'),''),
						NULLIF(trim(a.application->>'universityEmail'),''),''
					),
					COALESCE(
						NULLIF(trim(a.application->>'firstName'),''),
						NULLIF(split_part(trim(u.name),' ',1),''),'Hacker'
					),
					COALESCE(
						NOT u.is_fake AND NOT a.is_fake
						AND u.role='applicant'
						AND (
							(o.kind='signup' AND a.status::text='waitlisted'
								AND w.created_at=o.action_at
								AND w.signup_source=o.signup_source
								AND w.in_person_joined_at IS NULL)
							OR
							(o.kind='arrival' AND a.status::text='waitlist_confirmed'
								AND w.in_person_joined_at=o.action_at)
						),false
					)
				FROM waitlist_receipt_email_outbox o
				JOIN applications a
					ON a.user_id=o.user_id AND a.hackathon_id=o.hackathon_id
				JOIN users u ON u.id=o.user_id
				JOIN hackathons h ON h.id=o.hackathon_id
				LEFT JOIN waitlist w
					ON w.user_id=o.user_id AND w.hackathon_id=o.hackathon_id
				WHERE h.is_active AND o.hackathon_id='xii'
					AND o.sent_at IS NULL AND o.cancelled_at IS NULL
					AND o.next_attempt_at<=now()
				ORDER BY o.next_attempt_at,o.created_at,o.id
				LIMIT 1 FOR UPDATE OF a,o SKIP LOCKED
			`).Scan(&id, &kind, &source, &attempts, &recipient, &name, &eligible)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			if err != nil {
				return err
			}
			processed = true
			if !eligible {
				_, err := tx.Exec(ctx, `
					UPDATE waitlist_receipt_email_outbox
					SET cancelled_at=now(),last_error='Receipt no longer current'
					WHERE id=$1
				`, id)
				return err
			}

			var sendErr error
			if recipient == "" {
				sendErr = errors.New("recipient email is missing")
			} else {
				sendErr = send(ctx, recipient, name, kind, source)
			}
			if sendErr != nil {
				delay := time.Minute
				for n := 0; n < attempts && delay < time.Hour; n++ {
					delay *= 2
				}
				if delay > time.Hour {
					delay = time.Hour
				}
				_, err := tx.Exec(ctx, `
					UPDATE waitlist_receipt_email_outbox
					SET attempts=attempts+1,last_error=$2,next_attempt_at=$3
					WHERE id=$1
				`, id, sendErr.Error(), time.Now().UTC().Add(delay))
				return err
			}
			_, err = tx.Exec(ctx, `
				UPDATE waitlist_receipt_email_outbox
				SET sent_at=now(),attempts=attempts+1,last_error=NULL
				WHERE id=$1
			`, id)
			return err
		})
		if err != nil {
			return err
		}
		if !processed {
			return nil
		}
		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return nil
}
