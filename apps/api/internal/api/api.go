package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swamphacks/core/apps/api/internal/api/handlers"
	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type API struct {
	Router     *chi.Mux
	Logger     *zerolog.Logger
	Handlers   *handlers.Handlers
	Middleware *mw.Middleware
}

func NewAPI(logger *zerolog.Logger, handlers *handlers.Handlers, middleware *mw.Middleware) *API {
	api := &API{
		Router:     chi.NewRouter(),
		Logger:     logger,
		Handlers:   handlers,
		Middleware: middleware,
	}

	api.setupRoutes(api.Middleware)

	return api
}

func (api *API) setupRoutes(mw *mw.Middleware) {
	api.Router.Use(middleware.Logger)
	api.Router.Use(middleware.RealIP)
	api.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	api.Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		api.Logger.Trace().Str("method", r.Method).Str("path", r.URL.Path).Msg("Received ping.")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "6") // "pong!\n" is 6 bytes
		if _, err := w.Write([]byte("pong!\n")); err != nil {
			log.Err(err)
		}
	})

	// Auth routes
	api.Router.Route("/auth", func(r chi.Router) {
		r.Get("/callback", api.Handlers.Auth.OAuthCallback)

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequireAuth)
			r.Get("/me", api.Handlers.Auth.GetMe)
			r.Post("/logout", api.Handlers.Auth.Logout)
		})
	})

	// Event routes
	api.Router.Route("/event", func(r chi.Router) {
		r.Post("/", api.Handlers.Event.CreateEvent)
		r.Post("/{eventId}/interest", api.Handlers.EventInterest.AddEmailToEvent)
		r.Get("/{eventId}", api.Handlers.Event.GetEventByID)
		r.Patch("/{eventId}", api.Handlers.Event.UpdateEventById)
	})

	// Protected test routes
	api.Router.Route("/protected", func(r chi.Router) {
		r.Use(mw.Auth.RequireAuth)

		r.Get("/basic", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Welcome, arbitrarily roled user!\n")); err != nil {
				log.Err(err)
			}
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequirePlatformRole(sqlc.AuthUserRoleUser))
			r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("Welcome, user!\n")); err != nil {
					log.Err(err)
				}
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequirePlatformRole(sqlc.AuthUserRoleSuperuser))
			r.Get("/superuser", func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("Welcome, superuser!\n")); err != nil {
					log.Err(err)
				}
			})
		})
	})
}
