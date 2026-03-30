package bat

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/domains/email"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type BatService struct {
	applicationRepo *repository.ApplicationRepository
	hackathonRepo   *repository.HackathonRepository
	userRepo        *repository.UserRepository
	batRunsRepo     *repository.BatRunsRepository
	emailService    *email.EmailService
	txm             *database.TransactionManager
	taskQueue       *asynq.Client
	scheduler       *asynq.Scheduler
	config          *config.Config
	logger          zerolog.Logger
}

func NewBatService(
	applicationRepo *repository.ApplicationRepository, hackathonRepo *repository.HackathonRepository, userRepo *repository.UserRepository,
	batRunsRepo *repository.BatRunsRepository, emailService *email.EmailService, txm *database.TransactionManager,
	taskQueue *asynq.Client, scheduler *asynq.Scheduler, config *config.Config, logger zerolog.Logger) *BatService {
	return &BatService{
		taskQueue:       taskQueue,
		scheduler:       scheduler,
		applicationRepo: applicationRepo,
		hackathonRepo:   hackathonRepo,
		userRepo:        userRepo,
		batRunsRepo:     batRunsRepo,
		emailService:    emailService,
		txm:             txm,
		config:          config,
		logger:          logger.With().Str("service", "Bat Service").Logger(),
	}
}

var (
	ErrRunConflict    = errors.New("Run already exists for this event")
	ErrFailedToAddRun = errors.New("Failed to add run")
)

func (s *BatService) AddRun(ctx context.Context) (*sqlc.BatRun, error) {

	// TODO: don't hardcode the hackathonId
	run, err := s.batRunsRepo.AddRun(ctx, "xii")
	if err != nil && errors.Is(err, database.ErrDuplicateRun) {
		s.logger.Err(err).Msg("Could not insert result as it already exists.")
		return nil, ErrRunConflict
	} else if err != nil {
		s.logger.Err(err).Msg("An unknown error was caught!")
		return nil, ErrFailedToAddRun
	}

	return run, nil
}

func (s *BatService) GetRuns(ctx context.Context) (*[]sqlc.BatRun, error) {
	return s.batRunsRepo.GetRuns(ctx)
}

func (s *BatService) GetRunById(ctx context.Context, batId uuid.UUID) (sqlc.BatRun, error) {
	return s.batRunsRepo.GetRunById(ctx, batId)
}

func (s *BatService) UpdateRunById(ctx context.Context, params sqlc.UpdateBatRunByIdParams) (*sqlc.BatRun, error) {
	err := s.batRunsRepo.UpdateRunById(ctx, params)
	if err != nil {
		s.logger.Err(err).Msg("update run by id failed")
		return nil, err
	}

	run, err := s.batRunsRepo.GetRunById(ctx, params.ID)

	return &run, err
}

func (s *BatService) DeleteRunById(ctx context.Context, id uuid.UUID) error {
	err := s.batRunsRepo.DeleteRunById(ctx, id)
	if err != nil {
		s.logger.Err(err).Msg("delete run by id failed")
		return errors.New("Failed to delete run")
	}

	return err
}

