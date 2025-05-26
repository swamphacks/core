package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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

type OAuthCallbackRequest struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

func (h *AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	// This will be query params instead...
	var req OAuthCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sessionId, err := h.authService.AuthenticateWithOAuth(r.Context(), req.Code, req.Provider)
	if err != nil {
		http.Error(w, "code was invalid", http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "sh_session",
		Value:    sessionId,
		Path:     "/", // Available to entire site
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
