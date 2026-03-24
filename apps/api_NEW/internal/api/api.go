package api

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/domains/auth"
	"github.com/swamphacks/core/apps/api/internal/domains/hackathon"
	"github.com/swamphacks/core/apps/api/internal/domains/user"
	"github.com/swamphacks/core/apps/api/internal/logger"
)

func Run() {
	logger := logger.New()
	config := config.LoadConfig()

	db := database.NewDB(config.DatabaseURL)
	defer db.Close()

	txm := database.NewTransactionManager(db)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	mw := mw.NewMiddleware(db, logger, config)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	humaConfig := huma.DefaultConfig("SwampHacks API", "1.0.0")
	humaConfig.CreateHooks = nil

	// TODO: figure out a way to override the default schema name
	// humaConfig.OpenAPI.Components = &huma.Components{
	// 	Schemas: huma.NewMapRegistry("#/components/schemas/", func(t reflect.Type, hint string) string {
	// 		if t.Kind() == reflect.Pointer {
	// 			t = t.Elem()
	// 		}

	// 		bodyField, ok := t.FieldByName("Body")

	// 		if !ok {
	// 			return huma.DefaultSchemaNamer(t, hint)
	// 		}

	// 		if schemaName := bodyField.Tag.Get("schemaName"); schemaName != "" {
	// 			return schemaName
	// 		}

	// 		return huma.DefaultSchemaNamer(t, hint)
	// 	}),
	// }

	api := humachi.New(r, humaConfig)

	// Repositories Setup
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRespository(db)
	sessionRepo := repository.NewSessionRepository(db)
	hackathonRepo := repository.NewHackathonRepository(db)

	// Routes registrations
	authService := auth.NewService(userRepo, accountRepo, sessionRepo, txm, httpClient, logger, &config.Auth)
	authHandler := auth.NewHandler(authService, config, logger)
	auth.RegisterRoutes(authHandler, huma.NewGroup(api, "/auth"), mw)

	userService := user.NewService(userRepo, logger)
	userHandler := user.NewHandler(userService, config, logger)
	user.RegisterRoutes(userHandler, huma.NewGroup(api, "/user"), mw)

	hackathonService := hackathon.NewService(hackathonRepo, logger)
	hackathonHandler := hackathon.NewHandler(hackathonService, config, logger)
	hackathon.RegisterRoutes(hackathonHandler, huma.NewGroup(api, "/hackathon"), mw)

	logger.Info().Msgf("API listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Fatal().Msg("Failed to start server.")
	}
}
