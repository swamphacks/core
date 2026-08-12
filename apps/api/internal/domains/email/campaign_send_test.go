package email

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

// fakeCampaignStore is an in-memory campaignStore. Only the methods SendCampaign
// exercises carry behaviour; the rest exist to satisfy the interface.
type fakeCampaignStore struct {
	campaign        *sqlc.EmailCampaign
	getErr          error
	applicantEmails []string
	roleEmails      []string
	resolveErr      error
	statusCalls     []sqlc.UpdateEmailCampaignStatusParams
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

	if len(store.statusCalls) != 2 {
		t.Fatalf("expected 2 status updates, got %d", len(store.statusCalls))
	}
	if store.statusCalls[0].Status != sqlc.EmailCampaignStatusSending {
		t.Fatalf("expected first update to be sending, got %s", store.statusCalls[0].Status)
	}
	if store.statusCalls[1].Status != sqlc.EmailCampaignStatusSent {
		t.Fatalf("expected second update to be sent, got %s", store.statusCalls[1].Status)
	}
	if store.statusCalls[1].SentAt == nil {
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
	campaign := newTestCampaign(sqlc.EmailCampaignStatusDraft, sqlc.EmailCampaignFormatText, "interest_subscribers")
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

	if len(store.statusCalls) != 2 {
		t.Fatalf("expected sending then failed, got %d updates", len(store.statusCalls))
	}
	last := store.statusCalls[1]
	if last.Status != sqlc.EmailCampaignStatusFailed {
		t.Fatalf("expected final status failed, got %s", last.Status)
	}
	if last.LastError == nil || *last.LastError != queueErr.Error() {
		t.Fatalf("expected last_error to record the failure, got %v", last.LastError)
	}
}
