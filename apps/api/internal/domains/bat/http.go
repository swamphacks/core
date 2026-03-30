package bat

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func RegisterRoutes(batHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "get-bat-runs",
		Method:      http.MethodGet,
		Summary:     "Get Bat Runs",
		Description: "Returns all bat runs",
		Tags:        []string{"Bat"},
		Path:        "/runs",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusUnauthorized},
	}, batHandler.handleGetRuns)

	huma.Register(group, huma.Operation{
		OperationID: "delete-bat-run",
		Method:      http.MethodDelete,
		Summary:     "Delete Bat Run",
		Description: "Delete a bat run by id",
		Tags:        []string{"Bat"},
		Path:        "/runs/{runId}",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusUnauthorized},
	}, batHandler.handleDeleteRun)

	huma.Register(group, huma.Operation{
		OperationID: "queue-schedule-waitlist-transition-task",
		Method:      http.MethodPost,
		Summary:     "Queue Waitlist Transition Task",
		Description: "Queues an asynq task that transitions waitlisted applications, running every 3 days.",
		Tags:        []string{"Bat"},
		Path:        "/begin-waitlist-transition",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusUnauthorized},
	}, batHandler.handleQueueScheduleWaitlistTransitionTask)

	huma.Register(group, huma.Operation{
		OperationID: "queue-shutdown-waitlist-scheduler-task",
		Method:      http.MethodPost,
		Summary:     "Queue Shutdown Waitlist Scheduler",
		Description: "Shutsdown the scheduler used for the waitlist transition task. Error returned through logs if a scheduler is not active.",
		Tags:        []string{"Bat"},
		Path:        "/shutdown-waitlist-scheduler",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Errors:      []int{http.StatusInternalServerError, http.StatusUnauthorized},
	}, batHandler.handleQueueShutdownWaitlistSchedulerTask)
}

type handler struct {
	BatService *BatService
	logger     zerolog.Logger
}

func NewHandler(BatService *BatService, logger zerolog.Logger) *handler {
	return &handler{
		BatService: BatService,
		logger:     logger.With().Str("handler", "BatHandler").Str("domain", "bat").Logger(),
	}
}

type GetRunsOutput struct {
	Body *[]sqlc.BatRun
}

func (h *handler) handleGetRuns(ctx context.Context, input *struct{}) (*GetRunsOutput, error) {
	runs, err := h.BatService.GetRuns(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get bat runs")
	}

	return &GetRunsOutput{Body: runs}, nil
}

type DeleteRunOutput struct {
	Status int
}

func (h *handler) handleDeleteRun(ctx context.Context, input *struct {
	RunId string `path:"runId"`
}) (*DeleteRunOutput, error) {
	runId, err := uuid.Parse(input.RunId)

	if err != nil {
		return nil, huma.Error400BadRequest("Invalid run id")
	}

	err = h.BatService.DeleteRunById(ctx, runId)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete run by id")
	}

	return &DeleteRunOutput{Status: http.StatusOK}, nil
}

type QueueScheduleWaitlistTransitionTaskOutput struct {
	Status int
}

func (h *handler) handleQueueScheduleWaitlistTransitionTask(ctx context.Context, input *struct{}) (*QueueScheduleWaitlistTransitionTaskOutput, error) {
	err := h.BatService.QueueScheduleWaitlistTransitionTask(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to queue schedule waitlist transition task")
	}

	return &QueueScheduleWaitlistTransitionTaskOutput{Status: http.StatusOK}, nil
}

type QueueShutdownWaitlistSchedulerTaskOutput struct {
	Status int
}

func (h *handler) handleQueueShutdownWaitlistSchedulerTask(ctx context.Context, input *struct{}) (*QueueShutdownWaitlistSchedulerTaskOutput, error) {
	err := h.BatService.QueueShutdownWaitlistScheduler()

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to queue shutdown waitlist scheduler task")
	}

	return &QueueShutdownWaitlistSchedulerTaskOutput{Status: http.StatusOK}, nil
}
