package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

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

func (s *EmailService) SendEmail(recipient string, subject string, body []byte) error {
	cfg := config.Load()

	s.logger.Info().Msgf("Sending from %s to %s with subject '%s'", cfg.Smtp.Username, recipient, subject)

	smtpAuth := smtp.PlainAuth("", cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Hostname)

	err := smtp.SendMail(cfg.Smtp.ServerUrl, smtpAuth, cfg.Smtp.Username, []string{recipient}, body)

	if err != nil {
		s.logger.Err(err).Msg("Failed to send email")
		return err
	}

	return nil
}

func (s *EmailService) SendTextEmail(to []string, subject string, msg string) error {

	for _, recipient := range to {
		msg := fmt.Appendf(nil, "To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
			recipient, subject, msg)

		err := s.SendEmail(recipient, subject, msg)

		if err != nil {
			s.logger.Err(err).Msg("Failed to send text email")
			return err
		}
	}
	return nil
}

func (s *EmailService) SendHTMLEmail(to []string, subject string, html string) error {

	for _, recipient := range to {
		msg := fmt.Appendf(nil, "To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0 \r\nContent-type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
			recipient, subject, html)

		err := s.SendEmail(recipient, subject, msg)

		if err != nil {
			s.logger.Err(err).Msg("Failed to send html email")
			return err
		}
	}
	return nil
}

func (s *EmailService) SendConfirmationEmail(recipient string, name string) error {

	var body bytes.Buffer

	template, err := template.ParseFiles("../../internal/email/templates/ConfirmationEmail.html")
	if err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("Failed to parse email template for recipient '%s'.", recipient))
	}

	err = template.Execute(&body, struct{ Name string }{Name: name})
	if err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("Failed to inject template variables for recipient '%s'.", recipient))
	}

	err = s.SendHTMLEmail([]string{recipient}, "SwampHacks XI: we have recieved your application!", body.String())
	if err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("Failed to send confirmation email for recipient '%s'", recipient))
	}

	return nil
}

func (s *EmailService) QueueSendTextEmail(to []string, subject string, body string) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskSendTextEmail(tasks.SendTextEmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	})

	if err != nil {
		s.logger.Err(err).Msg("Failed to create SendTextEmail task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("email"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue SendTextEmail task")
		return nil, err
	}

	return info, nil
}

func (s *EmailService) QueueSendConfirmationEmail(to string, name string) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskSendConfirmationEmail(tasks.SendConfirmationEmailPayload{
		To:   to,
		Name: name,
	})

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pwd: " + pwd)

	if err != nil {
		s.logger.Err(err).Msg("Failed to create SendConfirmationEmail task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("email"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue SendConfirmationEmail task")
		return nil, err
	}

	return info, nil
}
