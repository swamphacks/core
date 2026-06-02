package email

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	//Reuses repository-level "not found" error
	ErrEmailCampaignNotFound = repository.ErrEmailCampaignNotFound

	//validation errors before writing error to db
	ErrEmailCampaignTitleRequired      = errors.New("email campaign title is required")
	ErrEmailCampaignSubjectRequired    = errors.New("email campaign subject is required")
	ErrEmailCampaignBodyRequired       = errors.New("email campaign body is required")
	ErrEmailCampaignRecipientsRequired = errors.New("email campaign recipients are required")

	ErrEmailCampaignCannotEdit = errors.New("email campaign cannot be edited")

	//status-specific validation errors
	ErrEmailCampaignScheduledAtRequired = errors.New("scheduled_at is required for scheduled campaigns")
	ErrEmailCampaignSentAtRequired      = errors.New("sent_at is required for sent campaigns")
)

// EmailCampaignService owns business rules for saved email campaigns.
type EmailCampaignService struct {
	emailCampaignRepo *repository.EmailCampaignRepository
	logger zerolog.Logger
}

// NewEmailCampaignService creates the service and stores its dependencies.
// This will eventually be called from api.go when wiring the app together.
func NewEmailCampaignService(
	emailCampaignRepo *repository.EmailCampaignRepository,
	logger zerolog.Logger,
) *EmailCampaignService {
	return &EmailCampaignService{
		emailCampaignRepo: emailCampaignRepo,
		logger:            logger.With().Str("service", "EmailCampaignService").Str("domain", "email").Logger(),
	}
}

// CreateCampaign validates required campaign fields, then stores a new campaign.
// The actual INSERT is handled by the repository/sqlc layer.
func (s *EmailCampaignService) CreateCampaign(
	ctx context.Context,
	params sqlc.CreateEmailCampaignParams,
) (*sqlc.EmailCampaign, error) {
	if err := validateCampaignContent(params.Title, params.Subject, params.Body, params.RecipientTypes); err != nil {
		return nil, err
	}

	return s.emailCampaignRepo.CreateEmailCampaign(ctx, params)
}


// GetCampaignByID fetches one campaign scoped to a hackathon.
// The hackathon scope prevents one event from reading another event's campaign.
func (s *EmailCampaignService) GetCampaignByID(
	ctx context.Context,
	params sqlc.GetEmailCampaignByIDParams,
) (*sqlc.EmailCampaign, error) {
	return s.emailCampaignRepo.GetEmailCampaignByID(ctx, params)
}

// ListCampaigns fetches all campaigns for one hackathon.
// Sorting is handled by the SQL query, currently newest first.
func (s *EmailCampaignService) ListCampaigns(
	ctx context.Context,
	hackathonID string,
) ([]sqlc.EmailCampaign, error) {
	return s.emailCampaignRepo.ListEmailCampaigns(ctx, hackathonID)
}

// UpdateCampaign updates editable campaign fields.
// It first loads the existing campaign so we can enforce status rules before updating.
func (s *EmailCampaignService) UpdateCampaign(
	ctx context.Context,
	params sqlc.UpdateEmailCampaignParams,
) (*sqlc.EmailCampaign, error) {
	existingCampaign, err := s.emailCampaignRepo.GetEmailCampaignByID(ctx, sqlc.GetEmailCampaignByIDParams{
		ID:          params.ID,
		HackathonID: params.HackathonID,
	})
	if err != nil {
		return nil, err
	}

	if !canEditCampaign(existingCampaign.Status) {
		return nil, ErrEmailCampaignCannotEdit
	}

	return s.emailCampaignRepo.UpdateEmailCampaign(ctx, params)
}

// UpdateCampaignStatus changes lifecycle fields such as draft -> scheduled or sending -> sent.
// The database also has constraints, but checking here gives cleaner service-level errors.
func (s *EmailCampaignService) UpdateCampaignStatus(
	ctx context.Context,
	params sqlc.UpdateEmailCampaignStatusParams,
) (*sqlc.EmailCampaign, error) {
	if params.Status == sqlc.EmailCampaignStatusScheduled && params.ScheduledAt == nil {
		return nil, ErrEmailCampaignScheduledAtRequired
	}

	if params.Status == sqlc.EmailCampaignStatusSent && params.SentAt == nil {
		return nil, ErrEmailCampaignSentAtRequired
	}

	return s.emailCampaignRepo.UpdateEmailCampaignStatus(ctx, params)
}

// validateCampaignContent checks fields that every campaign needs before it is saved.
// strings.TrimSpace prevents values like "   " from passing validation.
func validateCampaignContent(
	title string,
	subject string,
	body string,
	recipientTypes []sqlc.EmailRecipientType,
) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmailCampaignTitleRequired
	}

	if strings.TrimSpace(subject) == "" {
		return ErrEmailCampaignSubjectRequired
	}

	if strings.TrimSpace(body) == "" {
		return ErrEmailCampaignBodyRequired
	}

	if len(recipientTypes) == 0 {
		return ErrEmailCampaignRecipientsRequired
	}

	return nil
}

// canEditCampaign centralizes edit rules.
// Drafts are editable, and scheduled campaigns can still be adjusted before sending.
func canEditCampaign(status sqlc.EmailCampaignStatus) bool {
	return status == sqlc.EmailCampaignStatusDraft ||
		status == sqlc.EmailCampaignStatusScheduled
}