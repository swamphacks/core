package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
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
	if err != nil {
		return nil, err
	}
	return &redeemable, nil
}

func (r *RedeemablesRepository) RedeemRedeemable(ctx context.Context, params sqlc.RedeemRedeemableParams) (*sqlc.UserRedemption, error) {
	redemption, err := r.db.Query.RedeemRedeemable(ctx, params)
	if err != nil {
		return nil, err
	}
	return &redemption, nil
}

func (r *RedeemablesRepository) UpdateRedemption(ctx context.Context, params sqlc.UpdateRedemptionParams) error {
	err := r.db.Query.UpdateRedemption(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
