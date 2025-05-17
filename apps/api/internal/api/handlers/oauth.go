package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swamphacks/core/apps/api/internal/services"
)

type OAuthVerifyRequest struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

type OAuthHandler struct {
	authService *services.AuthService
}

func NewOAuthHandler(authService *services.AuthService) *OAuthHandler {
	return &OAuthHandler{
		authService: authService,
	}
}

// @Summary      Verify OAuth Code
// @Description  Verifies an OAuth code and sets a session token in a cookie
// @Tags         oauth
// @Accept       json
// @Produce      json
// @Param        request body OAuthVerifyRequest true "OAuth Verification Request"
// @Success      200 {string} string "Session token set in cookie"
// @Failure      400 {string} string "Invalid request body"
// @Failure      401 {string} string "Unauthorized or invalid code"
// @Router       /oauth/verify [post]
func (h *OAuthHandler) VerifyOAuth(w http.ResponseWriter, r *http.Request) {
	var req OAuthVerifyRequest
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
