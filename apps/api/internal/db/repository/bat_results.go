package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrDuplicateResult        = errors.New("result already exists in the database")
	ErrResultNotFound         = errors.New("result not found")
	ErrNoResultsDeleted       = errors.New("no results deleted")
	ErrMultipleResultsDeleted = errors.New("multiple results affected by delete query expecting to delete one")
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

func (r *BatResultsRepository) GetResultsByEvent(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetResultsByEventIdRow, error) {
	results, err := r.db.Query.GetResultsByEventId(ctx, eventId)
	return &results, err
}

func (r *BatResultsRepository) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteResultById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrResultNotFound
		}
	}
	if affectedRows == 0 {
		return ErrNoResultsDeleted
	} else if affectedRows > 1 {
		return ErrMultipleResultsDeleted
	}

	return err
}
