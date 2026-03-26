package email

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/emailutils"
)

func RegisterRoutes(emailHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "queue-text-email",
		Method:      http.MethodPost,
		Summary:     "Queue Text Email",
		Description: "Pushes a text email request to the task queue",
		Tags:        []string{"Email"},
		Path:        "/queue-text-email",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusUnauthorized},
	}, emailHandler.handleQueueTextEmail)

	huma.Register(group, huma.Operation{
		OperationID: "queue-confirmation-email",
		Method:      http.MethodPost,
		Summary:     "Queue Confirmation Email",
		Description: "Pushes a confirmation email request to the task queue",
		Tags:        []string{"Email"},
		Path:        "/queue-confirmation-email",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusUnauthorized},
	}, emailHandler.handleQueueConfirmationEmail)

	huma.Register(group, huma.Operation{
		OperationID: "queue-welcome-email",
		Method:      http.MethodPost,
		Summary:     "Queue Welcome Email",
		Description: "Pushes a welcome email request to the task queue",
		Tags:        []string{"Email"},
		Path:        "/queue-welcome-email",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusUnauthorized},
	}, emailHandler.handleQueueWelcomeEmail)

	huma.Register(group, huma.Operation{
		OperationID: "send-welcome-emails",
		Method:      http.MethodPost,
		Summary:     "Send Welcome Emails",
		Description: "Send welcome emails to all attendees",
		Tags:        []string{"Email"},
		Path:        "/send-welcome-emails",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusUnauthorized},
	}, emailHandler.handleSendWelcomeEmails)
}

type handler struct {
	emailService *EmailService
	logger       zerolog.Logger
}

func NewHandler(emailService *EmailService, logger zerolog.Logger) *handler {
	return &handler{
		emailService: emailService,
		logger:       logger.With().Str("handler", "EmailHandler").Str("domain", "email").Logger(),
	}
}

type QueueTextEmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject" minLength:"1"`
	Body    string   `json:"body" minLength:"1"`
}

type QueueTextEmailOutput struct {
	Status int
}

func (h *handler) handleQueueTextEmail(ctx context.Context, input *struct {
	Body QueueTextEmailRequest
}) (*QueueTextEmailOutput, error) {
	for _, to := range input.Body.To {
		if !emailutils.IsValidEmail(to) {
			return nil, huma.Error400BadRequest("Invalid email(s)")
		}
	}

	taskInfo, err := h.emailService.QueueSendTextEmail(input.Body.To, input.Body.Subject, input.Body.Body)

	if err != nil {
		h.logger.Err(err).Msg("Failed to queue SendTextEmail from EmailHandler")
		return nil, huma.Error500InternalServerError("Failed to queue text email")
	}

	h.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued SendTextEmail task!")

	return &QueueTextEmailOutput{Status: http.StatusOK}, nil
}

type QueueConfirmationEmailRequest struct {
	Email     string `json:"email" required:"true"`
	FirstName string `json:"firstName" required:"true"`
}

type QueueConfirmationEmailOutput struct {
	Status int
}

func (h *handler) handleQueueConfirmationEmail(ctx context.Context, input *struct {
	Body QueueConfirmationEmailRequest
}) (*QueueConfirmationEmailOutput, error) {
	err := h.emailService.QueueConfirmationEmail(input.Body.Email, input.Body.FirstName)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to queue confirmation email")
	}

	return &QueueConfirmationEmailOutput{Status: http.StatusOK}, nil
}

type QueueWelcomeEmailRequest struct {
	Email       string `json:"email" required:"true"`
	FirstName   string `json:"firstName" required:"true"`
	RecipientId string `json:"recipientId" required:"true"`
}

type QueueWelcomeEmailOutput struct {
	Status int
}

func (h *handler) handleQueueWelcomeEmail(ctx context.Context, input *struct {
	Body QueueWelcomeEmailRequest
}) (*QueueWelcomeEmailOutput, error) {
	recipientId, err := uuid.Parse(input.Body.RecipientId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid recipient id")
	}

	err = h.emailService.QueueWelcomeEmail(ctx, input.Body.Email, input.Body.FirstName, recipientId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to queue welcome email")
	}

	return &QueueWelcomeEmailOutput{Status: http.StatusOK}, nil
}

type SendWelcomeEmailsOutput struct {
	Status int
}

func (h *handler) handleSendWelcomeEmails(ctx context.Context, input *struct{}) (*SendWelcomeEmailsOutput, error) {
	err := h.emailService.SendWelcomeEmailToAttendees(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to send welcome emails")
	}

	return &SendWelcomeEmailsOutput{Status: http.StatusOK}, nil
}
