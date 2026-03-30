package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type AccountRepository struct {
	db *database.DB
}

func NewAccountRespository(db *database.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) NewTx(tx pgx.Tx) *AccountRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &AccountRepository{
		db: txDB,
	}
}

func (r *AccountRepository) Create(ctx context.Context, params sqlc.CreateAccountParams) (*sqlc.Account, error) {
	account, err := r.db.Query.CreateAccount(ctx, params)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetByProviderAndAccountID(ctx context.Context, params sqlc.GetByProviderAndAccountIDParams) (*sqlc.Account, error) {
	account, err := r.db.Query.GetByProviderAndAccountID(ctx, params)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetUserIDByDiscordAccountID(ctx context.Context, discordAccountID string) (*uuid.UUID, error) {
	userID, err := r.db.Query.GetUserIDByDiscordAccountID(ctx, discordAccountID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrAccountNotFound
		}
		return nil, err
	}

	return &userID, nil
}
