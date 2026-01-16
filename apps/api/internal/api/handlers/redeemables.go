package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type RedeemablesHandler struct {
	redeemablesService *services.RedeemablesService
	cfg                *config.Config
	logger             zerolog.Logger
}

func NewRedeemablesHandler(
	redeemablesService *services.RedeemablesService,
	cfg *config.Config,
	logger zerolog.Logger,
) *RedeemablesHandler {
	return &RedeemablesHandler{
		redeemablesService: redeemablesService,
		cfg:                cfg,
		logger:             logger,
	}
}

func (h *RedeemablesHandler) GetRedeemables(w http.ResponseWriter, r *http.Request) {
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
	redeemables, err := h.redeemablesService.GetRedeemablesByEventID(r.Context(), eventId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_server_error", "internal service error, failed to get redeemables by ID"))
		return
	}

	res.Send(w, http.StatusOK, redeemables)
}

type CreateRedeemableRequest struct {
	Name          string `json:"name"`
	Amount        int    `json:"amount"`
	MaxUserAmount int    `json:"max_user_amount"`
}

func (h *RedeemablesHandler) CreateRedeemable(w http.ResponseWriter, r *http.Request) {
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

	var req CreateRedeemableRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request_body", "The request body is invalid: "+err.Error()))
		return
	}

	redeemable, err := h.redeemablesService.CreateRedeemable(r.Context(), eventId, req.Name, req.Amount, req.MaxUserAmount)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_server_error", "internal service error, failed to create redeemable"))
		return
	}

	res.Send(w, http.StatusCreated, redeemable)
}

type UpdateRedeemableRequest struct {
	Name          *string `json:"name,omitempty"`
	Amount        *int    `json:"amount,omitempty"`
	MaxUserAmount *int    `json:"max_user_amount,omitempty"`
}

func (h *RedeemablesHandler) UpdateRedeemable(w http.ResponseWriter, r *http.Request) {
	redeemableIdStr := chi.URLParam(r, "redeemableId")
	if redeemableIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_redeemable_id", "The redeemable ID is missing from the URL!"))
		return
	}
	redeemableId, err := uuid.Parse(redeemableIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_redeemable_id", "The redeemable ID is not a valid UUID"))
		return
	}

	var req UpdateRedeemableRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents requests with extraneous fields
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request_body", "The request body is invalid: "+err.Error()))
		return
	}

	redeemable, err := h.redeemablesService.UpdateRedeemable(r.Context(), redeemableId, req.Name, req.Amount, req.MaxUserAmount)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_server_error", "internal service error, failed to update redeemable"))
		return
	}

	res.Send(w, http.StatusCreated, redeemable)
	return
	// decode request body
}

func (h *RedeemablesHandler) DeleteRedeemable(w http.ResponseWriter, r *http.Request) {
	redeemableIdStr := chi.URLParam(r, "redeemableId")
	if redeemableIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_redeemable_id", "The redeemable ID is missing from the URL!"))
		return
	}
	redeemableId, err := uuid.Parse(redeemableIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_redeemable_id", "The redeemable ID is not a valid UUID"))
		return
	}
	err = h.redeemablesService.DeleteRedeemable(r.Context(), redeemableId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_server_error", "internal service error, failed to delete redeemable"))
		return
	}
	res.Send(w, http.StatusNoContent, nil)
}

func (h *RedeemablesHandler) RedeemRedeemable(w http.ResponseWriter, r *http.Request) {
	// user id, redeemable id
	return
}

func (h *RedeemablesHandler) UpdateRedemption(w http.ResponseWriter, r *http.Request) {
	// user id, redeemable id
	return
}
