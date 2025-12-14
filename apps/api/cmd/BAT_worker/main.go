package main

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/bat"
	"github.com/swamphacks/core/apps/api/internal/config"
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

	batEngine, err := bat.NewBatEngine(0.5, 0.5)
	if err != nil {
		logger.Fatal().Err(err).Msg("Bat engine failed to initialize.")
	}

	batService := services.NewBatService(batEngine, nil, logger)

	BATWorker := workers.NewBATWorker(batService, logger)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeCalculateAdmissions, BATWorker.HandleCalculateAdmissionsTask)

	if err := srv.Run(mux); err != nil {
		logger.Fatal().Msg("Failed to run BAT worker")
	}

}
