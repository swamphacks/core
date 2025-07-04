package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

// MailingListRepository provides access to mailing list related queries.
// It acts as a thin wrapper around the auto-generated sqlc queries.
type MailingListRepository struct {
	db *db.DB
}

// NewMailingListRepository constructs a new MailingListRepository with the given *db.DB.
func NewMailingListRepository(db *db.DB) *MailingListRepository {
	return &MailingListRepository{
		db: db,
	}
}

// AddEmail adds a new email entry to the mailing list.
func (r *MailingListRepository) AddEmail(ctx context.Context, params sqlc.AddEmailParams) (sqlc.MailingListEmail, error) {
	// The return type is updated from sqlc.Email to sqlc.MailingListEmail
	return r.db.Query.AddEmail(ctx, params)
}

// UpdateEmailByID updates an existing email address by its primary key ID.
func (r *MailingListRepository) UpdateEmailByID(ctx context.Context, params sqlc.UpdateEmailByIDParams) (sqlc.MailingListEmail, error) {
	// Updated to call the new sqlc function and pass the correct params struct
	return r.db.Query.UpdateEmailByID(ctx, params)
}

// DeleteEmailByID deletes an email entry by its primary key ID.
func (r *MailingListRepository) DeleteEmailByID(ctx context.Context, id uuid.UUID) error {
	// The underlying sqlc function name was already correct
	return r.db.Query.DeleteEmailByID(ctx, id)
}

// DeleteEmailByUserAndEvent deletes an email entry using the logical composite key (user_id and event_id).
func (r *MailingListRepository) DeleteEmailByUserAndEvent(ctx context.Context, params sqlc.DeleteEmailByUserAndEventParams) error {
	// This is a new method wrapping the new sqlc function
	return r.db.Query.DeleteEmailByUserAndEvent(ctx, params)
}

// GetEmailsByEvent retrieves all emails for a specific event.
func (r *MailingListRepository) GetEmailsByEvent(ctx context.Context, eventID uuid.UUID) ([]sqlc.MailingListEmail, error) {
	// This is a new method wrapping the new sqlc function
	return r.db.Query.GetEmailsByEvent(ctx, eventID)
}

// GetEmailsByUser retrieves all mailing list entries for a specific user across all events.
func (r *MailingListRepository) GetEmailsByUser(ctx context.Context, userID uuid.UUID) ([]sqlc.MailingListEmail, error) {
	// This is a new method wrapping the new sqlc function
	return r.db.Query.GetEmailsByUser(ctx, userID)
}

// GetEmailByUserAndEvent retrieves a single email entry for a specific user and event.
func (r *MailingListRepository) GetEmailByUserAndEvent(ctx context.Context, params sqlc.GetEmailByUserAndEventParams) (sqlc.MailingListEmail, error) {
	// This is a new method wrapping the new sqlc function
	return r.db.Query.GetEmailByUserAndEvent(ctx, params)
}
