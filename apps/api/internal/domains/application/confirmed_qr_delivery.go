package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ConfirmedQRSender func(context.Context, uuid.UUID, string, string) (string, error)

func (s *ApplicationService) DeliverConfirmedQREmails(ctx context.Context, send ConfirmedQRSender) error {
	if send == nil {
		return errors.New("confirmed QR sender is required")
	}
	for i := 0; i < 20; i++ {
		processed := false
		err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
			var applicationID, userID uuid.UUID
			var status, role, recipient, name string
			var fake bool
			var attempts int
			err := tx.QueryRow(ctx, `
				SELECT o.application_id,o.user_id,a.status::text,u.role::text,
					(a.is_fake OR u.is_fake),o.attempts,
					COALESCE(NULLIF(btrim(u.preferred_email),''),
						NULLIF(btrim(u.email),''),
						NULLIF(btrim(a.application->>'preferredEmail'),''),
						NULLIF(btrim(a.application->>'universityEmail'),''),''),
					COALESCE(NULLIF(btrim(a.application->>'firstName'),''),
						NULLIF(split_part(btrim(u.name),' ',1),''),'Hacker')
				FROM confirmed_qr_email_outbox o
				JOIN applications a ON a.id=o.application_id
				JOIN users u ON u.id=o.user_id
				JOIN hackathons h ON h.id=o.hackathon_id
				WHERE h.is_active AND o.hackathon_id='xii'
					AND o.sent_at IS NULL AND o.cancelled_at IS NULL
					AND o.next_attempt_at<=now()
				ORDER BY o.next_attempt_at,o.created_at,o.application_id
				LIMIT 1 FOR UPDATE OF a,o SKIP LOCKED
			`).Scan(&applicationID, &userID, &status, &role, &fake,
				&attempts, &recipient, &name)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			if err != nil {
				return err
			}
			processed = true
			if fake || status != "confirmed" || role != "attendee" {
				_, err := tx.Exec(ctx, `
					UPDATE confirmed_qr_email_outbox
					SET cancelled_at=now(),last_error='Recipient no longer eligible'
					WHERE application_id=$1
				`, applicationID)
				return err
			}
			messageID, sendErr := send(ctx, userID, recipient, name)
			if sendErr != nil {
				delay := time.Minute
				for n := 0; n < attempts && delay < time.Hour; n++ {
					delay *= 2
				}
				if delay > time.Hour {
					delay = time.Hour
				}
				_, err := tx.Exec(ctx, `
					UPDATE confirmed_qr_email_outbox
					SET attempts=attempts+1,last_error=$2,next_attempt_at=$3
					WHERE application_id=$1
				`, applicationID, sendErr.Error(), time.Now().UTC().Add(delay))
				return err
			}
			_, err = tx.Exec(ctx, `
				UPDATE confirmed_qr_email_outbox
				SET sent_at=now(),message_id=$2,last_error=NULL,attempts=attempts+1
				WHERE application_id=$1
			`, applicationID, messageID)
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
