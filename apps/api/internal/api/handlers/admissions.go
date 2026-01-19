package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/web"
)

type AdmissionHandler struct {
	batService *services.BatService
	logger     zerolog.Logger
}

func NewAdmissionHandler(batService *services.BatService, logger zerolog.Logger) *AdmissionHandler {
	return &AdmissionHandler{
		batService: batService,
		logger:     logger.With().Str("handler", "AdmissionHandler").Logger(),
	}
}

func (h *AdmissionHandler) ReleaseDecisions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest,
			res.NewError("invalid_request", "Invalid or missing event_id."),
		)
		return
	}

	runId, err := web.PathParamToUUID(r, "runId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest,
			res.NewError("invalid_request", "Invalid or missing run_id."),
		)
		return
	}

	err = h.batService.ReleaseBatRunDecision(ctx, eventId, runId)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch {
	case errors.Is(err, services.ErrRunMismatch):
		res.SendError(w, http.StatusForbidden,
			res.NewError("run_mismatch", err.Error()),
		)

	case errors.Is(err, services.ErrRunStatusInvalid):
		res.SendError(w, http.StatusConflict,
			res.NewError("invalid_run_status", err.Error()),
		)

	case errors.Is(err, services.ErrNoAcceptedApplicants):
		res.SendError(w, http.StatusUnprocessableEntity,
			res.NewError("no_accepted_applicants", err.Error()),
		)

	case errors.Is(err, services.ErrCouldNotGetEventInfo):
		res.SendError(w, http.StatusNotFound,
			res.NewError("resource_not_found", err.Error()),
		)

	case errors.Is(err, services.ErrFailedToUpdateRun):
		res.SendError(w, http.StatusInternalServerError,
			res.NewError("update_failed", "Failed to release decisions."),
		)

	default:
		res.SendError(w, http.StatusInternalServerError,
			res.NewError("internal_error", "An unexpected error occurred."),
		)
	}
}

func (h *AdmissionHandler) HandleCalculateAdmissionsRequest(w http.ResponseWriter, r *http.Request) {
	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Invalid request body"))
		return
	}

	_, err = h.batService.QueueCalculateAdmissionsTask(r.Context(), eventId)
	if err != nil {
		res.Send(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went terribly wrong."))
		return
	}

	w.WriteHeader(http.StatusCreated)

}

type EventCheckInRequest struct {
	UserID uuid.UUID `json:"user_id"`
	RFID   *string   `json:"rfid"`
}

// Checks a user into an event. Staff only role.
//
//	@Summary		Check a user into an event
//	@Description	Staff route for checking a user to an event. The user to check in must be an attendee and have never been checked in yet.
//	@Tags			Admissions
//	@Param			sh_session_id	cookie	string				true	"The authenticated session token/id"
//	@Param			event_id		path	int					true	"The ID of the event"
//	@Param			request			body	EventCheckInRequest	true	"Event check in data"
//	@Success		204				"No Content"
//	@Failure		400				{object}	response.ErrorResponse	"Malformed request body."
//	@Failure		401				{object}	response.ErrorResponse	"Unauthenticated: Requester is not currently authenticated."
//	@Failure		500				{object}	response.ErrorResponse	"Something went seriously wrong."
//
//	@Router			/events/{eventId}/checkin [post]
func (h *AdmissionHandler) HandleEventCheckIn(w http.ResponseWriter, r *http.Request) {
	eventId, err := web.PathParamToUUID(r, "eventId")
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_param", "Missing event id"))
		return
	}

	var req EventCheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	if *req.RFID == "" {
		req.RFID = nil
	}

	err = h.batService.CheckInAttendee(r.Context(), eventId, req.UserID, req.RFID)
	if err != nil {
		res.Send(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went terribly wrong."))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
