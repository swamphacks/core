package tasks

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeCalculateAdmissions        = "admissions:calculate"
	TypeScheduleTransitionWaitlist = "waitlist:scheduletransition"
	TypeTransitionWaitlist         = "waitlist:transition"
	TypeShutdownScheduler          = "waitlist:shutdownscheduler"
)

type CalculateAdmissionsPayload struct {
	EventID  uuid.UUID
	BatRunID uuid.UUID
}

type ScheduleTransitionWaitlistPayload struct {
	EventID uuid.UUID
	Period  string
}

type TransitionWaitlistPayload struct {
	EventID                 uuid.UUID
	AcceptFromWaitlistCount uint32
	MaxAcceptedApplications uint32
}

func NewTaskCalculateAdmissions(payload CalculateAdmissionsPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeCalculateAdmissions, data), nil
}

func NewTaskScheduleTransitionWaitlist(payload ScheduleTransitionWaitlistPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeScheduleTransitionWaitlist, data), nil
}

func NewTaskTransitionWaitlist(payload TransitionWaitlistPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeTransitionWaitlist, data), nil
}

func NewTaskShutdownScheduler() (*asynq.Task, error) {
	return asynq.NewTask(TypeShutdownScheduler, nil), nil
}
