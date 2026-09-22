package email

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

// campaignStore is the data access EmailCampaignService needs.
// *repository.EmailCampaignRepository satisfies it in production; tests use a fake.
type campaignStore interface {
	CreateEmailCampaign(ctx context.Context, params sqlc.CreateEmailCampaignParams) (*sqlc.EmailCampaign, error)
	GetEmailCampaignByID(ctx context.Context, params sqlc.GetEmailCampaignByIDParams) (*sqlc.EmailCampaign, error)
	ListEmailCampaigns(ctx context.Context, hackathonID string) ([]sqlc.EmailCampaign, error)
	UpdateEmailCampaign(ctx context.Context, params sqlc.UpdateEmailCampaignParams) (*sqlc.EmailCampaign, error)
	UpdateEmailCampaignStatus(ctx context.Context, params sqlc.UpdateEmailCampaignStatusParams) (*sqlc.EmailCampaign, error)
	DeleteEmailCampaign(ctx context.Context, params sqlc.DeleteEmailCampaignParams) error
	ClaimCampaignForSending(ctx context.Context, params sqlc.ClaimCampaignForSendingParams) (*sqlc.EmailCampaign, error)
	ListDueScheduledCampaigns(ctx context.Context) ([]sqlc.EmailCampaign, error)
	GetApplicantContactEmailsByStatus(ctx context.Context, params sqlc.GetApplicantContactEmailsByStatusParams) ([]string, error)
	GetUserContactEmailsByRoles(ctx context.Context, roles []string) ([]string, error)
	GetInterestSubscriberEmails(ctx context.Context, hackathonID string) ([]string, error)
}

// campaignMailer is the send side EmailCampaignService needs.
// *EmailService satisfies it in production; tests use a fake.
type campaignMailer interface {
	QueueSendTextEmail(to []string, subject string, body string) (*asynq.TaskInfo, error)
	QueueSendRawHtmlEmail(to []string, subject string, body string) (*asynq.TaskInfo, error)
}

