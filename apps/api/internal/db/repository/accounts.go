package repository

import (
	"context"

	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type AccountRepository struct {
	db *db.DB
}

func NewAccountRespository(db *db.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) GetByProviderAndAccountID(ctx context.Context, params sqlc.GetByProviderAndAccountIDParams) (*sqlc.AuthAccount, error) {
	account, err := r.db.Query.GetByProviderAndAccountID(ctx, params)
	return &account, err
}
