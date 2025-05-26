package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swamphacks/core/apps/api/internal/api"
	"github.com/swamphacks/core/apps/api/internal/api/handlers"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/services"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg("Environment variables failed to load from .env file")
		os.Exit(1)
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	dbURL, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		logger.Error().Msg("DATABASE_URL not found")
		os.Exit(1)
	}

	// Init database connection
	db := db.NewDB(dbURL)

	// Injections into repositories
	userRepo := repository.NewUserRepository(db)

	// Injections into services
	authService := services.NewAuthService(userRepo)

	// Injections into handlers
	apiHandlers := handlers.NewHandlers(authService)

	api := api.NewAPI(&logger, apiHandlers)

	if err := http.ListenAndServe(":8080", api.Router); err != nil {
		log.Fatal().Msg("Failed to start server.")
	}
}
