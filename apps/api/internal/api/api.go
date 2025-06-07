package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/handlers"
)

type API struct {
	Router   *chi.Mux
	Logger   *zerolog.Logger
	Handlers *handlers.Handlers
}

func NewAPI(logger *zerolog.Logger, handlers *handlers.Handlers) *API {
	api := &API{
		Router:   chi.NewRouter(),
		Logger:   logger,
		Handlers: handlers,
	}

	api.setupRoutes()

	return api
}

func (api *API) setupRoutes() {
	api.Router.Use(middleware.Logger)
	api.Router.Use(middleware.RealIP)

	api.Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		api.Logger.Trace().Str("method", r.Method).Str("path", r.URL.Path).Msg("Received ping.")

		if _, err := w.Write([]byte("pong!\n")); err != nil {
			return
		}

	})

	api.Router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/callback", api.Handlers.Auth.OAuthCallback)
		})
	})
}
