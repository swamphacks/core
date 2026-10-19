package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *HackathonRepository) GetAcceptedApplicationConfirmationDeadline(
	ctx context.Context, userID uuid.UUID,
) (*time.Time, error) {
	var deadline *time.Time
	err := r.db.Pool.QueryRow(ctx, `
        SELECT COALESCE(o.confirmation_deadline, h.rsvp_deadline)
        FROM applications a
        JOIN hackathons h ON h.id = a.hackathon_id
        LEFT JOIN application_waitlist_offers o ON o.application_id = a.id
        WHERE a.user_id = $1 AND h.is_active = true
          AND a.status = 'accepted'
    `, userID).Scan(&deadline)
	return deadline, err
}
