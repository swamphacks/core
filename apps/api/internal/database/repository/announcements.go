package repository

import "github.com/swamphacks/core/apps/api/internal/database"

type AnnouncementsRepository struct {
	db *database.DB
}

func NewAnnouncementsRepository(db *database.DB) *AnnouncementsRepository {
	return &AnnouncementsRepository{db: db}
}