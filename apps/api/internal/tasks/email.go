package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSendTextEmail           = "textemail:send"
	TypeSendHtmlEmail           = "htmlemail:send"
	TypeSendRawHtmlEmail        = "rawhtmlemail:send"
	TypeSweepScheduledCampaigns = "emailcampaign:sweep_scheduled"
)

type SendTextEmailPayload struct {
	To      []string
	Subject string
	Body    string
}

type SendHtmlEmailPayload struct {
	To               string
	Subject          string
	TemplateData     interface{}
	TemplateFilePath string
}

type SendRawHtmlEmailPayload struct {
	To      []string
	Subject string
	Body    string // raw HTML content
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

func NewTaskSendRawHtmlEmail(payload SendRawHtmlEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendRawHtmlEmail, data), nil
}

// NewTaskSweepScheduledCampaigns builds the periodic sweep task. It carries no
// payload: the due campaigns are whatever the database says at run time.
func NewTaskSweepScheduledCampaigns() *asynq.Task {
	return asynq.NewTask(TypeSweepScheduledCampaigns, nil)
}
