package services

import (
	"bytes"
	"html/template"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type EmailService struct {
	logger    zerolog.Logger
	SESClient *email.SESClient
	taskQueue *asynq.Client
}

func NewEmailService(taskQueue *asynq.Client, SESClient *email.SESClient, logger zerolog.Logger) *EmailService {
	return &EmailService{
		logger:    logger.With().Str("service", "EmailService").Str("component", "email").Logger(),
		taskQueue: taskQueue,
		SESClient: SESClient,
	}
}

func (s *EmailService) QueueConfirmationEmail(recipient string, name string) error {
	cfg := config.Load()

	subject := "SwampHacks XI: we received your application!"
	templateEmailFilepath := cfg.EmailTemplateDirectory + "ConfirmationEmail.html"

	taskInfo, err := s.QueueSendHtmlEmailTask(recipient, subject, name, templateEmailFilepath)
	s.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued SendConfirmationEmail task!")

	if err != nil {
		s.logger.Err(err).Msg("Failed to send confirmation email to recipient")
		return err
	}

	return nil
}

func (s *EmailService) SendHtmlEmail(recipient string, subject string, name string, templateFilePath string) error {
	var body bytes.Buffer

	template, err := template.ParseFiles(templateFilePath)
	if err != nil {
		s.logger.Err(err).Msg("Failed to parse email template for recipient")
	}

	err = template.Execute(&body, struct{ Name string }{Name: name})
	if err != nil {
		s.logger.Err(err).Msg("Failed to inject template variables for recipient '%s'.")
	}

	err = s.SESClient.SendHTMLEmail([]string{recipient}, "noreply@swamphacks.com", subject, body.String())
	if err != nil {
		s.logger.Err(err).Msg("Failed to send html email to recipient")
		return err
	}
	s.logger.Info().Str("Template", templateFilePath).Msg("Sent email")

	return nil
}

// TODO: refactor other queue functions to use a similar naming scheme
func (s *EmailService) QueueSendHtmlEmailTask(to string, subject string, name string, templateFilePath string) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskSendHtmlEmail(tasks.SendHtmlEmailPayload{
		To:               to,
		Subject:          subject,
		Name:             name,
		TemplateFilePath: templateFilePath,
	})

	if err != nil {
		s.logger.Err(err).Msg("Failed to create SendHtmlEmail task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("email"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue SendHtmlEmail task")
		return nil, err
	}

	return info, nil
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
