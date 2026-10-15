package main

import (
	"context"
	"html"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/domains/application"
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
	defer func() {
		if err := taskQueueClient.Close(); err != nil {
			logger.Warn().Err(err).Msg("Failed to close email queue client")
		}
	}()

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

	admissionService := application.NewService(
		db, database.NewTransactionManager(db), nil, nil, nil,
		emailService, cfg, logger,
	)
	waitlistEnabled := os.Getenv("WAITLIST_DISPATCH_ENABLED") == "true"
	allowNewWaitlistOffers := os.Getenv("WAITLIST_NEW_OFFERS_ENABLED") != "false"

	mux := asynq.NewServeMux()

	qrEmailEnabled := os.Getenv("CONFIRMED_QR_EMAIL_ENABLED") == "true"

	receiptEnabled := os.Getenv("WAITLIST_RECEIPT_EMAIL_ENABLED") == "true"
	mux.HandleFunc("waitlistreceipt:sweep", func(ctx context.Context, task *asynq.Task) error {
		if !receiptEnabled {
			return nil
		}
		return admissionService.DeliverWaitlistReceipts(
			ctx,
			func(ctx context.Context, recipient, name, kind, source string) error {
				if err := ctx.Err(); err != nil {
					return err
				}
				subject := "You're on the SwampHacks XII waitlist"
				template := "WaitlistSignupReceiptEmail.html"
				if kind == "arrival" {
					subject = "Your SwampHacks XII standby arrival is recorded"
					template = "StandbyArrivalReceiptEmail.html"
				} else if source == "day_of" {
					subject = "Your SwampHacks XII day-of registration is complete"
					template = "DayOfRegistrationReceiptEmail.html"
				}
				return emailService.SendHtmlEmail(
					recipient, subject,
					map[string]string{"Name": html.EscapeString(name)},
					cfg.EmailTemplateDirectory+template,
				)
			},
		)
	})
	mux.HandleFunc("confirmedqr:sweep", func(ctx context.Context, task *asynq.Task) error {
		if !qrEmailEnabled {
			return nil
		}
		return admissionService.DeliverConfirmedQREmails(ctx, emailService.SendConfirmedQREmail)
	})

	mux.HandleFunc("waitlist:sweep_admissions", func(ctx context.Context, task *asynq.Task) error {
		if !waitlistEnabled {
			return nil
		}
		dispatched, err := admissionService.DispatchAdmissionWaitlistWithNewOffers(
			ctx, time.Now().UTC(), allowNewWaitlistOffers,
		)
		if err != nil {
			return err
		}
		delivered, err := admissionService.DeliverAdmissionWaitlistInvitations(
			ctx,
			func(ctx context.Context, recipient, name string, deadline time.Time) error {
				if err := ctx.Err(); err != nil {
					return err
				}
				location, err := time.LoadLocation("America/New_York")
				if err != nil {
					return err
				}
				return emailService.SendHtmlEmail(
					recipient,
					"You're off the waitlist! Confirm your SwampHacks XII attendance",
					map[string]string{
						"Name": html.EscapeString(name),
						"ConfirmationDeadline": deadline.In(location).Format(
							"Monday, January 2, 2006 at 3:04 PM MST",
						),
					},
					cfg.EmailTemplateDirectory+"ApplicationAcceptanceWithDeadlineEmail.html",
				)
			},
		)
		logger.Info().
			Int64("expired", dispatched.Expired).
			Int("offered", dispatched.Offered).
			Int("sent", delivered.Sent).
			Int("failed", delivered.Failed).
			Msg("Waitlist admission sweep")
		return err
	})

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

	if receiptEnabled {
		if _, err := scheduler.Register(
			"@every 1m", asynq.NewTask("waitlistreceipt:sweep", nil), asynq.Queue("email"),
		); err != nil {
			log.Fatalf("Failed to register waitlist receipt sweep: %v", err)
		}
	}

	if qrEmailEnabled {
		if _, err := scheduler.Register(
			"@every 1m", asynq.NewTask("confirmedqr:sweep", nil), asynq.Queue("email"),
		); err != nil {
			log.Fatalf("Failed to register confirmed QR sweep: %v", err)
		}
	}
	if _, err := scheduler.Register("@every 1m", tasks.NewTaskSweepScheduledCampaigns(), asynq.Queue("email")); err != nil {
		log.Fatalf("Failed to register scheduled campaign sweep: %v", err)
	}
	if waitlistEnabled {
		if _, err := scheduler.Register(
			"@every 1m",
			asynq.NewTask("waitlist:sweep_admissions", nil),
			asynq.Queue("email"),
		); err != nil {
			log.Fatalf("Failed to register waitlist sweep: %v", err)
		}
	}

	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer scheduler.Shutdown()

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Failed to run email worker")
	}
}
