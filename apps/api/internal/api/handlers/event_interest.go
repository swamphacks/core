package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type EventInterestHandler struct {
	eventInterestService *services.EventInterestService
	cfg                  *config.Config
	logger               zerolog.Logger
}

func NewEventInterestHandler(eventInterestService *services.EventInterestService, cfg *config.Config, logger zerolog.Logger) *EventInterestHandler {
	return &EventInterestHandler{
		eventInterestService: eventInterestService,
		cfg:                  cfg,
		logger:               logger.With().Str("handler", "EventInterestHandler").Str("component", "event_interest").Logger(),
	}
}

// AddEmailRequest is the expected payload for adding an email
type AddEmailRequest struct {
	Email  string  `json:"email"`
	Source *string `json:"source"`
}

func (h *EventInterestHandler) AddEmailToEvent(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event", "The event ID is missing from the URL!"))
		return
	}
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	// Parse JSON body
	var req AddEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	if !email.IsValidEmail(req.Email) {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_email", "Email is required"))
		return
	}

	// Call service to add email to event interest
	_, err = h.eventInterestService.CreateInterestSubmission(r.Context(), eventId, req.Email, req.Source)
	if err != nil {
		// Handle duplicate error or other service errors
		if err == services.ErrEmailConflict {
			res.SendError(w, http.StatusConflict, res.NewError("duplicate_email", "Email is already registered for this event"))
			return
		} else if err == services.ErrFailedToCreateSubmission {
			res.SendError(w, http.StatusInternalServerError, res.NewError("submission_error", "Failed to create event interest submission"))
			return
		}
		// generic server error
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
