package handlers

import (
	"errors"
	"net/http"

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
