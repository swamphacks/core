package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/email"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type EmailHandler struct {
	emailService *services.EmailService
	logger       zerolog.Logger
}

func NewEmailHandler(emailService *services.EmailService, logger zerolog.Logger) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
		logger:       logger.With().Str("handler", "EmailHandler").Str("component", "email").Logger(),
	}
}

type QueueEmailRequest struct {
	To   []string `json:"to"`
	From string   `json:"from"`
	Body string   `json:"body"`
}

// Queue an Email Request
//
//	@Summary		Queue an Email Request
//	@Description	Push an email request to the task queue
//	@Tags			Email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		QueueEmailRequest		true	"Email data"
//	@Success		201		{object}	string					"OK: Email request queued"
//	@Failure		400		{object}	response.ErrorResponse	"Bad request/Malformed request. The email request is potentially invalid."
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: The server went kaput while queueing email sending"
//	@Router			/email/queue [post]
func (h *EmailHandler) QueueEmail(w http.ResponseWriter, r *http.Request) {
	// TODO: remove "from" field as this will be from an env var
	var req QueueEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	if !email.IsValidEmail(req.From) {
		res.SendError(w, http.StatusBadRequest, res.NewError("malformed_email", "'From' email is malformed or missing"))
		return
	}
	for _, to := range req.To {
		if !email.IsValidEmail(to) {
			res.SendError(w, http.StatusBadRequest, res.NewError("malformed_email", "'To' email is malformed or missing"))
			return
		}
	}

	if req.Body == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_body", "Body is missing or is an empty string."))
		return
	}

	taskInfo, err := h.emailService.QueueSendEmail(req.To, req.From, req.Body)
	if err != nil {
		h.logger.Err(err).Msg("Failed to queue email sending from EmailHandler")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "The server went kaput while queueing email sending"))
		return
	}

	h.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued Send Email task!")

	w.WriteHeader(http.StatusCreated)
}
