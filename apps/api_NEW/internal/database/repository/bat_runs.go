package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrDuplicateRun        = errors.New("Run already exists in the database")
	ErrRunNotFound         = errors.New("Run not found")
	ErrNoRunsDeleted       = errors.New("No Runs deleted")
	ErrMultipleRunsDeleted = errors.New("Multiple Runs affected by delete query expecting to delete one")
)

type BatRunsRepository struct {
	db *database.DB
}

func NewBatRunsRepository(db *database.DB) *BatRunsRepository {
	return &BatRunsRepository{
		db: db,
	}
}

func (r *BatRunsRepository) AddRun(ctx context.Context) (*sqlc.BatRun, error) {
	run, err := r.db.Query.AddBatRun(ctx)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return nil, ErrDuplicateRun
		}
		return nil, err
	}
	return &run, nil
}

func (r *BatRunsRepository) GetRunById(ctx context.Context, id uuid.UUID) (sqlc.BatRun, error) {
	return r.db.Query.GetBatRunById(ctx, id)
}

func (r *BatRunsRepository) GetRuns(ctx context.Context) (*[]sqlc.BatRun, error) {
	runs, err := r.db.Query.GetBatRuns(ctx)
	return &runs, err
}

func (r *BatRunsRepository) UpdateRunById(ctx context.Context, params sqlc.UpdateBatRunByIdParams) error {
	err := r.db.Query.UpdateBatRunById(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRunNotFound
		}
	}
	return err
}

func (r *BatRunsRepository) DeleteRunById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteBatRunById(ctx, id)
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
