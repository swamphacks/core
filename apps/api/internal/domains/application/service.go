package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/domains/email"
	"github.com/swamphacks/core/apps/api/internal/storage"
	"golang.org/x/sync/errgroup"
)

type ApplicationService struct {
	db           *database.DB
	storage      storage.Storage
	buckets      *config.CoreBuckets
	txm          *database.TransactionManager
	scheduler    *asynq.Scheduler
	emailService *email.EmailService
	config       *config.Config
	logger       zerolog.Logger
}

func NewService(
	db *database.DB, txm *database.TransactionManager, storage storage.Storage, buckets *config.CoreBuckets,
	scheduler *asynq.Scheduler, emailService *email.EmailService, config *config.Config, logger zerolog.Logger,
) *ApplicationService {
	return &ApplicationService{
		db:           db,
		emailService: emailService,
		storage:      storage,
		buckets:      buckets,
		txm:          txm,
		scheduler:    scheduler,
		config:       config,
		logger:       logger.With().Str("service", "ApplicationService").Str("domain", "application").Logger(),
	}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, userID uuid.UUID) (*sqlc.Application, error) {
	hackathon, err := s.db.Query.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("Create application fail because can't retrieve hackathon")
		return nil, ErrFailedToGetHackathon
	}

	isEarly := false
	now := time.Now()

	if hackathon.AcceptEarlyApplications && (now.After(*hackathon.EarlyApplicationOpen) && now.Before(*hackathon.EarlyApplicationClose)) {
		isEarly = true
	}

	if !isEarly && (now.After(hackathon.ApplicationClose) || now.Before(hackathon.ApplicationOpen)) {
		return nil, ErrApplicationNotOpened
	}

	application, err := s.db.Query.CreateApplication(ctx, sqlc.CreateApplicationParams{
		UserID:      userID,
		HackathonID: hackathon.ID,
		IsEarly:     isEarly,
	})

	if err != nil {
		s.logger.Err(err).Msg("Failed to create application")
		return nil, ErrFailedToCreateApplication
	}

	return &application, nil
}

func (s *ApplicationService) GetApplicationById(ctx context.Context, id uuid.UUID) (*sqlc.Application, error) {
	application, err := s.db.Query.GetApplicationById(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrApplicationNotFound
		} else {
			s.logger.Err(err).Msg("Failed to get application by id")
			return nil, err
		}
	}

	return &application, nil
}

func (s *ApplicationService) GetApplicationByUserId(ctx context.Context, userID uuid.UUID) (*sqlc.Application, error) {
	application, err := s.db.Query.GetApplicationByUserId(ctx, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrApplicationNotFound
		} else {
			s.logger.Err(err).Msg("Failed to get application by user id")
			return nil, err
		}
	}

	return &application, nil
}

func (s *ApplicationService) UpdateApplicationById(ctx context.Context, req UpdateApplicationRequestDto) error {
	params := sqlc.UpdateApplicationByIdParams{
		ID: req.ApplicationID,
	}

	if req.Status != nil {
		params.StatusDoUpdate = true
		params.Status = *req.Status
	}

	return s.db.Query.UpdateApplicationById(ctx, params)
}

func (s *ApplicationService) GetExtendedApplicationById(ctx context.Context, id uuid.UUID) (*sqlc.GetExtendedApplicationByIdRow, error) {
	extendedApplication, err := s.db.Query.GetExtendedApplicationById(ctx, id)

	if err != nil {
		s.logger.Err(err).
			Str("ApplicationId", id.String()).
			Msg("GetExtendedApplicationById fail, unable to get extended application")
		return nil, err
	}

	return &extendedApplication, nil
}

