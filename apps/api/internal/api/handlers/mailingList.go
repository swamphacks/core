package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type MailingListHandler struct {
	mailingService *services.MailingListService
	cfg            *config.Config
	logger         zerolog.Logger
}

func NewMailingListHandler(mailingService *services.MailingListService, cfg *config.Config, logger zerolog.Logger) *MailingListHandler {
	return &MailingListHandler{
		mailingService: mailingService,
		cfg:            cfg,
		logger:         logger.With().Str("handler", "MailingListHandler").Str("component", "mailing_list").Logger(),
	}
}

// AddEmailRequest is the expected payload for adding an email
type AddEmailRequest struct {
	Email string `json:"email"`
}

// AddEmailToEvent adds a user's email to a specific event's mailing list.
// It expects 'eventID' and 'userID' to be present in the request context,
// placed there by a routing middleware.
func (h *MailingListHandler) AddEmailToEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Get IDs from request context (placed by a middleware)
	eventIDVal := r.Context().Value("eventID")
	eventIDStr, ok := eventIDVal.(string)
	if !ok {
		h.logger.Error().Msg("eventID not found in request context")
		res.SendError(w, http.StatusInternalServerError, res.NewError("server_error", "Could not process request."))
		return
	}
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("validation_error", "Invalid 'eventID' format in URL"))
		return
	}

	userIDVal := r.Context().Value("userID")
	userIDStr, ok := userIDVal.(string)
	if !ok {
		h.logger.Error().Msg("userID not found in request context")
		res.SendError(w, http.StatusInternalServerError, res.NewError("server_error", "Could not process request."))
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("validation_error", "Invalid 'userID' format in URL"))
		return
	}

	// 2. Decode email from JSON body
	var req AddEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request payload"))
		return
	}

	// 3. Call the service
	createdEmail, err := h.mailingService.AddEmail(r.Context(), eventID, userID, req.Email)
	if err != nil {
		if errors.Is(err, services.ErrEmailConflict) {
			res.SendError(w, http.StatusConflict, res.NewError("conflict", "This user is already on the mailing list for this event."))
		} else if errors.Is(err, services.ErrFailedToAddEmail) {
			h.logger.Error().Err(err).Msg("Service layer failed to add email")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to add email to the mailing list."))
		} else {
			h.logger.Error().Err(err).Msg("An unexpected error occurred while adding email")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "An unexpected server error occurred."))
		}
		return
	}

	// 4. Respond with the created resource
	res.Send(w, http.StatusCreated, createdEmail)
}

// DeleteEmailFromEvent removes a user's email from an event's mailing list.
func (h *MailingListHandler) DeleteEmailFromEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Get IDs from request context
	eventIDVal := r.Context().Value("eventID")
	eventIDStr, ok := eventIDVal.(string)
	if !ok {
		h.logger.Error().Msg("eventID not found in request context")
		res.SendError(w, http.StatusInternalServerError, res.NewError("server_error", "Could not process request."))
		return
	}
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("validation_error", "Invalid 'eventID' format in URL"))
		return
	}

	userIDVal := r.Context().Value("userID")
	userIDStr, ok := userIDVal.(string)
	if !ok {
		h.logger.Error().Msg("userID not found in request context")
		res.SendError(w, http.StatusInternalServerError, res.NewError("server_error", "Could not process request."))
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("validation_error", "Invalid 'userID' format in URL"))
		return
	}

	// 2. Call the service
	err = h.mailingService.DeleteEmailByUserAndEvent(r.Context(), eventID, userID)
	if err != nil {
		if errors.Is(err, services.ErrFailedToDeleteEmail) {
			h.logger.Error().Err(err).Msg("Service layer failed to delete email")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to delete email from the mailing list."))
		} else {
			h.logger.Error().Err(err).Msg("An unexpected error occurred while deleting email")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "An unexpected server error occurred."))
		}
		return
	}

	// 3. Respond with success
	w.WriteHeader(http.StatusNoContent)
}

// GetEmailsByEvent retrieves all emails for a specific event.
func (h *MailingListHandler) GetEmailsByEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("event_id")
	fmt.Println(eventIDStr)
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("validation_error", "Invalid 'eventID' format in URL"))
		return
	}

	emails, err := h.mailingService.GetEmailsByEvent(r.Context(), eventID)
	if err != nil {
		if errors.Is(err, services.ErrFailedToGetEmail) {
			h.logger.Error().Err(err).Msg("Service layer failed to get emails by event")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to retrieve emails for the event."))
		} else {
			h.logger.Error().Err(err).Msg("An unexpected error occurred while getting emails")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "An unexpected server error occurred."))
		}
		return
	}

	res.Send(w, http.StatusOK, emails)
}

// GetEmailsByUser retrieves all mailing list entries for a user.
func (h *MailingListHandler) GetEmailsByUser(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID")
	userIDStr, ok := userIDVal.(string)
	if !ok {
		h.logger.Error().Msg("userID not found in request context")
		res.SendError(w, http.StatusInternalServerError, res.NewError("server_error", "Could not process request."))
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("validation_error", "Invalid 'userID' format in URL"))
		return
	}

	emails, err := h.mailingService.GetEmailsByUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrFailedToGetEmail) {
			h.logger.Error().Err(err).Msg("Service layer failed to get emails by user")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "Failed to retrieve emails for the user."))
		} else {
			h.logger.Error().Err(err).Msg("An unexpected error occurred while getting emails")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_error", "An unexpected server error occurred."))
		}
		return
	}

	res.Send(w, http.StatusOK, emails)
}
