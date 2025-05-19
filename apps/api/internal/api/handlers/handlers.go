package handlers

import "github.com/swamphacks/core/apps/api/internal/services"

type Handlers struct {
	AuthHandler *AuthHandler
}

func NewHandlers(authService *services.AuthService) *Handlers {
	return &Handlers{
		AuthHandler: NewAuthHandler(authService),
	}
}
