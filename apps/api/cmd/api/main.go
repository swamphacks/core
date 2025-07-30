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
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/services"
)

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

	// Create asynq client
	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Msg("Failed to parse REDIS_URL")
	}
	taskQueueClient := asynq.NewClient(redisOpt)

	// Create new middleware injectable
	mw := middleware.NewMiddleware(database, logger, cfg)

	// Injections into repositories
	userRepo := repository.NewUserRepository(database)
	accountRepo := repository.NewAccountRespository(database)
	sessionRepo := repository.NewSessionRepository(database)
	eventInterestRepo := repository.NewEventInterestRepository(database)

	// Injections into services
	authService := services.NewAuthService(userRepo, accountRepo, sessionRepo, txm, client, logger, &cfg.Auth)
	eventInterestService := services.NewEventInterestService(eventInterestRepo, logger)
	emailService := services.NewEmailService(taskQueueClient, logger)

	// Injections into handlers
	apiHandlers := handlers.NewHandlers(authService, eventInterestService, emailService, cfg, logger)

	api := api.NewAPI(&logger, apiHandlers, mw)

	logger.Info().Msgf("API listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, api.Router); err != nil {
		log.Fatal().Msg("Failed to start server.")
	}
}
