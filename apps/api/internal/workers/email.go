package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/domains/email"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type EmailWorker struct {
	emailService *email.EmailService
	logger       zerolog.Logger
}

func NewEmailWorker(emailService *email.EmailService, logger zerolog.Logger) *EmailWorker {
	return &EmailWorker{
		emailService: emailService,
		logger:       logger.With().Str("worker", "EmailWorker").Logger(),
	}
}

func (w *EmailWorker) HandleSendTextEmailTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.SendTextEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		w.logger.Err(err)
		return fmt.Errorf("HandleSendTextEmailTask: json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := w.emailService.SendTextEmail(p.To, p.Subject, p.Body); err != nil {
		w.logger.Err(err).Msg("Failed to send text email from worker")
		return err
	}
	return nil
}

func (w *EmailWorker) HandleSendHtmlEmailTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.SendHtmlEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		w.logger.Err(err)
		return fmt.Errorf("HandleSendHtmlEmailTask: json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := w.emailService.SendHtmlEmail(p.To, p.Subject, p.TemplateData, p.TemplateFilePath); err != nil {
		w.logger.Err(err).Msg("Failed to send ConfirmationEmail from worker")
		return err
	}
	return nil
}

func (w *EmailWorker) HandleSendRawHtmlEmailTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.SendRawHtmlEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		w.logger.Err(err)
		return fmt.Errorf("HandleSendRawHtmlEmailTask: json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := w.emailService.SendRawHtmlEmail(p.To, p.Subject, p.Body); err != nil {
		w.logger.Err(err).Msg("Failed to send raw HTML email from worker")
		return err
	}
	return nil
}
