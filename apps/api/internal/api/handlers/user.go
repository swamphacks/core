package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/ptr"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/web"
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

type UpdateProfileRequest struct {
	Name           string `json:"name"`
	PreferredEmail string `json:"preferred_email"`
}

type UpdateEmailConsentRequest struct {
	EmailConsent bool `json:"email_consent"`
}

func (h *UserHandler) UpdateEmailConsent(w http.ResponseWriter, r *http.Request) {
	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	var req UpdateEmailConsentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	params := sqlc.UpdateUserParams{
		EmailConsentDoUpdate: true,
		EmailConsent:         req.EmailConsent,
	}

	err := h.userService.UpdateUser(r.Context(), *userId, params)
	if err != nil {
		h.logger.Err(err).Msg("failed to update email consent")
		if err == services.ErrUserNotFound {
			res.SendError(w, http.StatusNotFound, res.NewError("user_not_found", "User not found"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("update_failed", "Failed to update email consent"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	// Validate required fields
	if req.Name == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Name is required"))
		return
	}

	// Validate email format
	if req.PreferredEmail != "" {
		if !email.IsValidEmail(req.PreferredEmail) {
			res.SendError(w, http.StatusBadRequest, res.NewError("invalid_email", "Invalid email format"))
			return
		}
	}

	params := sqlc.UpdateUserParams{
		ID:                     *userId,
		NameDoUpdate:           true,
		Name:                   req.Name,
		PreferredEmailDoUpdate: true,
		PreferredEmail:         &req.PreferredEmail,
	}

	err := h.userService.UpdateUser(r.Context(), *userId, params)
	if err != nil {
		h.logger.Err(err).Msg("failed to update user settings")
		if err == services.ErrUserNotFound {
			res.SendError(w, http.StatusNotFound, res.NewError("user_not_found", "User not found"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("update_failed", "Failed to update user settings"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
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

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Parse search from query params, and limit, and offset
	searchTerm := web.ParseParamString(queryParams, "search", nil)
	limit, err := web.ParseParamInt32(queryParams, "limit", ptr.Int32ToPtr(50))
	if err != nil {
		h.logger.Err(err).Msg("Limit field was misconfigured. Please check your query parameters.")
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_query", "Your 'limit' query parameter was malformed."))
		return
	}
	offset, err := web.ParseParamInt32(queryParams, "offset", ptr.Int32ToPtr(0))
	if err != nil {

		h.logger.Err(err).Msg("Offset field was misconfigured. Please check your query parameters.")
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_query", "Your 'offset' query parameter was malformed."))
		return
	}

	users, err := h.userService.GetAllUsers(r.Context(), searchTerm, *limit, *offset)
	if err != nil {
		h.logger.Err(err).Msg("Failed to retrieve all users")
		res.SendError(w, http.StatusInternalServerError, res.NewError("update_failed", "Failed to complete onboarding"))
		return
	}

	res.Send(w, http.StatusOK, users)
}
