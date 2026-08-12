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
	ErrApiKeyNotFound = errors.New("api key not found")
)

type ApiKeysRepository struct {
	db *database.DB
}

func NewApiKeysRepository(db *database.DB) *ApiKeysRepository {
	return &ApiKeysRepository{db: db}
}

func (r *ApiKeysRepository) ListApiKeys(
	ctx context.Context,
) ([]sqlc.ListApiKeysRow, error) {
	return r.db.Query.ListApiKeys(ctx)
}

func (r *ApiKeysRepository) GetApiKeyByID(
	ctx context.Context,
	id uuid.UUID,
) (*sqlc.GetApiKeyByIdRow, error) {
	apiKey, err := r.db.Query.GetApiKeyById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrApiKeyNotFound
	}
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeysRepository) CreateApiKey(
	ctx context.Context,
	params sqlc.CreateApiKeyParams,
) (*sqlc.CreateApiKeyRow, error) {
	apiKey, err := r.db.Query.CreateApiKey(ctx, params)
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeysRepository) DeleteApiKeyByID(
	ctx context.Context,
	id uuid.UUID,
) error {
	err := r.db.Query.DeleteApiKey(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrApiKeyNotFound
	}
	return err
}