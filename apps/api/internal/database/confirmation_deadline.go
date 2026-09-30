package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (db *DB) GetApplicationConfirmationDeadline(
	ctx context.Context, applicationID uuid.UUID,
) (*time.Time, error) {
	var deadline *time.Time
	err := db.Pool.QueryRow(ctx, `
        SELECT COALESCE(o.confirmation_deadline, h.rsvp_deadline)
        FROM applications a
        JOIN hackathons h ON h.id = a.hackathon_id
        LEFT JOIN application_waitlist_offers o ON o.application_id = a.id
        WHERE a.id = $1
    `, applicationID).Scan(&deadline)
	return deadline, err
}
