package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/parse"
	. "github.com/swamphacks/core/apps/api/internal/parse"
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
	Name             string     `json:"name" validate:"required,min=5,max=30"`
	ApplicationOpen  time.Time  `json:"application_open" validate:"required"`
	ApplicationClose time.Time  `json:"application_close" validate:"required"`
	StartTime        time.Time  `json:"start_time" validate:"required"`
	EndTime          time.Time  `json:"end_time" validate:"required"`
	Description      *string    `json:"description"`
	Location         *string    `json:"location"`
	LocationUrl      *string    `json:"location_url"`
	MaxAttendees     *int32     `json:"max_attendees"`
	RsvpDeadline     *time.Time `json:"rsvp_deadline"`
	DecisionRelease  *time.Time `json:"decision_release"`
	WebsiteUrl       *string    `json:"website_url"`
	IsPublished      *bool      `json:"is_published"`
}

func (st CreateEventFields) ValidateTimeFields() bool {
	if st.ApplicationClose.Before(st.ApplicationOpen) || st.ApplicationClose.Equal(st.ApplicationOpen) {
		return false
	}
	if st.EndTime.Before(st.StartTime) || st.EndTime.Equal(st.StartTime) {
		return false
	}
	if st.ApplicationOpen.Before(time.Now()) ||
		st.ApplicationClose.Before(time.Now()) ||
		st.StartTime.Before(time.Now()) ||
		st.EndTime.Before(time.Now()) {
		return false
	}
	return true
}

