package repository

import (
	"context"
	"errors"

	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrDuplicateEmails = errors.New("email already exists in the database")
)

type InterestMailingListRepository struct {
	db *database.DB
}

func NewEventInterestRepository(db *database.DB) *InterestMailingListRepository {
	return &InterestMailingListRepository{
		db: db,
	}
}

func (r *InterestMailingListRepository) AddEmail(ctx context.Context, params sqlc.AddEmailParams) (*sqlc.InterestSubmission, error) {
	interestSubmission, err := r.db.Query.AddEmail(ctx, params)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return nil, ErrDuplicateEmails
		}

		return nil, err
	}

	return &interestSubmission, nil
}