func (s *BatService) QueueCalculateAdmissionsTask(ctx context.Context) (*asynq.TaskInfo, error) {
	newRun, err := s.AddRun(ctx)
	if err != nil {
		return nil, ErrFailedToAddRun
	}

	task, err := tasks.NewTaskCalculateAdmissions(tasks.CalculateAdmissionsPayload{
		BatRunID: newRun.ID,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to create CalculateAdmissions task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("bat"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue CalculateAdmissions task")
		return nil, err
	}

	return info, nil
}

func (s *BatService) CalculateAdmissions(ctx context.Context, batRunId uuid.UUID) error {
	s.logger.Debug().Str("RunID", batRunId.String()).Msg("")

	// check to make sure reviews are done, update state if true, return error if not
	nonReviewedApplicantUUIDs, err := s.applicationRepo.GetNonReviewedApplications(ctx)
	if err != nil {
		return errors.New("Could not get determine if application reviews have finished.")
	}

	areReviewsComplete := len(nonReviewedApplicantUUIDs) == 0

	if !areReviewsComplete {
		return errors.New("Please make sure reviews are finished before calculating application decisions.")
	}

	engine, err := NewBatEngine(0.6, 0.4)
	if err != nil {
		return err
	}

	// Aggregate data necessary
	applications, err := s.applicationRepo.ListAdmissionCandidates(ctx)
	if err != nil || len(applications) == 0 {
		return errors.New("Failed to retrieve applications")
	}

	admissionCandidates, err := s.mapToCandidates(engine, applications)
	if err != nil {
		return err
	}

	teams, idvs := engine.GroupCandidates(admissionCandidates)
	acceptedTeamMembers, remainder := engine.AcceptTeams(teams)

	idvs = append(idvs, remainder...)
	acceptedIdvs, rejected := engine.AcceptIndividuals(idvs)

	accepted := append(acceptedTeamMembers, acceptedIdvs...)

	acceptedIDs := make([]uuid.UUID, 0, len(accepted))
	rejectedIDs := make([]uuid.UUID, 0, len(rejected))

	for _, applicant := range accepted {
		acceptedIDs = append(acceptedIDs, applicant.UserID)
	}
	for _, applicant := range rejected {
		rejectedIDs = append(rejectedIDs, applicant.UserID)
	}

	params := sqlc.UpdateBatRunByIdParams{
		// TODO: add UF/other/early/late info?
		AcceptedApplicantsDoUpdate: true,
		RejectedApplicantsDoUpdate: true,
		StatusDoUpdate:             true,
		AcceptedApplicants:         acceptedIDs,
		RejectedApplicants:         rejectedIDs,
		Status:                     sqlc.BatRunStatusCompleted,
		ID:                         batRunId,
	}

	err = s.batRunsRepo.UpdateRunById(ctx, params)
	if err != nil {
		return errors.New("Failed to update run")
	}

	s.logger.Info().Int("Teams Members Accepted", int(len(acceptedTeamMembers))).Int("Accepted", int(engine.Quota.TotalAccepted)).Int("Rejected", len(rejected)).Msg("Finished Algo")

	return nil
}

// This could be moved into the engine instead, or some mapping function within the bat package.
func (s *BatService) mapToCandidates(engine *BatEngine, applications []sqlc.ListAdmissionCandidatesRow) ([]AdmissionCandidate, error) {
	var appAdmissionsData []AdmissionCandidate
	for _, app := range applications {
		if app.ExperienceRating == nil || app.PassionRating == nil {
			return []AdmissionCandidate{}, errors.New("Some applications are missing their review ratings")
		}

		var admissionContext AdmissionContext
		if err := json.Unmarshal(app.Application, &admissionContext); err != nil {
			s.logger.Debug().Bytes("App", app.Application).Msg("Application data")
			return []AdmissionCandidate{}, err
		}

		var teamId uuid.UUID
		if app.TeamID != nil {
			teamId = *app.TeamID
		}

		wScore, err := engine.CalculateWeightedScore(*app.PassionRating, *app.ExperienceRating)
		if err != nil {
			return []AdmissionCandidate{}, err
		}
		appAdmissionsData = append(appAdmissionsData, AdmissionCandidate{
			UserID: app.UserID,
			TeamID: uuid.NullUUID{
				UUID:  teamId,
				Valid: app.TeamID != nil,
			},
			WeightedScore: wScore,
			SortKey:       0.0,
			IsUFStudent:   admissionContext.School == "University of Florida",
			IsEarlyCareer: admissionContext.Year == "first_year" || admissionContext.Year == "second_year",
		})
	}

	return appAdmissionsData, nil
}

func (s *BatService) QueueScheduleWaitlistTransitionTask(ctx context.Context) error {
	task, err := tasks.NewTaskScheduleTransitionWaitlist(tasks.ScheduleTransitionWaitlistPayload{
		Period: s.config.AcceptFromWaitlistPeriod,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to create ScheduleTransitionWaitlist task")
		return err
	}

	_, err = s.taskQueue.Enqueue(task, asynq.Queue("bat"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue ScheduleTransitionWaitlist task")
		return err
	}
	s.logger.Info().Msg("Queued TransitionWaitlist task")

	return nil
}

func (s *BatService) QueueShutdownWaitlistScheduler() error {
	task, err := tasks.NewTaskShutdownScheduler()
	if err != nil {
		s.logger.Err(err).Msg("Failed to create ShutdownWaitlistScheduler task")
		return err
	}

	_, err = s.taskQueue.Enqueue(task, asynq.Queue("bat"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue ShutdownWaitlistScheduler task")
		return err
	}
	s.logger.Info().Msg("Queued ShutdownWaitlistScheduler task")

	return nil
}
