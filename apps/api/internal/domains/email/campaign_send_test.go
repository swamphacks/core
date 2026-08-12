package email

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

// fakeCampaignStore is an in-memory campaignStore. Only the methods SendCampaign
// exercises carry behaviour; the rest exist to satisfy the interface.
type fakeCampaignStore struct {
	campaign         *sqlc.EmailCampaign
	getErr           error
	applicantEmails  []string
	roleEmails       []string
	resolveErr       error
	statusCalls      []sqlc.UpdateEmailCampaignStatusParams
	deleted          bool
	dueCampaigns     []sqlc.EmailCampaign
	claimErr         error
	claimCalls       int
	subscriberEmails []string
}

func (f *fakeCampaignStore) GetEmailCampaignByID(ctx context.Context, params sqlc.GetEmailCampaignByIDParams) (*sqlc.EmailCampaign, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.campaign, nil
}

func (f *fakeCampaignStore) GetApplicantContactEmailsByStatus(ctx context.Context, params sqlc.GetApplicantContactEmailsByStatusParams) ([]string, error) {
	return f.applicantEmails, f.resolveErr
}

func (f *fakeCampaignStore) GetUserContactEmailsByRoles(ctx context.Context, roles []string) ([]string, error) {
	return f.roleEmails, f.resolveErr
}

func (f *fakeCampaignStore) UpdateEmailCampaignStatus(ctx context.Context, params sqlc.UpdateEmailCampaignStatusParams) (*sqlc.EmailCampaign, error) {
	f.statusCalls = append(f.statusCalls, params)
	updated := *f.campaign
	updated.Status = params.Status
	return &updated, nil
}

func (f *fakeCampaignStore) DeleteEmailCampaign(ctx context.Context, params sqlc.DeleteEmailCampaignParams) error {
	f.deleted = true
	return nil
}

func (f *fakeCampaignStore) ClaimCampaignForSending(ctx context.Context, params sqlc.ClaimCampaignForSendingParams) (*sqlc.EmailCampaign, error) {
	f.claimCalls++
	if f.claimErr != nil {
		return nil, f.claimErr
	}
	claimed := *f.campaign
	claimed.Status = sqlc.EmailCampaignStatusSending
	return &claimed, nil
}

func (f *fakeCampaignStore) ListDueScheduledCampaigns(ctx context.Context) ([]sqlc.EmailCampaign, error) {
	return f.dueCampaigns, nil
}

func (f *fakeCampaignStore) GetInterestSubscriberEmails(ctx context.Context, hackathonID string) ([]string, error) {
	return f.subscriberEmails, f.resolveErr
}

func (f *fakeCampaignStore) CreateEmailCampaign(ctx context.Context, params sqlc.CreateEmailCampaignParams) (*sqlc.EmailCampaign, error) {
	return nil, nil
}
func (f *fakeCampaignStore) ListEmailCampaigns(ctx context.Context, hackathonID string) ([]sqlc.EmailCampaign, error) {
	return nil, nil
}
func (f *fakeCampaignStore) UpdateEmailCampaign(ctx context.Context, params sqlc.UpdateEmailCampaignParams) (*sqlc.EmailCampaign, error) {
	return nil, nil
}

// fakeMailer records what would have been queued instead of touching Redis or SES.
type fakeMailer struct {
	textRecipients []string
	htmlRecipients []string
	err            error
}

func (f *fakeMailer) QueueSendTextEmail(to []string, subject string, body string) (*asynq.TaskInfo, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.textRecipients = append(f.textRecipients, to...)
	return &asynq.TaskInfo{}, nil
}

func (f *fakeMailer) QueueSendRawHtmlEmail(to []string, subject string, body string) (*asynq.TaskInfo, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.htmlRecipients = append(f.htmlRecipients, to...)
	return &asynq.TaskInfo{}, nil
}

