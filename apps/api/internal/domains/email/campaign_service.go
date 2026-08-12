package email

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/google/uuid"
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
	ErrEmailCampaignCannotSend          = errors.New("email campaign cannot be sent in its current status")
	ErrEmailCampaignNoRecipients        = errors.New("email campaign has no resolvable recipients")
	ErrUnsupportedRecipientType         = errors.New("unsupported recipient type")
)

// recipientStatuses maps an applicant recipient_type to the application statuses it covers.
var recipientStatuses = map[string][]string{
	"accepted_applicants":   {"accepted", "confirmed"},
	"rejected_applicants":   {"rejected"},
	"waitlisted_applicants": {"waitlisted"},
}

// EmailCampaignService owns business rules for saved email campaigns.
type EmailCampaignService struct {
	emailCampaignRepo *repository.EmailCampaignRepository
	emailService      *EmailService
	logger            zerolog.Logger
}

// NewEmailCampaignService creates the service and stores its dependencies.
// This will eventually be called from api.go when wiring the app together.
func NewEmailCampaignService(
	emailCampaignRepo *repository.EmailCampaignRepository,
	emailService *EmailService,
	logger zerolog.Logger,
) *EmailCampaignService {
	return &EmailCampaignService{
		emailCampaignRepo: emailCampaignRepo,
		emailService:      emailService,
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
	recipientTypes []string,
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

// canSendCampaign allows sending only from pre-send states.
// Drafts and still-scheduled campaigns can be sent; sending/sent/failed cannot.
func canSendCampaign(status sqlc.EmailCampaignStatus) bool {
	return status == sqlc.EmailCampaignStatusDraft ||
		status == sqlc.EmailCampaignStatusScheduled
}

// canEditCampaign centralizes edit rules.
// Drafts are editable, and scheduled campaigns can still be adjusted before sending.
func canEditCampaign(status sqlc.EmailCampaignStatus) bool {
	return status == sqlc.EmailCampaignStatusDraft ||
		status == sqlc.EmailCampaignStatusScheduled
}

// resolveRecipients turns a campaign's recipient_types into a deduplicated email list.
func (s *EmailCampaignService) resolveRecipients(ctx context.Context, campaign *sqlc.EmailCampaign) ([]string, error) {
	seen := make(map[string]struct{})
	emails := []string{}

	for _, rt := range campaign.RecipientTypes {
		statuses, ok := recipientStatuses[string(rt)]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedRecipientType, rt)
		}

		groupEmails, err := s.emailCampaignRepo.GetApplicantContactEmailsByStatus(ctx, sqlc.GetApplicantContactEmailsByStatusParams{
			HackathonID: campaign.HackathonID,
			Statuses:    statuses,
		})
		if err != nil {
			return nil, err
		}

		for _, email := range groupEmails {
			if _, exists := seen[email]; exists {
				continue
			}
			seen[email] = struct{}{}
			emails = append(emails, email)
		}
	}

	return emails, nil
}

// enqueueCampaignEmails sends one email per recipient so addresses are never
// exposed to each other, choosing the queue method by the campaign's format.
func (s *EmailCampaignService) enqueueCampaignEmails(campaign *sqlc.EmailCampaign, recipients []string) error {
	for _, recipient := range recipients {
		var err error
		switch campaign.Format {
		case sqlc.EmailCampaignFormatHtml:
			_, err = s.emailService.QueueSendRawHtmlEmail([]string{recipient}, campaign.Subject, campaign.Body)
		case sqlc.EmailCampaignFormatText:
			_, err = s.emailService.QueueSendTextEmail([]string{recipient}, campaign.Subject, campaign.Body)
		default:
			return fmt.Errorf("unsupported email format: %s", campaign.Format)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// SendCampaign loads a campaign, resolves its audience, enqueues the emails,
// and drives status: draft/scheduled -> sending -> sent (or failed).
func (s *EmailCampaignService) SendCampaign(
	ctx context.Context,
	campaignID uuid.UUID,
	hackathonID string,
	actorUserID uuid.UUID,
) (*sqlc.EmailCampaign, error) {
	campaign, err := s.emailCampaignRepo.GetEmailCampaignByID(ctx, sqlc.GetEmailCampaignByIDParams{
		ID:          campaignID,
		HackathonID: hackathonID,
	})
	if err != nil {
		return nil, err
	}

	if !canSendCampaign(campaign.Status) {
		return nil, ErrEmailCampaignCannotSend
	}

	recipients, err := s.resolveRecipients(ctx, campaign)
	if err != nil {
		return nil, err
	}
	if len(recipients) == 0 {
		return nil, ErrEmailCampaignNoRecipients
	}

	// Mark sending first, so a crash mid-send doesn't leave it looking like a draft.
	if _, err := s.emailCampaignRepo.UpdateEmailCampaignStatus(ctx, sqlc.UpdateEmailCampaignStatusParams{
		Status:                  sqlc.EmailCampaignStatusSending,
		UpdatedByUserIDDoUpdate: true,
		UpdatedByUserID:         actorUserID,
		ID:                      campaignID,
		HackathonID:             hackathonID,
	}); err != nil {
		return nil, err
	}

	if sendErr := s.enqueueCampaignEmails(campaign, recipients); sendErr != nil {
		errMsg := sendErr.Error()
		if _, uerr := s.emailCampaignRepo.UpdateEmailCampaignStatus(ctx, sqlc.UpdateEmailCampaignStatusParams{
			Status:                  sqlc.EmailCampaignStatusFailed,
			LastErrorDoUpdate:       true,
			LastError:               &errMsg,
			UpdatedByUserIDDoUpdate: true,
			UpdatedByUserID:         actorUserID,
			ID:                      campaignID,
			HackathonID:             hackathonID,
		}); uerr != nil {
			s.logger.Err(uerr).Msg("Failed to mark campaign as failed")
		}
		return nil, sendErr
	}

	now := time.Now()
	sent, err := s.emailCampaignRepo.UpdateEmailCampaignStatus(ctx, sqlc.UpdateEmailCampaignStatusParams{
		Status:                  sqlc.EmailCampaignStatusSent,
		SentAtDoUpdate:          true,
		SentAt:                  &now,
		UpdatedByUserIDDoUpdate: true,
		UpdatedByUserID:         actorUserID,
		ID:                      campaignID,
		HackathonID:             hackathonID,
	})
	if err != nil {
		return nil, err
	}

	return sent, nil
}
