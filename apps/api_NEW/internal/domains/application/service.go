package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/domains/bat"
	"github.com/swamphacks/core/apps/api/internal/domains/email"
	"github.com/swamphacks/core/apps/api/internal/storage"
	"golang.org/x/sync/errgroup"
)

var (
	ErrApplicationNotOpened = errors.New("Application not opened")
)

type ApplicationService struct {
	userRepo        *repository.UserRepository
	applicationRepo *repository.ApplicationRepository
	hackathonRepo   *repository.HackathonRepository
	storage         storage.Storage
	buckets         *config.CoreBuckets
	txm             *database.TransactionManager
	scheduler       *asynq.Scheduler
	emailService    *email.EmailService
	batService      *bat.BatService
	config          *config.Config
	logger          zerolog.Logger
}

func NewService(
	applicationRepo *repository.ApplicationRepository, userRepo *repository.UserRepository,
	hackathonRepo *repository.HackathonRepository, txm *database.TransactionManager, storage storage.Storage, buckets *config.CoreBuckets,
	scheduler *asynq.Scheduler, emailService *email.EmailService, batService *bat.BatService, config *config.Config, logger zerolog.Logger,
) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		userRepo:        userRepo,
		hackathonRepo:   hackathonRepo,
		emailService:    emailService, // TODO: is there anyway to structure this? I don't know if it's a good idea to depend on another service
		batService:      batService,
		storage:         storage,
		buckets:         buckets,
		txm:             txm,
		scheduler:       scheduler,
		config:          config,
		logger:          logger.With().Str("service", "ApplicationService").Str("domain", "application").Logger(),
	}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, userId uuid.UUID) (*sqlc.Application, error) {
	hackathon, err := s.hackathonRepo.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("Create application fail because can't retrieve hackathon")
		return nil, err
	}

	now := time.Now()
	isApplicationOpen := now.After(hackathon.ApplicationOpen) && now.Before(hackathon.ApplicationClose)

	if !isApplicationOpen {
		return nil, ErrApplicationNotOpened
	}

	application, err := s.applicationRepo.CreateApplication(ctx, userId)

	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	return application, nil
}

func (s *ApplicationService) GetApplicationByUserId(ctx context.Context, userId uuid.UUID) (*sqlc.Application, error) {
	application, err := s.applicationRepo.GetApplicationByUserId(ctx, userId)

	if err != nil {
		if errors.Is(err, database.ErrApplicationNotFound) {
			return nil, database.ErrApplicationNotFound
		} else {
			s.logger.Err(err).Msg("")
			return nil, err
		}
	}

	return application, nil
}

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

