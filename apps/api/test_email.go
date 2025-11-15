package main

import (
	"log"
	"github.com/hibiken/asynq"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)


func main() {
	cfg := config.Load()
	
	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL", err)
	}

	client := asynq.NewClient(redisOpt)
	defer client.Close()

	task, err := tasks.NewTaskSendConfirmationEmail(tasks.SendConfirmationEmailPayload{
		To: "kestanley101@gmail.com",
		Name: "Stanley Ke",
	})

	if err != nil {
		log.Fatal("Failed to create SendConfirmationEmail task", err)
	}

	info, err := client.Enqueue(task, asynq.Queue("email"))
	if err != nil {
		log.Fatal("Failed to enqueue task:", err)
	}

	log.Printf("Successfully enqueued task: ID=%s Queue=%s", info.ID, info.Queue)
}