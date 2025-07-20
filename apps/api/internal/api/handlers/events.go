package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type EventHandler struct {
	eventService *services.EventService
	cfg          *config.Config
	logger       zerolog.Logger
}

func NewEventHandler(eventService *services.EventService, cfg *config.Config, logger zerolog.Logger) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		cfg:          cfg,
		logger:       logger.With().Str("handler", "EventHandler").Str("component", "events").Logger(),
	}
}

type CreateEventRequest struct {
	Name             string    `json:"name"`
	ApplicationOpen  time.Time `json:"application_open"`
	ApplicationClose time.Time `json:"application_close"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// Parse JSON body
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	// Check for empty params
	if req.Name == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_param", "'name' must be given"))
		return
	}
	if req.ApplicationOpen.IsZero() {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_param", "'application_open' must be given"))
		return
	}
	if req.ApplicationClose.IsZero() {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_param", "'application_close' must be given"))
		return
	}
	if req.StartTime.IsZero() {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_param", "'start_time' must be given"))
		return
	}
	if req.EndTime.IsZero() {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_param", "'end_time' must be given"))
		return
	}

	// Time checks
	// Open/Close and Start/End blocks are allowed to overlap, as an event could have an open application while it is ongoing, so a comparison between the two is omitted.
	if !req.ApplicationClose.After(req.ApplicationOpen) || req.ApplicationClose.Equal(req.ApplicationOpen) {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_time", "'application_open' must be before 'application_close'"))
	}
	if !req.EndTime.After(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_time", "'start_time' must be before 'end_time'"))
	}
	if req.ApplicationOpen.Before(time.Now()) || req.ApplicationClose.Before(time.Now()) || req.StartTime.Before(time.Now()) || req.EndTime.Before(time.Now()) {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_time", "Event and application periods must not take place before the current time"))
	}

	_, err := h.eventService.CreateEvent(r.Context(), req.Name, req.ApplicationOpen, req.ApplicationClose, req.StartTime, req.EndTime)
	if err != nil {
		switch err {
		case services.ErrFailedToCreateEvent:
			res.SendError(w, http.StatusInternalServerError, res.NewError("creation_error", "Failed to create event"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusCreated)
}

type GetEventResponse struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Description      *string    `json:"description"`
	Location         *string    `json:"location"`
	LocationUrl      *string    `json:"location_url"`
	MaxAttendees     *int32     `json:"max_attendees"`
	ApplicationOpen  time.Time  `json:"application_open"`
	ApplicationClose time.Time  `json:"application_close"`
	RsvpDeadline     *time.Time `json:"rsvp_deadline"`
	DecisionRelease  *time.Time `json:"decision_release"`
	StartTime        time.Time  `json:"start_time"`
	EndTime          time.Time  `json:"end_time"`
	WebsiteUrl       *string    `json:"website_url"`
	IsPublished      *bool      `json:"is_published"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
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
	event, err := h.eventService.GetEventByID(r.Context(), eventId)

	if err != nil {
		switch err {
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	eventJson := GetEventResponse{ID: event.ID} // Compiler fills in the other parameters

	res.Send(w, http.StatusOK, eventJson)
	w.WriteHeader(http.StatusOK)
}
