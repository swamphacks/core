package repository

import (
	"context"
	"errors"

	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrDuplicateResult = errors.New("result already exists in the database")
)

type BatResultsRepository struct {
	db *db.DB
}

func NewBatResultsRepository(db *db.DB) *EventInterestRepository {
	return &EventInterestRepository{
		db: db,
	}
}

func (r *BatResultsRepository) AddResult(ctx context.Context, params sqlc.AddResultParams) (*sqlc.BatResult, error) {
	result, err := r.db.Query.AddResult(ctx, params)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil, ErrDuplicateResult
		}
		return nil, err
	}
	return &result, nil
}

func (r *BatResultsRepository) GetResultsByEvent(ctx context.Context, eventId uuid.UUID) (*[]sqlc.BatResult, error) {
	results, err := r.db.Query.GetResultsByEventId(ctx, eventId)
	return &results, err
}
