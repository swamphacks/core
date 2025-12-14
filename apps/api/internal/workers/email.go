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

func (w *EmailWorker) HandleSendConfirmationEmailTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.SendConfirmationEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		w.logger.Err(err)
		return fmt.Errorf("HandleSendConfirmationEmailTask: json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := w.emailService.SendConfirmationEmail(p.To, p.Name); err != nil {
		w.logger.Err(err).Msg("Failed to send ConfirmationEmail from worker")
		return err
	}
	return nil
}
