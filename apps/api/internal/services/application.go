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
	"golang.org/x/sync/errgroup"
)

var (
	ErrGetApplicationStatistics = errors.New("failed to aggregate application stats")
	ErrMismatchedReviewerCounts = errors.New("the total number of applications does not match the total number of assigned reviews")
	ErrWrongReviewerAssignment  = errors.New("an application has been assigned to a reviewer who is not authorized to review it")
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
	Country                 string `json:"country" validate:"required"`
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

	err = s.emailService.QueueConfirmationEmail(data.PreferredEmail, data.FirstName)

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

func (s *ApplicationService) SubmitApplicationReview(ctx context.Context) (*sqlc.Application, error) {
	// application, err := s.appRepo.GetAssignedApplicationByUserAndEventID(ctx, params)

	// if err != nil {
	// 	s.logger.Err(err).Msg(err.Error())
	// 	return nil, err
	// }

	return nil, nil
}

func (s *ApplicationService) DownloadResume(ctx context.Context, userId, eventId uuid.UUID, lifetimeSecs int64) (*storage.PresignedRequest, error) {
	presignableStorage, ok := s.storage.(storage.PresignableStorage)

	if !ok {
		err := errors.New("unable to type cast `Storage` to `PresignableStorage`")
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	if lifetimeSecs <= 0 {
		err := errors.New("invalid number of lifetime seconds")
		return nil, err
	}

	request, err := presignableStorage.PresignGetObject(ctx, s.buckets.ApplicationResumes, eventId.String()+"/"+userId.String(), lifetimeSecs)

	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	return request, nil
}

type ApplicationStatistics struct {
	GenderStatistics sqlc.GetApplicationGenderSplitRow   `json:"gender_stats"`
	AgeStatistics    sqlc.GetApplicationAgeSplitRow      `json:"age_stats"`
	RaceStatistics   []sqlc.GetApplicationRaceSplitRow   `json:"race_stats"`
	MajorStatistics  []sqlc.GetApplicationMajorSplitRow  `json:"major_stats"`
	SchoolStatistics []sqlc.GetApplicationSchoolSplitRow `json:"school_stats"`
	StatusStatistics sqlc.GetApplicationStatusSplitRow   `json:"status_stats"`
}

func (s *ApplicationService) GetApplicationStatistics(ctx context.Context, eventId uuid.UUID) (*ApplicationStatistics, error) {
	g, ctx := errgroup.WithContext(ctx)

	var genderStats sqlc.GetApplicationGenderSplitRow
	var ageStats sqlc.GetApplicationAgeSplitRow
	var raceStats []sqlc.GetApplicationRaceSplitRow
	var majorStats []sqlc.GetApplicationMajorSplitRow
	var schoolStats []sqlc.GetApplicationSchoolSplitRow
	var statusStats sqlc.GetApplicationStatusSplitRow

	g.Go(func() error {
		var err error
		genderStats, err = s.appRepo.GetSubmittedApplicationGenders(ctx, eventId)
		return err
	})

	g.Go(func() error {
		var err error
		ageStats, err = s.appRepo.GetSubmittedApplicationAges(ctx, eventId)
		return err
	})

	g.Go(func() error {
		var err error
		majorStats, err = s.appRepo.GetSubmittedApplicationMajors(ctx, eventId)
		return err
	})

	g.Go(func() error {
		var err error
		raceStats, err = s.appRepo.GetSubmittedApplicationRaces(ctx, eventId)
		return err
	})

	g.Go(func() error {
		var err error
		schoolStats, err = s.appRepo.GetSubmittedApplicationSchools(ctx, eventId)
		return err
	})

	g.Go(func() error {
		var err error
		statusStats, err = s.appRepo.GetApplicationStatuses(ctx, eventId)
		return err
	})

	if err := g.Wait(); err != nil {
		s.logger.Err(err).Msg("Something went wrong while getting application statistics")
		return nil, ErrGetApplicationStatistics
	}

	return &ApplicationStatistics{
		GenderStatistics: genderStats,
		AgeStatistics:    ageStats,
		RaceStatistics:   raceStats,
		MajorStatistics:  majorStats,
		SchoolStatistics: schoolStats,
		StatusStatistics: statusStats,
	}, nil

}

type ReviewerAssignment struct {
	ID     uuid.UUID `json:"id"`     // User/Reviewer ID
	Amount *int      `json:"amount"` // Number of applications assigned (nil if autoassign)
}

type ReviewerAllocation struct {
	ReviewerID             uuid.UUID   `json:"reviewer_id"`
	AssignedApplicationIDs []uuid.UUID `json:"assigned_application_ids"`
}

func (s *ApplicationService) AssignReviewers(ctx context.Context, eventId uuid.UUID, reviewers []ReviewerAssignment) error {

	//TODO: Must check if applications are closed, if we havent released decisions, and more.

	var fixedReviewers []ReviewerAssignment
	var autoReviewers []ReviewerAssignment
	var totalFixedAmount int

	for _, assignee := range reviewers {
		if assignee.Amount != nil {
			fixedReviewers = append(fixedReviewers, assignee)
			totalFixedAmount += *assignee.Amount
		} else {
			autoReviewers = append(autoReviewers, assignee)
		}
	}

	availableApplications, err := s.appRepo.ListAvailableApplicationForEvent(ctx, eventId)
	if err != nil {
		return err
	}

	totalAvailable := len(availableApplications)
	if totalAvailable == 0 {
		return nil
	}

	if totalFixedAmount > totalAvailable {
		return ErrMismatchedReviewerCounts
	}

	if totalAvailable > totalFixedAmount && len(autoReviewers) == 0 {
		return ErrMismatchedReviewerCounts
	}

	var appIndex int = 0
	var finalAllocations []ReviewerAllocation

	// Assign fixed
	for _, assignee := range fixedReviewers {
		amountToAssign := *assignee.Amount

		assignedSlice := availableApplications[appIndex : appIndex+amountToAssign]
		finalAllocations = append(finalAllocations, ReviewerAllocation{
			ReviewerID:             assignee.ID,
			AssignedApplicationIDs: assignedSlice,
		})

		appIndex += amountToAssign
	}

	remainingApps := totalAvailable - appIndex
	if remainingApps > 0 && len(autoReviewers) > 0 {
		baseShare := remainingApps / len(autoReviewers)
		remainder := remainingApps % len(autoReviewers)

		for index, assignee := range autoReviewers {
			assignCount := baseShare
			if index < remainder {
				assignCount++
			}

			assignedSlice := availableApplications[appIndex : appIndex+assignCount]
			finalAllocations = append(finalAllocations, ReviewerAllocation{
				ReviewerID:             assignee.ID,
				AssignedApplicationIDs: assignedSlice,
			})
			appIndex += assignCount
		}
	}

	for _, allocation := range finalAllocations {
		s.logger.Info().Str("ReviewerID", allocation.ReviewerID.String()).Int("AssignedCount", len(allocation.AssignedApplicationIDs)).Msg("Reviewer assigned applications")
	}

	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.appRepo.NewTx(tx)
		txEventRepo := s.eventsService.eventRepo.NewTx(tx)

		for _, allocation := range finalAllocations {
			err := txAppRepo.AssignApplicationToReviewByEvent(ctx, allocation.ReviewerID, eventId, allocation.AssignedApplicationIDs)
			if err != nil {
				s.logger.Err(err).Msg(err.Error())
				return err
			}
		}

		return txEventRepo.UpdateEventById(ctx, sqlc.UpdateEventByIdParams{
			ID:                               eventId,
			ApplicationReviewStartedDoUpdate: true,
			ApplicationReviewStarted:         true,
		})
	})
}

