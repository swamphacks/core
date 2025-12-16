package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type BatHandler struct {
	BatService *services.BatService
	logger     zerolog.Logger
}

func NewBatHandler(BatService *services.BatService, logger zerolog.Logger) *BatHandler {
	return &BatHandler{
		BatService: BatService,
		logger:     logger.With().Str("handler", "BatRunHandler").Str("component", "event_interest").Logger(),
	}
}

// Get BatRuns
//
//	@Summary		Get BatRuns
//	@Description	Gets BatRuns.
//	@Tags			Bat
//	@Accept         json
//	@Produce		json
//	@Success		200		{array}	sqlc.GetBatRunsWithUserInfoRow	"OK: BatRuns returned"
//	@Router			/BatRuns [get]
func (h *BatHandler) GetRunsByEventId(w http.ResponseWriter, r *http.Request) {
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

	runs, err := h.BatService.GetRunsByEventId(r.Context(), eventId)
	if errors.Is(err, services.ErrMissingFields) {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_fields", "Missing/malformed query. Available parameters:  eventId"))
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
	if err := json.NewEncoder(w).Encode(runs); err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong encoding response"))
		return
	}
}
