package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/cookie"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	cfg         *config.Config
	logger      zerolog.Logger
}

func NewAuthHandler(authService *services.AuthService, cfg *config.Config, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		logger:      logger.With().Str("handler", "AuthHandler").Str("component", "auth").Logger(),
	}
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetMe(r.Context())
	if err != nil {
		res.SendError(w, http.StatusNotFound, res.NewError("no_user", "Your profile could not be loaded."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went seriously wrong."))
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.authService.Logout(r.Context())
	if err != nil && errors.Is(err, services.ErrFetchSessionContextFailed) {
		res.SendError(w, http.StatusUnauthorized, res.NewError("no_auth", "You are not authorized."))
		return
	} else if err != nil && errors.Is(err, services.ErrInvalidateSessionFailed) {
		res.SendError(w, http.StatusInternalServerError, res.NewError("logout_err", "Failed to logout of your session"))
		return
	} else if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went seriously wrong."))
		return
	}

	// Invalidate cookie
	cookie.ClearSessionCookie(w, h.cfg.Cookie)
}

// OAuth Callbacks
type OAuthState struct {
	Nonce    string `json:"nonce"`
	Provider string `json:"provider"`
	Redirect string `json:"redirect"`
}

func ensureLeadingSlash(s string) string {
	if len(s) == 0 || s[0] != '/' {
		return "/" + s
	}
	return s
}

func (h *AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	codeParam := q.Get("code")
	stateParam := q.Get("state")

	// User Agent + IpAddress for session
	var ipAddress *string
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && ip != "" {
		ipAddress = &ip
	}

	var userAgent *string
	ua := r.Header.Get("User-Agent")
	if ua != "" {
		userAgent = &ua
	}

	// Empty parameters
	if codeParam == "" || stateParam == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_callback", "This callback was invalid. Please try again."))
		return
	}

	decodedStateBytes, err := base64.URLEncoding.DecodeString(stateParam)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_callback", "This callback was invalid. Please try again."))
		return
	}

	var state OAuthState
	if err := json.Unmarshal(decodedStateBytes, &state); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_callback", "This callback was invalid. Please try again."))
		return
	}

	nonceCookie, err := r.Cookie("sh_auth_nonce")
	if err != nil {
		if err == http.ErrNoCookie {
			res.SendError(w, http.StatusForbidden, res.NewError("auth_error", "Failed to authenticate. Please try again."))
			return
		}

		res.SendError(w, http.StatusBadRequest, res.NewError("bad_cookie", "The cookie jar spilled over 😔"))
		return
	}

	if nonceCookie.Value != state.Nonce {
		res.SendError(w, http.StatusUnauthorized, res.NewError("auth_error", "Failed to authenticate. Please try again."))
		return
	}

	// Delete nonce cookie!
	cookie.ExpireCookie(w, h.cfg.Cookie, "sh_auth_nonce")

	// At this point, nonce has matched, proceed with remaining authentication services
	session, err := h.authService.AuthenticateWithOAuth(r.Context(), codeParam, state.Provider, ipAddress, userAgent)
	if err != nil {
		switch err {
		case services.ErrProviderUnsupported:
			res.SendError(w, http.StatusNotImplemented, res.NewError("provider_error", "This provider is not supported... are you sure you're supposed to be here?"))
			return
		case services.ErrAuthenticationFailed:
			res.SendError(w, http.StatusNotImplemented, res.NewError("auth_err", "Failed to authenticate the user."))
			return
		default:
			h.logger.Err(err).Msg("Something unexpected happened.")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went horribly wrong!"))
			return
		}
	}

	redirectPath := ensureLeadingSlash(state.Redirect)

	cookie.SetSessionCookie(w, session.ID, session.ExpiresAt, h.cfg.Cookie)
	http.Redirect(w, r, h.cfg.ClientUrl+redirectPath, http.StatusSeeOther)
}
