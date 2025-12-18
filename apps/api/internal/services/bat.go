package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/bat"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/tasks"
	"golang.org/x/sync/errgroup"
)

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
)

type BatService struct {
	appRepo      *repository.ApplicationRepository
	eventRepo    *repository.EventRepository
	userRepo     *repository.UserRepository
	batRunsRepo  *repository.BatRunsRepository
	emailService *EmailService
	txm          *db.TransactionManager
	taskQueue    *asynq.Client
	logger       zerolog.Logger
}

func NewBatService(appRepo *repository.ApplicationRepository, eventRepo *repository.EventRepository, userRepo *repository.UserRepository, batRunsRepo *repository.BatRunsRepository, emailService *EmailService, txm *db.TransactionManager, taskQueue *asynq.Client, logger zerolog.Logger) *BatService {
	return &BatService{
		taskQueue:    taskQueue,
		appRepo:      appRepo,
		eventRepo:    eventRepo,
		userRepo:     userRepo,
		batRunsRepo:  batRunsRepo,
		emailService: emailService,
		txm:          txm,
		logger:       logger.With().Str("service", "Bat Service").Str("component", "admissions").Logger(),
	}
}

func (s *BatService) ReleaseBatRunDecision(ctx context.Context, eventId, batRunId uuid.UUID) error {
	// Retrieve Bat Run AND Event
	g, egCtx := errgroup.WithContext(ctx)

	var event sqlc.Event
	g.Go(func() error {
		eventPtr, err := s.eventRepo.GetEventByID(egCtx, eventId)
		if err != nil {
			return err
		}

		event = *eventPtr
		return nil
	})

	var batRun sqlc.BatRun
	g.Go(func() error {
		run, err := s.batRunsRepo.GetRunById(egCtx, batRunId)
		if err != nil {
			return err
		}

		batRun = run
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if event.ID != batRun.EventID {
		return ErrRunMismatch
	}

	if batRun.Status.Valid != true || (batRun.Status.Valid == true && batRun.Status.BatRunStatus != sqlc.BatRunStatusCompleted) {
		return ErrRunStatusInvalid
	}

	if len(batRun.AcceptedApplicants) == 0 {
		return ErrNoAcceptedApplicants
	}

	err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txAppRepo := s.appRepo.NewTx(tx)

		err := txAppRepo.UpdateApplicationStatusByEventId(ctx, sqlc.ApplicationStatusAccepted, event.ID, batRun.AcceptedApplicants)
		if err != nil {
			return err
		}

		return txAppRepo.UpdateApplicationStatusByEventId(ctx, sqlc.ApplicationStatusRejected, event.ID, batRun.RejectedApplicants)
	})
	if err != nil {
		return err
	}

	err = s.SendDecisionEmails(ctx, batRun)
	if err != nil {
		return ErrFailedToSendDecisionEmails
	}

	return nil
}

func (s *BatService) SendDecisionEmails(ctx context.Context, batRun sqlc.BatRun) error {
	cfg := config.Load()
	accepetedEmailTemplatePath := cfg.EmailTemplateDirectory + "ApplicationAcceptedEmail.html"
	rejectedEmailTemplatePath := cfg.EmailTemplateDirectory + "ApplicationRejectedEmail.html"
	acceptedEmailSubject := "Congratulations on being accepted to hack in SwampHacks XI!"
	rejectedEmailSubject := "Update on Your SwampHacks XI Application"

	for _, uuid := range batRun.AcceptedApplicants {
		emailInfo, err := s.userRepo.GetUserEmailInfoById(ctx, uuid)
		if err != nil {
			return ErrCouldNotGetEmailInfo
		}

		contactEmail, ok := emailInfo.ContactEmail.(string)
		if !ok {
			return ErrFailedToGetContactEmail
		}
		taskInfo, err := s.emailService.QueueSendHtmlEmailTask(contactEmail, acceptedEmailSubject, emailInfo.Name, accepetedEmailTemplatePath)
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
		taskInfo, err := s.emailService.QueueSendHtmlEmailTask(contactEmail, rejectedEmailSubject, emailInfo.Name, rejectedEmailTemplatePath)
		s.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued rejection email")
	}

	return nil
}

func (s *BatService) AddRun(ctx context.Context, eventId uuid.UUID) (*sqlc.BatRun, error) {

	run, err := s.batRunsRepo.AddRun(ctx, eventId)
	if err != nil && errors.Is(err, repository.ErrDuplicateRun) {
		s.logger.Err(err).Msg("Could not insert result as it already exists.")
		return nil, ErrRunConflict
	} else if err != nil {
		s.logger.Err(err).Msg("An unknown error was caught!")
		return nil, ErrFailedToAddRun
	}

	return run, nil
}

