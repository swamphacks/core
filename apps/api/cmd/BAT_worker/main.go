package main

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/workers"
)

// Entrypoint for the BAT worker which handles
// hackathon application review and admissions.
func main() {
	logger := logger.New()
	cfg := config.Load()

	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Msg("failed to parse REDIS_URL")
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 1,
			Queues: map[string]int{
				"bat": 1,
			},
			TaskCheckInterval:        10 * time.Second,
			DelayedTaskCheckInterval: time.Minute,
			HealthCheckInterval:      2 * time.Minute,
			JanitorInterval:          time.Hour,
			JanitorBatchSize:         100,
		},
	)

	_ = workers.NewBATWorker(logger)

	mux := asynq.NewServeMux()

	if err := srv.Run(mux); err != nil {
		logger.Fatal().Msg("Failed to run BAT worker")
	}

}
