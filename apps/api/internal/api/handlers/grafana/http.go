package grafana

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	auth "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type handler struct {
	logger zerolog.Logger
}

func RegisterRoutes(grafanaHandler *handler, mw *middleware.Middleware, router *chi.Mux) {
	grafanaURL, err := url.Parse("http://grafana:3000")
	if err != nil {
		grafanaHandler.logger.Error().Err(err).Msg("Failed to parse grafana URL")
		return
	}
	grafanaProxy := httputil.NewSingleHostReverseProxy(grafanaURL)

	handler := mw.Auth.RequireAuth(
		mw.Auth.RequireRoles([]sqlc.UserRole{sqlc.UserRoleAdmin})(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set header to mark identity for Grafana
				if userCtx, ok := r.Context().Value(auth.UserContextKey).(*auth.UserContext); ok {
					r.Header.Set("X-WEBAUTH-USER", userCtx.UserID.String())
				}

				grafanaProxy.ServeHTTP(w, r)
			}),
		),
	)

	router.Handle("/grafana", handler)
	router.Handle("/grafana/*", handler)
}

func NewHandler(logger zerolog.Logger) *handler {
	return &handler{
		logger: logger.With().Str("handler", "GrafanaHandler").Logger(),
	}
}
