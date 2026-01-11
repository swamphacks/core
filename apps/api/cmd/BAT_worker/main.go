package main

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/tasks"
	"github.com/swamphacks/core/apps/api/internal/workers"
)

/*
                 -.                       .-
              _..-'(                       )`-.._
           ./'. '||\\.      _ _ /| HACK! .//||` .`\.
        ./'.|'.'||||\\|..   \'o.O' /  ..|//||||`.`|.`\.
     ./'..|'.|| |||||\``````=(___)=''''''/||||| ||.`|..`\.
   ./'.||'.|||| ||||||||||||.  U  .|||||||||||| ||||.`||.`\.
  /'|||'.|||||| |||||||||||||     ||||||||||||| ||||||.`|||`\
 '.|||'.||||||| |||||||||||||     ||||||||||||| |||||||.`|||.`
'.||| ||||||||| |/'   ``\||``     ''||/''   `\| ||||||||| |||.`
|/' \./'     `\./         \!|\   /|!/         \./'     `\./ `\|
V    V         V          }' `\ /' `{          V         V    V
`    `         `               V               '         '    '

	Entrypoint for the BAT worker which handles hackathon application review
	and admissions tasks.
*/

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

	schedulerLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	scheduler := asynq.NewScheduler(
		redisOpt,
		&asynq.SchedulerOpts{
			Location: schedulerLocation,
		},
	)

	taskQueueClient := asynq.NewClient(redisOpt)

	database := db.NewDB(cfg.DatabaseURL)
	defer database.Close()

	txm := db.NewTransactionManager(database)

	applicationRepo := repository.NewApplicationRepository(database)
	eventRepo := repository.NewEventRespository(database)
	userRepo := repository.NewUserRepository(database)
	eventService := services.NewEventService(eventRepo, userRepo, nil, nil, logger)
	batRunsRepo := repository.NewBatRunsRepository(database)

	sesClient := email.NewSESClient(cfg.AWS.AccessKey, cfg.AWS.AccessKeySecret, cfg.AWS.Region, logger)
	emailService := services.NewEmailService(taskQueueClient, sesClient, logger)
	batService := services.NewBatService(applicationRepo, eventRepo, userRepo, batRunsRepo, emailService, txm, nil, logger)
	applicationService := services.NewApplicationService(applicationRepo, userRepo, eventService, emailService, txm, nil, nil, scheduler, logger)

	BATWorker := workers.NewBATWorker(batService, applicationService, scheduler, logger)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeCalculateAdmissions, BATWorker.HandleCalculateAdmissionsTask)
	mux.HandleFunc(tasks.TypeTransitionWaitlist, BATWorker.HandleTransitionWaitlistTask)
	mux.HandleFunc(tasks.TypeScheduleTransitionWaitlist, BATWorker.HandleScheduleTransitionWaitlistTask)

	if err := srv.Run(mux); err != nil {
		logger.Fatal().Msg("Failed to run BAT worker")
	}

}
