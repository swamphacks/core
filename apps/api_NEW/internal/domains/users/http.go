package users

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/emailutils"
)

func RegisterRoutes(userHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "get-me",
		Method:      http.MethodGet,
		Summary:     "Get Me",
		Description: "Returns the authenticated user's profile",
		Tags:        []string{"Users"},
		Path:        "/me",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleGetMe)

	huma.Register(group, huma.Operation{
		OperationID: "get-users",
		Method:      http.MethodGet,
		Summary:     "Get Users",
		Description: "Get or search for users by name or email. If no search term is provided, returns all users with pagination.",
		Tags:        []string{"Users"},
		Path:        "",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleGetUsers)

	huma.Register(group, huma.Operation{
		OperationID: "get-user-by-id",
		Method:      http.MethodGet,
		Summary:     "Get User By Id",
		Description: "Returns the user associated with the user id",
		Tags:        []string{"Users"},
		Path:        "/userid/{userId}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleGetUserById)

	huma.Register(group, huma.Operation{
		OperationID: "get-user-by-email",
		Method:      http.MethodGet,
		Summary:     "Get User By Email",
		Description: "Returns the user associated with the email",
		Tags:        []string{"Users"},
		Path:        "/email/{email}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleGetUserByEmail)

	huma.Register(group, huma.Operation{
		OperationID: "get-user-by-rfid",
		Method:      http.MethodGet,
		Summary:     "Get User By RFID",
		Description: "Returns the user associated with the RFID",
		Tags:        []string{"Users"},
		Path:        "/rfid/{rfid}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleGetUserByRFID)

	huma.Register(group, huma.Operation{
		OperationID:   "update-email-consent",
		Method:        http.MethodPatch,
		Summary:       "Update Email Consent",
		Description:   "Updates the user's email consent setting",
		Tags:          []string{"Users"},
		Path:          "/me/email-consent",
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
		Tags:          []string{"Users"},
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
		Tags:          []string{"Users"},
		Path:          "/me/onboarding",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, userHandler.handleOnboarding)

	huma.Register(group, huma.Operation{
		OperationID: "assign-role",
		Method:      http.MethodPost,
		Summary:     "Assign Role",
		Description: "Assigns/modify a user's role",
		Tags:        []string{"Users"},
		Path:        "/roles/assign",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleAssignRole)

	huma.Register(group, huma.Operation{
		OperationID: "batch-assign-roles",
		Method:      http.MethodPost,
		Summary:     "Batch Assign Roles",
		Description: "Batch assign/modify multiple users' roles",
		Tags:        []string{"Users"},
		Path:        "/roles/batch-assign",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleBatchAssignRoles)

	huma.Register(group, huma.Operation{
		OperationID: "revoke-role",
		Method:      http.MethodPost,
		Summary:     "Revoke Role",
		Description: "Remove a user's role",
		Tags:        []string{"Users"},
		Path:        "/roles/revoke/{userId}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, userHandler.handleRevokeEventRole)
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

type GetMeOutput struct {
	Body *middleware.UserContext
}

func (h *handler) handleGetMe(ctx context.Context, input *struct{}) (*GetMeOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	return &GetMeOutput{Body: userCtx}, nil
}

type GetUserByEmailOutput struct {
	Body *sqlc.User
}

func (h *handler) handleGetUserByEmail(ctx context.Context, input *struct {
	Email string `path:"email"`
}) (*GetUserByEmailOutput, error) {
	if !emailutils.IsValidEmail(input.Email) {
		return nil, huma.Error400BadRequest("Invalid email")
	}

	user, err := h.userService.GetUserByEmail(ctx, input.Email)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user by email")
	}

	if user == nil {
		return nil, huma.Error404NotFound("User not found")
	}

	return &GetUserByEmailOutput{Body: user}, nil
}

type GetUserByIdOutput struct {
	Body *sqlc.User
}

func (h *handler) handleGetUserById(ctx context.Context, input *struct {
	UserId string `path:"userId"`
}) (*GetUserByIdOutput, error) {
	userId, err := uuid.Parse(input.UserId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid user id")
	}

	user, err := h.userService.GetUserById(ctx, userId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user by id")
	}

	if user == nil {
		return nil, huma.Error404NotFound("User not found")
	}

	return &GetUserByIdOutput{Body: user}, nil
}

type GetUserByRFIDOutput struct {
	Body *sqlc.User
}

func (h *handler) handleGetUserByRFID(ctx context.Context, input *struct {
	RFID string `path:"rfid"`
}) (*GetUserByRFIDOutput, error) {
	user, err := h.userService.GetUserByRFID(ctx, input.RFID)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user by rfid")
	}

	if user == nil {
		return nil, huma.Error404NotFound("User not found")
	}

	return &GetUserByRFIDOutput{Body: user}, nil
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
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	if input.Body.Name == "" {
		return nil, huma.Error400BadRequest("Name is required")
	}

	if input.Body.PreferredEmail != "" && !emailutils.IsValidEmail(input.Body.PreferredEmail) {
		return nil, huma.Error400BadRequest("Invalid email format")
	}

	// TODO: Allow/add more fields here
	params := sqlc.UpdateUserParams{
		ID:                     userCtx.UserID,
		NameDoUpdate:           true,
		Name:                   input.Body.Name,
		PreferredEmailDoUpdate: true,
		PreferredEmail:         &input.Body.PreferredEmail,
	}

	err := h.userService.UpdateUser(ctx, userCtx.UserID, params)
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
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	params := sqlc.UpdateUserParams{
		EmailConsentDoUpdate: true,
		EmailConsent:         input.Body.EmailConsent,
	}

	err := h.userService.UpdateUser(ctx, userCtx.UserID, params)

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
	userCtx := ctxutils.GetUserFromCtx(ctx)

	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	if input.Body.Name == "" || input.Body.PreferredEmail == "" {
		return nil, huma.Error400BadRequest("Name and preferred email are required")
	}

	if input.Body.PreferredEmail != "" && !emailutils.IsValidEmail(input.Body.PreferredEmail) {
		return nil, huma.Error400BadRequest("Invalid email format")
	}

	err := h.userService.CompleteOnboarding(ctx, userCtx.UserID, input.Body.Name, input.Body.PreferredEmail)

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
	Body *[]sqlc.User
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

type AssignRoleRequest struct {
	Email  *string       `json:"email"`
	UserID *string       `json:"user_id"`
	Role   sqlc.UserRole `json:"role"`
}

type AssignRoleOutput struct {
	Status int
}

func (h *handler) handleAssignRole(ctx context.Context, input *struct {
	Body AssignRoleRequest
}) (*AssignRoleOutput, error) {
	var userId *uuid.UUID

	if input.Body.UserID == nil {
		userId = nil
	} else {
		userIdTemp, err := uuid.Parse(*input.Body.UserID)

		if err != nil {
			userId = nil
		} else {
			userId = &userIdTemp
		}
	}

	err := h.userService.AssignRole(ctx, userId, input.Body.Email, input.Body.Role)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to assign role")
	}

	return &AssignRoleOutput{Status: http.StatusOK}, nil
}

type AssignRoleBatchRequest struct {
	Assignments []AssignRoleRequest `json:"assignments"`
}

type BatchAssignRolesOutput struct {
	Status int
}

func (h *handler) handleBatchAssignRoles(ctx context.Context, input *struct {
	Body AssignRoleBatchRequest
}) (*BatchAssignRolesOutput, error) {
	for _, assignment := range input.Body.Assignments {
		userId := ParseUUIDOrNil(assignment.UserID)

		err := h.userService.AssignRole(ctx, userId, assignment.Email, assignment.Role)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to batch assign roles")
		}
	}

	return &BatchAssignRolesOutput{Status: http.StatusOK}, nil
}

type RevokeEventRoleOutput struct {
	Status int
}

func (h *handler) handleRevokeEventRole(ctx context.Context, input *struct {
	UserId string `path:"userId"`
}) (*RevokeEventRoleOutput, error) {
	userId, err := uuid.Parse(input.UserId)

	if err != nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	err = h.userService.RevokeRole(ctx, userId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to revoke role")
	}

	return &RevokeEventRoleOutput{Status: http.StatusOK}, nil
}

func ParseUUIDOrNil(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}
