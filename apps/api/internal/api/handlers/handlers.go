package handlers

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type Handlers struct {
	Auth          *AuthHandler
	User          *UserHandler
	EventInterest *EventInterestHandler
	Event         *EventHandler
	Email         *EmailHandler
	Application   *ApplicationHandler
	Teams         *TeamHandler
	Admission     *AdmissionHandler
	Bat           *BatHandler
	Redeemables   *RedeemablesHandler
	Discord       *DiscordHandler
}

func NewHandlers(
	authService *services.AuthService,
	userService *services.UserService,
	eventInterestService *services.EventInterestService,
	eventService *services.EventService,
	emailService *services.EmailService,
	appService *services.ApplicationService,
	teamService *services.TeamService,
	batService *services.BatService,
	redeemablesService *services.RedeemablesService,
	discordService *services.DiscordService,
	cfg *config.Config,
	logger zerolog.Logger,
) *Handlers {
	return &Handlers{
		Auth:          NewAuthHandler(authService, cfg, logger),
		User:          NewUserHandler(userService, logger),
		EventInterest: NewEventInterestHandler(eventInterestService, cfg, logger),
		Event:         NewEventHandler(eventService, cfg, logger),
		Email:         NewEmailHandler(emailService, logger),
		Application:   NewApplicationHandler(appService),
		Teams:         NewTeamHandler(teamService, logger),
		Admission:     NewAdmissionHandler(batService, logger),
		Bat:           NewBatHandler(batService, logger),
		Redeemables:   NewRedeemablesHandler(redeemablesService, cfg, logger),
		Discord:       NewDiscordHandler(discordService, logger),
	}
}
