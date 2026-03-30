package main

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/domains/application"
	"github.com/swamphacks/core/apps/api/internal/domains/bat"
	"github.com/swamphacks/core/apps/api/internal/domains/email"
	"github.com/swamphacks/core/apps/api/internal/domains/hackathon"
	"github.com/swamphacks/core/apps/api/internal/emailutils"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/storage"
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
	cfg := config.LoadConfig()

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
	defer taskQueueClient.Close()

	db := database.NewDB(cfg.DatabaseURL)
	defer db.Close()

	txm := database.NewTransactionManager(db)

	hackathonRepo := repository.NewHackathonRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	userRepo := repository.NewUserRepository(db)
	batRunsRepo := repository.NewBatRunsRepository(db)
	eventInterestsRepo := repository.NewEventInterestsRepository(db)

	sesClient := emailutils.NewSESClient(cfg.AWS.AccessKey, cfg.AWS.AccessKeySecret, cfg.AWS.Region, logger)
	emailService := email.NewEmailService(hackathonRepo, userRepo, taskQueueClient, sesClient, nil, logger, cfg)
	batService := bat.NewBatService(applicationRepo, hackathonRepo, userRepo, batRunsRepo, emailService, txm, nil, scheduler, cfg, logger)
	applicationService := application.NewService(applicationRepo, userRepo, hackathonRepo, txm, nil, nil, scheduler, emailService, batService, cfg, logger)
	hackathonService := hackathon.NewService(hackathonRepo, userRepo, eventInterestsRepo, &storage.R2Client{}, nil, logger)

	BATWorker := workers.NewBATWorker(batService, applicationService, hackathonService, scheduler, taskQueueClient, cfg, logger)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeCalculateAdmissions, BATWorker.HandleCalculateAdmissionsTask)
	mux.HandleFunc(tasks.TypeTransitionWaitlist, BATWorker.HandleTransitionWaitlistTask)
	mux.HandleFunc(tasks.TypeScheduleTransitionWaitlist, BATWorker.HandleScheduleTransitionWaitlistTask)
	mux.HandleFunc(tasks.TypeShutdownScheduler, BATWorker.HandleShutdownScheduler)

	if err := srv.Run(mux); err != nil {
		logger.Fatal().Msg("Failed to run BAT worker")
	}

}
