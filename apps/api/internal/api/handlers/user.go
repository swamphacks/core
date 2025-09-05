package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type UserHandler struct {
	userService *services.UserService
	logger      zerolog.Logger
}

func NewUserHandler(userService *services.UserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger.With().Str("handler", "UserHandler").Str("component", "user").Logger(),
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	user, err := h.userService.GetUser(r.Context(), *userId)
	if err != nil {
		h.logger.Err(err).Msg("failed to get user profile")
		if err == services.ErrUserNotFound {
			res.SendError(w, http.StatusNotFound, res.NewError("user_not_found", "User profile not found"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went seriously wrong."))
		}
		return
	}

	res.Send(w, http.StatusOK, user)
}

type CompleteOnboardingRequest struct {
	Name           string `json:"name"`
	PreferredEmail string `json:"preferred_email"`
}

func (h *UserHandler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	var req CompleteOnboardingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	// Validate required fields
	if req.Name == "" || req.PreferredEmail == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Name and preferred email are required"))
		return
	}

	// Validate email format
	if !email.IsValidEmail(req.PreferredEmail) {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_email", "Invalid email format"))
		return
	}

	err := h.userService.CompleteOnboarding(r.Context(), *userId, req.Name, req.PreferredEmail)
	if err != nil {
		h.logger.Err(err).Msg("failed to complete onboarding")
		if err == services.ErrUserNotFound {
			res.SendError(w, http.StatusNotFound, res.NewError("user_not_found", "User not found"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("update_failed", "Failed to complete onboarding"))
		}
		return
	}
}
