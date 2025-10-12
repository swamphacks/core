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
	taskQueue *asynq.Client
}

func NewEmailService(taskQueue *asynq.Client, logger zerolog.Logger) *EmailService {
	return &EmailService{
		logger:    logger.With().Str("service", "EmailService").Str("component", "email").Logger(),
		taskQueue: taskQueue,
	}
}

func (s *EmailService) SendEmail(to []string, subject string, body string) error {
	cfg := config.Load()

	s.logger.Info().Msgf("Sending from %s to %s with subject '%s' and body '%s'", cfg.Smtp.Username, to, subject, body)

	smtpAuth := smtp.PlainAuth("", cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Hostname)

	for _, recipient := range to {
		msg := fmt.Appendf(nil, "To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
			recipient, subject, body)

		fmt.Println(recipient, string(msg))

		err := smtp.SendMail(cfg.Smtp.ServerUrl, smtpAuth, cfg.Smtp.Username, to, msg)
		if err != nil {
			s.logger.Err(err).Msg("Failed to send email")
			return err
		}

	}

	return nil
}

func (s *EmailService) QueueSendEmail(to []string, subject string, body string) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskSendEmail(tasks.SendEmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
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
