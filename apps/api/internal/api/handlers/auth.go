package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
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
		res := response.NewError("A100", "invalid_callback", "Missing or malformed code or state in callback")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	decodedStateBytes, err := base64.URLEncoding.DecodeString(stateParam)
	if err != nil {
		res := response.NewError("A100", "invalid_callback", "Missing or malformed code or state in callback")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var state OAuthState
	if err := json.Unmarshal(decodedStateBytes, &state); err != nil {
		res := response.NewError("A100", "invalid_callback", "Missing or malformed code or state in callback")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	nonceCookie, err := r.Cookie("sh_auth_nonce")
	if err != nil {
		if err == http.ErrNoCookie {
			res := response.NewError("A102", "nonce_missing", "Nonce cookie missing")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	if nonceCookie.Value != state.Nonce {
		res := response.NewError("A101", "nonce_mismatch", "Nonce in state does not match cookie. Possible CSRF attack...")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// At this point, nonce has matched, proceed with remaining authentication services
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("Try again soon~!")); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