func newTestCampaign(status sqlc.EmailCampaignStatus, format sqlc.EmailCampaignFormat, types ...sqlc.EmailRecipientType) *sqlc.EmailCampaign {
	return &sqlc.EmailCampaign{
		ID:             uuid.New(),
		HackathonID:    "swamphacks-xii",
		Title:          "Announcement",
		Subject:        "Hello",
		Body:           "Body",
		Format:         format,
		RecipientTypes: types,
		Status:         status,
	}
}

func newTestService(store *fakeCampaignStore, mailer *fakeMailer) *EmailCampaignService {
	return &EmailCampaignService{emailCampaignRepo: store, emailService: mailer}
}

func TestSendCampaignQueuesOneEmailPerRecipient(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "accepted_applicants")
	store := &fakeCampaignStore{campaign: campaign, applicantEmails: []string{"a@ufl.edu", "b@ufl.edu"}}
	mailer := &fakeMailer{}

	sent, err := newTestService(store, mailer).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New())
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if len(mailer.textRecipients) != 2 {
		t.Fatalf("expected 2 queued emails, got %d (%v)", len(mailer.textRecipients), mailer.textRecipients)
	}
	if len(mailer.htmlRecipients) != 0 {
		t.Fatalf("text campaign must not queue html emails, got %v", mailer.htmlRecipients)
	}
	if sent.Status != sqlc.EmailCampaignStatusSent {
		t.Fatalf("expected final status sent, got %s", sent.Status)
	}
}

func TestSendCampaignMarksSendingBeforeSent(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "accepted_applicants")
	store := &fakeCampaignStore{campaign: campaign, applicantEmails: []string{"a@ufl.edu"}}

	if _, err := newTestService(store, &fakeMailer{}).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New()); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	// The move into "sending" is the atomic claim, not a plain status update.
	if store.claimCalls != 1 {
		t.Fatalf("expected the campaign to be claimed exactly once, got %d", store.claimCalls)
	}
	if len(store.statusCalls) != 1 {
		t.Fatalf("expected 1 status update after the claim, got %d", len(store.statusCalls))
	}
	if store.statusCalls[0].Status != sqlc.EmailCampaignStatusSent {
		t.Fatalf("expected final status sent, got %s", store.statusCalls[0].Status)
	}
	if store.statusCalls[0].SentAt == nil {
		t.Fatal("expected sent_at to be set when marking sent")
	}
}

func TestSendCampaignUsesHtmlQueueForHtmlFormat(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatHtml, "accepted_applicants")
	store := &fakeCampaignStore{campaign: campaign, applicantEmails: []string{"a@ufl.edu"}}
	mailer := &fakeMailer{}

	if _, err := newTestService(store, mailer).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New()); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if len(mailer.htmlRecipients) != 1 || len(mailer.textRecipients) != 0 {
		t.Fatalf("html campaign routed wrong: html=%v text=%v", mailer.htmlRecipients, mailer.textRecipients)
	}
}

func TestSendCampaignDeduplicatesAcrossGroups(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "accepted_applicants", "admins")
	store := &fakeCampaignStore{
		campaign:        campaign,
		applicantEmails: []string{"shared@ufl.edu", "applicant@ufl.edu"},
		roleEmails:      []string{"shared@ufl.edu", "admin@ufl.edu"},
	}
	mailer := &fakeMailer{}

	if _, err := newTestService(store, mailer).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New()); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if len(mailer.textRecipients) != 3 {
		t.Fatalf("expected 3 unique recipients, got %d (%v)", len(mailer.textRecipients), mailer.textRecipients)
	}
}

func TestSendCampaignRejectsAlreadySentCampaign(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusSent, sqlc.EmailCampaignFormatText, "accepted_applicants")
	store := &fakeCampaignStore{campaign: campaign, applicantEmails: []string{"a@ufl.edu"}}
	mailer := &fakeMailer{}

	_, err := newTestService(store, mailer).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New())
	if !errors.Is(err, ErrEmailCampaignCannotSend) {
		t.Fatalf("expected %v, got %v", ErrEmailCampaignCannotSend, err)
	}
	if len(mailer.textRecipients) != 0 || len(store.statusCalls) != 0 {
		t.Fatal("a non-sendable campaign must not queue mail or change status")
	}
}

