package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/cookie"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type ctxKey string

// Use this variable to retrieve the user object later from context!
const UserContextKey ctxKey = "user"
const SessionContextKey ctxKey = "session"

type AuthMiddleware struct {
	db     *db.DB
	logger zerolog.Logger
	cfg    *config.Config
}

type UserContext struct {
	UserID    uuid.UUID         `json:"userId"`
	Name      string            `json:"name"`
	Onboarded bool              `json:"onboarded"`
	Image     *string           `json:"image,omitempty"` // omit if nil
	Role      sqlc.AuthUserRole `json:"role"`
}

type SessionContext struct {
	SessionID uuid.UUID
}

func NewAuthMiddleware(db *db.DB, logger zerolog.Logger, cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		db:     db,
		logger: logger.With().Str("middleware", "AuthMiddleware").Str("component", "api").Logger(),
		cfg:    cfg,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Trace().Msg("Checking auth status")
		cookie, err := r.Cookie("sh_session_id")
		if err != nil {
			m.logger.Warn().Msg("Missing session cookie.")
			response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized"))
			return
		}

		m.logger.Trace().Msg("Checking session id")
		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			m.logger.Warn().Msg("Session cookie was unparsable into UUID")
			response.SendError(w, http.StatusBadRequest, response.NewError("bad_cookie", "The cookie went bad... blegh"))
			return
		}

		user, err := m.db.Query.GetActiveSessionUserInfo(r.Context(), sessionID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			m.logger.Info().Msg("Session is no longer valid or does not exist.")
			response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized"))
			return
		} else if err != nil {
			m.logger.Err(err).Msg("Something went wrong getting active session user info.")
			response.SendError(w, http.StatusInternalServerError, response.NewError("internal_err", "Something went horrible wrong!"))
			return
		}

		userContext := UserContext{
			UserID:    user.UserID,
			Name:      user.Name,
			Image:     user.Image,
			Onboarded: user.Onboarded,
			Role:      user.Role,
		}

		sessionContext := SessionContext{
			SessionID: sessionID,
		}

		m.checkLastUsedAt(w, r, sessionID, user.LastUsedAt)

		ctx := context.WithValue(r.Context(), UserContextKey, &userContext)
		ctx = context.WithValue(ctx, SessionContextKey, &sessionContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequirePlatformRole(roles []sqlc.AuthUserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user from context
			userCtx, ok := r.Context().Value(UserContextKey).(*UserContext)
			if !ok {
				m.logger.Warn().Msg("No user context found.")
				response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized."))
				return
			}

			// check if user role matches required role
			if !slices.Contains(roles, userCtx.Role) {
				m.logger.Warn().Msgf("User tried to access %s with insufficient permissions as role %s", r.URL.Path, string(userCtx.Role))
				response.SendError(w, http.StatusForbidden, response.NewError("forbidden", "You are forbidden from this resource."))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// If lastUsedAt is more than a day ago from now, update using TouchSession (rolling session expiration)
// Also make sure to reflect on the cookie!
func (m *AuthMiddleware) checkLastUsedAt(w http.ResponseWriter, r *http.Request, sessionID uuid.UUID, lastUsedAt time.Time) {

	// Was last used within a day, do not update
	if lastUsedAt.After(time.Now().Add(-24 * time.Hour)) {
		return
	}

	newExpiration := time.Now().AddDate(0, 1, 0) // In 30 days
	err := m.db.Query.TouchSession(r.Context(), sqlc.TouchSessionParams{
		ID:        sessionID,
		ExpiresAt: newExpiration,
	})
	if err != nil {
		m.logger.Err(err).Msg("Failed to update session expiration.")
		return
	}

	// Update user's cookie
	cookie.SetSessionCookie(w, sessionID, newExpiration, m.cfg.Cookie)
}
