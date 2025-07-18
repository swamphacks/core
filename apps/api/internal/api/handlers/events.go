package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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

	// TODO: create checks for correct application open/close times and start/end times

	// event, err = h.eventService.CreateEvent(r.Context(), req.Name, req.ApplicationOpen, req.ApplicationClose, req.StartTime, req.EndTime)
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
