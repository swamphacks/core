package workers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
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
	logger             zerolog.Logger
}

func NewBATWorker(batService *services.BatService, applicationService *services.ApplicationService, logger zerolog.Logger) *BATWorker {
	return &BATWorker{
		batService: batService,
		logger:     logger.With().Str("worker", "BATWorker").Str("component", "BAT").Logger(),
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
func (w *BATWorker) HandleTransitionWaitlistTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.TransitionWaitlistPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		w.logger.Err(err).Msg("Failed to unmarshal payload.")
		return err
	}

	var acceptanceCount uint32 = 50
	var acceptanceQuota uint32 = 500
	err := w.applicationService.TransitionWaitlistedApplications(ctx, payload.EventID, acceptanceCount, acceptanceQuota)
	if err != nil {
		w.logger.Err(err)
		return nil
	}

	return nil
}
