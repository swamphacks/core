package handlers

import (
	"encoding/json"
	"net/http"

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

// @Summary      Get authenticated user's information
// @Description  Retrieves current user's user information from database
// @Tags         auth
// @Produce      json
// @Success      200 {object} sqlc.AuthUser "User information"
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal server error"
// @Router       /auth/me [get]
func (h *AuthHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) {
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
