package tasks

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeWaitlistTransitionStatuses = "waitlist:transitionstatuses"
)

type TransitionStatusesPayload struct {
	EventID uuid.UUID
}

func NewTaskTransitionStatuses(payload TransitionStatusesPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeWaitlistTransitionStatuses, data), nil
}