func TestSendCampaignFailsWhenNoRecipientsResolve(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "accepted_applicants")
	store := &fakeCampaignStore{campaign: campaign, applicantEmails: []string{}}
	mailer := &fakeMailer{}

	_, err := newTestService(store, mailer).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New())
	if !errors.Is(err, ErrEmailCampaignNoRecipients) {
		t.Fatalf("expected %v, got %v", ErrEmailCampaignNoRecipients, err)
	}
	if len(store.statusCalls) != 0 {
		t.Fatal("an empty audience must not move the campaign to sending")
	}
}

func TestSendCampaignRejectsUnsupportedRecipientType(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "attendees_not_a_real_group")
	store := &fakeCampaignStore{campaign: campaign}

	_, err := newTestService(store, &fakeMailer{}).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New())
	if !errors.Is(err, ErrUnsupportedRecipientType) {
		t.Fatalf("expected %v, got %v", ErrUnsupportedRecipientType, err)
	}
}

func TestSendCampaignMarksFailedWhenQueueingFails(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "accepted_applicants")
	store := &fakeCampaignStore{campaign: campaign, applicantEmails: []string{"a@ufl.edu"}}
	queueErr := errors.New("redis unavailable")

	_, err := newTestService(store, &fakeMailer{err: queueErr}).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New())
	if !errors.Is(err, queueErr) {
		t.Fatalf("expected the queue error to surface, got %v", err)
	}

	if store.claimCalls != 1 {
		t.Fatalf("expected the campaign to be claimed before dispatch, got %d claims", store.claimCalls)
	}
	if len(store.statusCalls) != 1 {
		t.Fatalf("expected a single failed update after the claim, got %d", len(store.statusCalls))
	}
	last := store.statusCalls[0]
	if last.Status != sqlc.EmailCampaignStatusFailed {
		t.Fatalf("expected final status failed, got %s", last.Status)
	}
	if last.LastError == nil || *last.LastError != queueErr.Error() {
		t.Fatalf("expected last_error to record the failure, got %v", last.LastError)
	}
}

func TestDeleteCampaignRejectsSentCampaign(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusSent, sqlc.EmailCampaignFormatText, "accepted_applicants")
	store := &fakeCampaignStore{campaign: campaign}

	err := newTestService(store, &fakeMailer{}).DeleteCampaign(context.Background(), sqlc.DeleteEmailCampaignParams{
		ID:          campaign.ID,
		HackathonID: campaign.HackathonID,
	})

	if !errors.Is(err, ErrEmailCampaignCannotDelete) {
		t.Fatalf("expected %v, got %v", ErrEmailCampaignCannotDelete, err)
	}
	if store.deleted {
		t.Fatal("a sent campaign must never be removed from the record")
	}
}

func TestDeleteCampaignAllowsDraftAndFailed(t *testing.T) {
	for _, status := range []sqlc.EmailCampaignStatus{
		sqlc.EmailCampaignStatusDraft,
		sqlc.EmailCampaignStatusScheduled,
		sqlc.EmailCampaignStatusFailed,
	} {
		campaign := newTestCampaign(status, sqlc.EmailCampaignFormatText, "accepted_applicants")
		store := &fakeCampaignStore{campaign: campaign}

		err := newTestService(store, &fakeMailer{}).DeleteCampaign(context.Background(), sqlc.DeleteEmailCampaignParams{
			ID:          campaign.ID,
			HackathonID: campaign.HackathonID,
		})

		if err != nil {
			t.Fatalf("status %s: expected delete to succeed, got %v", status, err)
		}
		if !store.deleted {
			t.Fatalf("status %s: expected the campaign to be deleted", status)
		}
	}
}

func scheduledAt(d time.Duration) *time.Time {
	t := time.Now().Add(d)
	return &t
}

