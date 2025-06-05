package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/swamphacks/core/apps/api/internal/api"
	"github.com/swamphacks/core/apps/api/internal/api/handlers"
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
	db := db.NewDB(cfg.DatabaseURL)
	defer db.Close()

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

	// Injections into repositories
	userRepo := repository.NewUserRepository(db)

	// Injections into services
	authService := services.NewAuthService(userRepo, client)

	// Injections into handlers
	apiHandlers := handlers.NewHandlers(authService, logger)

	api := api.NewAPI(&logger, apiHandlers)

	logger.Info().Msgf("API listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, api.Router); err != nil {
		log.Fatal().Msg("Failed to start server.")
	}
}
