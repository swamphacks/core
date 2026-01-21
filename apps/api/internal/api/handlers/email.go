package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

type QueueTextEmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// Queue an Email Request
//
//	@Summary		Queue an Email Request
//	@Description	Push an email request to the task queue
//	@Tags			Email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		QueueTextEmailRequest	true	"Email data"
//	@Success		201		{object}	string					"OK: Email request queued"
//	@Failure		400		{object}	response.ErrorResponse	"Bad request/Malformed request. The email request is potentially invalid."
//	@Failure		500		{object}	response.ErrorResponse	"Server Error: The server went kaput while queueing email sending"
//	@Router			/email/queue [post]
func (h *EmailHandler) QueueTextEmail(w http.ResponseWriter, r *http.Request) {
	var req QueueTextEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	for _, to := range req.To {
		if !email.IsValidEmail(to) {
			res.SendError(w, http.StatusBadRequest, res.NewError("malformed_email", "'To' email is malformed or missing"))
			return
		}
	}

	if req.Subject == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_subject", "Subject is missing or is an empty string."))
		return
	}

	if req.Body == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_body", "Body is missing or is an empty string."))
		return
	}

	taskInfo, err := h.emailService.QueueSendTextEmail(req.To, req.Subject, req.Body)
	if err != nil {
		h.logger.Err(err).Msg("Failed to queue SendTextEmail from EmailHandler")
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "The server went kaput while queueing email sending"))
		return
	}

	h.logger.Info().Str("TaskID", taskInfo.ID).Str("Task Queue", taskInfo.Queue).Str("Task Type", taskInfo.Type).Msg("Queued SendTextEmail task!")

	w.WriteHeader(http.StatusCreated)
}

type QueueConfirmationEmailFields struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
}

type QueueWelcomeEmailFields struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	UserId    string `json:userId validate:"required"`
}

func (h *EmailHandler) QueueConfirmationEmail(w http.ResponseWriter, r *http.Request) {
	var req QueueConfirmationEmailFields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", err.Error()))
	}

	err = h.emailService.QueueConfirmationEmail(req.Email, req.FirstName)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Confirmation email could not be queued."))
	}

	res.Send(w, http.StatusOK, nil)
}

func (h *EmailHandler) QueueWelcomeEmail(w http.ResponseWriter, r *http.Request) {
	var req QueueWelcomeEmailFields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", err.Error()))
	}

	parsedUserId, err := uuid.Parse(req.UserId)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "userId must be of type uuid"))
		return
	}

	err = h.emailService.QueueWelcomeEmail(req.Email, req.FirstName, parsedUserId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Hacker email could not be queued."))
	}

	res.Send(w, http.StatusOK, nil)
}
