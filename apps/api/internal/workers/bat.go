package workers

import "github.com/rs/zerolog"

// BAT Worker
// The BAT worker runs the background execution pipeline for our
// Balanced Admissions Thresher (BAT). This worker processes
// applicant-selection tasks using a combination of review data,
// randomized selection mechanisms, team-formation logic, and
// other decision heuristics. It operates asynchronously to ensure
// fair, consistent, and scalable admissions handling.
type BATWorker struct {
	logger zerolog.Logger
}

func NewBATWorker(logger zerolog.Logger) *BATWorker {
	return &BATWorker{
		logger: logger.With().Str("worker", "BATWorker").Str("component", "BAT").Logger(),
	}
}
