package workers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

// BAT Worker
// The BAT worker runs the background execution pipeline for our
// Balanced Admissions Thresher (BAT). This worker processes
// applicant-selection tasks using a combination of review data,
// randomized selection mechanisms, team-formation logic, and
// other decision heuristics. It operates asynchronously to ensure
// fair, consistent, and scalable admissions handling.
type BATWorker struct {
	batService         *services.BatService
	applicationService *services.ApplicationService
	scheduler          *asynq.Scheduler
	taskQueue          *asynq.Client
	logger             zerolog.Logger
}

func NewBATWorker(batService *services.BatService, applicationService *services.ApplicationService, scheduler *asynq.Scheduler, taskQueue *asynq.Client, logger zerolog.Logger) *BATWorker {
	return &BATWorker{
		batService:         batService,
		applicationService: applicationService,
		logger:             logger.With().Str("worker", "BATWorker").Str("component", "BAT").Logger(),
		scheduler:          scheduler,
		taskQueue:          taskQueue,
	}
}

func (w *BATWorker) HandleCalculateAdmissionsTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.CalculateAdmissionsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		w.logger.Err(err).Msg("Failed to unmarshal payload.")
		return err
	}

	err := w.batService.CalculateAdmissions(ctx, p.EventID, p.BatRunID)

	// On error, mark run as failed.
	// To be fair, this should be handled more granulary. Like for example,
	// what if the original run was never created? Just food for thought for now.
	if err != nil {
		w.logger.Err(err).Msg("Something went wrong calculating admissions.")
		_, _ = w.batService.UpdateRunById(ctx, sqlc.UpdateRunByIdParams{
			ID:             p.BatRunID,
			StatusDoUpdate: true,
			Status: sqlc.NullBatRunStatus{
				BatRunStatus: sqlc.BatRunStatusFailed,
				Valid:        true,
			},
		})

		return err
	}

	return nil
}
func (w *BATWorker) HandleScheduleTransitionWaitlistTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.ScheduleTransitionWaitlistPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		w.logger.Err(err).Msg("Failed to unmarshal payload.")
		return err
	}

	cfg := config.Load()
	task, err := tasks.NewTaskTransitionWaitlist(tasks.TransitionWaitlistPayload{
		EventID:                 payload.EventID,
		AcceptFromWaitlistCount: cfg.AcceptFromWaitlistCount,
		MaxAcceptedApplications: cfg.MaxAcceptedApplications,
	})

	// The scheduler will make its first run after the period cycles once. So we queue our task immediately as well.
	_, err = w.taskQueue.Enqueue(task, asynq.Queue("bat"))

	w.scheduler.Start()
	_, err = w.scheduler.Register(payload.Period, task, asynq.Queue("bat"))
	if err != nil {
		w.logger.Err(err)
		return nil
	}

	return nil
}

func (w *BATWorker) HandleTransitionWaitlistTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.TransitionWaitlistPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		w.logger.Err(err).Msg("Failed to unmarshal payload.")
		return err
	}

	err := w.applicationService.TransitionWaitlistedApplications(ctx, payload.EventID, payload.AcceptFromWaitlistCount, payload.MaxAcceptedApplications)
	if err != nil {
		w.logger.Err(err)
		return nil
	}

	return nil
}

func (w *BATWorker) HandleShutdownScheduler(ctx context.Context, t *asynq.Task) error {
	w.scheduler.Shutdown()
	// Error returned by logging.

	return nil
}
