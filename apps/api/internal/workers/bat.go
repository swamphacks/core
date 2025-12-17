package workers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
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
	batService *services.BatService
	logger     zerolog.Logger
}

func NewBATWorker(batService *services.BatService, logger zerolog.Logger) *BATWorker {
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

	// TODO
	// Check to make sure all reviews are done

	// Fire off bat service
	err := w.batService.CalculateAdmissions(ctx, p.EventID)
	if err != nil {
		w.logger.Err(err).Msg("Something went wrong calculating admissions.")
		return err
	}

	return nil
}
