package handlers

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type Handlers struct {
	Auth          *AuthHandler
	EventInterest *EventInterestHandler
	Email         *EmailHandler
}

func NewHandlers(authService *services.AuthService, eventInterestService *services.EventInterestService, emailService *services.EmailService, cfg *config.Config, logger zerolog.Logger) *Handlers {
	return &Handlers{
		Auth:          NewAuthHandler(authService, cfg, logger),
		EventInterest: NewEventInterestHandler(eventInterestService, cfg, logger),
		Email:         NewEmailHandler(emailService, logger),
	}
}
