package email

import (
	"context"
	"errors"
	"testing"

	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func TestValidateCampaignContent(t *testing.T) {
	tests := []struct {
		name           string
		title          string
		subject        string
		body           string
		recipientTypes []string
		expectedError  error
	}{
		{
			name:           "valid campaign",
			title:          "Welcome",
			subject:        "Welcome to SwampHacks",
			body:           "Campaign body",
			recipientTypes: []string{"admins"},
			expectedError:  nil,
		},
		{
			name:           "missing title",
			title:          "   ",
			subject:        "Subject",
			body:           "Body",
			recipientTypes: []string{"admins"},
			expectedError:  ErrEmailCampaignTitleRequired,
		},
		{
			name:           "missing subject",
			title:          "Title",
			subject:        "",
			body:           "Body",
			recipientTypes: []string{"admins"},
			expectedError:  ErrEmailCampaignSubjectRequired,
		},
		{
			name:           "missing body",
			title:          "Title",
			subject:        "Subject",
			body:           "   ",
			recipientTypes: []string{"admins"},
			expectedError:  ErrEmailCampaignBodyRequired,
		},
		{
			name:           "missing recipients",
			title:          "Title",
			subject:        "Subject",
			body:           "Body",
			recipientTypes: nil,
			expectedError:  ErrEmailCampaignRecipientsRequired,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateCampaignContent(
				test.title,
				test.subject,
				test.body,
				test.recipientTypes,
			)

			if !errors.Is(err, test.expectedError) {
				t.Fatalf("expected error %v, got %v", test.expectedError, err)
			}
		})
	}
}

func TestCanEditCampaign(t *testing.T) {
	tests := []struct {
		status   sqlc.EmailCampaignStatus
		expected bool
	}{
		{status: sqlc.EmailCampaignStatusDraft, expected: true},
		{status: sqlc.EmailCampaignStatusScheduled, expected: true},
		{status: sqlc.EmailCampaignStatusSending, expected: false},
		{status: sqlc.EmailCampaignStatusSent, expected: false},
		{status: sqlc.EmailCampaignStatusFailed, expected: false},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			result := canEditCampaign(test.status)

			if result != test.expected {
				t.Fatalf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestUpdateCampaignStatusRequiresScheduledAt(t *testing.T) {
	service := &EmailCampaignService{}

	_, err := service.UpdateCampaignStatus(
		context.Background(),
		sqlc.UpdateEmailCampaignStatusParams{
			Status:      sqlc.EmailCampaignStatusScheduled,
			ScheduledAt: nil,
		},
	)

	if !errors.Is(err, ErrEmailCampaignScheduledAtRequired) {
		t.Fatalf(
			"expected %v, got %v",
			ErrEmailCampaignScheduledAtRequired,
			err,
		)
	}
}

func TestUpdateCampaignStatusRequiresSentAt(t *testing.T) {
	service := &EmailCampaignService{}

	_, err := service.UpdateCampaignStatus(
		context.Background(),
		sqlc.UpdateEmailCampaignStatusParams{
			Status: sqlc.EmailCampaignStatusSent,
			SentAt: nil,
		},
	)

	if !errors.Is(err, ErrEmailCampaignSentAtRequired) {
		t.Fatalf(
			"expected %v, got %v",
			ErrEmailCampaignSentAtRequired,
			err,
		)
	}
}