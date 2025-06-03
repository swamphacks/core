package main

import (
	"net/http"

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

	// Injections into repositories
	userRepo := repository.NewUserRepository(db)

	// Injections into services
	authService := services.NewAuthService(userRepo)

	// Injections into handlers
	apiHandlers := handlers.NewHandlers(authService, logger)

	api := api.NewAPI(&logger, apiHandlers)

	logger.Info().Msgf("API listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, api.Router); err != nil {
		log.Fatal().Msg("Failed to start server.")
	}
}
