package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrInPersonWaitlistEligibility = errors.New("only rejected or waitlisted applicants can join the in-person waitlist")
var ErrInPersonWaitlistClosed = errors.New("the in-person waitlist is not currently open")
var ErrInPersonWaitlistPermission = errors.New("an administrator must record in-person arrival")

func (s *ApplicationService) RecordInPersonAdmissionWaitlist(
	ctx context.Context, userID, recordedBy uuid.UUID, now time.Time,
) error {
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		var admin bool
		if err := tx.QueryRow(ctx, `
            SELECT EXISTS (SELECT 1 FROM users WHERE id=$1 AND role='admin')
        `, recordedBy).Scan(&admin); err != nil {
			return err
		}
		if !admin {
			return ErrInPersonWaitlistPermission
		}

		var eventID string
		var opensAt, closesAt time.Time
		err := tx.QueryRow(ctx, `
            SELECT h.id, p.in_person_opens_at, p.invitations_close_at
            FROM hackathons h
            JOIN waitlist_dispatch_policies p ON p.hackathon_id=h.id
            WHERE h.is_active AND p.enabled
            FOR UPDATE OF h, p
        `).Scan(&eventID, &opensAt, &closesAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInPersonWaitlistClosed
		}
		if err != nil {
			return err
		}
		if now.Before(opensAt) || !now.Before(closesAt) {
			return ErrInPersonWaitlistClosed
		}

		var applicationID uuid.UUID
		var status string
		err = tx.QueryRow(ctx, `
            SELECT a.id, a.status::text
            FROM applications a JOIN users u ON u.id=a.user_id
            WHERE a.hackathon_id=$1 AND a.user_id=$2
              AND u.role='applicant' AND NOT u.is_fake AND NOT a.is_fake
            FOR UPDATE OF a
        `, eventID, userID).Scan(&applicationID, &status)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInPersonWaitlistEligibility
		}
		if err != nil {
			return err
		}
		if status != "rejected" && status != "waitlisted" {
			return ErrInPersonWaitlistEligibility
		}

		if _, err := tx.Exec(ctx, `
            INSERT INTO waitlist
                (hackathon_id, user_id, in_person_joined_at, in_person_recorded_by)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (hackathon_id, user_id) DO UPDATE
            SET in_person_joined_at=COALESCE(waitlist.in_person_joined_at, EXCLUDED.in_person_joined_at),
                in_person_recorded_by=COALESCE(waitlist.in_person_recorded_by, EXCLUDED.in_person_recorded_by)
        `, eventID, userID, now, recordedBy); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
            UPDATE applications SET status='waitlisted', updated_at=now()
            WHERE id=$1
        `, applicationID)
		return err
	})
}
