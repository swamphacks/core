package tasks

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeCalculateAdmissions = "admissions:calculate"
	TypeTransitionWaitlist  = "waitlist:transition"
)

type CalculateAdmissionsPayload struct {
	EventID  uuid.UUID
	BatRunID uuid.UUID
}

type TransitionWaitlistPayload struct {
	EventID uuid.UUID
}

func NewTaskCalculateAdmissions(payload CalculateAdmissionsPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeCalculateAdmissions, data), nil
}

func NewTaskTransitionWaitlist(payload TransitionWaitlistPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeTransitionWaitlist, data), nil
}
