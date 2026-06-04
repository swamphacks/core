package email

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func RegisterCampaignRoutes(emailCampaignHandler *emailCampaignHandler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID:   "create-email-campaign",
		Method:        http.MethodPost,
		Summary:       "Create Email Campaign",
		Description:   "Creates a saved email campaign draft for a hackathon.",
		Tags:          []string{"Email Campaigns"},
		Path:          "/campaigns",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusCreated,
	}, emailCampaignHandler.handleCreateCampaign)

	huma.Register(group, huma.Operation{
		OperationID:   "list-email-campaigns",
		Method:        http.MethodGet,
		Summary:       "List Email Campaigns",
		Description:   "Returns all saved email campaigns for a hackathon.",
		Tags:          []string{"Email Campaigns"},
		Path:          "/campaigns",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, emailCampaignHandler.handleListCampaigns)

	huma.Register(group, huma.Operation{
		OperationID:   "get-email-campaign",
		Method:        http.MethodGet,
		Summary:       "Get Email Campaign",
		Description:   "Returns one saved email campaign by id.",
		Tags:          []string{"Email Campaigns"},
		Path:          "/campaigns/{campaignId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, emailCampaignHandler.handleGetCampaign)

	huma.Register(group, huma.Operation{
		OperationID:   "update-email-campaign",
		Method:        http.MethodPatch,
		Summary:       "Update Email Campaign",
		Description:   "Updates editable fields on a draft or scheduled email campaign.",
		Tags:          []string{"Email Campaigns"},
		Path:          "/campaigns/{campaignId}",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, emailCampaignHandler.handleUpdateCampaign)

	huma.Register(group, huma.Operation{
		OperationID:   "update-email-campaign-status",
		Method:        http.MethodPatch,
		Summary:       "Update Email Campaign Status",
		Description:   "Updates lifecycle fields such as status, scheduled_at, sent_at, and last_error.",
		Tags:          []string{"Email Campaigns"},
		Path:          "/campaigns/{campaignId}/status",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, emailCampaignHandler.handleUpdateCampaignStatus)
}

type emailCampaignHandler struct {
	emailCampaignService *EmailCampaignService
	logger               zerolog.Logger
}

func NewCampaignHandler(emailCampaignService *EmailCampaignService, logger zerolog.Logger) *emailCampaignHandler {
	return &emailCampaignHandler{
		emailCampaignService: emailCampaignService,
		logger:               logger.With().Str("handler", "EmailCampaignHandler").Str("domain", "email").Logger(),
	}
}

type CreateEmailCampaignRequest struct {
	HackathonID    string                    `json:"hackathonId" required:"true"`
	Title          string                    `json:"title" minLength:"1"`
	Description    *string                   `json:"description,omitempty"`
	Subject        string                    `json:"subject" minLength:"1"`
	Body           string                    `json:"body" minLength:"1"`
	Format         sqlc.EmailCampaignFormat  `json:"format" required:"true"`
	RecipientTypes []sqlc.EmailRecipientType `json:"recipientTypes" minItems:"1"`
	ScheduledAt    *time.Time                `json:"scheduledAt,omitempty"`
}

type UpdateEmailCampaignRequest struct {
	Title          *string                    `json:"title,omitempty"`
	Description    *string                    `json:"description,omitempty"`
	Subject        *string                    `json:"subject,omitempty"`
	Body           *string                    `json:"body,omitempty"`
	Format         *sqlc.EmailCampaignFormat  `json:"format,omitempty"`
	RecipientTypes *[]sqlc.EmailRecipientType `json:"recipientTypes,omitempty"`
	ScheduledAt    *time.Time                 `json:"scheduledAt,omitempty"`
}

type UpdateEmailCampaignStatusRequest struct {
	Status      sqlc.EmailCampaignStatus `json:"status" required:"true"`
	ScheduledAt *time.Time               `json:"scheduledAt,omitempty"`
	SentAt      *time.Time               `json:"sentAt,omitempty"`
	LastError   *string                  `json:"lastError,omitempty"`
}

type EmailCampaignOutput struct {
	Body *sqlc.EmailCampaign
}

type ListEmailCampaignsOutput struct {
	Body []sqlc.EmailCampaign
}

func (h *emailCampaignHandler) handleCreateCampaign(ctx context.Context, input *struct {
	Body CreateEmailCampaignRequest
}) (*EmailCampaignOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	campaign, err := h.emailCampaignService.CreateCampaign(ctx, sqlc.CreateEmailCampaignParams{
		HackathonID:     input.Body.HackathonID,
		Title:           input.Body.Title,
		Description:     input.Body.Description,
		Subject:         input.Body.Subject,
		Body:            input.Body.Body,
		Format:          input.Body.Format,
		RecipientTypes:  input.Body.RecipientTypes,
		ScheduledAt:     input.Body.ScheduledAt,
		CreatedByUserID: &userCtx.UserID,
		UpdatedByUserID: &userCtx.UserID,
	})
	if err != nil {
		return nil, campaignHTTPError(err, "Failed to create email campaign")
	}

	return &EmailCampaignOutput{Body: campaign}, nil
}

func (h *emailCampaignHandler) handleListCampaigns(ctx context.Context, input *struct {
	HackathonID string `query:"hackathonId" required:"true"`
}) (*ListEmailCampaignsOutput, error) {
	campaigns, err := h.emailCampaignService.ListCampaigns(ctx, input.HackathonID)
	if err != nil {
		return nil, campaignHTTPError(err, "Failed to list email campaigns")
	}

	return &ListEmailCampaignsOutput{Body: campaigns}, nil
}

func (h *emailCampaignHandler) handleGetCampaign(ctx context.Context, input *struct {
	CampaignID  string `path:"campaignId"`
	HackathonID string `query:"hackathonId" required:"true"`
}) (*EmailCampaignOutput, error) {
	campaignID, err := uuid.Parse(input.CampaignID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid campaign id")
	}

	campaign, err := h.emailCampaignService.GetCampaignByID(ctx, sqlc.GetEmailCampaignByIDParams{
		ID:          campaignID,
		HackathonID: input.HackathonID,
	})
	if err != nil {
		return nil, campaignHTTPError(err, "Failed to get email campaign")
	}

	return &EmailCampaignOutput{Body: campaign}, nil
}

func (h *emailCampaignHandler) handleUpdateCampaign(ctx context.Context, input *struct {
	CampaignID  string `path:"campaignId"`
	HackathonID string `query:"hackathonId" required:"true"`
	Body        UpdateEmailCampaignRequest
}) (*EmailCampaignOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	campaignID, err := uuid.Parse(input.CampaignID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid campaign id")
	}
	if input.Body.Title != nil && *input.Body.Title == "" {
		return nil, huma.Error400BadRequest(ErrEmailCampaignTitleRequired.Error())
	}
	if input.Body.Subject != nil && *input.Body.Subject == "" {
		return nil, huma.Error400BadRequest(ErrEmailCampaignSubjectRequired.Error())
	}
	if input.Body.Body != nil && *input.Body.Body == "" {
		return nil, huma.Error400BadRequest(ErrEmailCampaignBodyRequired.Error())
	}
	if input.Body.RecipientTypes != nil && len(*input.Body.RecipientTypes) == 0 {
		return nil, huma.Error400BadRequest(ErrEmailCampaignRecipientsRequired.Error())
	}

	params := sqlc.UpdateEmailCampaignParams{
		TitleDoUpdate:           input.Body.Title != nil,
		DescriptionDoUpdate:     input.Body.Description != nil,
		SubjectDoUpdate:         input.Body.Subject != nil,
		BodyDoUpdate:            input.Body.Body != nil,
		FormatDoUpdate:          input.Body.Format != nil,
		RecipientTypesDoUpdate:  input.Body.RecipientTypes != nil,
		ScheduledAtDoUpdate:     input.Body.ScheduledAt != nil,
		UpdatedByUserIDDoUpdate: true,
		UpdatedByUserID:         userCtx.UserID,
		ID:                      campaignID,
		HackathonID:             input.HackathonID,
	}

	if input.Body.Title != nil {
		params.Title = *input.Body.Title
	}
	if input.Body.Description != nil {
		params.Description = input.Body.Description
	}
	if input.Body.Subject != nil {
		params.Subject = *input.Body.Subject
	}
	if input.Body.Body != nil {
		params.Body = *input.Body.Body
	}
	if input.Body.Format != nil {
		params.Format = *input.Body.Format
	}
	if input.Body.RecipientTypes != nil {
		params.RecipientTypes = *input.Body.RecipientTypes
	}
	if input.Body.ScheduledAt != nil {
		params.ScheduledAt = input.Body.ScheduledAt
	}

	campaign, err := h.emailCampaignService.UpdateCampaign(ctx, params)
	if err != nil {
		return nil, campaignHTTPError(err, "Failed to update email campaign")
	}

	return &EmailCampaignOutput{Body: campaign}, nil
}

func (h *emailCampaignHandler) handleUpdateCampaignStatus(ctx context.Context, input *struct {
	CampaignID  string `path:"campaignId"`
	HackathonID string `query:"hackathonId" required:"true"`
	Body        UpdateEmailCampaignStatusRequest
}) (*EmailCampaignOutput, error) {
	userCtx := ctxutils.GetUserFromCtx(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Failed to get current user info")
	}

	campaignID, err := uuid.Parse(input.CampaignID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid campaign id")
	}

	campaign, err := h.emailCampaignService.UpdateCampaignStatus(ctx, sqlc.UpdateEmailCampaignStatusParams{
		Status:                  input.Body.Status,
		ScheduledAtDoUpdate:     input.Body.ScheduledAt != nil,
		ScheduledAt:             input.Body.ScheduledAt,
		SentAtDoUpdate:          input.Body.SentAt != nil,
		SentAt:                  input.Body.SentAt,
		LastErrorDoUpdate:       input.Body.LastError != nil,
		LastError:               input.Body.LastError,
		UpdatedByUserIDDoUpdate: true,
		UpdatedByUserID:         userCtx.UserID,
		ID:                      campaignID,
		HackathonID:             input.HackathonID,
	})
	if err != nil {
		return nil, campaignHTTPError(err, "Failed to update email campaign status")
	}

	return &EmailCampaignOutput{Body: campaign}, nil
}

func campaignHTTPError(err error, fallback string) error {
	if errors.Is(err, ErrEmailCampaignNotFound) {
		return huma.Error404NotFound("Email campaign not found")
	}

	if errors.Is(err, ErrEmailCampaignCannotEdit) ||
		errors.Is(err, ErrEmailCampaignTitleRequired) ||
		errors.Is(err, ErrEmailCampaignSubjectRequired) ||
		errors.Is(err, ErrEmailCampaignBodyRequired) ||
		errors.Is(err, ErrEmailCampaignRecipientsRequired) ||
		errors.Is(err, ErrEmailCampaignScheduledAtRequired) ||
		errors.Is(err, ErrEmailCampaignSentAtRequired) {
		return huma.Error400BadRequest(err.Error())
	}

	return huma.Error500InternalServerError(fallback)
}
