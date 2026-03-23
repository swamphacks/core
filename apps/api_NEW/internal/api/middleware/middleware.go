package middleware

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
)

type Middleware struct {
	Auth *AuthMiddleware
}

func NewMiddleware(db *database.DB, logger zerolog.Logger, config *config.Config) *Middleware {
	return &Middleware{
		Auth: NewAuthMiddleware(db, logger, config),
	}
}
