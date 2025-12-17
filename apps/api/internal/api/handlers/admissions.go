package handlers

import (
	"net/http"

	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/services"
	"github.com/swamphacks/core/apps/api/internal/web"
)

type AdmissionHandler struct {
	batService *services.BatService
}

func NewAdmissionHandler(batService *services.BatService) *AdmissionHandler {
	return &AdmissionHandler{
		batService: batService,
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