var (
	//Reuses repository-level "not found" error
	ErrEmailCampaignNotFound = repository.ErrEmailCampaignNotFound

	//validation errors before writing error to db
	ErrEmailCampaignTitleRequired      = errors.New("email campaign title is required")
	ErrEmailCampaignSubjectRequired    = errors.New("email campaign subject is required")
	ErrEmailCampaignBodyRequired       = errors.New("email campaign body is required")
	ErrEmailCampaignRecipientsRequired = errors.New("email campaign recipients are required")

	ErrEmailCampaignCannotEdit   = errors.New("email campaign cannot be edited")
	ErrEmailCampaignCannotDelete = errors.New("email campaign cannot be deleted after it has been sent")

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

// recipientRoles maps a role-based recipient_type to the user roles it covers.
// Roles live on the users table globally, so these groups are not hackathon-scoped.
var recipientRoles = map[string][]string{
	"admins":   {"admin"},
	"staff":    {"staff"},
	"visitors": {"visitor"},
}

// recipientTypeInterestSubscribers reads from the public interest form rather
// than from users or applications, so it gets its own branch.
const recipientTypeInterestSubscribers = "interest_subscribers"

// claimableSendStatuses are the states a campaign can be claimed for sending from.
var claimableSendStatuses = []string{"draft", "scheduled"}

// scheduledSendGrace is how late a scheduled campaign may be and still go out.
// Beyond this it is failed instead, so an outage cannot trigger a surprise blast.
const scheduledSendGrace = 2 * time.Hour

// EmailCampaignService owns business rules for saved email campaigns.
type EmailCampaignService struct {
	emailCampaignRepo campaignStore
	emailService      campaignMailer
	logger            zerolog.Logger
}

// NewEmailCampaignService creates the service and stores its dependencies.
// This will eventually be called from api.go when wiring the app together.
func NewEmailCampaignService(
	emailCampaignRepo campaignStore,
	emailService campaignMailer,
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

	// A draft has no send time, so unscheduling clears it. The handler cannot do
	// this itself: ScheduledAt is a *time.Time, so an omitted field and an explicit
	// null both arrive as nil and it has to assume "leave it alone".
	if params.Status == sqlc.EmailCampaignStatusDraft {
		params.ScheduledAtDoUpdate = true
		params.ScheduledAt = nil
	}

	return s.emailCampaignRepo.UpdateEmailCampaignStatus(ctx, params)
}

// DeleteCampaign removes a campaign after checking it is still deletable.
// It loads the campaign first so the status guard runs before anything is destroyed.
func (s *EmailCampaignService) DeleteCampaign(
	ctx context.Context,
	params sqlc.DeleteEmailCampaignParams,
) error {
	existingCampaign, err := s.emailCampaignRepo.GetEmailCampaignByID(ctx, sqlc.GetEmailCampaignByIDParams{
		ID:          params.ID,
		HackathonID: params.HackathonID,
	})
	if err != nil {
		return err
	}

	if !canDeleteCampaign(existingCampaign.Status) {
		return ErrEmailCampaignCannotDelete
	}

	return s.emailCampaignRepo.DeleteEmailCampaign(ctx, params)
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

// canDeleteCampaign keeps delivered mail auditable: a campaign that went out
// (or is going out) stays on the record, so only pre-send and failed ones are removable.
func canDeleteCampaign(status sqlc.EmailCampaignStatus) bool {
	return status == sqlc.EmailCampaignStatusDraft ||
		status == sqlc.EmailCampaignStatusScheduled ||
		status == sqlc.EmailCampaignStatusFailed
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
		var groupEmails []string
		var err error

		switch {
		case recipientStatuses[string(rt)] != nil:
			groupEmails, err = s.emailCampaignRepo.GetApplicantContactEmailsByStatus(ctx, sqlc.GetApplicantContactEmailsByStatusParams{
				HackathonID: campaign.HackathonID,
				Statuses:    recipientStatuses[string(rt)],
			})
		case recipientRoles[string(rt)] != nil:
			groupEmails, err = s.emailCampaignRepo.GetUserContactEmailsByRoles(ctx, recipientRoles[string(rt)])
		case string(rt) == recipientTypeInterestSubscribers:
			groupEmails, err = s.emailCampaignRepo.GetInterestSubscriberEmails(ctx, campaign.HackathonID)
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedRecipientType, rt)
		}

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

// SendCampaign sends a campaign on an admin's behalf.
// The campaign is claimed atomically, so a double-click or two concurrent
// requests cannot both dispatch it.
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

	// Claim before dispatching: this is the only write that can race, and losing
	// it means someone else already started this send.
	claimed, err := s.emailCampaignRepo.ClaimCampaignForSending(ctx, sqlc.ClaimCampaignForSendingParams{
		ID:                      campaignID,
		HackathonID:             hackathonID,
		FromStatuses:            claimableSendStatuses,
		UpdatedByUserIDDoUpdate: true,
		UpdatedByUserID:         actorUserID,
	})
	if errors.Is(err, repository.ErrEmailCampaignNotClaimable) {
		return nil, ErrEmailCampaignCannotSend
	}
	if err != nil {
		return nil, err
	}

	return s.dispatchClaimedCampaign(ctx, claimed, recipients, &actorUserID)
}

// dispatchClaimedCampaign enqueues the mail for a campaign already moved into
// "sending", then records the outcome. Both the admin path and the scheduler
// share it so the status bookkeeping lives in exactly one place.
func (s *EmailCampaignService) dispatchClaimedCampaign(
	ctx context.Context,
	campaign *sqlc.EmailCampaign,
	recipients []string,
	actorUserID *uuid.UUID,
) (*sqlc.EmailCampaign, error) {
	if sendErr := s.enqueueCampaignEmails(campaign, recipients); sendErr != nil {
		if _, uerr := s.markCampaignFailed(ctx, campaign, sendErr.Error(), actorUserID); uerr != nil {
			s.logger.Err(uerr).Msg("Failed to mark campaign as failed")
		}
		return nil, sendErr
	}

	now := time.Now()
	sent, err := s.emailCampaignRepo.UpdateEmailCampaignStatus(ctx, sqlc.UpdateEmailCampaignStatusParams{
		Status:                  sqlc.EmailCampaignStatusSent,
		SentAtDoUpdate:          true,
		SentAt:                  &now,
		UpdatedByUserIDDoUpdate: actorUserID != nil,
		UpdatedByUserID:         derefUserID(actorUserID),
		ID:                      campaign.ID,
		HackathonID:             campaign.HackathonID,
	})
	if err != nil {
		return nil, err
	}

	return sent, nil
}

// markCampaignFailed records a failure and the reason on the campaign row.
func (s *EmailCampaignService) markCampaignFailed(
	ctx context.Context,
	campaign *sqlc.EmailCampaign,
	reason string,
	actorUserID *uuid.UUID,
) (*sqlc.EmailCampaign, error) {
	return s.emailCampaignRepo.UpdateEmailCampaignStatus(ctx, sqlc.UpdateEmailCampaignStatusParams{
		Status:                  sqlc.EmailCampaignStatusFailed,
		LastErrorDoUpdate:       true,
		LastError:               &reason,
		UpdatedByUserIDDoUpdate: actorUserID != nil,
		UpdatedByUserID:         derefUserID(actorUserID),
		ID:                      campaign.ID,
		HackathonID:             campaign.HackathonID,
	})
}

// derefUserID keeps the generated params happy for scheduler-driven sends,
// which have no acting user; the paired DoUpdate flag is false in that case.
func derefUserID(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.UUID{}
	}
	return *id
}

