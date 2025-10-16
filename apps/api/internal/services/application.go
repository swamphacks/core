package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/ptr"
	"github.com/swamphacks/core/apps/api/internal/storage"
)

// TODO: figure out a way to create the submission fields dynamically using the json form files with proper validation.
// these fields are only applicable to swamphacks xi, not other events
type ApplicationSubmissionFields struct {
	FirstName               string `json:"firstName" validate:"required,max=50"`
	LastName                string `json:"lastName" validate:"required,max=50"`
	Age                     int    `json:"age" validate:"required,min=0,max=99"`
	Phone                   string `json:"phone" validate:"required,len=10"`
	PreferredEmail          string `json:"preferredEmail" validate:"required,email"`
	UniversityEmail         string `json:"universityEmail" validate:"required,email"`
	Gender                  string `json:"gender"`
	GenderOther             string `json:"gender-other"`
	Pronouns                string `json:"pronouns"`
	Race                    string `json:"race"`
	RaceOther               string `json:"race-other"`
	Orientation             string `json:"orientation"`
	Linkedin                string `json:"linkedin" validate:"required,url"`
	Github                  string `json:"github" validate:"required,url"`
	AgeCertification        bool   `json:"ageCertification" validate:"required,boolean"`
	School                  string `json:"school" validate:"required"`
	Level                   string `json:"level" validate:"required"`
	LevelOther              string `json:"level-other"`
	Year                    string `json:"year" validate:"required"`
	YearOther               string `json:"year-other"`
	GraduationYear          string `json:"graduationYear" validate:"required"`
	Majors                  string `json:"majors" validate:"required"`
	Minors                  string `json:"minors"`
	Experience              string `json:"experience" validate:"required"`
	UfHackathonExp          string `json:"ufHackathonExp" validate:"required"`
	ProjectExperience       string `json:"projectExperience" validate:"required"`
	ShirtSize               string `json:"shirtSize" validate:"required"`
	Diet                    string `json:"diet"`
	Essay1                  string `json:"essay1" validate:"required"`
	Essay2                  string `json:"essay2" validate:"required"`
	Referral                string `json:"referral" validate:"required"`
	PictureConsent          string `json:"pictureConsent" validate:"required"`
	InPersonAcknowledgement string `json:"inpersonAcknowledgement" validate:"required"`
	AgreeToConduct          string `json:"agreeToConduct" validate:"required"`
	InfoShareAuthorization  string `json:"infoShareAuthorization" validate:"required"`
	AgreeToMLHEmails        string `json:"agreeToMLHEmails"`
}

var (
	ErrApplicationDeadlinePassed = errors.New("the application deadline has passed")
	ErrApplicationUnavailable    = errors.New("unable to access the application")
	ErrApplicationCannotSave     = errors.New("unable to save the application")
	ErrApplicationPastSubmitted  = errors.New("application has already been submitted and cannot be modified")
)

type ApplicationService struct {
	appRepo       *repository.ApplicationRepository
	eventsService *EventService
	emailService  *EmailService
	storage       storage.Storage
	buckets       *config.CoreBuckets
	txm           *db.TransactionManager
	logger        zerolog.Logger
}

func NewApplicationService(appRepo *repository.ApplicationRepository, eventsService *EventService, emailService *EmailService, txm *db.TransactionManager, storage storage.Storage, buckets *config.CoreBuckets, logger zerolog.Logger) *ApplicationService {
	return &ApplicationService{
		appRepo:       appRepo,
		eventsService: eventsService,
		emailService:  emailService,
		storage:       storage,
		buckets:       buckets,
		txm:           txm,
		logger:        logger,
	}
}

func (s *ApplicationService) GetApplicationByUserAndEventID(ctx context.Context, params sqlc.GetApplicationByUserAndEventIDParams) (*sqlc.Application, error) {
	application, err := s.appRepo.GetApplicationByUserAndEventID(ctx, params)

	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	return application, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, params sqlc.CreateApplicationParams) (*sqlc.Application, error) {
	canCreateApplication, err := s.eventsService.IsApplicationsOpen(ctx, params.EventID)

	if err != nil {
		return nil, err
	}

	if !canCreateApplication {
		return nil, nil
	}

	application, err := s.appRepo.CreateApplication(ctx, params)

	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	return application, nil
}

func (s *ApplicationService) SubmitApplication(ctx context.Context, data ApplicationSubmissionFields, resume []byte, userId uuid.UUID, eventId uuid.UUID) error {
	canSubmitApplication, err := s.eventsService.IsApplicationsOpen(ctx, eventId)

	if err != nil {
		return err
	}

	if !canSubmitApplication {
		return ErrApplicationDeadlinePassed
	}

	// Submitting application is an atomic operation
	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.appRepo.NewTx(tx)

		err := txAppRepo.SubmitApplication(ctx, data, userId, eventId)

		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		contentType := "application/pdf"
		err = s.storage.Store(ctx, s.buckets.ApplicationResumes, eventId.String()+"/"+userId.String(), resume, &contentType)

		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		err = s.eventsService.AssignEventRole(ctx, ptr.UUIDToPtr(userId), nil, eventId, sqlc.EventRoleTypeApplicant)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		return nil
	})

	taskInfo, err := s.emailService.QueueSendConfirmationEmail(data.PreferredEmail, data.FirstName)
	s.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued SendConfirmationEmail task!")

	// Non-blocking error
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
	}

	return nil
}

func (s *ApplicationService) SaveApplication(ctx context.Context, data any, userId, eventId uuid.UUID) error {
	// Guard clauses to ensure application can be saved
	// 1) Check if applications are open for the event
	// 2) Ensure application status is "started" (Reject all other statuses)
	canSaveApplication, err := s.eventsService.IsApplicationsOpen(ctx, eventId)
	if err != nil {
		return err
	}

	if !canSaveApplication {
		return ErrApplicationCannotSave
	}

	application, err := s.GetApplicationByUserAndEventID(ctx, sqlc.GetApplicationByUserAndEventIDParams{
		UserID:  userId,
		EventID: eventId,
	})
	if err != nil {
		return err
	}

	// This check should almost never fail, but just in case
	if application == nil {
		return ErrApplicationUnavailable
	}

	if application.Status.ApplicationStatus != sqlc.ApplicationStatusStarted {
		return ErrApplicationPastSubmitted
	}

	err = s.appRepo.SaveApplication(ctx, data, userId, eventId)

	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return err
	}

	return nil
}
