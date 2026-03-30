package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type BatRunsRepository struct {
	db *database.DB
}

func NewBatRunsRepository(db *database.DB) *BatRunsRepository {
	return &BatRunsRepository{
		db: db,
	}
}

func (r *BatRunsRepository) AddRun(ctx context.Context, hackathonID string) (*sqlc.BatRun, error) {
	run, err := r.db.Query.AddBatRun(ctx, hackathonID)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return nil, database.ErrDuplicateRun
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
			return database.ErrRunNotFound
		}
	}
	return err
}

func (r *BatRunsRepository) DeleteRunById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteBatRunById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.ErrRunNotFound
		}
	}
	if affectedRows == 0 {
		return database.ErrNoRunsDeleted
	} else if affectedRows > 1 {
		return database.ErrMultipleRunsDeleted
	}

	return err
}
