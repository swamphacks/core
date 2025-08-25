package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/parse"
	"github.com/swamphacks/core/apps/api/internal/ptr"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/web"
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
		if err == services.ErrFailedToCreateEvent {
			res.SendError(w, http.StatusInternalServerError, res.NewError("creation_error", "Failed to create event"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	res.Send(w, http.StatusCreated, event)
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
		case services.ErrFailedToGetEvent:
			res.SendError(w, http.StatusNotFound, res.NewError("no_event", "Event not found"))
		default:
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	res.Send(w, http.StatusOK, event)
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	// Refactorme: could be improved by unmarshalling values into a generic that can include nil information
	// Todo: make sure that non nullable values can't be updated to null
	// Todo: Time validation
	req.NameDoUpdate = reflect.ValueOf(req.Name).IsValid()
	req.DescriptionDoUpdate = reflect.ValueOf(req.Description).IsValid()
	req.LocationDoUpdate = reflect.ValueOf(req.Location).IsValid()
	req.LocationUrlDoUpdate = reflect.ValueOf(req.LocationUrl).IsValid()
	req.MaxAttendeesDoUpdate = reflect.ValueOf(req.MaxAttendees).IsValid()
	req.ApplicationOpenDoUpdate = reflect.ValueOf(req.ApplicationOpen).IsValid()
	req.ApplicationCloseDoUpdate = reflect.ValueOf(req.ApplicationClose).IsValid()
	req.RsvpDeadlineDoUpdate = reflect.ValueOf(req.RsvpDeadline).IsValid()
	req.DecisionReleaseDoUpdate = reflect.ValueOf(req.DecisionRelease).IsValid()
	req.StartTimeDoUpdate = reflect.ValueOf(req.StartTime).IsValid()
	req.EndTimeDoUpdate = reflect.ValueOf(req.EndTime).IsValid()
	req.WebsiteUrlDoUpdate = reflect.ValueOf(req.WebsiteUrl).IsValid()
	req.IsPublishedDoUpdate = reflect.ValueOf(req.IsPublished).IsValid()
	req.ID = eventId

	event, err := h.eventService.UpdateEventById(r.Context(), req)

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

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	// include_unpublished="true,false", default: false
	queryParams := r.URL.Query()
	includeUnpublished, err := web.ParseParamBoolean(queryParams, "include_unpublished", ptr.BoolToPtr(false))
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Missing/malformed query. Available parameters: include_unpublished=true,false"))
		return
	}

	events, err := h.eventService.GetEvents(r.Context(), *includeUnpublished)
	if err == services.ErrMissingFields {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Missing/malformed query. Available parameteres: status=all,published"))
		return
	}

	if err == services.ErrMissingPerms {
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

// func (h *EventHandler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
// 	// filename := r.URL.Query().Get("filename")

// 	// if filename == "" {
// 	// 	res.SendError(w, http.StatusBadRequest, res.NewError("a", "no file name"))
// 	// 	return
// 	// }

// 	// // Ensure directory exists
// 	// if err := os.MkdirAll("./uploads", 0755); err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// out, err := os.Create("./uploads/" + filename)

// 	// if err != nil {
// 	// 	fmt.Println(err)

// 	// 	res.SendError(w, http.StatusBadRequest, res.NewError("a", "can't create file"))
// 	// 	return
// 	// }
// 	// defer out.Close()
// 	// _, err = io.Copy(out, r.Body)
// 	// if err != nil {
// 	// 	res.SendError(w, http.StatusBadRequest, res.NewError("a", "can't copy file"))
// 	// 	return
// 	// }

// 	// w.Header().Set("Content-Type", "text/plain")
// 	// fmt.Fprintln(w, "Uploaded:", filename)
// 	sdkConfig, err := awsConfig.LoadDefaultConfig(context.TODO())

// 	if err != nil {
// 		res.SendError(w, http.StatusInternalServerError, res.NewError("a", "load sdk config"))
// 		return
// 	}

// 	s3Client := s3.NewFromConfig(sdkConfig)
// 	presignClient := s3.NewPresignClient(s3Client)
// 	presigner := Presigner{PresignClient: presignClient}

// 	presignedPutRequest, err := presigner.PutObject(r.Context(), "hieu-infra-test", "testKey", 60)

// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("Got a presigned %v request to URL:\n\t%v\n", presignedPutRequest.Method,
// 		presignedPutRequest.URL)

// 	res.Send(w, http.StatusOK, presignedPutRequest)
// }

// type Presigner struct {
// 	PresignClient *s3.PresignClient
// }

// // PutObject makes a presigned request that can be used to put an object in a bucket.
// // The presigned request is valid for the specified number of seconds.
// func (presigner Presigner) PutObject(
// 	ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {

// 	ContentType := "application/pdf"
// 	request, err := presigner.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
// 		Bucket:      aws.String(bucketName),
// 		Key:         aws.String(objectKey),
// 		ContentType: &ContentType,
// 	}, func(opts *s3.PresignOptions) {
// 		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
// 			bucketName, objectKey, err)
// 	}
// 	return request, err
// }
