package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

const EventRoleContextKey ctxKey = "event_role"

type EventMiddleware struct {
	db     *db.DB
	logger zerolog.Logger
	cfg    *config.Config
}

type UserRoleContext struct {
	UserID  uuid.UUID          `json:"userId"`
	EventID uuid.UUID          `json:"eventId"`
	Role    sqlc.EventRoleType `json:"role"`
}

func NewEventMiddleware(db *db.DB, logger zerolog.Logger, cfg *config.Config) *EventMiddleware {
	return &EventMiddleware{
		db:     db,
		logger: logger.With().Str("middleware", "EventMiddleware").Str("component", "api").Logger(),
		cfg:    cfg,
	}
}

func (m *EventMiddleware) RequireEventRole(role sqlc.EventRoleType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user from context
			userCtx, ok := r.Context().Value(EventRoleContextKey).(*UserRoleContext)
			if !ok {
				m.logger.Warn().Msg("No user context found.")
				response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized."))
				return
			}

			// check if user role matches required role
			if userCtx.Role != role {
				m.logger.Warn().Msgf("User tried to access %s with insufficient permissions as role %s", r.URL.Path, string(userCtx.Role))
				response.SendError(w, http.StatusForbidden, response.NewError("forbidden", "You are forbidden from this resource."))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
