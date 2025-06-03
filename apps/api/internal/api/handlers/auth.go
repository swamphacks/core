package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      zerolog.Logger
}

func NewAuthHandler(authService *services.AuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.With().Str("handler", "AuthHandler").Str("component", "auth").Logger(),
	}
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// If no userId here, throw
	user, err := h.authService.GetMe(r.Context())
	if err != nil {
		http.Error(w, "not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "something went wrong retrieving user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// OAuth Callbacks
type OAuthState struct {
	Nonce    string `json:"nonce"`
	Provider string `json:"provider"`
	Redirect string `json:"redirect"`
}

func (h *AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	log.Trace().Str("method", r.Method).Str("path", r.URL.Path).Msg("[HIT]")

	q := r.URL.Query()

	codeParam := q.Get("code")
	stateParam := q.Get("state")

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

	// At this point, nonce has matched, proceed with remaining authentication services
	session, err := h.authService.AuthenticateWithOAuth(r.Context(), codeParam, state.Provider)
	if err != nil {
		switch err {
		case services.ErrProviderUnsupported:
			res.SendError(w, http.StatusNotImplemented, res.NewError("provider_error", "This provider is not supported... are you sure you're supposed to be here?"))
			return
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went horribly wrong!"))
			return
		}
	}

	cookie := &http.Cookie{
		Name:     "sh_session_id",
		Value:    session.Token,
		Domain:   "localhost",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Only for dev
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "http://localhost:3000"+state.Redirect, http.StatusSeeOther)
}
