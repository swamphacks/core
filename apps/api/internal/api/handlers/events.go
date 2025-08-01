package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
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

// Be very careful with the types in this struct. If a type is not a pointer (pointer types allow a null value), and the field is not present in the json body, its default value will be passed to the SQL query, and be a non-null value will be put into coalese(), which will then make a NULL value impossible and instead make the default value the type's zero value in Go.
type CreateEventFields struct {
	Name             string     `json:"name" tag:"required"`
	ApplicationOpen  time.Time  `json:"application_open" tag:"required"`
	ApplicationClose time.Time  `json:"application_close" tag:"required"`
	StartTime        time.Time  `json:"start_time" tag:"required"`
	EndTime          time.Time  `json:"end_time" tag:"required"`
	Description      *string    `json:"description"`
	Location         *string    `json:"location"`
	LocationUrl      *string    `json:"location_url"`
	MaxAttendees     *int32     `json:"max_attendees"`
	RsvpDeadline     *time.Time `json:"rsvp_deadline"`
	DecisionRelease  *time.Time `json:"decision_release"`
	WebsiteUrl       *string    `json:"website_url"`
	IsPublished      *bool      `json:"is_published"`
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	// Parse JSON body
	var req CreateEventFields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	// This will also throw an error for empty values for fields which correspond to types that cannot convert an empty string to a zero value (e.g. time.Time)
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", err.Error())) //TODO: change back
		return
	}

	// == Validation ==
	// While encoding/json can accept omitempty and omitzero tags, there is no tag for denoting that a field is required.
	// This returns error when field containing `tag"required"` is set to a zero value or is missing (e.g. empty string).
	fields := reflect.ValueOf(&req).Elem()
	fmt.Printf("%v\n", fields.NumField())
	for i := 0; i < fields.NumField(); i++ {
		tags := fields.Type().Field(i).Tag.Get("tag")
		if strings.Contains(tags, "required") && fields.Field(i).IsZero() {
			res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", fmt.Sprintf("Required field is empty or missing: %v", fields.Type().Field(i).Tag.Get("json"))))
			return
		}
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

	params := sqlc.CreateEventParams{
		Name:             req.Name,
		ApplicationOpen:  req.ApplicationOpen,
		ApplicationClose: req.ApplicationClose,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Description:      req.Description,
		Location:         req.Location,
		LocationUrl:      req.LocationUrl,
		MaxAttendees:     req.MaxAttendees,
		RsvpDeadline:     req.RsvpDeadline,
		DecisionRelease:  req.DecisionRelease,
		WebsiteUrl:       req.WebsiteUrl,
		IsPublished:      req.IsPublished,
	}

	fmt.Printf("req.Description: %v\n", req.Description)

	_, err := h.eventService.CreateEvent(r.Context(), params)
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

	eventJson := sqlc.Event{
		ID:               event.ID,
		Name:             event.Name,
		Description:      event.Description,
		Location:         event.Location,
		LocationUrl:      event.LocationUrl,
		MaxAttendees:     event.MaxAttendees,
		ApplicationOpen:  event.ApplicationOpen,
		ApplicationClose: event.ApplicationClose,
		RsvpDeadline:     event.RsvpDeadline,
		DecisionRelease:  event.DecisionRelease,
		StartTime:        event.StartTime,
		EndTime:          event.EndTime,
		WebsiteUrl:       event.WebsiteUrl,
		IsPublished:      event.IsPublished,
		CreatedAt:        event.CreatedAt,
		UpdatedAt:        event.UpdatedAt,
	}

	res.Send(w, http.StatusOK, eventJson)
	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) UpdateEventById(w http.ResponseWriter, r *http.Request) {
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

	var req sqlc.UpdateEventByIdParams
	// Because sqlc.UpdateEventById() must take in parameters through a sqlc.UpdateEventByIdParams struct, we must inject the eventId from the URL parameters into the struct.
	// This will override an eventId that is added by the request body, if they mistakenly try to add one there. But if a client fails to add it as a URL parameter, they will get an error anyways.
	req.ID = eventId

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	err = h.eventService.UpdateEventById(r.Context(), req)

	if err != nil {
		switch err {
		case services.ErrFailedToUpdateEvent:
			res.SendError(w, http.StatusInternalServerError, res.NewError("patch_error", "Failed to update event"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusNoContent) // We return a 204 for http PATCH requests since nothing is being returned, and we would like to indicate success.
}

func (h *EventHandler) DeleteEventById(w http.ResponseWriter, r *http.Request) {
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
	err = h.eventService.DeleteEventById(r.Context(), eventId)

	// TODO: throw error on trying to delete a non-existant event
	fmt.Printf("%v\n", err)

	if err != nil {
		switch err {
		case services.ErrFailedToDeleteEvent:
			res.SendError(w, http.StatusInternalServerError, res.NewError("delete_error", "Failed to delete event"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusNoContent) // We return a 204 for http DELETE requests since nothing is being returned, and we would like to indicate success.
}