func (s *ApplicationService) SearchApplications(ctx context.Context, limit, offset int32, search string) (*int64, []sqlc.SearchApplicationsWithUserInfoRow, error) {
	hackathon, err := s.db.Query.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("SearchApplications fail because can't retrieve hackathon")
		return nil, nil, ErrFailedToGetHackathon
	}

	g, ctx := errgroup.WithContext(ctx)

	var applicationsCount int64
	var applications []sqlc.SearchApplicationsWithUserInfoRow

	g.Go(func() error {
		var err error
		applicationsCount, err = s.db.Query.GetApplicationsCount(ctx, hackathon.ID)
		return err
	})

	g.Go(func() error {
		var err error
		applications, err = s.db.Query.SearchApplicationsWithUserInfo(ctx, sqlc.SearchApplicationsWithUserInfoParams{
			HackathonID: hackathon.ID,
			Offset:      int32(offset * limit),
			Limit:       int32(limit),
			Search:      &search,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		s.logger.Err(err).Msg("SearchApplications fail")
		return nil, nil, err
	}

	return &applicationsCount, applications, nil
}

func (s *ApplicationService) SubmitApplication(ctx context.Context, data ApplicationSubmissionFields, resume []byte, userID uuid.UUID) (*time.Time, error) {
	hackathon, err := s.db.Query.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("Submit application fail because can't retrieve hackathon")
		return nil, err
	}

	now := time.Now()
	isEarly := false

	if hackathon.AcceptEarlyApplications && (now.After(*hackathon.EarlyApplicationOpen) && now.Before(*hackathon.EarlyApplicationClose)) {
		isEarly = true
	}

	if !isEarly && (now.After(hackathon.ApplicationClose) || now.Before(hackathon.ApplicationOpen)) {
		return nil, ErrApplicationNotOpened
	}

	dataJSON, err := json.Marshal(data)

	if err != nil {
		return nil, errors.New("Failed to parse application data")
	}

	// Submitting application is an atomic operation
	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txDB := s.db.NewTX(tx)

		err := txDB.Query.UpdateApplicationByUserId(ctx, sqlc.UpdateApplicationByUserIdParams{
			UserID:              userID,
			StatusDoUpdate:      true,
			Status:              sqlc.ApplicationStatusSubmitted,
			ApplicationDoUpdate: true,
			Application:         dataJSON,
			SubmittedAtDoUpdate: true,
			SubmittedAt:         now,
			SavedAtDoUpdate:     true,
			SavedAt:             now,
			IsEarlyDoUpdate:     true,
			IsEarly:             isEarly,
		})

		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		contentType := "application/pdf"
		err = s.storage.Store(ctx, s.buckets.ApplicationResumes, hackathon.ID+"/"+userID.String(), resume, &contentType)

		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		err = txDB.Query.UpdateRole(ctx, sqlc.UpdateRoleParams{
			UserID: userID,
			Role:   sqlc.UserRoleApplicant,
		})
		if err != nil {
			s.logger.Err(err).Msg("submit application assign role fail")
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	err = s.emailService.QueueApplicationConfirmationEmail(data.PreferredEmail, data.FirstName)

	// Non-blocking error
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
	}

	return &now, nil
}

func (s *ApplicationService) SaveApplication(ctx context.Context, data any, userID uuid.UUID) error {
	// Guard clauses to ensure application can be saved
	// 1) Check if applications are open for the event
	// 2) Ensure application status is "started" (Reject all other statuses)
	if err := s.isApplicationOpen(ctx); err != nil {
		return ErrApplicationNotOpened
	}

	application, err := s.GetApplicationByUserId(ctx, userID)
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

	dataJSON, err := json.Marshal(data)

	if err != nil {
		return errors.New("Failed to parse application data")
	}

	err = s.db.Query.UpdateApplicationByUserId(ctx, sqlc.UpdateApplicationByUserIdParams{
		UserID:              userID,
		StatusDoUpdate:      true,
		Status:              sqlc.ApplicationStatusStarted,
		ApplicationDoUpdate: true,
		Application:         dataJSON,
		SavedAtDoUpdate:     true,
		SavedAt:             time.Now(),
	})

	if err != nil {
		s.logger.Err(err).Msg("Save application fail")
		return err
	}

	return nil
}

func (s *ApplicationService) GetApplicationResumeURL(ctx context.Context, userID uuid.UUID, lifetimeSecs int64) (*storage.PresignedRequest, error) {
	hackathon, err := s.db.Query.GetHackathon(ctx)

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

	request, err := presignableStorage.PresignGetObject(ctx, s.buckets.ApplicationResumes, hackathon.ID+"/"+userID.String(), lifetimeSecs)

	if err != nil {
		s.logger.Err(err).Msg("fail presign get object")
		return nil, err
	}

	return request, nil
}

// ReplaceResume overwrites the resume of an already-submitted application without
// touching any of the question responses. Hackers sometimes submit the wrong resume
// and need to swap it out after the fact.
func (s *ApplicationService) ReplaceResume(ctx context.Context, userID uuid.UUID, resume []byte) error {
	hackathon, err := s.db.Query.GetHackathon(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Replace resume fail because can't retrieve hackathon")
		return ErrFailedToGetHackathon
	}

	application, err := s.GetApplicationByUserId(ctx, userID)
	if err != nil {
		return err
	}

	// Only allow replacing the resume once the application has actually been submitted.
	// Before submission the resume is handled as part of the normal submit flow.
	if application.Status == sqlc.ApplicationStatusStarted {
		return ErrCannotReplaceResume
	}

	contentType := "application/pdf"
	if err := s.storage.Store(ctx, s.buckets.ApplicationResumes, hackathon.ID+"/"+userID.String(), resume, &contentType); err != nil {
		s.logger.Err(err).Msg("Replace resume fail while storing resume")
		return err
	}

	return nil
}

func (s *ApplicationService) GetDownloadResumeURL(ctx context.Context, userID uuid.UUID, lifetimeSecs int64) (*storage.PresignedRequest, error) {
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

	hackathon, err := s.db.Query.GetHackathon(ctx)
	if err != nil {
		s.logger.Err(err).Msg("download resume fail because can't retrieve hackathon")
		return nil, ErrFailedToGetHackathon
	}

	// Resumes are stored under hackathonID/userID (see SubmitApplication and ReplaceResume),
	// so the presigned download key must match that prefix.
	request, err := presignableStorage.PresignGetObject(ctx, s.buckets.ApplicationResumes, hackathon.ID+"/"+userID.String(), lifetimeSecs)

	if err != nil {
		s.logger.Err(err).Msg("fail presign get object")
		return nil, err
	}

	return request, nil
}

func (s *ApplicationService) GetApplicationStatistics(ctx context.Context) (*ApplicationStatisticsDto, error) {
	g, ctx := errgroup.WithContext(ctx)

	var genderStats sqlc.GetSubmittedApplicationGendersRow
	var ageStats sqlc.GetSubmittedApplicationAgesRow
	var raceStats []sqlc.GetSubmittedApplicationRacesRow
	var majorStats []sqlc.GetSubmittedApplicationMajorsRow
	var schoolStats []sqlc.GetSubmittedApplicationSchoolsRow
	var statusStats sqlc.GetApplicationStatusesRow

	g.Go(func() error {
		var err error
		genderStats, err = s.db.Query.GetSubmittedApplicationGenders(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		ageStats, err = s.db.Query.GetSubmittedApplicationAges(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		majorStats, err = s.db.Query.GetSubmittedApplicationMajors(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		raceStats, err = s.db.Query.GetSubmittedApplicationRaces(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		schoolStats, err = s.db.Query.GetSubmittedApplicationSchools(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		statusStats, err = s.db.Query.GetApplicationStatuses(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		s.logger.Err(err).Msg("Something went wrong while getting application statistics")
		return nil, errors.New("Get application stats error")
	}

	return &ApplicationStatisticsDto{
		GenderStatistics: genderStats,
		AgeStatistics:    ageStats,
		RaceStatistics:   raceStats,
		MajorStatistics:  majorStats,
		SchoolStatistics: schoolStats,
		StatusStatistics: statusStats,
	}, nil

}

// func (s *ApplicationService) JoinWaitlist(ctx context.Context, userID uuid.UUID) error {
// 	err := s.db.Query.WaitlistApplicationByUserId(ctx, userID)
// 	if err != nil {
// 		s.logger.Err(err).Msg("Join waitlist fail")
// 		return err
// 	}
// 	return nil
// }

func (s *ApplicationService) WithdrawApplication(ctx context.Context, userID uuid.UUID) error {
	// Make atomic
	err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txDB := s.db.NewTX(tx)

		if err := txDB.Query.UpdateApplicationByUserId(ctx, sqlc.UpdateApplicationByUserIdParams{
			UserID:         userID,
			StatusDoUpdate: true,
			Status:         sqlc.ApplicationStatusWithdrawn,
		}); err != nil {
			return err
		}

		return txDB.Query.UpdateRole(ctx,
			sqlc.UpdateRoleParams{
				UserID: userID,
				Role:   sqlc.UserRoleApplicant,
			},
		)
	})
	if err != nil {
		s.logger.Err(err).Str("userID", userID.String()).Msg("WithdrawAttendance fail")
		return err
	}
	return nil
}

func (s *ApplicationService) ConfirmAttendance(ctx context.Context, userID uuid.UUID) error {
	// Atomic
	err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txDB := s.db.NewTX(tx)

		application, err := s.db.Query.GetApplicationByUserId(ctx, userID)

		if err != nil {
			s.logger.Err(err).Msg("ConfirmAttendance fail, unable to retrieve user application")
			return err
		}

		if application.Status != sqlc.ApplicationStatusAccepted {
			err = errors.New("User is not accepted to hack")
			s.logger.Err(err).Msg(fmt.Sprintf("ConfirmAttendance fail, application is not accepted, status: %s", application.Status))
			return err
		}

		if err := txDB.Query.UpdateApplicationByUserId(ctx, sqlc.UpdateApplicationByUserIdParams{
			UserID:         userID,
			StatusDoUpdate: true,
			Status:         sqlc.ApplicationStatusConfirmed,
		}); err != nil {
			return err
		}

		return txDB.Query.UpdateRole(ctx,
			sqlc.UpdateRoleParams{
				UserID: userID,
				Role:   sqlc.UserRoleAttendee,
			},
		)
	})

	if err != nil {
		s.logger.Err(err).Str("userID", userID.String()).Msg("ConfirmAttendance fail")
		return err
	}
	return nil
}

// ============================== APPLICATION REVIEW FUNCTIONS ==============================

func (s *ApplicationService) UpdateApplicationReviewStatusForHackathon(ctx context.Context, started bool) error {
	err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txDB := s.db.NewTX(tx)

		err := txDB.Query.UpdateHackathon(ctx, sqlc.UpdateHackathonParams{
			ApplicationReviewStartedDoUpdate: true,
			ApplicationReviewStarted:         started,
		})

		if err != nil {
			s.logger.Err(err).Msg("UpdateApplicationReviewStatusForHackathon fail because can't update hackathon")
			return err
		}

		if started {
			err = txDB.Query.MarkSubmittedApplicationsAsUnderReview(ctx)
		} else {
			err = txDB.Query.ResetApplicationsToSubmitted(ctx)
		}

		if err != nil {
			s.logger.Err(err).Msg("UpdateApplicationReviewStatusForHackathon fail because can't update applications status")
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Err(err).Msg("StartApplicationReview fail")
		return err
	}

	return nil
}

func (s *ApplicationService) AssignReviewersToApplications(ctx context.Context, assignments []ReviewerAssignmentRequestDto) error {
	// TODO: Must check if applications are closed, if we havent released decisions, and more.
	type ReviewerAllocation struct {
		ReviewerID             uuid.UUID   `json:"reviewerIdd"`
		AssignedApplicationIDs []uuid.UUID `json:"assignedApplicationIds"`
	}

	hackathon, err := s.db.Query.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("AssignReviewerToApplications fail because can't retrieve hackathon")
		return err
	}

	if !hackathon.ApplicationReviewStarted {
		return errors.New("Application Review has not started")
	}

	var fixedAssignments []ReviewerAssignmentRequestDto
	var autoAssignments []ReviewerAssignmentRequestDto
	var totalFixedAmount int

	for _, assignment := range assignments {
		if assignment.Amount != nil {
			// Guard negative/zero amounts
			if *assignment.Amount <= 0 {
				return errors.New("assignment amount must be positive")
			}
			fixedAssignments = append(fixedAssignments, assignment)
			totalFixedAmount += *assignment.Amount
		} else {
			autoAssignments = append(autoAssignments, assignment)
		}
	}

	availableApplications, err := s.db.Query.ListUnderReviewApplicationIds(ctx, hackathon.ID)
	if err != nil {
		return err
	}

	totalAvailable := len(availableApplications)
	if totalAvailable == 0 {
		s.logger.Info().Msg("no available applications to assign")
		return nil
	}

	if totalFixedAmount > totalAvailable {
		return errors.New("total fixed assignment amount exceeds available applications")
	}

	if totalAvailable > totalFixedAmount && len(autoAssignments) == 0 {
		return errors.New("not enough assignment slots: some applications would remain unassigned and no auto assignments provided")
	}

	// helper to take a slice safely
	take := func(start, count int) ([]uuid.UUID, error) {
		if count == 0 {
			return []uuid.UUID{}, nil
		}
		if start < 0 || count < 0 || start+count > totalAvailable {
			return nil, errors.New("allocation out of bounds")
		}
		// copy to avoid referencing the backing array
		out := make([]uuid.UUID, count)
		copy(out, availableApplications[start:start+count])
		return out, nil
	}

	var appIndex int
	var finalAllocations []ReviewerAllocation

	// Assign fixed amounts
	for _, assignment := range fixedAssignments {
		amountToAssign := *assignment.Amount
		if appIndex+amountToAssign > totalAvailable {
			return errors.New("not enough available applications to satisfy fixed assignments")
		}

		assignedSlice, err := take(appIndex, amountToAssign)
		if err != nil {
			return err
		}

		finalAllocations = append(finalAllocations, ReviewerAllocation{
			ReviewerID:             assignment.ID,
			AssignedApplicationIDs: assignedSlice,
		})
		appIndex += amountToAssign
	}

	// Auto-assign remaining applications evenly
	remainingApps := totalAvailable - appIndex
	if remainingApps > 0 && len(autoAssignments) > 0 {
		baseShare := remainingApps / len(autoAssignments)
		remainder := remainingApps % len(autoAssignments)

		for i, assignment := range autoAssignments {
			assignCount := baseShare
			if i < remainder {
				assignCount++
			}

			if assignCount == 0 {
				// nothing to assign to this reviewer
				finalAllocations = append(finalAllocations, ReviewerAllocation{
					ReviewerID:             assignment.ID,
					AssignedApplicationIDs: []uuid.UUID{},
				})
				continue
			}

			if appIndex+assignCount > totalAvailable {
				return errors.New("auto assignment would exceed available applications")
			}

			assignedSlice, err := take(appIndex, assignCount)
			if err != nil {
				return err
			}

			finalAllocations = append(finalAllocations, ReviewerAllocation{
				ReviewerID:             assignment.ID,
				AssignedApplicationIDs: assignedSlice,
			})

			appIndex += assignCount
		}
	}

	for _, allocation := range finalAllocations {
		s.logger.Info().Str("ReviewerID", allocation.ReviewerID.String()).Int("AssignedCount", len(allocation.AssignedApplicationIDs)).Msg("Reviewer assigned applications")
	}

	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txDB := s.db.NewTX(tx)

		err = txDB.Query.DeleteAllApplicationReviews(ctx)

		if err != nil {
			s.logger.Err(err).Msg("unable to reset all application reviews before assigning")
			return err
		}

		err = txDB.Query.DeleteAllAutoDecisionRequests(ctx)

		if err != nil {
			s.logger.Err(err).Msg("unable to delete all decision requests before assigning")
			return err
		}

		for _, allocation := range finalAllocations {
			if len(allocation.AssignedApplicationIDs) == 0 {
				continue
			}

			err := txDB.Query.AssignReviewerToApplications(ctx, sqlc.AssignReviewerToApplicationsParams{
				ReviewerID:     allocation.ReviewerID,
				ApplicationIds: allocation.AssignedApplicationIDs,
			})

			if err != nil {
				s.logger.Err(err).Msg("assign application to reviewer failed while allocating")
				return err
			}
		}

		return nil
	})
}

func (s *ApplicationService) DeleteAllApplicationReviews(ctx context.Context) error {
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txDB := s.db.NewTX(tx)

		err := txDB.Query.DeleteAllApplicationReviews(ctx)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return err
		}

		return nil
	})
}

