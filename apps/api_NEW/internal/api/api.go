package api

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/domains/application"
	"github.com/swamphacks/core/apps/api/internal/domains/auth"
	"github.com/swamphacks/core/apps/api/internal/domains/bat"
	"github.com/swamphacks/core/apps/api/internal/domains/email"
	"github.com/swamphacks/core/apps/api/internal/domains/hackathon"
	"github.com/swamphacks/core/apps/api/internal/domains/redeemables"
	"github.com/swamphacks/core/apps/api/internal/domains/team"
	"github.com/swamphacks/core/apps/api/internal/domains/user"
	"github.com/swamphacks/core/apps/api/internal/emailutils"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/storage"
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

	// Create SES Client for email service
	sesClient := emailutils.NewSESClient(config.AWS.AccessKey, config.AWS.AccessKeySecret, config.AWS.Region, logger)

	// Create asynq client
	redisOpt, err := asynq.ParseRedisURI(config.RedisURL)
	if err != nil {
		logger.Fatal().Msg("Failed to parse REDIS_URL")
	}
	taskQueueClient := asynq.NewClient(redisOpt)

	r2Client, err := storage.NewR2Client(config.CF.AccountID, config.CF.AccessKeyId, config.CF.AccessKeySecret, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create R2 client")
	}

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
	eventRolesRepo := repository.NewEventRolesRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	teamRepo := repository.NewTeamRespository(db)
	teamMemberRepo := repository.NewTeamMemberRespository(db)
	teamJoinRequestRepo := repository.NewTeamJoinRequestRepository(db)
	redeemablesRepo := repository.NewRedeemablesRepository(db)
	batRunsRepo := repository.NewBatRunsRepository(db)

	mw := mw.NewMiddleware(eventRolesRepo, db, logger, config)

	// Routes registrations
	authService := auth.NewService(userRepo, accountRepo, sessionRepo, txm, httpClient, logger, &config.Auth)
	authHandler := auth.NewHandler(authService, config, logger)
	auth.RegisterRoutes(authHandler, huma.NewGroup(api, "/auth"), mw)

	userService := user.NewService(userRepo, eventRolesRepo, logger)
	userHandler := user.NewHandler(userService, config, logger)
	user.RegisterRoutes(userHandler, huma.NewGroup(api, "/user"), mw)

	hackathonService := hackathon.NewService(hackathonRepo, eventRolesRepo, logger)
	hackathonHandler := hackathon.NewHandler(hackathonService, config, logger)
	hackathon.RegisterRoutes(hackathonHandler, huma.NewGroup(api, "/hackathon"), mw)

	emailService := email.NewEmailService(hackathonRepo, userRepo, taskQueueClient, sesClient, r2Client, logger, config)
	batService := bat.NewBatService(applicationRepo, hackathonRepo, userRepo, batRunsRepo, emailService, txm, taskQueueClient, nil, config, logger)
	applicationService := application.NewService(applicationRepo, userRepo, hackathonRepo, eventRolesRepo, txm, r2Client, &config.CoreBuckets, nil, emailService, batService, config, logger)
	applicationHandler := application.NewHandler(applicationService, batService, config, logger)
	application.RegisterRoutes(applicationHandler, huma.NewGroup(api, "/application"), mw)

	teamService := team.NewService(teamRepo, teamMemberRepo, teamJoinRequestRepo, hackathonRepo, eventRolesRepo, txm, logger)
	teamHandler := team.NewHandler(teamService, logger)
	team.RegisterRoutes(teamHandler, huma.NewGroup(api, "/team"), mw)

	redeemablesService := redeemables.NewService(redeemablesRepo, logger)
	redeemablesHandler := redeemables.NewHandler(redeemablesService, config, logger)
	redeemables.RegisterRoutes(redeemablesHandler, huma.NewGroup(api, "/redeemables"), mw)

	logger.Info().Msgf("API listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Fatal().Msg("Failed to start server.")
	}
}
