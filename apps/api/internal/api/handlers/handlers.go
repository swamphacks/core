package handlers

import "github.com/swamphacks/core/apps/api/internal/services"

type Handlers struct {
	Auth *AuthHandler
}

func NewHandlers(authService *services.AuthService) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(authService),
	}
}