func (s *BatService) GetRunsByEventId(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetRunsByEventIdRow, error) {
	return s.batRunsRepo.GetRunsByEventId(ctx, eventId)
}

func (s *BatService) UpdateRunById(ctx context.Context, params sqlc.UpdateRunByIdParams) (*sqlc.BatRun, error) {
	err := s.batRunsRepo.UpdateRunById(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrRunNotFound) {
			s.logger.Err(err).Msg(repository.ErrRunNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToUpdateRun
	}

	run, err := s.batRunsRepo.GetRunById(ctx, params.ID)

	return &run, err
}

func (s *BatService) DeleteRunById(ctx context.Context, id uuid.UUID) error {
	err := s.batRunsRepo.DeleteRunById(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrRunNotFound:
			s.logger.Err(err).Msg(repository.ErrRunNotFound.Error())
		case repository.ErrNoRunsDeleted:
			s.logger.Err(err).Msg(repository.ErrNoRunsDeleted.Error())
		case repository.ErrMultipleRunsDeleted:
			s.logger.Err(err).Msg(repository.ErrMultipleRunsDeleted.Error())
		default:
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return ErrFailedToDeleteRun
	}

	return err
}

func (s *BatService) CheckApplicationReviewsComplete(ctx context.Context, eventId uuid.UUID) (bool, error) {
	nonReviewedApplicantUUIDs, err := s.appRepo.GetNonReviewedApplications(ctx, eventId)
	if err != nil {
		return false, ErrFailedToCheckAppReviewsComplete
	}

	return len(nonReviewedApplicantUUIDs) == 0, nil
}

func (s *BatService) QueueCalculateAdmissionsTask(ctx context.Context, eventId uuid.UUID) (*asynq.TaskInfo, error) {
	newRun, err := s.AddRun(ctx, eventId)
	if err != nil {
		return nil, ErrFailedToAddRun
	}

	task, err := tasks.NewTaskCalculateAdmissions(tasks.CalculateAdmissionsPayload{
		EventID:  eventId,
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

func (s *BatService) CalculateAdmissions(ctx context.Context, eventId, batRunId uuid.UUID) error {
	s.logger.Debug().Str("RunID", batRunId.String()).Msg("")

	// check to make sure reviews are done, update state if true, return error if not
	reviewStatus, err := s.CheckApplicationReviewsComplete(ctx, eventId)
	if err != nil {
		return ErrFailedToCheckAppReviewsComplete
	}
	if reviewStatus == false {
		return ErrReviewsNotComplete
	}

	engine, err := bat.NewBatEngine(0.6, 0.4)
	if err != nil {
		return err
	}

	// Aggregate data necessary
	applications, err := s.appRepo.ListAdmissionCandidatesByEvent(ctx, eventId)
	if err != nil || len(applications) == 0 {
		return ErrListApplicationsFailure
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

	params := sqlc.UpdateRunByIdParams{
		// TODO: add UF/other/early/late info?
		AcceptedApplicantsDoUpdate: true,
		RejectedApplicantsDoUpdate: true,
		StatusDoUpdate:             true,
		AcceptedApplicants:         acceptedIDs,
		RejectedApplicants:         rejectedIDs,
		Status: sqlc.NullBatRunStatus{
			BatRunStatus: sqlc.BatRunStatusCompleted,
			Valid:        true,
		},
		ID: batRunId,
	}

	err = s.batRunsRepo.UpdateRunById(ctx, params)
	if err != nil {
		return ErrFailedToUpdateRun
	}

	s.logger.Info().Int("Teams Members Accepted", int(len(acceptedTeamMembers))).Int("Accepted", int(engine.Quota.TotalAccepted)).Int("Rejected", len(rejected)).Msg("Finished Algo")

	return nil
}

// This could be moved into the engine instead, or some mapping function within the bat package.
func (s *BatService) mapToCandidates(engine *bat.BatEngine, applications []sqlc.ListAdmissionCandidatesByEventRow) ([]bat.AdmissionCandidate, error) {
	var appAdmissionsData []bat.AdmissionCandidate
	for _, app := range applications {
		if app.ExperienceRating == nil || app.PassionRating == nil {
			return []bat.AdmissionCandidate{}, ErrMissingRatings
		}

		var admissionContext bat.AdmissionContext
		if err := json.Unmarshal(app.Application, &admissionContext); err != nil {
			s.logger.Debug().Bytes("App", app.Application).Msg("Application data")
			return []bat.AdmissionCandidate{}, err
		}

		var teamId uuid.UUID
		if app.TeamID != nil {
			teamId = *app.TeamID
		}

		wScore, err := engine.CalculateWeightedScore(*app.PassionRating, *app.ExperienceRating)
		if err != nil {
			return []bat.AdmissionCandidate{}, err
		}
		appAdmissionsData = append(appAdmissionsData, bat.AdmissionCandidate{
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
