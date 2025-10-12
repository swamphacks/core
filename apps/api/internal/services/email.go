package services

import (
	"fmt"
	"net/smtp"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type EmailService struct {
	logger    zerolog.Logger
	smtpAuth  *smtp.Auth
	taskQueue *asynq.Client
}

func NewEmailService(taskQueue *asynq.Client, smtpAuth *smtp.Auth, logger zerolog.Logger) *EmailService {
	return &EmailService{
		logger:    logger.With().Str("service", "EmailService").Str("component", "email").Logger(),
		smtpAuth:  smtpAuth,
		taskQueue: taskQueue,
	}
}

func (s *EmailService) SendEmail(to []string, from, body string) error {
	s.logger.Info().Msgf("Sending from %s to %s with body %s", from, to, body)

	cfg := config.Load()
	msg := fmt.Appendf(nil, "To: %s\r\n"+
		"Subject: test\r\n"+
		"\r\n"+
		"This is the email body.\r\n",
		to)

	err := smtp.SendMail(cfg.Smtp.ServerUrl, *s.smtpAuth, from, to, msg)
	if err != nil {
		s.logger.Err(err).Msg("Failed to send email")
		return err
	}

	return nil
}

func (s *EmailService) QueueSendEmail(to []string, from, body string) (*asynq.TaskInfo, error) {
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

// func (s *EmailService) QueueConfirmationEmail() (*asynq.TaskInfo, error) {
// 	// get to name, get from name, generate body
// 	info, err := s.QueueSendEmail()

// 	return info, nil
// }
