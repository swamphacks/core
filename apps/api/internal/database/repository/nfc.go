package repository

import (
	"context"

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

func (r *NFCRepository) GetMeals(ctx context.Context) ([]sqlc.GetMealsRow, error) {
	meals, err := r.db.Query.GetMeals(ctx)
	if err != nil {
		return nil, err
	}

	return meals, nil
}

func (r *NFCRepository) GetTshirts(ctx context.Context) ([]sqlc.GetTshirtsRow, error) {
	tshirts, err := r.db.Query.GetTshirts(ctx)
	if err != nil {
		return nil, err
	}

	return tshirts, nil
}

func (r *NFCRepository) CheckinUser(ctx context.Context, params sqlc.CheckinUserParams) (*sqlc.NfcTagsUser, error) {
	user, err := r.db.Query.CheckinUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *NFCRepository) TagToWorkshop(ctx context.Context, params sqlc.TagToWorkshopParams) (int64, error) {
	workshop, err := r.db.Query.TagToWorkshop(ctx, params)
	if err != nil {
		return 0, err
	}

	return workshop, nil
}

func (r *NFCRepository) TagToRedeemable(ctx context.Context, params sqlc.TagToRedeemableParams) error {
	return r.db.Query.TagToRedeemable(ctx, params)
}

func (r *NFCRepository) GetSocials(ctx context.Context) ([]sqlc.GetSocialsRow, error) {
	socials, err := r.db.Query.GetSocials(ctx)
	if err != nil {
		return nil, err
	}

	return socials, nil
}

func (r *NFCRepository) GetWorkshops(ctx context.Context) ([]sqlc.GetWorkshopsRow, error) {
	workshops, err := r.db.Query.GetWorkshops(ctx)
	if err != nil {
		return nil, err
	}

	return workshops, nil
}