// SendDueScheduledCampaigns is the scheduler sweep: it finds campaigns whose
// send time has arrived and dispatches each one. It returns the number sent.
//
// One campaign failing never aborts the sweep, otherwise a single bad campaign
// would block every later one indefinitely.
func (s *EmailCampaignService) SendDueScheduledCampaigns(ctx context.Context) (int, error) {
	due, err := s.emailCampaignRepo.ListDueScheduledCampaigns(ctx)
	if err != nil {
		return 0, err
	}

	sentCount := 0
	for i := range due {
		campaign := due[i]
		if err := s.sendScheduledCampaign(ctx, &campaign); err != nil {
			s.logger.Err(err).
				Str("campaignID", campaign.ID.String()).
				Str("hackathonID", campaign.HackathonID).
				Msg("Failed to send scheduled campaign")
			continue
		}
		sentCount++
	}

	return sentCount, nil
}

// sendScheduledCampaign handles one due campaign: it drops campaigns that are
// too far past their slot, resolves recipients, claims the row, then dispatches.
func (s *EmailCampaignService) sendScheduledCampaign(ctx context.Context, campaign *sqlc.EmailCampaign) error {
	if campaign.ScheduledAt == nil {
		return nil
	}

	// A long outage must not trigger a surprise blast hours after the fact.
	if lateness := time.Since(*campaign.ScheduledAt); lateness > scheduledSendGrace {
		reason := fmt.Sprintf("missed its scheduled send by %s", lateness.Round(time.Minute))
		if _, err := s.markCampaignFailed(ctx, campaign, reason, nil); err != nil {
			return err
		}
		s.logger.Warn().
			Str("campaignID", campaign.ID.String()).
			Dur("lateBy", lateness).
			Msg("Scheduled campaign expired past its grace window")
		return nil
	}

	recipients, err := s.resolveRecipients(ctx, campaign)
	if err != nil {
		if _, ferr := s.markCampaignFailed(ctx, campaign, err.Error(), nil); ferr != nil {
			s.logger.Err(ferr).Msg("Failed to mark campaign as failed")
		}
		return err
	}
	if len(recipients) == 0 {
		if _, ferr := s.markCampaignFailed(ctx, campaign, ErrEmailCampaignNoRecipients.Error(), nil); ferr != nil {
			s.logger.Err(ferr).Msg("Failed to mark campaign as failed")
		}
		return ErrEmailCampaignNoRecipients
	}

	// Claiming is what makes an overlapping tick harmless: only one wins the row.
	claimed, err := s.emailCampaignRepo.ClaimCampaignForSending(ctx, sqlc.ClaimCampaignForSendingParams{
		ID:           campaign.ID,
		HackathonID:  campaign.HackathonID,
		FromStatuses: []string{"scheduled"},
	})
	if errors.Is(err, repository.ErrEmailCampaignNotClaimable) {
		// Another tick already picked it up; nothing to do and nothing wrong.
		return nil
	}
	if err != nil {
		return err
	}

	_, err = s.dispatchClaimedCampaign(ctx, claimed, recipients, nil)
	return err
}
