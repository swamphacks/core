package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSendEmail = "email:send"
)

type SendEmailPayload struct {
	To      []string
	Subject string
	Body    string
}

func NewTaskSendEmail(payload SendEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendEmail, data), nil
}
