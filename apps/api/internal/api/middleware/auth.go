package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
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
const userContextKey ctxKey = "user"

type AuthMiddleware struct {
	db     *db.DB
	logger zerolog.Logger
	cfg    *config.Config
}

// Copied from sqlc.GetActiveSessionUserInfoRow for readability when casting in other files
type UserContext struct {
	UserID    uuid.UUID
	Name      string
	Onboarded bool
	Image     *string
	Role      sqlc.AuthUserRole
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
		cookie, err := r.Cookie("sh_session_id")
		if err != nil {
			m.logger.Warn().Msg("Missing session cookie.")
			response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized"))
			return
		}

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

		m.checkLastUsedAt(w, r, sessionID, user.LastUsedAt)

		ctx := context.WithValue(r.Context(), userContextKey, userContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequirePlatformRole(role sqlc.AuthUserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user from context
			userCtx, ok := r.Context().Value(userContextKey).(UserContext)
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
