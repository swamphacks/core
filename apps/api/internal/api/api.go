package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/swamphacks/core/apps/api/internal/api/handlers"
	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
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
	var (
		// Both requireXXRole functions automatically allow superusers
		ensureSuperuser  = mw.Auth.RequirePlatformRole([]sqlc.AuthUserRole{sqlc.AuthUserRoleSuperuser})
		ensureEventAdmin = mw.Event.RequireEventRole([]sqlc.EventRoleType{sqlc.EventRoleTypeAdmin})
	)

	AllowedOrigins := config.Load().AllowedOrigins

	api.Router.Use(middleware.Logger)
	api.Router.Use(middleware.RealIP)
	api.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api.Router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "SwampHacks API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})

	// Health check
	api.Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		api.Logger.Trace().Str("method", r.Method).Str("path", r.URL.Path).Msg("Received ping.")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "6")
		if _, err := w.Write([]byte("pong!\n")); err != nil {
			log.Err(err)
		}
	})

	// --- Auth routes ---
	api.Router.Route("/auth", func(r chi.Router) {
		r.Get("/callback", api.Handlers.Auth.OAuthCallback)

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequireAuth)
			r.Get("/me", api.Handlers.Auth.GetMe)
			r.Post("/logout", api.Handlers.Auth.Logout)
		})
	})

	// --- User routes ---
	api.Router.Route("/users", func(r chi.Router) {
		r.Use(mw.Auth.RequireAuth)
		r.Get("/", api.Handlers.User.GetUsers)
		r.Get("/me", api.Handlers.User.GetProfile)
		r.Patch("/me", api.Handlers.User.UpdateUser)
		r.Patch("/me/email-consent", api.Handlers.User.UpdateEmailConsent)
		r.Patch("/me/onboarding", api.Handlers.User.CompleteOnboarding)
	})

	// --- Event routes ---
	api.Router.Route("/events", func(r chi.Router) {

		// Superuser-only
		r.With(mw.Auth.RequireAuth, ensureSuperuser).Post("/", api.Handlers.Event.CreateEvent)

		// Authenticated
		r.With(mw.Auth.RequireAuth).Get("/", api.Handlers.Event.GetEvents)

		r.Post("/{eventId}/interest", api.Handlers.EventInterest.AddEmailToEvent) // Unprotected

		// Event-specific routes
		r.Route("/{eventId}", func(r chi.Router) {
			r.Use(mw.Auth.RequireAuth) // routes below this are protected

			// General access
			r.Get("/", api.Handlers.Event.GetEventByID)
			r.Get("/role", api.Handlers.Event.GetEventRole)

			// Admin-only
			r.With(ensureEventAdmin).Patch("/", api.Handlers.Event.UpdateEventById)
			r.With(ensureEventAdmin).Post("/banner", api.Handlers.Event.UploadEventBanner)
			r.With(ensureEventAdmin).Delete("/banner", api.Handlers.Event.DeleteBanner)
			r.With(ensureEventAdmin).Get("/staff", api.Handlers.Event.GetEventStaffUsers)
			r.With(ensureEventAdmin).Post("/roles", api.Handlers.Event.AssignEventRole)
			r.With(ensureEventAdmin).Delete("/roles/{userId}", api.Handlers.Event.RevokeEventRole)
			r.With(ensureEventAdmin).Post("/roles/batch", api.Handlers.Event.BatchAssignEventRoles)

			// Superuser-only
			r.With(ensureSuperuser).Delete("/", api.Handlers.Event.DeleteEventById)

			// Application routes
			r.Route("/application", func(r chi.Router) {
				r.Use(mw.Auth.RequireAuth)
				r.Get("/", api.Handlers.Application.GetApplicationByUserAndEventID)
				r.Post("/submit", api.Handlers.Application.SubmitApplication)
				r.Post("/save", api.Handlers.Application.SaveApplication)
			})
		})
	})

	// Email routes
	api.Router.Route("/email", func(r chi.Router) {
		r.Post("/queue", api.Handlers.Email.QueueEmail)
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
			r.Use(mw.Auth.RequirePlatformRole([]sqlc.AuthUserRole{sqlc.AuthUserRoleUser}))
			r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("Welcome, user!\n")); err != nil {
					log.Err(err)
				}
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.RequirePlatformRole([]sqlc.AuthUserRole{sqlc.AuthUserRoleSuperuser}))
			r.Get("/superuser", func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("Welcome, superuser!\n")); err != nil {
					log.Err(err)
				}
			})
		})
	})
}
