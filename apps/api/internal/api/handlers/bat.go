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
//	@Accept         json
//	@Produce		json
//	@Success		200		{array}	sqlc.GetBatRunsWithUserInfoRow	"OK: BatRuns returned"
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

// Update event application_reviews_finished status
//
//	@Summary		Update an event's application_reviews_finished status
//	@Description	Update an event's application_reviews_finished status
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	string	true	"Event ID"	Format(uuid)
//	@Success		204		"OK - Event updated (patched)"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/app-review-decision-status [post]
func (h *BatHandler) UpdateEventApplicationReviewsFinishedStatus(w http.ResponseWriter, r *http.Request) {
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

	reviewsNotFinished, err := h.BatService.UpdateEventApplicationReviewsFinishedStatus(r.Context(), eventId)
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
	if err := json.NewEncoder(w).Encode(reviewsNotFinished); err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
		return
	}
}

// Delete an event
//
//	@Summary		Delete an event
//	@Description	Delete an existing event
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
