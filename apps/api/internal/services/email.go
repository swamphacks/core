package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/storage"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type EmailService struct {
	logger    zerolog.Logger
	taskQueue *asynq.Client
	SESClient *email.SESClient
	storage   storage.Storage
}

func NewEmailService(taskQueue *asynq.Client, SESClient *email.SESClient, storage storage.Storage, logger zerolog.Logger) *EmailService {
	return &EmailService{
		logger:    logger.With().Str("service", "EmailService").Str("component", "email").Logger(),
		taskQueue: taskQueue,
		SESClient: SESClient,
		storage:   storage,
	}
}

func (s *EmailService) QueueConfirmationEmail(recipient string, name string) error {
	cfg := config.Load()

	subject := "SwampHacks XI: we received your application!"
	templateEmailFilepath := cfg.EmailTemplateDirectory + "ConfirmationEmail.html"

	type emailTemplateData struct {
		Name string
	}
	_, err := s.QueueSendHtmlEmailTask(recipient, subject, emailTemplateData{Name: name}, templateEmailFilepath)

	if err != nil {
		s.logger.Err(err).Msg("Failed to send confirmation email to recipient")
		return err
	}

	return nil
}

func (s *EmailService) QueueWelcomeEmail(ctx context.Context, recipient string, name string, userId uuid.UUID) error {
	cfg := config.Load()

	qrString := fmt.Sprintf("IDENT::%s", userId)
	qrPng, err := qrcode.Encode(qrString, qrcode.Medium, 256)
	if err != nil {
		s.logger.Err(err).Msg("Failed to generate QR code png")
		return err
	}

	contentType := "image/png"
	if s.storage == nil {
		s.logger.Err(err).Msg("A R2 client must be connected for this function to run")
		return err
	}
	err = s.storage.Store(ctx, cfg.CoreBuckets.QRCodes, userId.String(), qrPng, &contentType)
	if err != nil {
		s.logger.Err(err).Msg("Failed to upload QR code to R2")
		return err
	}

	qrPngLink := fmt.Sprintf("%s/%s", cfg.CoreBuckets.QRCodesBaseUrl, userId.String())

	subject := "SwampHacks XI – A welcome from our Organizers!"
	templateEmailFilepath := cfg.EmailTemplateDirectory + "WelcomeEmail.html"

	type emailTemplateData struct {
		Name      string
		QRPngLink string
	}
	_, err = s.QueueSendHtmlEmailTask(recipient, subject, emailTemplateData{Name: name, QRPngLink: qrPngLink}, templateEmailFilepath)

	if err != nil {
		s.logger.Err(err).Msgf("Failed to send welcome email to recipient with userId %s", userId.String())
		return err
	}

	return nil
}

func (s *EmailService) QueueWaitlistAcceptanceEmail(recipient string, name string) error {
	cfg := config.Load()

	subject := "Congratulations! You're in – confirm in 72 hours to keep your spot in SwampHacks XI"
	templateEmailFilepath := cfg.EmailTemplateDirectory + "WaitlistAcceptanceEmail.html"

	type emailTemplateData struct {
		Name string
	}
	_, err := s.QueueSendHtmlEmailTask(recipient, subject, emailTemplateData{Name: name}, templateEmailFilepath)

	if err != nil {
		s.logger.Err(err).Msg("Failed to send waitlist acceptance email to recipient")
		return err
	}

	return nil
}

// SendHtmlEmail
//
//	  templateData: a struct holding the data which should replace {{}} tags inside of an html template.
//		 For example, if an email template uses the tag {{ .Name }}, then the templateData struct would look like
//	     type templateData struct {
//	         Name string
//	     }
func (s *EmailService) SendHtmlEmail(recipient string, subject string, templateData interface{}, templateFilePath string) error {
	var body bytes.Buffer

	template, err := template.ParseFiles(templateFilePath)
	if err != nil {
		s.logger.Err(err).Msg("Failed to parse email template for recipient")
	}

	err = template.Execute(&body, templateData)
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

func (s *EmailService) QueueSendHtmlEmailTask(to string, subject string, templateData interface{}, templateFilePath string) (*asynq.TaskInfo, error) {
	if len(to) == 0 {
		s.logger.Warn().Msgf("No recipient email found for email being sent from template '%s'", templateFilePath)
	}

	task, err := tasks.NewTaskSendHtmlEmail(tasks.SendHtmlEmailPayload{
		To:               to,
		Subject:          subject,
		TemplateData:     templateData,
		TemplateFilePath: templateFilePath,
	})

	if err != nil {
		s.logger.Err(err).Msg("Failed to create SendHtmlEmail task")
		return nil, err
	}

	taskInfo, err := s.taskQueue.Enqueue(task, asynq.Queue("email"))
	s.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued SendHtmlEmail task!")
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue SendHtmlEmail task")
		return nil, err
	}

	return taskInfo, nil
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
