package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (r *AccountRepository) NewTx(tx pgx.Tx) *AccountRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &AccountRepository{
		db: txDB,
	}
}

func (r *AccountRepository) Create(ctx context.Context, params sqlc.CreateAccountParams) (*sqlc.AuthAccount, error) {
	account, err := r.db.Query.CreateAccount(ctx, params)
	return &account, err
}

func (r *AccountRepository) GetByProviderAndAccountID(ctx context.Context, params sqlc.GetByProviderAndAccountIDParams) (*sqlc.AuthAccount, error) {
	account, err := r.db.Query.GetByProviderAndAccountID(ctx, params)
	return &account, err
}
