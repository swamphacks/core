package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strings"
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

// UserContext represents the authenticated user in API requests.
// @Description Information about the current user session.
type UserContext struct {
	// Unique identifier for the user
	UserID uuid.UUID `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`

	// Primary email address (nullable)
	Email *string `json:"email" example:"user@example.com"`

	// Preferred email address for communications
	PreferredEmail *string `json:"preferredEmail" example:"user.alt@example.com"`

	// Full display name
	Name string `json:"name" example:"Jane Doe"`

	// Whether the user completed onboarding
	Onboarded bool `json:"onboarded" example:"true"`

	// Optional profile image URL
	Image *string `json:"image,omitempty" example:"https://cdn.example.com/avatar.png" extensions:"nullable"`

	// Role assigned to the user
	Role sqlc.AuthUserRole `json:"role"`

	// Whether the user agreed to receive emails
	EmailConsent bool `json:"emailConsent" example:"false"`
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

func (m *AuthMiddleware) RequireMobileAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Trace().Msg("Incoming mobile request")

		auth := r.Header.Get("Authorization")
		if auth == "" {
			m.logger.Warn().Msg("Authorization header is empty.")
			response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized"))
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 {
			m.logger.Warn().Msg("Authorization header is malformed > 2 parts")
			response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized"))
			return
		}

		if parts[0] != "Key" || parts[1] != m.cfg.MobileAuthKey {
			m.logger.Warn().Msg("Authorization header is not a valid value")
			response.SendError(w, http.StatusUnauthorized, response.NewError("no_auth", "You are not authorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
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
			UserID:         user.UserID,
			Name:           user.Name,
			Email:          user.Email,
			PreferredEmail: user.PreferredEmail,
			Image:          user.Image,
			Onboarded:      user.Onboarded,
			Role:           user.Role,
			EmailConsent:   user.EmailConsent,
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

			if userCtx.Role == sqlc.AuthUserRoleSuperuser {
				next.ServeHTTP(w, r)
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
