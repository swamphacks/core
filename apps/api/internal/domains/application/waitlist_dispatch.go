package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WaitlistDispatchResult struct {
	Expired int64
	Offered int
}

// The event and policy locks serialize dispatchers for the same event.
// Candidate application locks serialize offers with joins, leaves and confirmations.
func (s *ApplicationService) DispatchAdmissionWaitlist(
	ctx context.Context, now time.Time,
) (WaitlistDispatchResult, error) {
	var result WaitlistDispatchResult
	err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		var eventID string
		var capacity *int
		var opensAt, inPersonAt, closesAt time.Time
		err := tx.QueryRow(ctx, `
            SELECT h.id, h.max_attendees, p.invitations_open_at,
                   p.in_person_opens_at, p.invitations_close_at
            FROM hackathons h
            JOIN waitlist_dispatch_policies p ON p.hackathon_id = h.id
            WHERE h.is_active AND p.enabled
            FOR UPDATE OF h, p
        `).Scan(&eventID, &capacity, &opensAt, &inPersonAt, &closesAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if now.Before(opensAt) {
			return nil
		}
		if capacity == nil || *capacity <= 0 {
			return errors.New("waitlist dispatch requires a positive event capacity")
		}

		// Release current offers that could not be sent before invitations closed.
		if !now.Before(closesAt) {
			if _, err := tx.Exec(ctx, `
                UPDATE applications a
                SET status='withdrawn', updated_at=now()
                WHERE a.hackathon_id=$1 AND a.status='accepted'
                  AND EXISTS (
                    SELECT 1 FROM waitlist_invitation_outbox pending
                    JOIN application_waitlist_offers current_offer
                      ON current_offer.application_id=pending.application_id
                     AND current_offer.offered_at=pending.offered_at
                    WHERE pending.application_id=a.id
                      AND pending.sent_at IS NULL
                      AND pending.cancelled_at IS NULL
                  )
            `, eventID); err != nil {
				return err
			}
		}

		// Confirmed attendance is never expired here.
		expired, err := tx.Exec(ctx, `
            UPDATE applications a
            SET status = 'withdrawn', updated_at = now()
            FROM hackathons h
            WHERE a.hackathon_id = h.id AND h.id = $1
              AND a.status = 'accepted'
              AND NOT EXISTS (
                  SELECT 1 FROM waitlist_invitation_outbox pending
                  JOIN application_waitlist_offers current_offer
                    ON current_offer.application_id = pending.application_id
                   AND current_offer.offered_at = pending.offered_at
                  WHERE pending.application_id = a.id
                    AND pending.sent_at IS NULL
                    AND pending.cancelled_at IS NULL
              )
              AND COALESCE(
                  (SELECT o.confirmation_deadline
                   FROM application_waitlist_offers o
                   WHERE o.application_id = a.id),
                  h.rsvp_deadline
              ) <= $2
        `, eventID, now)
		if err != nil {
			return err
		}
		result.Expired = expired.RowsAffected()

		if _, err := tx.Exec(ctx, `
            UPDATE waitlist_invitation_outbox o
            SET cancelled_at = $2
            FROM applications a
            WHERE o.application_id = a.id
              AND a.hackathon_id = $1
              AND o.sent_at IS NULL AND o.cancelled_at IS NULL
              AND (
                  a.status <> 'accepted'
                  OR NOT EXISTS (
                      SELECT 1 FROM application_waitlist_offers current_offer
                      WHERE current_offer.application_id = o.application_id
                        AND current_offer.offered_at = o.offered_at
                  )
              )
        `, eventID, now); err != nil {
			return err
		}

		if !now.Before(closesAt) {
			return nil
		}

		var reserved int
		if err := tx.QueryRow(ctx, `
            SELECT count(*) FROM applications
            WHERE hackathon_id = $1 AND status IN ('accepted', 'confirmed')
        `, eventID).Scan(&reserved); err != nil {
			return err
		}
		available := *capacity - reserved
		if available <= 0 {
			return nil
		}
		// Limit each sweep to a batch of 50 invitations.
		if available > 50 {
			available = 50
		}

		rows, err := tx.Query(ctx, `
            SELECT a.id,
                   COALESCE(NULLIF(btrim(u.preferred_email), ''),
                            NULLIF(btrim(u.email), '')),
                   COALESCE(NULLIF(a.application->>'firstName', ''),
                            NULLIF(split_part(btrim(u.name), ' ', 1), ''),
                            'Hacker')
            FROM waitlist w
            JOIN applications a
              ON a.user_id = w.user_id AND a.hackathon_id = w.hackathon_id
            JOIN users u ON u.id = a.user_id
            WHERE w.hackathon_id = $1 AND a.status = 'waitlisted'
              AND u.role = 'applicant'
              AND NOT u.is_fake AND NOT a.is_fake
              AND COALESCE(NULLIF(btrim(u.preferred_email), ''),
                           NULLIF(btrim(u.email), '')) IS NOT NULL
            ORDER BY
              CASE WHEN $2::boolean AND w.in_person_joined_at IS NOT NULL
                   THEN 0 ELSE 1 END,
              CASE WHEN $2::boolean THEN w.in_person_joined_at END NULLS LAST,
              w.created_at, w.user_id
            LIMIT $3
            FOR UPDATE OF a
        `, eventID, !now.Before(inPersonAt), available)
		if err != nil {
			return err
		}
		type candidate struct {
			id              uuid.UUID
			recipient, name string
		}
		var candidates []candidate
		for rows.Next() {
			var c candidate
			if err := rows.Scan(&c.id, &c.recipient, &c.name); err != nil {
				rows.Close()
				return err
			}
			candidates = append(candidates, c)
		}
		rowsErr := rows.Err()
		rows.Close()
		if rowsErr != nil {
			return rowsErr
		}

		for _, c := range candidates {
			if _, err := tx.Exec(ctx, `
                UPDATE applications
                SET status = 'accepted', updated_at = now()
                WHERE id = $1
            `, c.id); err != nil {
				return err
			}

			// Use the same clock for the offer and its deadline.
			if _, err := tx.Exec(ctx, `
                UPDATE application_waitlist_offers
                SET offered_at = $2::timestamptz, confirmation_deadline = $2::timestamptz + interval '48 hours'
                WHERE application_id = $1
            `, c.id, now); err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, `
                INSERT INTO waitlist_invitation_outbox
                    (application_id, offered_at, confirmation_deadline,
                     recipient, first_name)
                SELECT application_id, offered_at, confirmation_deadline, $2, $3
                FROM application_waitlist_offers
                WHERE application_id = $1
            `, c.id, c.recipient, c.name); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `
                UPDATE users SET has_seen_new_application_status = false
                WHERE id = (SELECT user_id FROM applications WHERE id = $1)
            `, c.id); err != nil {
				return err
			}
			result.Offered++
		}
		return nil
	})
	if err != nil {
		return WaitlistDispatchResult{}, err
	}
	return result, nil
}
