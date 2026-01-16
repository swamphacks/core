package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type RedeemablesRepository struct {
	db *db.DB
}

func NewRedeemablesRepository(db *db.DB) *RedeemablesRepository {
	return &RedeemablesRepository{
		db: db,
	}
}

func (r *RedeemablesRepository) GetRedeemablesByEventID(ctx context.Context, eventID uuid.UUID) (*[]sqlc.GetRedeemablesByEventIDRow, error) {
	redeemables, err := r.db.Query.GetRedeemablesByEventID(ctx, eventID)
	if err != nil {
		return nil, err // TODO: more robust error handling
	}
	return &redeemables, nil
}

func (r *RedeemablesRepository) CreateRedeemable(ctx context.Context, params sqlc.CreateRedeemableParams) (*sqlc.Redeemable, error) {
	redeemable, err := r.db.Query.CreateRedeemable(ctx, params)
	if err != nil {
		return nil, err // TODO: more robust error handling
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