type ApplicationReviewStatus string

const (
	ApplicationReviewStatusInProgress ApplicationReviewStatus = "in_progress"
	ApplicationReviewStatusCompleted  ApplicationReviewStatus = "completed"
)

func (s *ApplicationService) GetReviewsForReviewer(ctx context.Context, reviewerId uuid.UUID) ([]sqlc.ListReviewsByReviewerIdRow, error) {
	reviews, err := s.db.Query.ListReviewsByReviewerId(ctx, reviewerId)

	if err != nil {
		s.logger.Err(err).Msg("get assigned applications and progress fail because get applications by reviewer failed")
		return nil, err
	}

	return reviews, nil
}

func (s *ApplicationService) GetAllReviewersAndProgress(ctx context.Context) ([]sqlc.ListReviewersAndProgressRow, error) {
	results, err := s.db.Query.ListReviewersAndProgress(ctx)

	if err != nil {
		s.logger.Err(err).Msg("GetAllReviewersAndProgress fail")
		return nil, err
	}

	return results, nil
}

func (s *ApplicationService) GetReviewById(ctx context.Context, reviewId uuid.UUID) (*sqlc.GetReviewByIdRow, *storage.PresignedRequest, error) {
	review, err := s.db.Query.GetReviewById(ctx, reviewId)

	if err != nil {
		s.logger.Err(err).Msg("GetReviewById fail, unable to get review")
		return nil, nil, err
	}

	resumeRequest, err := s.GetApplicationResumeURL(ctx, review.UserID, 600)

	if err != nil {
		s.logger.Err(err).Msg("GetReviewById fail, unable to get resume")
		return nil, nil, err
	}

	return &review, resumeRequest, nil
}

