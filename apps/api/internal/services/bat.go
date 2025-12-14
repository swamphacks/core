package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/bat"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

var (
	ErrListApplicationsFailure = errors.New("Failed to retrieve applications")
)

type BatService struct {
	engine    *bat.BatEngine
	appRepo   *repository.ApplicationRepository
	eventRepo *repository.EventRepository
	taskQueue *asynq.Client
	logger    zerolog.Logger
}

func NewBatService(engine *bat.BatEngine, appRepo *repository.ApplicationRepository, eventRepo *repository.EventRepository, taskQueue *asynq.Client, logger zerolog.Logger) *BatService {
	return &BatService{
		engine:    engine,
		taskQueue: taskQueue,
		logger:    logger.With().Str("service", "Bat Service").Str("component", "admissions").Logger(),
	}
}

func (s *BatService) QueueCalculateAdmissionsTask(eventId uuid.UUID) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskCalculateAdmissions(tasks.CalculateAdmissionsPayload{
		EventID: eventId,
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

type ApplicantAdmissionsData struct {
	ID            uuid.UUID
	TeamID        uuid.NullUUID
	WeightedScore float32
	IsUFStudent   bool // Is a University of Florida student?
	IsEarlyCareer bool // Is a freshman to sophomore?
}

func (s *BatService) CalculateAdmissions(ctx context.Context, eventId uuid.UUID) error {
	// First aggregate data necessary
	applications, err := s.appRepo.ListApplicationsByEventAndStatus(ctx, eventId, sqlc.ApplicationStatusUnderReview)
	if err != nil || len(applications) == 0 {
		return ErrListApplicationsFailure
	}

	event, err := s.eventRepo.GetEventByID(ctx, eventId)
	if err != nil {
		return err
	}

	maxAttendees := int32(500)
	if event.MaxAttendees != nil {
		maxAttendees = *event.MaxAttendees
	}

	var appAdmissionsData []ApplicantAdmissionsData
	for _, app := range applications {
		appAdmissionsData = append(appAdmissionsData, ApplicantAdmissionsData{
			ID: app.UserID,
			TeamID: ,
		})
	}

	return nil
}
