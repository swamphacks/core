package repository

import "github.com/swamphacks/core/apps/api/internal/db"

type ApplicationRepository struct {
	db *db.DB
}

func NewApplicationRepository(db *db.DB) *ApplicationRepository {
	return &ApplicationRepository{
		db: db,
	}
}