func (s *ApplicationService) SaveApplicationReview(ctx context.Context, req SaveReviewRequestDto, reviewerId uuid.UUID, reviewerRole sqlc.UserRole) error {
	// Log everything for debug
	s.logger.Debug().Str("ReviewerId", reviewerId.String()).
		Str("ApplicantId", req.ApplicationId.String()).
		Int32("Passion Rating", int32(req.PassionRating)).
		Int32("Experiene Rating", int32(req.ExperienceRating)).Msg("Saving app review.")

	if reviewerRole == sqlc.UserRoleStaff {
		reviewerIds, err := s.db.Query.ListApplicationReviewersById(ctx, req.ApplicationId)
		if err != nil {
			s.logger.Err(err).
				Str("ApplicationId", req.ApplicationId.String()).
				Msg("SaveApplicationReview fail, unable to get reviewers for application")
			return err
		}

		if !slices.Contains(reviewerIds, reviewerId) {
			s.logger.
				Warn().
				Str("ReviewID", reviewerId.String()).
				Msg("Cannot review this application. either the assigned review is different")
			return errors.New("An application has been assigned to a reviewer who is not authorized to review it")
		}
	}

	newExperienceRating := int32(req.ExperienceRating)
	newPassionRating := int32(req.PassionRating)

	if err := s.db.Query.UpdateApplicationReview(ctx, sqlc.UpdateApplicationReviewParams{
		ID:                       req.ReviewId,
		ReviewerID:               reviewerId,
		ExperienceRatingDoUpdate: true,
		ExperienceRating:         &newExperienceRating,
		PassionRatingDoUpdate:    true,
		PassionRating:            &newPassionRating,
		NotesDoUpdate:            true,
		Notes:                    &req.Notes,
		UpdatedByDoUpdate:        true,
		UpdatedBy:                &reviewerId,
	}); err != nil {
		s.logger.Err(err).Msg("Something went wrong updating the application review")
	}

	return nil
}

