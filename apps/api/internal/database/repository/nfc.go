package repository

import (
	"context"


	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type NFCRepository struct {
	db *database.DB
}

func NewNFCRepository(db *database.DB) *NFCRepository {
	return &NFCRepository{
		db: db,
	}
}



func (r *NFCRepository) GetMeals(ctx context.Context) ([]sqlc.getMealsRow, error) {
	meals, err := r.db.Query.getMeals(ctx)
	if err != nil {
		return nil, err
	}

	return meals, nil
}

func (r *NFCRepository) GetTshirts(ctx context.Context) ([]sqlc.getTshirtsRow, error) {
	tshirts, err := r.db.Query.getTshirts(ctx)
	if err != nil {
		return nil, err
	}

	return tshirts, nil
}

func (r *NFCRepository) CheckinUser(ctx context.Context, params sqlc.checkinUserParams) (*sqlc.NfcTagsUser, error) {
	user, err := r.db.Query.checkinUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *NFCRepository) TagToWorkshop(ctx context.Context, params sqlc.TagToWorkshopParams) (int64, error) {
	workshop, err := r.db.Query.tagToWorkshop(ctx, params)
	if err != nil {
		return 0, err
	}

	return &workshop, nil
}

func (r *NFCRepository) TagToRedeemable(ctx context.Context, params sqlc.TagToRedeemableParams) (error) {
	return r.db.Query.tagToRedeemable(ctx, params)
}

func (r *NFCRepository) GetSocials(ctx context.Context) ([]sqlc.getSocialsRow, error) {
	socials, err := r.db.Query.getSocials(ctx)
	if err != nil {
		return nil, err
	}

	return socials, nil
}

func (r *NFCRepository) GetWorkshops(ctx context.Context) ([]sqlc.getWorkshopsRow, error) {
	workshops, err := r.db.Query.getWorkshops(ctx)
	if err != nil {
		return nil, err
	}

	return workshops, nil
}

