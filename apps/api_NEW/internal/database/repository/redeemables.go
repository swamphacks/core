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

func (r *RedeemablesRepository) CreateRedeemable(ctx context.Context, name string, amount int, maxUserAmount int) (*sqlc.Redeemable, error) {
	params := sqlc.CreateRedeemableParams{
		Name:          name,
		Amount:        int32(amount),
		MaxUserAmount: int32(maxUserAmount),
	}
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
func (r *RedeemablesRepository) UpdateRedeemable(ctx context.Context, redeemableID uuid.UUID, name *string, amount *int, maxUserAmount *int) (*sqlc.Redeemable, error) {
	var amount32 *int32
	if amount != nil {
		v := int32(*amount)
		amount32 = &v
	}
	var maxUserAmount32 *int32
	if maxUserAmount != nil {
		v := int32(*maxUserAmount)
		maxUserAmount32 = &v
	}
	params := sqlc.UpdateRedeemableParams{
		ID:            redeemableID,
		Name:          name,
		Amount:        amount32,
		MaxUserAmount: maxUserAmount32,
	}
	redeemable, err := r.db.Query.UpdateRedeemable(ctx, params)
	if err != nil {
		return nil, err
	}
	return &redeemable, nil
}

func (r *RedeemablesRepository) RedeemRedeemable(ctx context.Context, redeemableID uuid.UUID, userID uuid.UUID) (*sqlc.UserRedemption, error) {
	params := sqlc.RedeemRedeemableParams{
		RedeemableID: redeemableID,
		UserID:       userID,
	}
	redemption, err := r.db.Query.RedeemRedeemable(ctx, params)
	if err != nil {
		return nil, err
	}
	return &redemption, nil
}

func (r *RedeemablesRepository) UpdateRedemption(ctx context.Context, redeemableID uuid.UUID, userID uuid.UUID, amount int) error {
	err := r.db.Query.UpdateRedemption(ctx, sqlc.UpdateRedemptionParams{
		RedeemableID: redeemableID,
		UserID:       userID,
		Amount:       int32(amount),
	})
	if err != nil {
		return err
	}
	return nil
}
