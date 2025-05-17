package api

import (
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/handlers"
	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type API struct {
	Router *chi.Mux
	Logger *zerolog.Logger
	Wg     *sync.WaitGroup
}

func NewAPI(log *zerolog.Logger, db *db.DB) {
	// Obviously don't make a new client here..., this is just for testing
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ur := repository.NewUserRepository(db)
	as := services.NewAuthService(ur)
	hl := handlers.NewHandlers(as)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/oauth/verify", hl.OAuthHandler.VerifyOAuth)

		r.Route("/auth", func(r chi.Router) {
			r.Use(mw.SessionMiddleware(rdb, db))
			r.Get("/me", hl.AuthHandler.HandleGetMe)
		})
	})
}
