package handlers

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type Handlers struct {
	Auth *AuthHandler
}

func NewHandlers(authService *services.AuthService, logger zerolog.Logger) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(authService, logger),
	}
}
