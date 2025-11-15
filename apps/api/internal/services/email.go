package services

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
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

// TODO: Make this generic for any event
func (s *EmailService) SendConfirmationEmail(recipient string, name string) error {

	var body bytes.Buffer

	template, err := template.ParseFiles("/app/internal/email/templates/ConfirmationEmail.html")
	if err != nil {
		s.logger.Err(err).Msg("Failed to parse email template for recipient")
	}

	err = template.Execute(&body, struct{ Name string }{Name: name})
	if err != nil {
		s.logger.Err(err).Msg("Failed to inject template variables for recipient '%s'.")
	}
	err = s.SESClient.SendHTMLEmail([]string{recipient}, "noreply@swamphacks.com", "SwampHacks XI: we received your application!", body.String())
	if err != nil {
		s.logger.Err(err).Msg("Failed to send confirmation email to recipient")
		return err
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

func (s *EmailService) SendTeamInvitationEmail(recipient string, teamName string, inviterName string, eventName string, inviteLink string) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskSendTeamInvitation(tasks.SendTeamInvitationPayload{
		To: to,
		TeamName: teamName,
		InviterName: inviterName,
		EventName: eventName,
		InviteLink: inviteLink,
	})

	if err != nil {
		s.logger.Err(err).Msg("Failed to create SendTeamInvitation task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("email"))
	if err != nil {
		s.logger.Err(err).Msg("failed to queue SendTeamInvitation task")
		return nil, err
	}

	return info, nil
}

func (s *EmailService) QueueSendTeamInvitation(to string, teamName string, inviterName string, eventName string, inviteLink string) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskSendTeamInvitation(tasks.sendTeamInvitationPayload{
		To: to,
		TeamName: teamName,
		InviterName: inviterName,
		EventName: eventName,
		InviteLink: inviteLink,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to create SendTeamInvitation task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("email"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue SendTeamInvitation task")
		return nil, err
	}

	return info, nil
}