func (s *ApplicationService) ResetApplicationReviews(ctx context.Context, eventId uuid.UUID) error {
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.appRepo.NewTx(tx)
		txEventRepo := s.eventsService.eventRepo.NewTx(tx)

		err := txAppRepo.ResetApplicationReviewsForEvent(ctx, eventId)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		return txEventRepo.UpdateEventById(ctx, sqlc.UpdateEventByIdParams{
			ID:                               eventId,
			ApplicationReviewStartedDoUpdate: true,
			ApplicationReviewStarted:         false,
		})
	})

}

type ApplicationReviewStatus string

const (
	ApplicationReviewStatusInProgress ApplicationReviewStatus = "in_progress"
	ApplicationReviewStatusCompleted  ApplicationReviewStatus = "completed"
)

type AssignedApplication struct {
	UserID uuid.UUID               `json:"user_id"`
	Status ApplicationReviewStatus `json:"status"`
}

func (s *ApplicationService) GetAssignedApplicationsAndProgress(ctx context.Context, reviewerId, eventId uuid.UUID) ([]AssignedApplication, error) {
	applications, err := s.appRepo.ListApplicationByReviewerAndEvent(ctx, reviewerId, eventId)
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	var assignedApps []AssignedApplication
	for _, app := range applications {
		status := ApplicationReviewStatusInProgress
		if app.ExperienceRating != nil && app.PassionRating != nil {
			status = ApplicationReviewStatusCompleted
		}

		assignedApps = append(assignedApps, AssignedApplication{
			UserID: app.UserID,
			Status: status,
		})
	}

	return assignedApps, nil
}

