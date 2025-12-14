package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/web"
)

const EventContextKey ctxKey = "event"

type EventMiddleware struct {
	db               *db.DB
	logger           zerolog.Logger
	cfg              *config.Config
	eventRespository *repository.EventRepository
}

type EventContext struct {
	EventRole *sqlc.EventRole
}

func NewEventMiddleware(db *db.DB, logger zerolog.Logger, cfg *config.Config) *EventMiddleware {
	return &EventMiddleware{
		db:               db,
		logger:           logger.With().Str("middleware", "EventMiddleware").Str("component", "api").Logger(),
		cfg:              cfg,
		eventRespository: repository.NewEventRespository(db),
	}
}

func (m *EventMiddleware) AttachEventRoleToContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := r.Context().Value(UserContextKey).(*UserContext)
			if !ok {
				m.logger.Warn().Msg("No event role context found.")
				response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized."))
				return
			}

			eventId, err := web.PathParamToUUID(r, "eventId")
			if err != nil {
				response.SendError(w, http.StatusBadRequest, response.NewError("invalid_event_id", "The event ID is not a valid UUID"))
				return
			}

			eventRole, err := m.eventRespository.GetEventRoleByIds(r.Context(), userCtx.UserID, eventId)

			eventRoleContext := EventContext{
				EventRole: eventRole,
			}

			ctx := context.WithValue(r.Context(), EventContextKey, &eventRoleContext)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *EventMiddleware) RequireEventRole(eventRoles []sqlc.EventRoleType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user from context
			userCtx, ok := r.Context().Value(UserContextKey).(*UserContext)
			if !ok {
				m.logger.Warn().Msg("No event role context found.")
				response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized."))
				return
			}

			eventIdStr := chi.URLParam(r, "eventId")
			eventId, err := uuid.Parse(eventIdStr)
			if err != nil {
				response.SendError(w, http.StatusBadRequest, response.NewError("invalid_event_id", "The event ID is not a valid UUID"))
				return
			}

			if userCtx.Role == sqlc.AuthUserRoleSuperuser {
				next.ServeHTTP(w, r)
				return
			}

			userEventRole, err := m.eventRespository.GetEventRoleByIds(r.Context(), userCtx.UserID, eventId)
			if err != nil {
				// TODO: Will throw if user doesn't have permission, but how should we handle that with other possible errors?
				m.logger.Warn().Msgf("Error while trying to access %s with insufficient permissions (userId: %s, eventId: %s)", r.URL.Path, userCtx.UserID, eventId)
				response.SendError(w, http.StatusForbidden, response.NewError("forbidden", "You are forbidden from this resource."))
				return

			}
			if !slices.Contains(eventRoles, userEventRole.Role) {
				m.logger.Warn().Msgf("User tried to access %s with insufficient permissions (eventRole: %s)", r.URL.Path, string(userEventRole.Role))
				response.SendError(w, http.StatusForbidden, response.NewError("forbidden", "You are forbidden from this resource."))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
