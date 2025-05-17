package handlers

import "github.com/swamphacks/core/apps/api/internal/services"

type Handlers struct {
	OAuthHandler *OAuthHandler
	AuthHandler  *AuthHandler
}

func NewHandlers(authService *services.AuthService) *Handlers {
	return &Handlers{
		OAuthHandler: NewOAuthHandler(authService),
		AuthHandler:  NewAuthHandler(authService),
	}
}
