package repository

import (
	"context"

	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type EventInterestsRepository struct {
	db *database.DB
}

func NewEventInterestsRepository(db *database.DB) *EventInterestsRepository {
	return &EventInterestsRepository{
		db: db,
	}
}

func (r *EventInterestsRepository) AddEmail(ctx context.Context, params sqlc.AddEmailParams) (*sqlc.InterestSubmission, error) {
	interestSubmission, err := r.db.Query.AddEmail(ctx, params)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return nil, database.ErrDuplicateEmails
		}

		return nil, err
	}

	return &interestSubmission, nil
}
