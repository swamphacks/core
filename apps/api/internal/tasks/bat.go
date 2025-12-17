package tasks

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeCalculateAdmissions = "admissions:calculate"
)

type CalculateAdmissionsPayload struct {
	EventID  uuid.UUID
	BatRunID uuid.UUID
}

func NewTaskCalculateAdmissions(payload CalculateAdmissionsPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeCalculateAdmissions, data), nil
}
