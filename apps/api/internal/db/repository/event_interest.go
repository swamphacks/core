package repository

import (
	"context"
	"errors"

	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrDuplicateEmails = errors.New("email already exists in the database")
)

type EventInterestRepository struct {
	db *db.DB
}

func NewEventInterestRepository(db *db.DB) *EventInterestRepository {
	return &EventInterestRepository{
		db: db,
	}
}

func (r *EventInterestRepository) AddEmail(ctx context.Context, params sqlc.AddEmailParams) (*sqlc.EventInterestSubmission, error) {
	interestSubmission, err := r.db.Query.AddEmail(ctx, params)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil, ErrDuplicateEmails
		}

		return nil, err
	}

	return &interestSubmission, nil
}