func (s *ApplicationService) SubmitApplication(ctx context.Context, data ApplicationSubmissionFields, resume []byte, userId uuid.UUID) error {
	hackathon, err := s.hackathonRepo.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("Submit application fail because can't retrieve hackathon")
		return err
	}

	now := time.Now()
	isApplicationOpen := now.After(hackathon.ApplicationOpen) && now.Before(hackathon.ApplicationClose)

	if !isApplicationOpen {
		return ErrApplicationNotOpened
	}

	// Submitting application is an atomic operation
	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.applicationRepo.NewTx(tx)

		err := txAppRepo.SubmitApplication(ctx, data, userId)

		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		contentType := "application/pdf"
		err = s.storage.Store(ctx, s.buckets.ApplicationResumes, userId.String(), resume, &contentType)

		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		err = s.userRepo.UpdateRole(ctx, sqlc.UpdateRoleParams{
			UserID: userId,
			Role:   sqlc.UserRoleApplicant,
		})
		if err != nil {
			s.logger.Err(err).Msg("submit application assign role fail")
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

func (s *ApplicationService) SaveApplication(ctx context.Context, data any, userId uuid.UUID) error {
	// Guard clauses to ensure application can be saved
	// 1) Check if applications are open for the event
	// 2) Ensure application status is "started" (Reject all other statuses)
	if err := s.isApplicationOpen(ctx); err != nil {
		return ErrApplicationNotOpened
	}

	application, err := s.GetApplicationByUserId(ctx, userId)
	if err != nil {
		return err
	}

	// This check should almost never fail, but just in case
	if application == nil {
		return errors.New("Application not found when saving the application")
	}

	if application.Status != sqlc.ApplicationStatusStarted {
		return errors.New("application has already been submitted and cannot be modified")
	}

	err = s.applicationRepo.SaveApplication(ctx, data, userId)

	if err != nil {
		s.logger.Err(err).Msg("Save application fail")
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

func (s *ApplicationService) GetDownloadResumeURL(ctx context.Context, userId uuid.UUID, lifetimeSecs int64) (*storage.PresignedRequest, error) {
	presignableStorage, ok := s.storage.(storage.PresignableStorage)

	if !ok {
		err := errors.New("unable to type cast `Storage` to `PresignableStorage`")
		s.logger.Err(err).Msg("download resume fail storage setup")
		return nil, err
	}

	if lifetimeSecs <= 0 {
		err := errors.New("invalid number of lifetime seconds")
		return nil, err
	}

	request, err := presignableStorage.PresignGetObject(ctx, s.buckets.ApplicationResumes, userId.String(), lifetimeSecs)

	if err != nil {
		s.logger.Err(err).Msg("fail presign get object")
		return nil, err
	}

	return request, nil
}

type ApplicationStatistics struct {
	GenderStatistics sqlc.GetApplicationGenderSplitRow   `json:"genderStats"`
	AgeStatistics    sqlc.GetApplicationAgeSplitRow      `json:"ageStats"`
	RaceStatistics   []sqlc.GetApplicationRaceSplitRow   `json:"raceStats"`
	MajorStatistics  []sqlc.GetApplicationMajorSplitRow  `json:"majorStats"`
	SchoolStatistics []sqlc.GetApplicationSchoolSplitRow `json:"schoolStats"`
	StatusStatistics sqlc.GetApplicationStatusSplitRow   `json:"statusStats"`
}

func (s *ApplicationService) GetApplicationStatistics(ctx context.Context) (*ApplicationStatistics, error) {
	g, ctx := errgroup.WithContext(ctx)

	var genderStats sqlc.GetApplicationGenderSplitRow
	var ageStats sqlc.GetApplicationAgeSplitRow
	var raceStats []sqlc.GetApplicationRaceSplitRow
	var majorStats []sqlc.GetApplicationMajorSplitRow
	var schoolStats []sqlc.GetApplicationSchoolSplitRow
	var statusStats sqlc.GetApplicationStatusSplitRow

	g.Go(func() error {
		var err error
		genderStats, err = s.applicationRepo.GetSubmittedApplicationGenders(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		ageStats, err = s.applicationRepo.GetSubmittedApplicationAges(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		majorStats, err = s.applicationRepo.GetSubmittedApplicationMajors(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		raceStats, err = s.applicationRepo.GetSubmittedApplicationRaces(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		schoolStats, err = s.applicationRepo.GetSubmittedApplicationSchools(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		statusStats, err = s.applicationRepo.GetApplicationStatuses(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		s.logger.Err(err).Msg("Something went wrong while getting application statistics")
		return nil, errors.New("Get application stats error")
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
	ID     uuid.UUID `json:"userId"` // User/Reviewer ID
	Amount *int      `json:"amount"` // Number of applications assigned (nil if autoassign)
}

type ReviewerAllocation struct {
	ReviewerID             uuid.UUID   `json:"reviewerIdd"`
	AssignedApplicationIDs []uuid.UUID `json:"assignedApplicationIds"`
}

func (s *ApplicationService) AssignReviewers(ctx context.Context, reviewers []ReviewerAssignment) error {

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

	availableApplications, err := s.applicationRepo.ListAvailableApplications(ctx)
	if err != nil {
		return err
	}

	totalAvailable := len(availableApplications)
	if totalAvailable == 0 {
		return nil
	}

	if totalFixedAmount > totalAvailable {
		return errors.New("the total number of applications does not match the total number of assigned reviews")
	}

	if totalAvailable > totalFixedAmount && len(autoReviewers) == 0 {
		return errors.New("the total number of applications does not match the total number of assigned reviews")
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
		txAppRepo := s.applicationRepo.NewTx(tx)
		txHackathonRepo := s.hackathonRepo.NewTx(tx)

		for _, allocation := range finalAllocations {
			err := txAppRepo.AssignApplicationToReview(ctx, allocation.ReviewerID, allocation.AssignedApplicationIDs)
			if err != nil {
				s.logger.Err(err).Msg("assign applicattion to review fail while allocating")
				return err
			}
		}

		return txHackathonRepo.UpdateHackathon(ctx, sqlc.UpdateHackathonParams{
			ApplicationReviewStartedDoUpdate: true,
			ApplicationReviewStarted:         true,
		})
	})
}

func (s *ApplicationService) ResetApplicationReviews(ctx context.Context) error {
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.applicationRepo.NewTx(tx)
		txHackathonRepo := s.hackathonRepo.NewTx(tx)

		err := txAppRepo.ResetApplicationReviews(ctx)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		return txHackathonRepo.UpdateHackathon(ctx, sqlc.UpdateHackathonParams{
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
	UserID uuid.UUID               `json:"applicantId"`
	Status ApplicationReviewStatus `json:"status"`
}

func (s *ApplicationService) GetAssignedApplicationsAndProgress(ctx context.Context, reviewerId uuid.UUID) ([]AssignedApplication, error) {
	applications, err := s.applicationRepo.ListApplicationByReviewer(ctx, reviewerId)
	if err != nil {
		s.logger.Err(err).Msg("get assigned applications and progress fail because list application by reviewer failed")
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

func (s *ApplicationService) SaveApplicationReview(ctx context.Context, reviewerId, applicantId uuid.UUID, experienceRating, passionRating int) error {
	// Log everything for debug
	s.logger.Debug().Str("ReviewerId", reviewerId.String()).Str("ApplicantId", applicantId.String()).Int32("Passion Rating", int32(passionRating)).Int32("Experiene Rating", int32(experienceRating)).Msg("Saving app review.")

	// Get the assigned application
	application, err := s.applicationRepo.GetApplicationByUserId(ctx, applicantId)
	if err != nil {
		s.logger.Err(err).Msg("SaveApplicationReview fail, unable to get application for user")
		return err
	}

	// Ensure the application is assigned to the reviewer
	if application.AssignedReviewerID == nil || *application.AssignedReviewerID != reviewerId {
		s.logger.
			Warn().
			Str("AssignedReviewID", application.AssignedReviewerID.String()).
			Str("ReviewID", reviewerId.String()).
			Msg("Cannot review this application. either the assigned review is different or is nil.")
		return errors.New("an application has been assigned to a reviewer who is not authorized to review it")
	}

	if err = s.applicationRepo.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		UserID:                   applicantId,
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

func (s *ApplicationService) CheckApplicationReviewsComplete(ctx context.Context) (bool, error) {
	nonReviewedApplicantUUIDs, err := s.applicationRepo.GetNonReviewedApplications(ctx)
	if err != nil {
		return false, errors.New("Failed to check application reviews status")
	}

	return len(nonReviewedApplicantUUIDs) == 0, nil
}

func (s *ApplicationService) JoinWaitlist(ctx context.Context, userId uuid.UUID) error {
	err := s.applicationRepo.JoinWaitlist(ctx, userId)
	if err != nil {
		s.logger.Err(err).Msg("Join waitlist fail")
		return err
	}
	return nil
}

func (s *ApplicationService) WithdrawAcceptance(ctx context.Context, userId uuid.UUID) error {
	err := s.applicationRepo.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		UserID:         userId,
		StatusDoUpdate: true,
		Status:         sqlc.ApplicationStatusWithdrawn,
	})
	if err != nil {
		s.logger.Err(err).Msg("WithdrawAcceptance fail, unable to update application")
		return err
	}
	return nil
}

func (s *ApplicationService) WithdrawAttendance(ctx context.Context, userId uuid.UUID) error {
	// Make atomic
	err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.applicationRepo.NewTx(tx)
		txUserRepo := s.userRepo.NewTx(tx)

		if err := txAppRepo.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
			UserID:         userId,
			StatusDoUpdate: true,
			Status:         sqlc.ApplicationStatusWithdrawn,
		}); err != nil {
			return err
		}

		return txUserRepo.UpdateRole(ctx,
			sqlc.UpdateRoleParams{
				UserID: userId,
				Role:   sqlc.UserRoleApplicant,
			},
		)
	})
	if err != nil {
		s.logger.Err(err).Msg("WithdrawAttendance fail")
		return err
	}
	return nil
}

func (s *ApplicationService) AcceptApplicationAcceptance(ctx context.Context, userId uuid.UUID) error {
	// is a check for a user being accepted necessary here? or is the frontend enough

	err := s.userRepo.UpdateRole(ctx,
		sqlc.UpdateRoleParams{
			UserID: userId,
			Role:   sqlc.UserRoleAttendee,
		},
	)
	if err != nil {
		s.logger.Err(err).Msg("AcceptApplicationAcceptance fail, unable to update role")
		return err
	}
	return nil
}

func (s *ApplicationService) TransitionWaitlistedApplications(ctx context.Context, acceptanceCount uint32, acceptanceQuota uint32) error {
	var acceptedUserIds []uuid.UUID

	ErrEventAlreadyStarted := errors.New("the event has already started")
	ErrFailedToGetContactEmail := errors.New("Failed to get contact email")

	hackathon, err := s.hackathonRepo.GetHackathon(ctx)
	if err != nil {
		s.logger.Err(err).Msg("TransitionWaitlistedApplications fail, unable to get hackathon")
		return err
	}
	currentTime := time.Now()
	if currentTime.After(hackathon.StartTime) {
		s.logger.Err(ErrEventAlreadyStarted).Msg("Could not transition waitlisted applications: the event has already started.")
		return ErrEventAlreadyStarted
	}

	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.applicationRepo.NewTx(tx)

		err := txAppRepo.TransitionAcceptedApplicationsToWaitlist(ctx)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		attendeeCount, err := s.applicationRepo.GetAttendeeCount(ctx)
		if err != nil {
			s.logger.Err(err).Msg("Failed to get total accepted application amount.")
		}
		if (acceptanceQuota - attendeeCount) <= acceptanceCount {
			s.logger.Info().Msgf("Acceptance quota is close, shutting down waitlist transition scheduler. Remaining acceptances: %v - %v <= %v", acceptanceQuota, attendeeCount, acceptanceCount)
			if s.scheduler != nil {
				// The API also uses this file, and this function can be run from an endpoint so we have to check that the scheduler exists.
				// Technically the task should be removed from the scheduler via an scheduler ENTRY_ID. However the scheduler is only running for this task.
				s.scheduler.Shutdown()
			}
			acceptanceCount = acceptanceQuota - attendeeCount
		}

		s.logger.Info().Msgf("Acceptance count: %v", acceptanceCount)
		acceptedUserIds, err = txAppRepo.TransitionWaitlistedApplicationsToAccepted(ctx, int32(acceptanceCount))
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		s.logger.Debug().Msgf("Statuses transitioned: %s", acceptedUserIds)
		return nil
	})

	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return err
	}

	for _, userId := range acceptedUserIds {
		userContactInfo, err := s.userRepo.GetUserEmailInfoById(ctx, userId)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		contactEmail, ok := userContactInfo.ContactEmail.(string)
		if !ok {
			return ErrFailedToGetContactEmail
		}

		err = s.emailService.QueueWaitlistAcceptanceEmail(contactEmail, userContactInfo.Name)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}
	}

	return nil
}

// Checks if application is opened
//
// Returns nil if yes, otherwise returns error
func (s *ApplicationService) isApplicationOpen(ctx context.Context) error {
	hackathon, err := s.hackathonRepo.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("Submit application fail because can't retrieve hackathon")
		return err
	}

	now := time.Now()
	isApplicationOpen := now.After(hackathon.ApplicationOpen) && now.Before(hackathon.ApplicationClose)

	if !isApplicationOpen {
		return ErrApplicationNotOpened
	}

	return nil
}

func (s *ApplicationService) ReleaseDecisions(ctx context.Context, batRunId uuid.UUID) error {
	batRun, err := s.batService.GetRunById(ctx, batRunId)

	if batRun.Status != sqlc.BatRunStatusCompleted {
		return errors.New("This run status is not valid for this action.")
	}

	if len(batRun.AcceptedApplicants) == 0 {
		return errors.New("No applicants marked as accepted.")
	}

	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.applicationRepo.NewTx(tx)

		err := txAppRepo.UpdateApplicationsStatuses(ctx, sqlc.ApplicationStatusAccepted, batRun.AcceptedApplicants)
		if err != nil {
			return err
		}

		return txAppRepo.UpdateApplicationsStatuses(ctx, sqlc.ApplicationStatusRejected, batRun.RejectedApplicants)
	})
	if err != nil {
		return err
	}

	err = s.emailService.SendDecisionEmails(ctx, batRun)
	if err != nil {
		return errors.New("Failed to send decision emails")
	}

	return nil
}
