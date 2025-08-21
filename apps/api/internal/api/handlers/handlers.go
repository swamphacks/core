package handlers

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type Handlers struct {
	Auth          *AuthHandler
	EventInterest *EventInterestHandler
	Event         *EventHandler
	Email         *EmailHandler
	Application   *ApplicationHandler
}

func NewHandlers(
	authService *services.AuthService,
	eventInterestService *services.EventInterestService,
	eventService *services.EventService,
	emailService *services.EmailService,
	appService *services.ApplicationService,
	cfg *config.Config,
	logger zerolog.Logger,
) *Handlers {
	return &Handlers{
		Auth:          NewAuthHandler(authService, cfg, logger),
		EventInterest: NewEventInterestHandler(eventInterestService, cfg, logger),
		Event:         NewEventHandler(eventService, cfg, logger),
		Email:         NewEmailHandler(emailService, logger),
		Application:   NewApplicationHandler(appService),
	}
}