func (s *ApplicationService) SaveApplicationReview(ctx context.Context, reviewerId, userId, eventId uuid.UUID, experienceRating, passionRating int) error {
	// Log everything for debug
	s.logger.Debug().Str("ReviewerId", reviewerId.String()).Str("UserId", userId.String()).Str("eventId", eventId.String()).Int32("Passion Rating", int32(passionRating)).Int32("Experiene Rating", int32(experienceRating)).Msg("Saving app review.")

	// Get the assigned application
	application, err := s.appRepo.GetApplicationByUserAndEventID(ctx, sqlc.GetApplicationByUserAndEventIDParams{
		UserID:  userId,
		EventID: eventId,
	})
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return err
	}

	// Ensure the application is assigned to the reviewer
	if application.AssignedReviewerID == nil || *application.AssignedReviewerID != reviewerId {
		s.logger.
			Warn().
			Str("AssignedReviewID", application.AssignedReviewerID.String()).
			Str("ReviewID", reviewerId.String()).
			Msg("Cannot review this application. either the assigned review is different or is nil.")
		return ErrWrongReviewerAssignment
	}

	if err = s.appRepo.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		UserID:                   userId,
		EventID:                  eventId,
		ExperienceRatingDoUpdate: true,
		ExperienceRating:         int32(experienceRating),
		PassionRatingDoUpdate:    true,
		PassionRating:            int32(passionRating),

		//TODO: Make it so I don't have to set this!
		StatusDoUpdate: false,
		Status:         sqlc.ApplicationStatusUnderReview,
	}); err != nil {
		s.logger.Err(err).Msg("Something went wrong Updating the application")
	}

	return nil
}

func (s *ApplicationService) JoinWaitlist(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) error {
	err := s.appRepo.JoinWaitlist(ctx, userId, eventId)
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return err
	}
	return nil
}

func (s *ApplicationService) WithdrawAcceptance(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) error {
	err := s.appRepo.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		UserID:  userId,
		EventID: eventId,
		//TODO: Make it so I don't have to set this!
		StatusDoUpdate: true,
		Status:         sqlc.ApplicationStatusRejected,
	})
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return err
	}
	return nil
}

func (s *ApplicationService) AcceptApplicationAcceptance(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) error {
	// is a check for a user being accepted necessary here? or is the frontend enough

	err := s.eventsService.eventRepo.UpdateRole(ctx,
		userId,
		eventId,
		sqlc.EventRoleTypeAttendee,
	)
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return err
	}
	return nil
}

func (s *ApplicationService) TransitionWaitlistedApplications(ctx context.Context, eventId uuid.UUID, acceptanceCount uint32) error {
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.appRepo.NewTx(tx)

		err := txAppRepo.TransitionAcceptedApplicationsToWaitlistByEventID(ctx, eventId)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		acceptedUserIds, err := txAppRepo.TransitionWaitlistedApplicationsToAcceptedByEventID(ctx, eventId, acceptanceCount)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		// sendWaitlistAcceptanceEmails([]uuid.UUID)
		s.logger.Debug().Msgf("Statuses transitioned: %s", acceptedUserIds)
		return nil
	})
}
