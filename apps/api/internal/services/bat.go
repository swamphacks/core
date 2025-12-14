package services

import (
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/bat"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

type BatService struct {
	engine    *bat.BatEngine
	taskQueue *asynq.Client
	logger    zerolog.Logger
}

func NewBatService(engine *bat.BatEngine, taskQueue *asynq.Client, logger zerolog.Logger) *BatService {
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

func (s *BatService) CalculateAdmissions() {
	// First calculate scores for each applicant
	score, err := s.engine.CalculateWeightedScore(1, 5)
	if err != nil {
		s.logger.Err(err).Msg("Failed to calculate weighted score.")
	}

	s.logger.Info().Float64("Weighted score", score).Msg("Successed in calculating score.")
}
