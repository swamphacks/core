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

func (h *EmailHandler) QueueEmail(w http.ResponseWriter, r *http.Request) {
	var req QueueEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	if req.Body == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_body", "Body is missing or is an empty string."))
		return
	}

	for _, to := range req.To {
		if !email.IsValidEmail(to) || !email.IsValidEmail(req.From) {
			res.SendError(w, http.StatusBadRequest, res.NewError("malformed_email", "To and/or From email is malformed or missing"))
			return
		}

		taskInfo, err := h.emailService.QueueSendEmail(to, req.From, req.Body)
		if err != nil {
			h.logger.Err(err).Msg("Failed to queue email sending from EmailHandler")
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "The server went kaput while queueing email sending"))
			return
		}

		h.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued Send Email task!")
	}

	w.WriteHeader(http.StatusCreated)
}
