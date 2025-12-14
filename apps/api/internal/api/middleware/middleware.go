package middleware

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
)

type Middleware struct {
	Auth  *AuthMiddleware
	Event *EventMiddleware
}

func NewMiddleware(db *db.DB, logger zerolog.Logger, cfg *config.Config) *Middleware {
	return &Middleware{
		Auth:  NewAuthMiddleware(db, logger, cfg),
		Event: NewEventMiddleware(db, logger, cfg),
	}
}
