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
	ErrDuplicateRun        = errors.New("Run already exists in the database")
	ErrRunNotFound         = errors.New("Run not found")
	ErrNoRunsDeleted       = errors.New("No Runs deleted")
	ErrMultipleRunsDeleted = errors.New("Multiple Runs affected by delete query expecting to delete one")
)

type BatRunsRepository struct {
	db *db.DB
}

func NewBatRunsRepository(db *db.DB) *BatRunsRepository {
	return &BatRunsRepository{
		db: db,
	}
}

func (r *BatRunsRepository) AddRun(ctx context.Context, eventId uuid.UUID) (*sqlc.BatRun, error) {
	run, err := r.db.Query.AddRun(ctx, eventId)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil, ErrDuplicateRun
		}
		return nil, err
	}
	return &run, nil
}

func (r *BatRunsRepository) GetRunById(ctx context.Context, id uuid.UUID) (sqlc.BatRun, error) {
	return r.db.Query.GetRunById(ctx, id)
}

func (r *BatRunsRepository) GetRunsByEventId(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetRunsByEventIdRow, error) {
	runs, err := r.db.Query.GetRunsByEventId(ctx, eventId)
	return &runs, err
}

func (r *BatRunsRepository) UpdateRunById(ctx context.Context, params sqlc.UpdateRunByIdParams) error {
	err := r.db.Query.UpdateRunById(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrEventNotFound
		}
	}
	return err
}

func (r *BatRunsRepository) DeleteRunById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteRunById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRunNotFound
		}
	}
	if affectedRows == 0 {
		return ErrNoRunsDeleted
	} else if affectedRows > 1 {
		return ErrMultipleRunsDeleted
	}

	return err
}
