package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrWaitlistEligibility = errors.New("only rejected applicants can join the waitlist")
var ErrWithdrawalEligibility = errors.New("only accepted or confirmed applicants can withdraw")

// Lock the active hackathon's application before checking its status.
// This serializes joins and withdrawals with concurrent decision updates.
func (s *ApplicationService) changeAdmissionStatus(
	ctx context.Context, userID uuid.UUID, join bool,
) error {
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		var applicationID uuid.UUID
		var hackathonID, status string

		err := tx.QueryRow(ctx, `
			SELECT a.id, a.hackathon_id, a.status::text
			FROM applications a
			JOIN hackathons h ON h.id = a.hackathon_id
			WHERE a.user_id = $1 AND h.is_active = true
			FOR UPDATE OF a
		`, userID).Scan(&applicationID, &hackathonID, &status)
		if errors.Is(err, pgx.ErrNoRows) {
			if join {
				return ErrWaitlistEligibility
			}
			return ErrWithdrawalEligibility
		}
		if err != nil {
			return err
		}

		if join {
			if status != "rejected" && status != "waitlisted" {
				return ErrWaitlistEligibility
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO waitlist (hackathon_id, user_id)
				VALUES ($1, $2)
				ON CONFLICT (hackathon_id, user_id) DO NOTHING
			`, hackathonID, userID)
			if err != nil {
				return err
			}
			_, err = tx.Exec(ctx, `
				UPDATE applications
				SET status = 'waitlisted', updated_at = now()
				WHERE id = $1
			`, applicationID)
			return err
		}

		if status != "accepted" && status != "confirmed" && status != "withdrawn" {
			return ErrWithdrawalEligibility
		}
		_, err = tx.Exec(ctx, `
			DELETE FROM waitlist
			WHERE hackathon_id = $1 AND user_id = $2
		`, hackathonID, userID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			UPDATE applications
			SET status = 'withdrawn', updated_at = now()
			WHERE id = $1
		`, applicationID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			UPDATE users
			SET role = 'applicant', updated_at = now()
			WHERE id = $1 AND role = 'attendee'
		`, userID)
		return err
	})
}

func (s *ApplicationService) JoinAdmissionWaitlist(
	ctx context.Context, userID uuid.UUID,
) error {
	return s.changeAdmissionStatus(ctx, userID, true)
}