// Create a new event
//
//	@Summary		Create a new event
//	@Description	Create a new event with the provided details
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateEventFields		true	"Event creation data"
//	@Success		201		{object}	sqlc.Event				"OK: Event created"
//	@Failure		400		{object}	response.ErrorResponse	"Bad request/Malformed request"
//	@Failure		409		{object}	response.ErrorResponse	"endTime is before startTime or applicationClose is before applicationOpen"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events [post]
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	// Parse JSON body
	var req CreateEventFields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	// This will also throw an error for empty values for fields which correspond to types that cannot convert an empty string to a zero value (e.g. time.Time)
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", err.Error()))
	}

	if !req.ValidateTimeFields() {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_time", "Time fields must be sequential and not in the past."))
		return
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

	event, err := h.eventService.CreateEvent(r.Context(), params)
	if err != nil {
		if errors.Is(err, services.ErrFailedToCreateEvent) {
			res.SendError(w, http.StatusInternalServerError, res.NewError("creation_error", "Failed to create event"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	res.Send(w, http.StatusCreated, event)
}

// Get an event
//
//	@Summary		Get an event
//	@Description	Get a specific event by ID
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path		string					true	"Event ID"	Format(uuid)
//	@Success		201		{object}	sqlc.Event				"OK - Event received"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId} [get]
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
		case services.ErrFailedToGetEvent:
			res.SendError(w, http.StatusNotFound, res.NewError("no_event", "Event not found"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	res.Send(w, http.StatusOK, event)
}

type UpdateEventFields struct {
	Name             Optional[string]     `json:"name"`
	Description      Optional[*string]    `json:"description"`
	Location         Optional[*string]    `json:"location"`
	LocationUrl      Optional[*string]    `json:"location_url"`
	MaxAttendees     Optional[*int32]     `json:"max_attendees"`
	ApplicationOpen  Optional[time.Time]  `json:"application_open"`
	ApplicationClose Optional[time.Time]  `json:"application_close"`
	RsvpDeadline     Optional[*time.Time] `json:"rsvp_deadline"`
	DecisionRelease  Optional[*time.Time] `json:"decision_release"`
	StartTime        Optional[time.Time]  `json:"start_time"`
	EndTime          Optional[time.Time]  `json:"end_time"`
	WebsiteUrl       Optional[*string]    `json:"website_url"`
	IsPublished      Optional[bool]       `json:"is_published"`
}

// Update an event
//
//	@Summary		Update an event
//	@Description	Update an existing event
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	string	true	"Event ID"	Format(uuid)
//	@Success		204		"OK - Event updated (patched)"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId} [patch]
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

	var req UpdateEventFields

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	var params = sqlc.UpdateEventByIdParams{
		NameDoUpdate: req.Name.Present,
		Name:         req.Name.Value,

		DescriptionDoUpdate: req.Description.Present,
		Description:         req.Description.Value,

		LocationDoUpdate: req.Location.Present,
		Location:         req.Location.Value,

		LocationUrlDoUpdate: req.LocationUrl.Present,
		LocationUrl:         req.LocationUrl.Value,

		MaxAttendeesDoUpdate: req.MaxAttendees.Present,
		MaxAttendees:         req.MaxAttendees.Value,

		ApplicationOpenDoUpdate: req.ApplicationOpen.Present,
		ApplicationOpen:         req.ApplicationOpen.Value,

		ApplicationCloseDoUpdate: req.ApplicationClose.Present,
		ApplicationClose:         req.ApplicationClose.Value,

		RsvpDeadlineDoUpdate: req.RsvpDeadline.Present,
		RsvpDeadline:         req.RsvpDeadline.Value,

		DecisionReleaseDoUpdate: req.DecisionRelease.Present,
		DecisionRelease:         req.DecisionRelease.Value,

		StartTimeDoUpdate: req.StartTime.Present,
		StartTime:         req.StartTime.Value,

		EndTimeDoUpdate: req.EndTime.Present,
		EndTime:         req.EndTime.Value,

		WebsiteUrlDoUpdate: req.WebsiteUrl.Present,
		WebsiteUrl:         req.WebsiteUrl.Value,

		IsPublishedDoUpdate: req.IsPublished.Present,
		IsPublished:         &req.IsPublished.Value,

		BannerDoUpdate: false, // Banners are uploaded using a separate endpoint
		Banner:         nil,

		ID: eventId,
	}

	event, err := h.eventService.UpdateEventById(r.Context(), params)

	if err != nil {
		switch err {
		case services.ErrFailedToUpdateEvent:
			res.SendError(w, http.StatusInternalServerError, res.NewError("patch_error", "Failed to update event"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	res.Send(w, http.StatusOK, event)
}

type NullableEventRole struct {
	UserID     uuid.UUID           `json:"user_id"`
	EventID    uuid.UUID           `json:"event_id"`
	Role       *sqlc.EventRoleType `json:"role"`
	AssignedAt *time.Time          `json:"assigned_at"`
}

// Get the current user's event role for an event
//
//	@Summary		Get the current user's event role for an event
//	@Description	Get current user's role for a specific event
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path		string					true	"Event ID"	Format(uuid)
//	@Success		200		{object}	NullableEventRole		"OK - Return role"
//	@Failure		400		{object}	response.ErrorResponse	"Not Found - Role not found"
//	@Failure		404		{object}	response.ErrorResponse	"Not Found - Role not found"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/role [get]
func (h *EventHandler) GetEventRole(w http.ResponseWriter, r *http.Request) {
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

	userId := ctxutils.GetUserIdFromCtx(r.Context())
	if userId == nil {
		res.SendError(w, http.StatusUnauthorized, res.NewError("unauthorized", "User not authenticated"))
		return
	}

	eventRole, err := h.eventService.GetEventRoleByIds(r.Context(), *userId, eventId)
	if err != nil {
		if errors.Is(err, repository.ErrEventRoleNotFound) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(NullableEventRole{
				UserID:     *userId,
				EventID:    eventId,
				Role:       nil,
				AssignedAt: nil,
			})
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went terribly wrong."))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NullableEventRole{
		UserID:     eventRole.UserID,
		EventID:    eventRole.EventID,
		Role:       &eventRole.Role,
		AssignedAt: eventRole.AssignedAt,
	})
}

// Delete an event
//
//	@Summary		Delete an event
//	@Description	Delete an existing event
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	string	true	"Event ID"	Format(uuid)
//	@Success		204		"OK - Event deleted"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId} [delete]
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

	if err != nil {
		switch err {
		case services.ErrFailedToDeleteEvent:
			res.SendError(w, http.StatusInternalServerError, res.NewError("delete_error", "Failed to delete event"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get events
//
//	@Summary		Get events
//	@Description	Gets events with a nullable event role for authenticated users.
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			scope	query	string							false	"Can be scoped to either published, scoped, or all. Scoped means admins and staff can see unpublished events"	default("published")
//	@Success		200		{array}	sqlc.GetEventsWithUserInfoRow	"OK: Events returned"
//	@Router			/events [get]
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	scope, err := parse.ParseGetEventScopeType(q.Get("scope"))

	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Missing/malformed query. Available parameters: include_unpublished=published, scoped, all, or none (default to published)"))
		return
	}

	events, err := h.eventService.GetEvents(r.Context(), scope)
	if errors.Is(err, services.ErrMissingFields) {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Missing/malformed query. Available parameters:  include_unpublished=published, scoped, all, or none (default to published)"))
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
	if err := json.NewEncoder(w).Encode(events); err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
		return
	}
}

// Get all staff users for an event
//
//	@Summary		Get all staff users for an event
//	@Description	Gets all users with role STAFF or ADMIN
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path		string					true	"Event ID"	Format(uuid)
//	@Success		200		{array}		sqlc.GetEventStaffRow	"OK - Return users"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/staff [get]
func (h *EventHandler) GetEventStaffUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.eventService.GetEventStaffUsers(r.Context(), eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
		return
	}
}

type AssignRoleFields struct {
	Email  *string            `json:"email"`
	UserID *string            `json:"user_id"`
	Role   sqlc.EventRoleType `json:"role"`
}

func (r *AssignRoleFields) Validate() error {
	if r.Email != nil && !email.IsValidEmail(*r.Email) {
		return fmt.Errorf("invalid email: %s", *r.Email)
	}

	switch r.Role {
	case sqlc.EventRoleTypeAdmin, sqlc.EventRoleTypeStaff, sqlc.EventRoleTypeAttendee, sqlc.EventRoleTypeApplicant:
		return nil
	default:
		return fmt.Errorf("invalid role: %q", r.Role)
	}
}

// Change or add event role of a user
//
//	@Summary		Change or add event role of a user
//	@Description	Modify user's role for a specific event
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	string				true	"Event ID"	Format(uuid)
//	@Param			request	body	AssignRoleFields	true	"Event role data"
//	@Success		200		"OK - Role updated"
//	@Failure		404		{object}	response.ErrorResponse	"Not Found - User not found"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/roles [post]
func (h *EventHandler) AssignEventRole(w http.ResponseWriter, r *http.Request) {
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

	var input AssignRoleFields
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_body", err.Error()))
		return
	}

	if err := input.Validate(); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_body", err.Error()))
		return
	}

	userId := parse.ParseUUIDOrNil(input.UserID)
	email := parse.ParseStrToPtr(input.Email)

	err = h.eventService.AssignEventRole(r.Context(), userId, email, eventId, input.Role)
	if err != nil {
		switch err {
		case repository.ErrUserNotFound:
			res.SendError(w, http.StatusNotFound, res.NewError("user_missing", "The user does not exist"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong on our end"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

type AssignRoleBatch struct {
	Assignments []AssignRoleFields `json:"assignments"`
}

func (b *AssignRoleBatch) Validate() error {
	if len(b.Assignments) == 0 {
		return fmt.Errorf("at least one assignment is required")
	}
	for i, a := range b.Assignments {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("assignment[%d]: %w", i, err)
		}
	}
	return nil
}

// Change or add event role of a user in batch
//
//	@Summary		Change or add event role of a user in batch
//	@Description	Modify users' role for a specific event
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	string			true	"Event ID"	Format(uuid)
//	@Param			request	body	AssignRoleBatch	true	"Event roles data"
//	@Success		200		"OK - Roles updated"
//	@Failure		404		{object}	response.ErrorResponse	"Not Found - User not found"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/roles/batch [post]
func (h *EventHandler) BatchAssignEventRoles(w http.ResponseWriter, r *http.Request) {
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

	var input AssignRoleBatch
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_body", err.Error()))
		return
	}

	if err := input.Validate(); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_body", err.Error()))
		return
	}

	for _, assignment := range input.Assignments {
		userId := parse.ParseUUIDOrNil(assignment.UserID)
		email := parse.ParseStrToPtr(assignment.Email)

		err = h.eventService.AssignEventRole(r.Context(), userId, email, eventId, assignment.Role)
		if err != nil {
			switch err {
			case repository.ErrUserNotFound:
				res.SendError(w, http.StatusNotFound, res.NewError("user_missing", fmt.Sprintf("The user %v does not exist", assignment.UserID)))
			default:
				res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong on our end"))
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Revoke event role of a user
//
//	@Summary		Revoke event role of a user
//	@Description	Remove user's role for a specific event
//	@Tags			Event
//	@Produce		json
//	@Param			eventId	path	string	true	"Event ID"	Format(uuid)
//	@Param			userId	path	string	true	"User ID"	Format(uuid)
//	@Success		200		"OK - Role revoked"
//	@Failure		404		{object}	response.ErrorResponse	"Not Found - User not found"
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: Something went terribly wrong on our end."
//	@Router			/events/{eventId}/roles/{userId} [delete]
func (h *EventHandler) RevokeEventRole(w http.ResponseWriter, r *http.Request) {
	eventIdStr := chi.URLParam(r, "eventId")
	userIdStr := chi.URLParam(r, "userId")
	if eventIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_event_id", "The event ID is missing from the URL!"))
		return
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_event_id", "The event ID is not a valid UUID"))
		return
	}

	if userIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_user_id", "The user ID is missing from the URL!"))
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "The user ID is not a valid UUID"))
		return
	}

	err = h.eventService.RevokeEventRole(r.Context(), userId, eventId)
	if err != nil {
		switch err {
		case repository.ErrUserNotFound:
			res.SendError(w, http.StatusNotFound, res.NewError("user_missing", "The user does not exist"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong on our end"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deferredCloser(c io.Closer, name string) func() {
	return func() {
		if err := c.Close(); err != nil {
			log.Err(err).Msg("Failed to close " + name)
		}
	}
}

const maxBannerUploadSize = 5 << 20 // 5 Mb

type EventBannerUploadResponse struct {
	BannerUrl string `json:"banner_url"`
}

func (h *EventHandler) UploadEventBanner(w http.ResponseWriter, r *http.Request) {
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

	if err := r.ParseMultipartForm(maxBannerUploadSize); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "could not parse multipart form"))
		return
	}

	bannerFile, header, err := r.FormFile("image")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "invalid resume file"))
		return
	}
	defer deferredCloser(bannerFile, "banner file")

	url, err := h.eventService.UploadBanner(r.Context(), eventId, bannerFile, header)
	switch err {
	case services.ErrFailedToUploadBanner:
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong on our end"))
		return
	case services.ErrUnexpectedFileType:
		res.SendError(w, http.StatusBadRequest, res.NewError("file_error", err.Error()))
		return
	case nil:
		// Continue
	default:
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong on our end"))
		return
	}

	res.Send(w, http.StatusOK, EventBannerUploadResponse{
		BannerUrl: *url,
	})

}

func (h *EventHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
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

	err = h.eventService.DeleteBanner(r.Context(), eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong on our end"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
