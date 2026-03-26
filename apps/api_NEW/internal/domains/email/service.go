package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/emailutils"
	"github.com/swamphacks/core/apps/api/internal/storage"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type EmailService struct {
	hackathonRepo *repository.HackathonRepository
	userRepo      *repository.UserRepository
	logger        zerolog.Logger
	taskQueue     *asynq.Client
	SESClient     *emailutils.SESClient
	storage       storage.Storage
	config        *config.Config
}

func NewEmailService(
	hackathonRepo *repository.HackathonRepository, userRepo *repository.UserRepository,
	taskQueue *asynq.Client, SESClient *emailutils.SESClient, storage storage.Storage,
	logger zerolog.Logger, config *config.Config,
) *EmailService {
	return &EmailService{
		hackathonRepo: hackathonRepo,
		userRepo:      userRepo,
		logger:        logger.With().Str("service", "EmailService").Str("component", "email").Logger(),
		taskQueue:     taskQueue,
		SESClient:     SESClient,
		storage:       storage,
		config:        config,
	}
}

func (s *EmailService) QueueConfirmationEmail(recipient string, name string) error {
	subject := "SwampHacks XII: we received your application!"
	templateEmailFilepath := s.config.EmailTemplateDirectory + "ConfirmationEmail.html"

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
	err = s.storage.Store(ctx, s.config.CoreBuckets.QRCodes, userId.String(), qrPng, &contentType)
	if err != nil {
		s.logger.Err(err).Msg("Failed to upload QR code to R2")
		return err
	}

	qrPngLink := fmt.Sprintf("%s/%s", s.config.CoreBuckets.QRCodesBaseUrl, userId.String())

	subject := "SwampHacks XII – A welcome from our Organizers!"
	templateEmailFilepath := s.config.EmailTemplateDirectory + "WelcomeEmail.html"

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
	subject := "Congratulations! You're in – confirm in 72 hours to keep your spot in SwampHacks XII"
	templateEmailFilepath := s.config.EmailTemplateDirectory + "WaitlistAcceptanceEmail.html"

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

func (s *EmailService) SendWelcomeEmailToAttendees(ctx context.Context) error {
	attendees, err := s.hackathonRepo.GetAttendeeUserIds(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Could not get attendee user ids")
		return err
	}

	s.logger.Info().Msgf("Sending welcome emails to %v attendees", len(attendees))

	for _, userId := range attendees {
		contactInfo, err := s.userRepo.GetUserEmailInfoById(ctx, userId)
		if err != nil {
			s.logger.Err(err).Msgf("Could not get contact info for user with id %s", userId)
			return err
		}
		contactEmail, ok := contactInfo.ContactEmail.(string)
		if !ok {
			s.logger.Err(err).Msgf("could got convert id %s", userId)
			continue
		}
		if contactEmail == "" {
			s.logger.Err(err).Msgf("empty contact email found for user with id %s", userId)
			continue
		}

		err = s.QueueWelcomeEmail(ctx, contactEmail, contactInfo.Name, userId)
		if err != nil {
			s.logger.Err(err).Msgf("Could not queue welcome email for user with id %s", userId)
			return err
		}
	}
	return nil
}

var (
	ErrListApplicationsFailure         = errors.New("Failed to retrieve applications")
	ErrMissingRatings                  = errors.New("Some applications are missing their review ratings")
	ErrRunConflict                     = errors.New("Run already exists for this event")
	ErrFailedToAddRun                  = errors.New("Failed to add run")
	ErrFailedToDeleteRun               = errors.New("Failed to delete run")
	ErrFailedToUpdateRun               = errors.New("Failed to update run")
	ErrCouldNotGetEventInfo            = errors.New("Could not retreive event info.")
	ErrReviewsNotFinished              = errors.New("Please make sure reviews are finished before calculating application decisions.")
	ErrRunMismatch                     = errors.New("That bat run does not belong to this event.")
	ErrRunStatusInvalid                = errors.New("This run status is not valid for this action.")
	ErrNoAcceptedApplicants            = errors.New("No applicants marked as accepted.")
	ErrFailedToCheckAppReviewsComplete = errors.New("Could not get determine if application reviews have finished.")
	ErrReviewsNotComplete              = errors.New("Please make sure reviews are finished before calculating application decisions.")
	ErrCouldNotGetEmailInfo            = errors.New("Could not get email info for applicant.")
	ErrParseTemplateFilepathFailed     = errors.New("Could not parse filepath for template.")
	ErrFailedToSendDecisionEmails      = errors.New("Failed to send decision emails")
	ErrTestErr                         = errors.New("Err while testing")
	ErrFailedToGetContactEmail         = errors.New("Failed to get contact email")
	ErrUserNotAttendee                 = errors.New("user is not an attendee")
	ErrUserCheckedIn                   = errors.New("user already checked in")
)

func (s *EmailService) SendDecisionEmails(ctx context.Context, batRun sqlc.BatRun) error {
	accepetedEmailTemplatePath := s.config.EmailTemplateDirectory + "ApplicationAcceptedEmail.html"
	rejectedEmailTemplatePath := s.config.EmailTemplateDirectory + "ApplicationRejectedEmail.html"
	acceptedEmailSubject := "Congratulations on being accepted to hack in SwampHacks XII!"
	rejectedEmailSubject := "Update on Your SwampHacks XII Application"

	for _, uuid := range batRun.AcceptedApplicants {
		emailInfo, err := s.userRepo.GetUserEmailInfoById(ctx, uuid)
		if err != nil {
			return ErrCouldNotGetEmailInfo
		}

		contactEmail, ok := emailInfo.ContactEmail.(string)
		if !ok {
			return ErrFailedToGetContactEmail
		}
		type emailTemplateData struct {
			Name string
		}
		taskInfo, err := s.QueueSendHtmlEmailTask(contactEmail, acceptedEmailSubject, emailTemplateData{Name: emailInfo.Name}, accepetedEmailTemplatePath)
		s.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued acceptance email")
	}

	for _, uuid := range batRun.RejectedApplicants {
		emailInfo, err := s.userRepo.GetUserEmailInfoById(ctx, uuid)
		if err != nil {
			return ErrCouldNotGetEmailInfo
		}

		contactEmail, ok := emailInfo.ContactEmail.(string)
		if !ok {
			return ErrFailedToGetContactEmail
		}
		type emailTemplateData struct {
			Name string
		}
		taskInfo, err := s.QueueSendHtmlEmailTask(contactEmail, rejectedEmailSubject, emailTemplateData{emailInfo.Name}, rejectedEmailTemplatePath)
		s.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued rejection email")
	}

	return nil
}
