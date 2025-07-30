package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type EmailWorker struct {
	emailService *services.EmailService
	logger       zerolog.Logger
}

func NewEmailWorker(emailService *services.EmailService, logger zerolog.Logger) *EmailWorker {
	return &EmailWorker{
		emailService: emailService,
		logger:       logger.With().Str("worker", "EmailWorker").Str("component", "email").Logger(),
	}
}

func (w *EmailWorker) HandleSendEmailTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.SendEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		w.logger.Err(err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := w.emailService.SendEmail(p.To, p.From, p.Body); err != nil {
		w.logger.Err(err).Msg("Failed to send email from worker")
	}
	return nil
}
