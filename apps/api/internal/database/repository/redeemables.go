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
	ErrNoActiveHackathon  = errors.New("no active hackathon")
	ErrRedeemableNotFound = errors.New("redeemable not found")
)

type RedeemablesRepository struct {
	db *database.DB
}

func NewRedeemablesRepository(db *database.DB) *RedeemablesRepository {
	return &RedeemablesRepository{
		db: db,
	}
}

func (r *RedeemablesRepository) GetRedeemables(ctx context.Context) (*[]sqlc.GetRedeemablesRow, error) {
	redeemables, err := r.db.Query.GetRedeemables(ctx)
	if err != nil {
		return nil, err
	}
	return &redeemables, nil
}

func (r *RedeemablesRepository) CreateRedeemable(ctx context.Context, params sqlc.CreateRedeemableParams) (*sqlc.Redeemable, error) {
	redeemable, err := r.db.Query.CreateRedeemable(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNoActiveHackathon
	}
	if err != nil {
		return nil, err
	}
	return &redeemable, nil
}

func (r *RedeemablesRepository) DeleteRedeemable(ctx context.Context, redeemableID uuid.UUID) error {
	err := r.db.Query.DeleteRedeemable(ctx, redeemableID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedeemablesRepository) UpdateRedeemable(ctx context.Context, params sqlc.UpdateRedeemableParams) (*sqlc.Redeemable, error) {
	redeemable, err := r.db.Query.UpdateRedeemable(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRedeemableNotFound
	}
	if err != nil {
		return nil, err
	}
	return &redeemable, nil
}

// Returns pgx.ErrNoRows when nothing was inserted; GetRedeemEligibility says why.
func (r *RedeemablesRepository) RedeemForUser(ctx context.Context, params sqlc.RedeemForUserParams) (*sqlc.RedeemForUserRow, error) {
	row, err := r.db.Query.RedeemForUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *RedeemablesRepository) GetRedeemEligibility(ctx context.Context, params sqlc.GetRedeemEligibilityParams) (*sqlc.GetRedeemEligibilityRow, error) {
	row, err := r.db.Query.GetRedeemEligibility(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRedeemableNotFound
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}
