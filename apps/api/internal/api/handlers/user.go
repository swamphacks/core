package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
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

// GetProfile returns the current user's profile information
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

// UpdateProfileRequest represents the request body for updating user profile
type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// UpdateProfile updates the current user's profile information
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	var req UpdateProfileRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	// Validate that at least one field is provided
	if req.Name == nil && req.Email == nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "At least one field (name or email) must be provided"))
		return
	}

	// Validate email format if provided
	if req.Email != nil && *req.Email != "" {
		// Basic email validation - you might want to use a more robust validation library
		if !isValidEmail(*req.Email) {
			res.SendError(w, http.StatusBadRequest, res.NewError("invalid_email", "Invalid email format"))
			return
		}
	}

	err := h.userService.UpdateUserProfile(r.Context(), *userId, req.Name, req.Email)
	if err != nil {
		h.logger.Err(err).Msg("failed to update user profile")
		if err == services.ErrUserNotFound {
			res.SendError(w, http.StatusNotFound, res.NewError("user_not_found", "User not found"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("update_failed", "Failed to update user profile"))
		}
		return
	}

	// Return the updated user profile
	user, err := h.userService.GetUser(r.Context(), *userId)
	if err != nil {
		h.logger.Err(err).Msg("failed to get updated user profile")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Profile updated but failed to retrieve updated data"))
		return
	}

	res.Send(w, http.StatusOK, user)
}

type UpdateOnboardedRequest struct {
	Onboarded bool `json:"onboarded"`
}

func (h *UserHandler) UpdateOnboarded(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	var req UpdateOnboardedRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	if req.Onboarded {
		err := h.userService.UpdateUserOnboarded(r.Context(), *userId)
		if err != nil {
			h.logger.Err(err).Msg("failed to update user onboarded status")
			if err == services.ErrUserNotFound {
				res.SendError(w, http.StatusNotFound, res.NewError("user_not_found", "User not found"))
			} else {
				res.SendError(w, http.StatusInternalServerError, res.NewError("update_failed", "Failed to update onboarded status"))
			}
			return
		}
	} else {
		params := sqlc.UpdateUserParams{
			ID:                *userId,
			OnboardedDoUpdate: true,
			Onboarded:         false,
		}
		err := h.userService.UpdateUser(r.Context(), *userId, params)
		if err != nil {
			h.logger.Err(err).Msg("failed to update user onboarded status")
			if err == services.ErrUserNotFound {
				res.SendError(w, http.StatusNotFound, res.NewError("user_not_found", "User not found"))
			} else {
				res.SendError(w, http.StatusInternalServerError, res.NewError("update_failed", "Failed to update onboarded status"))
			}
			return
		}
	}

	// Return the updated user profile
	user, err := h.userService.GetUser(r.Context(), *userId)
	if err != nil {
		h.logger.Err(err).Msg("failed to get updated user profile")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Status updated but failed to retrieve updated data"))
		return
	}

	res.Send(w, http.StatusOK, user)
}

// Very basic email validation, checking for @ and .
func isValidEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	atCount := 0
	dotAfterAt := false
	for _, char := range email {
		if char == '@' {
			atCount++
			if atCount > 1 {
				return false // Multiple @ symbols
			}
		} else if char == '.' && atCount == 1 {
			dotAfterAt = true
		}
	}
	return atCount == 1 && dotAfterAt && len(email) > 5
}
