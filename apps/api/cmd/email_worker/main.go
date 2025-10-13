package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
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
			TaskCheckInterval:        10 * time.Second,
			DelayedTaskCheckInterval: 2 * time.Minute,
			HealthCheckInterval:      2 * time.Minute,
			JanitorInterval:          time.Hour,
			JanitorBatchSize:         100,
		},
	)

	emailService := services.NewEmailService(nil, logger)
	emailWorker := workers.NewEmailWorker(emailService, logger)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeSendTextEmail, emailWorker.HandleSendTextEmailTask)
	mux.HandleFunc(tasks.TypeSendConfirmationEmail, emailWorker.HandleSendConfirmationEmailTask)
	fmt.Println("Starting email worker")

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Failed to run email worker")
	}
}
