package main

import (
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/tasks"
	"github.com/swamphacks/core/apps/api/internal/workers"
)

func main() {
	logger := logger.New()
	cfg := config.Load()

	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Msg("Failed to parse REDIS_URL")
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 1,
			Queues: map[string]int{
				"email": 1,
			},
			TaskCheckInterval:        5 * time.Second,
			DelayedTaskCheckInterval: time.Minute,
			HealthCheckInterval:      2 * time.Minute,
			JanitorInterval:          time.Hour,
			JanitorBatchSize:         100,
		},
	)

	// Create ses client
	sesClient := email.NewSESClient(cfg.AWS.AccessKey, cfg.AWS.AccessKeySecret, cfg.AWS.Region, logger)

	emailService := services.NewEmailService(nil, sesClient, logger)
	emailWorker := workers.NewEmailWorker(emailService, logger)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeSendHtmlEmail, emailWorker.HandleSendHtmlEmailTask)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to run email worker (could not get working directory)")
	}

	logger.Info().Str("Working dir", wd).Msg("Starting email worker")

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Failed to run email worker")
	}
}
