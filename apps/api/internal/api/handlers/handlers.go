package handlers

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type Handlers struct {
	Auth        *AuthHandler
	MailingList *MailingListHandler
}

func NewHandlers(authService *services.AuthService, mailingListService *services.MailingListService, cfg *config.Config, logger zerolog.Logger) *Handlers {
	return &Handlers{
		Auth:        NewAuthHandler(authService, cfg, logger),
		MailingList: NewMailingListHandler(mailingListService, cfg, logger),
	}
}
