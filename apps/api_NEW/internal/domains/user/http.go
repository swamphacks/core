package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/email"
)

func RegisterRoutes(userHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "update-email-consent",
		Method:        http.MethodPatch,
		Summary:       "Update Email Consent",
		Description:   "Updates the user's email consent setting",
		Tags:          []string{"User"},
		Path:          "/email-consent",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, userHandler.handleUpdateEmailConsent)

	huma.Register(group, huma.Operation{
		OperationID:   "update-user",
		Method:        http.MethodPatch,
		Summary:       "Update User",
		Description:   "Updates information of the authenticated user",
		Tags:          []string{"User"},
		Path:          "/me",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, userHandler.handleUpdateUser)

	huma.Register(group, huma.Operation{
		OperationID:   "onboard-user",
		Method:        http.MethodPatch,
		Summary:       "Onboard User",
		Description:   "Allows the user to submit information such as name and preferred email, and complete the onboarding process",
		Tags:          []string{"User"},
		Path:          "/onboarding",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, userHandler.handleOnboarding)

	huma.Register(group, huma.Operation{
		OperationID: "get-users",
		Method:      http.MethodGet,
		Summary:     "Get Users",
		Description: "Get or search for users by name or email. If no search term is provided, returns all users with pagination.",
		Tags:        []string{"User"},
		Path:        "/search",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleGetUsers)
}

type handler struct {
	userService *UserService
	config      *config.Config
	logger      zerolog.Logger
}

func NewHandler(userService *UserService, config *config.Config, logger zerolog.Logger) *handler {
	return &handler{
		userService: userService,
		config:      config,
		logger:      logger.With().Str("handler", "UserHandler").Str("domain", "user").Logger(),
	}
}

type UpdateUserOutput struct {
	Status int
}

type UpdateUserRequest struct {
	Name           string `json:"name"`
	PreferredEmail string `json:"preferred_email"`
}

func (h *handler) handleUpdateUser(ctx context.Context, input *struct {
	Body UpdateUserRequest
}) (*UpdateUserOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	if input.Body.Name == "" {
		return nil, huma.Error400BadRequest("Name is required")
	}

	if input.Body.PreferredEmail != "" && !email.IsValidEmail(input.Body.PreferredEmail) {
		return nil, huma.Error400BadRequest("Invalid email format")
	}

	params := sqlc.UpdateUserParams{
		ID:                     *userId,
		NameDoUpdate:           true,
		Name:                   input.Body.Name,
		PreferredEmailDoUpdate: true,
		PreferredEmail:         &input.Body.PreferredEmail,
	}

	err := h.userService.UpdateUser(ctx, *userId, params)
	if err != nil {
		h.logger.Err(err).Msg("failed to update user")
		if errors.Is(err, ErrUserNotFound) {
			return nil, huma.Error404NotFound("User not found")
		} else {
			return nil, huma.Error500InternalServerError("failed to update user")
		}
	}

	res := &UpdateUserOutput{
		Status: http.StatusOK,
	}

	return res, nil
}

type UpdateEmailConsentOutput struct {
	Status int
}

type UpdateEmailConsentRequest struct {
	EmailConsent bool `json:"email_consent"`
}

func (h *handler) handleUpdateEmailConsent(ctx context.Context, input *struct {
	Body UpdateEmailConsentRequest
}) (*UpdateEmailConsentOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	params := sqlc.UpdateUserParams{
		EmailConsentDoUpdate: true,
		EmailConsent:         input.Body.EmailConsent,
	}

	err := h.userService.UpdateUser(ctx, *userId, params)

	if err != nil {
		h.logger.Err(err).Msg("failed to update email consent")
		if errors.Is(err, ErrUserNotFound) {
			return nil, huma.Error404NotFound("User not found")
		} else {
			return nil, huma.Error500InternalServerError("Failed to update email consent")
		}
	}

	res := &UpdateEmailConsentOutput{
		Status: http.StatusOK,
	}

	return res, nil
}

type OnboardingOutput struct {
	Status int
}

type OnboardingRequest struct {
	Name           string `json:"name"`
	PreferredEmail string `json:"preferred_email"`
}

func (h *handler) handleOnboarding(ctx context.Context, input *struct {
	Body OnboardingRequest
}) (*OnboardingOutput, error) {
	userId := ctxutils.GetUserIdFromCtx(ctx)

	if userId == nil {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	if input.Body.Name == "" || input.Body.PreferredEmail == "" {
		return nil, huma.Error400BadRequest("Name and preferred email are required")
	}

	if input.Body.PreferredEmail != "" && !email.IsValidEmail(input.Body.PreferredEmail) {
		return nil, huma.Error400BadRequest("Invalid email format")
	}

	err := h.userService.CompleteOnboarding(ctx, *userId, input.Body.Name, input.Body.PreferredEmail)

	if err != nil {
		h.logger.Err(err).Msg("failed to complete onboarding")
		if errors.Is(err, ErrUserNotFound) {
			return nil, huma.Error404NotFound("User not found")
		} else {
			return nil, huma.Error500InternalServerError("failed to complete onboarding")
		}
	}

	res := &OnboardingOutput{
		Status: http.StatusOK,
	}

	return res, nil
}

type GetUsersOutput struct {
	Body *[]sqlc.AuthUser
}

func (h *handler) handleGetUsers(ctx context.Context, input *struct {
	Search string `query:"search"`
	Limit  int    `query:"limit" default:"50"`
	Offset int    `query:"offset" default:"0"`
}) (*GetUsersOutput, error) {
	var searchTerm *string
	if input.Search == "" {
		searchTerm = nil
	} else {
		searchTerm = &input.Search
	}

	users, err := h.userService.GetAllUsers(ctx, searchTerm, int32(input.Limit), int32(input.Offset))

	if err != nil {
		h.logger.Err(err).Msg("Failed to retrieve users")
		return nil, huma.Error500InternalServerError("Failed to retrieve users")
	}

	res := &GetUsersOutput{Body: &users}

	return res, nil
}
