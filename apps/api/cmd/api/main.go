package main

import (
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/swamphacks/core/apps/api/internal/api"
	"github.com/swamphacks/core/apps/api/internal/api/handlers"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/storage"
)

//	@title			SwampHacks Test API
//	@version		1.0
//	@description	This is SwampHacks' OpenAPI documentation.
//	@termsOfService	TODO

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	logger := logger.New()
	cfg := config.Load()

	// Init database connection and defer close
	database := db.NewDB(cfg.DatabaseURL)
	defer database.Close()

	// Create transaction manager
	txm := db.NewTransactionManager(database)

	// Create injectable http client
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// Create SES Client for email service
	sesClient := email.NewSESClient(cfg.AWS.AccessKey, cfg.AWS.AccessKeySecret, cfg.AWS.Region, logger)

	// Create asynq client
	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Msg("Failed to parse REDIS_URL")
	}
	taskQueueClient := asynq.NewClient(redisOpt)

	// Create new middleware injectable
	mw := middleware.NewMiddleware(database, logger, cfg)

	// Initialize storage clients
	r2Client, err := storage.NewR2Client(cfg.CF.AccountID, cfg.CF.AccessKeyId, cfg.CF.AccessKeySecret, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create R2 client")
	}

	// Injections into repositories
	userRepo := repository.NewUserRepository(database)
	accountRepo := repository.NewAccountRespository(database)
	sessionRepo := repository.NewSessionRepository(database)
	eventInterestRepo := repository.NewEventInterestRepository(database)
	eventRepo := repository.NewEventRespository(database)
	applicationRepo := repository.NewApplicationRepository(database)
	teamRepo := repository.NewTeamRespository(database)
	teamMemberRepo := repository.NewTeamMemberRespository(database)

	// Injections into services
	authService := services.NewAuthService(userRepo, accountRepo, sessionRepo, txm, client, logger, &cfg.Auth)
	userService := services.NewUserService(userRepo, logger)
	eventInterestService := services.NewEventInterestService(eventInterestRepo, logger)
	eventService := services.NewEventService(eventRepo, userRepo, r2Client, &cfg.CoreBuckets, logger)
	emailService := services.NewEmailService(taskQueueClient, sesClient, logger)
	applicationService := services.NewApplicationService(applicationRepo, eventService, emailService, txm, r2Client, &cfg.CoreBuckets, logger)
	teamService := services.NewTeamService(teamRepo, teamMemberRepo, txm, logger)

	// Injections into handlers
	apiHandlers := handlers.NewHandlers(authService, userService, eventInterestService, eventService, emailService, applicationService, teamService, cfg, logger)

	api := api.NewAPI(&logger, apiHandlers, mw)

	logger.Info().Msgf("API listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, api.Router); err != nil {
		log.Fatal().Msg("Failed to start server.")
	}
}
