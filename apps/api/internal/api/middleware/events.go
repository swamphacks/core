package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type EventMiddleware struct {
	db               *db.DB
	logger           zerolog.Logger
	cfg              *config.Config
	eventRespository *repository.EventRepository
}

func NewEventMiddleware(db *db.DB, logger zerolog.Logger, cfg *config.Config) *EventMiddleware {
	return &EventMiddleware{
		db:               db,
		logger:           logger.With().Str("middleware", "EventMiddleware").Str("component", "api").Logger(),
		cfg:              cfg,
		eventRespository: repository.NewEventRespository(db),
	}
}

func (m *EventMiddleware) RequireEventRole(eventRole sqlc.EventRoleType) func(http.Handler) http.Handler {
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

			userEventRole, err := m.eventRespository.GetEventRoleByIds(r.Context(), userCtx.UserID, eventId)
			if err != nil {
				// TODO: Will throw if user doesn't have permission, but how should we handle that with other possible errors?
				m.logger.Warn().Msgf("Error while trying to access %s with insufficient permissions (userId: %s, eventId: %s)", r.URL.Path, userCtx.UserID, eventId, err.Error())
				response.SendError(w, http.StatusForbidden, response.NewError("forbidden", "You are forbidden from this resource."))
				return

			}
			if userEventRole.Role != eventRole {
				m.logger.Warn().Msgf("User tried to access %s with insufficient permissions as role %s", r.URL.Path, string(userCtx.Role))
				response.SendError(w, http.StatusForbidden, response.NewError("forbidden", "You are forbidden from this resource."))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
