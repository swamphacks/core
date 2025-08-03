package services

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type EmailService struct {
	logger    zerolog.Logger
	taskQueue *asynq.Client
	// Add SMTP client later
}

func NewEmailService(taskQueue *asynq.Client, logger zerolog.Logger) *EmailService {
	return &EmailService{
		logger:    logger.With().Str("service", "EmailService").Str("component", "email").Logger(),
		taskQueue: taskQueue,
	}
}

func (s *EmailService) SendEmail(to, from, body string) error {
	s.logger.Info().Msgf("Sending from %s to %s with body %s", from, to, body)
	return nil
}

func (s *EmailService) QueueSendEmail(to, from, body string) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskSendEmail(tasks.SendEmailPayload{
		To:   to,
		From: from,
		Body: body,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to create Send Email task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("email"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue Send Email task")
		return nil, err
	}

	return info, nil
}
