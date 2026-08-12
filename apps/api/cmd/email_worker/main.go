package main

import (
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/domains/email"
	"github.com/swamphacks/core/apps/api/internal/emailutils"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/tasks"
	"github.com/swamphacks/core/apps/api/internal/workers"
)

func main() {
	logger := logger.New()
	cfg := config.LoadConfig()

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

	taskQueueClient := asynq.NewClient(redisOpt)
	defer taskQueueClient.Close()

	db := database.NewDB(cfg.DatabaseURL)
	defer db.Close()

	hackathonRepo := repository.NewHackathonRepository(db)
	userRepo := repository.NewUserRepository(db)
	emailCampaignRepo := repository.NewEmailCampaignRepository(db)

	// Create ses client
	sesClient := emailutils.NewSESClient(cfg.AWS.AccessKey, cfg.AWS.AccessKeySecret, cfg.AWS.Region, logger)

	emailService := email.NewEmailService(hackathonRepo, userRepo, taskQueueClient, sesClient, nil, logger, cfg)
	emailCampaignService := email.NewEmailCampaignService(emailCampaignRepo, emailService, logger)
	emailWorker := workers.NewEmailWorker(emailService, emailCampaignService, logger)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeSendTextEmail, emailWorker.HandleSendTextEmailTask)
	mux.HandleFunc(tasks.TypeSendHtmlEmail, emailWorker.HandleSendHtmlEmailTask)
	mux.HandleFunc(tasks.TypeSendRawHtmlEmail, emailWorker.HandleSendRawHtmlEmailTask)
	mux.HandleFunc(tasks.TypeSweepScheduledCampaigns, emailWorker.HandleSweepScheduledCampaignsTask)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to run email worker (could not get working directory)")
	}

	logger.Info().Str("Working dir", wd).Msg("Starting email worker")

	// The scheduler enqueues the sweep task on an interval; the server below
	// consumes it like any other task, so scheduled sends share the same queue.
	scheduler := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{})
	if _, err := scheduler.Register("@every 1m", tasks.NewTaskSweepScheduledCampaigns(), asynq.Queue("email")); err != nil {
		log.Fatalf("Failed to register scheduled campaign sweep: %v", err)
	}
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer scheduler.Shutdown()

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Failed to run email worker")
	}
}
