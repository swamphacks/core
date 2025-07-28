package main

import (
	"fmt"
	"log"

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

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisURL},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"email": 10,
			},
		},
	)

	emailService := services.NewEmailService(nil, logger)
	emailWorker := workers.NewEmailWorker(emailService, logger)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeSendEmail, emailWorker.HandleSendEmailTask)
	fmt.Println("Starting email worker")

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Failed to run email worker")
	}
}
