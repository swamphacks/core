package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSendTextEmail = "textemail:send"
	TypeSendHtmlEmail = "htmlemail:send"
)

type SendTextEmailPayload struct {
	To      []string
	Subject string
	Body    string
}

type SendHtmlEmailPayload struct {
	To               string
	Name             string
	Subject          string
	TemplateFilePath string
}

func NewTaskSendTextEmail(payload SendTextEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendTextEmail, data), nil
}

func NewTaskSendHtmlEmail(payload SendHtmlEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendHtmlEmail, data), nil
}
