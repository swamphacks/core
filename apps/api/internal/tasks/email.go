package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSendTextEmail         = "textemail:send"
	TypeSendConfirmationEmail = "confirmationemail:send"
)

type SendTextEmailPayload struct {
	To      []string
	Subject string
	Body    string
}

type SendConfirmationEmailPayload struct {
	To   string
	Name string
}

func NewTaskSendTextEmail(payload SendTextEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendTextEmail, data), nil
}

func NewTaskSendConfirmationEmail(payload SendConfirmationEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendConfirmationEmail, data), nil
}
