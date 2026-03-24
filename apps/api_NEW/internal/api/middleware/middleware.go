package middleware

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
)

type Middleware struct {
	Auth *AuthMiddleware
}

func NewMiddleware(eventRolesRepo *repository.EventRolesRepository, db *database.DB, logger zerolog.Logger, config *config.Config) *Middleware {
	return &Middleware{
		Auth: NewAuthMiddleware(eventRolesRepo, db, logger, config),
	}
}