func TestSweepSendsDueScheduledCampaign(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusScheduled, sqlc.EmailCampaignFormatText, "accepted_applicants")
	campaign.ScheduledAt = scheduledAt(-5 * time.Minute)
	store := &fakeCampaignStore{
		campaign:        campaign,
		dueCampaigns:    []sqlc.EmailCampaign{*campaign},
		applicantEmails: []string{"a@ufl.edu", "b@ufl.edu"},
	}
	mailer := &fakeMailer{}

	sent, err := newTestService(store, mailer).SendDueScheduledCampaigns(context.Background())
	if err != nil {
		t.Fatalf("expected sweep to succeed, got %v", err)
	}
	if sent != 1 {
		t.Fatalf("expected 1 campaign sent, got %d", sent)
	}
	if len(mailer.textRecipients) != 2 {
		t.Fatalf("expected 2 queued emails, got %v", mailer.textRecipients)
	}
}

func TestSweepFailsCampaignPastGraceWindow(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusScheduled, sqlc.EmailCampaignFormatText, "accepted_applicants")
	campaign.ScheduledAt = scheduledAt(-6 * time.Hour)
	store := &fakeCampaignStore{
		campaign:        campaign,
		dueCampaigns:    []sqlc.EmailCampaign{*campaign},
		applicantEmails: []string{"a@ufl.edu"},
	}
	mailer := &fakeMailer{}

	sent, err := newTestService(store, mailer).SendDueScheduledCampaigns(context.Background())
	if err != nil {
		t.Fatalf("expected sweep to succeed, got %v", err)
	}
	if sent != 1 {
		t.Fatalf("an expired campaign should be handled, not error: got %d", sent)
	}
	if len(mailer.textRecipients) != 0 {
		t.Fatalf("expired campaign must not send mail, got %v", mailer.textRecipients)
	}
	if store.claimCalls != 0 {
		t.Fatal("expired campaign must never be claimed")
	}
	if len(store.statusCalls) != 1 || store.statusCalls[0].Status != sqlc.EmailCampaignStatusFailed {
		t.Fatalf("expected the campaign to be marked failed, got %+v", store.statusCalls)
	}
	if store.statusCalls[0].LastError == nil {
		t.Fatal("expected last_error to explain the missed schedule")
	}
}

func TestSweepSkipsCampaignClaimedByAnotherTick(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusScheduled, sqlc.EmailCampaignFormatText, "accepted_applicants")
	campaign.ScheduledAt = scheduledAt(-1 * time.Minute)
	store := &fakeCampaignStore{
		campaign:        campaign,
		dueCampaigns:    []sqlc.EmailCampaign{*campaign},
		applicantEmails: []string{"a@ufl.edu"},
		claimErr:        repository.ErrEmailCampaignNotClaimable,
	}
	mailer := &fakeMailer{}

	if _, err := newTestService(store, mailer).SendDueScheduledCampaigns(context.Background()); err != nil {
		t.Fatalf("losing the claim race is not an error, got %v", err)
	}
	if len(mailer.textRecipients) != 0 {
		t.Fatalf("a campaign claimed elsewhere must not be sent again, got %v", mailer.textRecipients)
	}
}

func TestSendCampaignSurfacesLostClaimAsCannotSend(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "accepted_applicants")
	store := &fakeCampaignStore{
		campaign:        campaign,
		applicantEmails: []string{"a@ufl.edu"},
		claimErr:        repository.ErrEmailCampaignNotClaimable,
	}

	_, err := newTestService(store, &fakeMailer{}).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New())
	if !errors.Is(err, ErrEmailCampaignCannotSend) {
		t.Fatalf("expected %v, got %v", ErrEmailCampaignCannotSend, err)
	}
}

func TestSendCampaignResolvesInterestSubscribers(t *testing.T) {
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "interest_subscribers")
	store := &fakeCampaignStore{campaign: campaign, subscriberEmails: []string{"curious@gmail.com", "maybe@gmail.com"}}
	mailer := &fakeMailer{}

	if _, err := newTestService(store, mailer).SendCampaign(context.Background(), campaign.ID, campaign.HackathonID, uuid.New()); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if len(mailer.textRecipients) != 2 {
		t.Fatalf("expected 2 subscriber emails, got %v", mailer.textRecipients)
	}
}