func (s *ApplicationService) UpdateApplicationReview(ctx context.Context, req UpdateReviewRequestDto, reviewerId uuid.UUID) error {
	params := sqlc.UpdateApplicationReviewParams{
		ID:         req.ReviewId,
		ReviewerID: reviewerId,
	}

	if req.ExperienceRating != nil {
		newExperienceRating := int32(*req.ExperienceRating)
		params.ExperienceRatingDoUpdate = true
		params.ExperienceRating = &newExperienceRating
	}

	if req.PassionRating != nil {
		newPassionRating := int32(*req.PassionRating)
		params.PassionRatingDoUpdate = true
		params.PassionRating = &newPassionRating
	}

	if req.Notes != nil {
		params.NotesDoUpdate = true
		params.Notes = req.Notes
	}

	params.UpdatedByDoUpdate = true
	params.UpdatedBy = &reviewerId

	return s.db.Query.UpdateApplicationReview(ctx, params)
}

func (s *ApplicationService) SearchAutoDecisionRequests(ctx context.Context, req SearchAutoDecisionRequestsDto) (*int64, []sqlc.SearchAutoDecisionRequestsRow, error) {
	g, ctx := errgroup.WithContext(ctx)

	var requestsCount int64
	var requests []sqlc.SearchAutoDecisionRequestsRow

	g.Go(func() error {
		var err error
		requestsCount, err = s.db.Query.GetAutoDecisionRequestsCount(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		var decision sqlc.NullApplicationAutoDecisionType

		if req.Decision == "all" {
			decision.Valid = false
		} else {
			decision.Valid = true
			decision.ApplicationAutoDecisionType = sqlc.ApplicationAutoDecisionType(req.Decision)
		}

		requests, err = s.db.Query.SearchAutoDecisionRequests(ctx, sqlc.SearchAutoDecisionRequestsParams{
			Offset:   int32(req.Offset * req.Limit),
			Limit:    int32(req.Limit),
			Search:   &req.Search,
			Approved: req.Approved,
			Decision: decision,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		s.logger.Err(err).Msg("SearchAutoDecisionRequests fail")
		return nil, nil, err
	}

	return &requestsCount, requests, nil
}

func (s *ApplicationService) GetAllAutoDecisionRequests(ctx context.Context) ([]sqlc.ListAutoDecisionRequestsRow, error) {
	requests, err := s.db.Query.ListAutoDecisionRequests(ctx)

	if err != nil {
		s.logger.Err(err).Msg("GetAllAutoDecisionRequests fail")
		return nil, err
	}

	return requests, nil
}

func (s *ApplicationService) RequestAutoDecision(ctx context.Context, request CreateAutoDecisionRequestDto, reviewerId uuid.UUID, reviewerRole sqlc.UserRole) (*sqlc.ApplicationAutoDecisionRequest, error) {
	hackathon, err := s.db.Query.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("RequestAutoDecision fail because can't retrieve hackathon")
		return nil, err
	}

	if !hackathon.ApplicationReviewStarted {
		return nil, errors.New("Application Review has not started")
	}

	params := sqlc.RequestAutoDecisionParams{
		ApplicationID:     request.ApplicationID,
		ReviewerID:        reviewerId,
		RequestedDecision: (sqlc.ApplicationAutoDecisionType)(request.RequestedDecision),
		Justification:     request.Justification,
	}

	if reviewerRole == sqlc.UserRoleAdmin {
		approved := true
		params.Approved = &approved
		params.DecidedBy = &reviewerId
	}

	req, err := s.db.Query.RequestAutoDecision(ctx, params)

	if err != nil {
		s.logger.Err(err).Msg("RequestAutoDecision fail")
		return nil, err
	}

	return &req, nil
}

func (s *ApplicationService) DeleteAutoDecisionRequest(ctx context.Context, requestId, reviewerId uuid.UUID) error {
	err := s.db.Query.DeleteAutoDecisionRequest(ctx, sqlc.DeleteAutoDecisionRequestParams{
		ID:         requestId,
		ReviewerID: reviewerId,
	})

	if err != nil {
		s.logger.Err(err).Msg("DeleteAutoDecisionRequest fail")
		return err
	}

	return nil
}

func (s *ApplicationService) UpdateAutoDecisionRequest(ctx context.Context, req UpdateAutoDecisionRequestDto, approverId uuid.UUID) error {
	err := s.db.Query.UpdateAutoDecisionRequest(ctx, sqlc.UpdateAutoDecisionRequestParams{
		ID:                 req.RequestID,
		ApprovedDoUpdate:   true,
		Approved:           &req.Approved,
		ApprovedByDoUpdate: true,
		DecidedBy:          &approverId,
	})

	if err != nil {
		s.logger.Err(err).Msg("UpdateAutoDecisionRequest fail")
		return err
	}

	return nil
}

func (s *ApplicationService) CheckApplicationReviewsComplete(ctx context.Context) (bool, error) {
	hackathon, err := s.db.Query.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("CheckApplicationReviewsComplete fail because can't retrieve hackathon")
		return false, err
	}

	underReviewApplicationIds, err := s.db.Query.ListUnderReviewApplicationIds(ctx, hackathon.ID)
	if err != nil {
		return false, errors.New("Failed to check application reviews status")
	}

	return len(underReviewApplicationIds) == 0, nil
}

func (s *ApplicationService) isApplicationOpen(ctx context.Context) error {
	hackathon, err := s.db.Query.GetHackathon(ctx)

	if err != nil {
		s.logger.Err(err).Msg("IsApplicationsOpen fail because can't retrieve hackathon")
		return err
	}

	now := time.Now()
	isEarly := false

	if hackathon.AcceptEarlyApplications && (now.After(*hackathon.EarlyApplicationOpen) && now.Before(*hackathon.EarlyApplicationClose)) {
		isEarly = true
	}

	if !isEarly && (now.After(hackathon.ApplicationClose) || now.Before(hackathon.ApplicationOpen)) {
		return ErrApplicationNotOpened
	}

	return nil
}

// func (s *ApplicationService) TransitionWaitlistedApplications(ctx context.Context, acceptanceCount uint32, acceptanceQuota uint32) error {
// 	var acceptedUserIds []uuid.UUID

// 	ErrEventAlreadyStarted := errors.New("the event has already started")
// 	ErrFailedToGetContactEmail := errors.New("Failed to get contact email")

// 	hackathon, err := s.db.Query.GetHackathon(ctx)
// 	if err != nil {
// 		s.logger.Err(err).Msg("TransitionWaitlistedApplications fail, unable to get hackathon")
// 		return err
// 	}
// 	currentTime := time.Now()
// 	if currentTime.After(hackathon.StartTime) {
// 		s.logger.Err(ErrEventAlreadyStarted).Msg("Could not transition waitlisted applications: the event has already started.")
// 		return ErrEventAlreadyStarted
// 	}

// 	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
// 		txDB := s.db.NewTX(tx)

// 		err := txDB.Query.TransitionAcceptedApplicationsToWaitlist(ctx)
// 		if err != nil {
// 			s.logger.Err(err).Msg(err.Error())
// 			return err
// 		}

// 		attendeeCount, err := s.db.Query.GetAttendeeCount(ctx)
// 		if err != nil {
// 			s.logger.Err(err).Msg("Failed to get total accepted application amount.")
// 		}
// 		if (acceptanceQuota - attendeeCount) <= acceptanceCount {
// 			s.logger.Info().Msgf("Acceptance quota is close, shutting down waitlist transition scheduler. Remaining acceptances: %v - %v <= %v", acceptanceQuota, attendeeCount, acceptanceCount)
// 			if s.scheduler != nil {
// 				// The API also uses this file, and this function can be run from an endpoint so we have to check that the scheduler exists.
// 				// Technically the task should be removed from the scheduler via an scheduler ENTRY_ID. However the scheduler is only running for this task.
// 				s.scheduler.Shutdown()
// 			}
// 			acceptanceCount = acceptanceQuota - attendeeCount
// 		}

// 		s.logger.Info().Msgf("Acceptance count: %v", acceptanceCount)
// 		acceptedUserIds, err = txDB.Query.TransitionWaitlistedApplicationsToAccepted(ctx, int32(acceptanceCount))
// 		if err != nil {
// 			s.logger.Err(err).Msg(err.Error())
// 			return err
// 		}

// 		s.logger.Debug().Msgf("Statuses transitioned: %s", acceptedUserIds)
// 		return nil
// 	})

// 	if err != nil {
// 		s.logger.Err(err).Msg(err.Error())
// 		return err
// 	}

// 	for _, userID := range acceptedUserIds {
// 		userContactInfo, err := s.db.Query.GetUserEmailInfoById(ctx, userID)
// 		if err != nil {
// 			s.logger.Err(err).Msg(err.Error())
// 			return err
// 		}

// 		contactEmail, ok := userContactInfo.ContactEmail.(string)
// 		if !ok {
// 			return ErrFailedToGetContactEmail
// 		}

// 		err = s.emailService.QueueWaitlistAcceptanceEmail(contactEmail, userContactInfo.Name)
// 		if err != nil {
// 			s.logger.Err(err).Msg(err.Error())
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *ApplicationService) ReleaseDecisions(ctx context.Context, batRunId uuid.UUID) error {
// 	batRun, err := s.batService.GetRunById(ctx, batRunId)

// 	if batRun.Status != sqlc.BatRunStatusCompleted {
// 		return errors.New("This run status is not valid for this action.")
// 	}

// 	if len(batRun.AcceptedApplicants) == 0 {
// 		return errors.New("No applicants marked as accepted.")
// 	}

// 	err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
// 		txDB := s.db.NewTX(tx)

// 		err := txDB.Query.UpdateApplicationStatusForMultipleUserIds(ctx, sqlc.UpdateApplicationStatusForMultipleUserIdsParams{
// 			Status:  sqlc.ApplicationStatusAccepted,
// 			UserIds: batRun.AcceptedApplicants,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return txDB.Query.UpdateApplicationStatusForMultipleUserIds(ctx, sqlc.UpdateApplicationStatusForMultipleUserIdsParams{
// 			Status:  sqlc.ApplicationStatusRejected,
// 			UserIds: batRun.RejectedApplicants,
// 		})
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	err = s.emailService.SendDecisionEmails(ctx, batRun)
// 	if err != nil {
// 		return errors.New("Failed to send decision emails")
// 	}

// 	return nil
// }
