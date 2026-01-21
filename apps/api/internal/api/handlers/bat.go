package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/web"
)

type BatHandler struct {
	BatService *services.BatService
	logger     zerolog.Logger
}

func NewBatHandler(BatService *services.BatService, logger zerolog.Logger) *BatHandler {
	return &BatHandler{
		BatService: BatService,
		logger:     logger.With().Str("handler", "BatRunHandler").Str("component", "event_interest").Logger(),
	}
}

// Get BatRuns
//
//	@Summary		Get BatRuns
//	@Description	Gets BatRuns.
//	@Tags			Bat
//	@Accept			json
//	@Produce		json
//	@Success		200 {array} sqlc.GetRunsByEventIdRow	"OK: BatRuns returned"
//	@Router			/events/{eventId}/bat-runs [get]
func (h *BatHandler) GetRunsByEventId(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}
	eventId, err := uuid.Parse(eventIdStr)

	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	runs, err := h.BatService.GetRunsByEventId(r.Context(), eventId)
	if errors.Is(err, services.ErrMissingFields) {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Missing/malformed query. Available parameters:  eventId"))
		return
	}

	if errors.Is(err, services.ErrMissingPerms) {
		res.SendError(w, http.StatusForbidden, res.NewError("forbidden", "You are forbidden from this resource."))
		return
	}

	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(runs); err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
		return
	}

}

// Check application reviews complete
//
//	@Summary		Check if application reviews complete
//	@Description	Check if application reviews complete
//	@Tags			Bat
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	string	true	"Event ID"	Format(uuid)
//	@Success		200		"OK"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/review-status [get]
func (h *BatHandler) CheckApplicationReviewsComplete(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}
	eventId, err := uuid.Parse(eventIdStr)

	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	reviewsComplete, err := h.BatService.CheckApplicationReviewsComplete(r.Context(), eventId)
	if errors.Is(err, services.ErrMissingFields) {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Missing/malformed query. Available parameters:  eventId"))
		return
	}

	if errors.Is(err, services.ErrMissingPerms) {
		res.SendError(w, http.StatusForbidden, res.NewError("forbidden", "You are forbidden from this resource."))
		return
	}

	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reviewsComplete); err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
		return
	}
}

// Delete a run
//
//	@Summary		Delete a run
//	@Description	Delete an existing BAT run
//	@Tags			Bat
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	string	true	"Run ID"	Format(uuid)
//	@Success		204		"OK - Run deleted"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/bat-runs [delete]
func (h *BatHandler) DeleteRunById(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}
	err = h.BatService.DeleteRunById(r.Context(), eventId)

	if err != nil {
		switch err {
		case services.ErrFailedToDeleteRun:
			res.SendError(w, http.StatusInternalServerError, res.NewError("delete_error", "Failed to delete event"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

//	 Queue transition waitlist task
//
//		@Summary		Queues a waitlist transition task
//		@Description	Queues an asynq task that transitions waitlisted applications, running every 3 days.
//		@Tags
//
//		@Param			eventId	path	string	true	"ID of the event to join the waitlist for"
//		@Success		200		"Transitioned application statuses successfully"
//		@Failure		400		{object}	res.ErrorResponse	"Bad request: invalid event ID"
//		@Failure		500		{object}	res.ErrorResponse	"Server error: failed to transition application statuses"
//		@Router			/events/{eventId}/queue-transition-waitlist-task [post]
func (h *BatHandler) QueueScheduleWaitlistTransitionTask(w http.ResponseWriter, r *http.Request) {
	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not valid."))
		return
	}

	err = h.BatService.QueueScheduleWaitlistTransitionTask(r.Context(), eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Failed to create ScheduleWaitlistTransition task."))
	}

	res.Send(w, http.StatusCreated, nil)
}

//	 Queue Shutdown scheduler task
//
//		@Summary		Shutsdown an asynq scheduler
//		@Description	Shutsdown the scheduler used for the waitlist transition task. Error returned through logs if a scheduler is not active.
//		@Tags
//
//		@Success		200		"Scheduler shutdown successfully"
//		@Failure		500		{object}	res.ErrorResponse	"Server error: failed to shutdown scheduler"
//		@Router			/events/{eventId}/queue-transition-waitlist-task [post]
func (h *BatHandler) QueueShutdownWaitlistSchedulerTask(w http.ResponseWriter, r *http.Request) {
	err := h.BatService.QueueShutdownWaitlistScheduler()
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Failed to shutdown scheduler."))
	}

	res.Send(w, http.StatusOK, nil)
}

//	 Send welcome emails
//
//		@Summary		Sends welcome emails to attendees
//		@Description
//		@Tags
//
//		@Param			eventId	path	string	true	"ID of the event"
//		@Success		200		"Welcome emails began to queue successfully"
//		@Failure		400		{object}	res.ErrorResponse	"Bad request: invalid event ID"
//		@Failure		500		{object}	res.ErrorResponse	"Server error: failed to begin queuing welcome emails"
//		@Router			/events/{eventId}/send-welcome-emails [post]
func (h *BatHandler) SendWelcomeEmails(w http.ResponseWriter, r *http.Request) {
	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not valid."))
		return
	}

	err = h.BatService.SendWelcomeEmailToAttendees(r.Context(), eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Failed to send welcome emails."))
	}

	res.Send(w, http.StatusCreated, nil)
}
