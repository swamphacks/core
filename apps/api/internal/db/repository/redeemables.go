package repository

import "github.com/swamphacks/core/apps/api/internal/db"

type RedeemablesRepository struct {
	db *db.DB
}

func NewRedeemablesRepository(db *db.DB) *RedeemablesRepository {
	return &RedeemablesRepository{
		db: db,
	}
}
