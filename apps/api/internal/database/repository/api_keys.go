package repository

import "github.com/swamphacks/core/apps/api/internal/database"

type ApiKeysRepository struct {
	db *database.DB
}

func NewApiKeysRepository(db *database.DB) *ApiKeysRepository {
	return &ApiKeysRepository{db: db}
}