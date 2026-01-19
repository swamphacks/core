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

// GetRedeemables
//
//	@Summary      Get all redeemables for an event
//	@Description   Retrieve a list of all redeemable items associated with a specific event ID.
//	@Tags         Redeemables
//	@Accept       json
//	@Produce      json
//	@Param        eventId path      string  true  "Event ID (UUID)"
//	@Success      200     {array}   sqlc.Redeemable
//	@Failure      400     {object}  response.ErrorResponse "Missing or invalid Event ID"
//	@Failure      500     {object}  response.ErrorResponse "Internal Server Error"
//	@Router       /events/{eventId}/redeemables [get]
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

// CreateRedeemable
//
//	@Summary      Create a new redeemable
//	@Description   Create a new redeemable item for a specific event.
//	@Tags         Redeemables
//	@Accept       json
//	@Produce      json
//	@Param        eventId path      string                   true  "Event ID (UUID)"
//	@Param        request body      CreateRedeemableRequest  true  "Redeemable creation data"
//	@Success      201     {object}  sqlc.Redeemable
//	@Failure      400     {object}  response.ErrorResponse "Invalid request body or ID"
//	@Failure      500     {object}  response.ErrorResponse "Internal Server Error"
//	@Router       /events/{eventId}/redeemables [post]
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
	Amount        *int    `json:"total_stock,omitempty"`
	MaxUserAmount *int    `json:"max_user_amount,omitempty"`
}

// UpdateRedeemable
//
//	@Summary      Update an existing redeemable
//	@Description   Update specific fields (name, stock, max per user) of a redeemable.
//	@Tags         Redeemables
//	@Accept       json
//	@Produce      json
//	@Param        redeemableId path      string                   true  "Redeemable ID (UUID)"
//	@Param        request      body      UpdateRedeemableRequest  true  "Redeemable update data (partial fields allowed)"
//	@Success      201          {object}  sqlc.Redeemable
//	@Failure      400          {object}  response.ErrorResponse "Invalid ID or request body"
//	@Failure      500          {object}  response.ErrorResponse "Internal Server Error"
//	@Router       /redeemables/{redeemableId} [patch]
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
}

// DeleteRedeemable
//
//	@Summary      Delete a redeemable
//	@Description   Permanently delete a redeemable item by ID.
//	@Tags         Redeemables
//	@Param        redeemableId path  string  true  "Redeemable ID (UUID)"
//	@Success      204          "No Content"
//	@Failure      400          {object}  response.ErrorResponse "Invalid Redeemable ID"
//	@Failure      500          {object}  response.ErrorResponse "Internal Server Error"
//	@Router       /redeemables/{redeemableId} [delete]
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

// RedeemRedeemable
//
//	@Summary      Redeem an item for a user
//	@Description   Create a redemption record linking a specific user to a redeemable item.
//	@Tags         Redeemables
//	@Param        redeemableId path  string  true  "Redeemable ID (UUID)"
//	@Param        userId       path  string  true  "User ID (UUID)"
//	@Success      204          "No Content"
//	@Failure      400          {object}  response.ErrorResponse "Invalid IDs"
//	@Failure      500          {object}  response.ErrorResponse "Internal Server Error"
//	@Router       /redeemables/{redeemableId}/users/{userId} [post]
func (h *RedeemablesHandler) RedeemRedeemable(w http.ResponseWriter, r *http.Request) {
	// user id, redeemable id
	userIdStr := chi.URLParam(r, "userId")
	if userIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_user_id", "The user ID is missing from the URL!"))
		return
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "The user ID is not a valid UUID"))
		return
	}
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
	err = h.redeemablesService.RedeemRedeemable(r.Context(), redeemableId, userId)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_server_error", "internal service error, failed to redeem redeemable"))
		return
	}
	res.Send(w, http.StatusNoContent, nil)
}

type UpdateRedemptionRequest struct {
	amount int `json:"new_amount,omitempty"`
}

func (h *RedeemablesHandler) UpdateRedemption(w http.ResponseWriter, r *http.Request) {
	// user id, redeemable id
	userIdStr := chi.URLParam(r, "userId")
	if userIdStr == "" {
		res.SendError(w, http.StatusBadRequest, res.NewError("missing_user_id", "The user ID is missing from the URL!"))
		return
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_user_id", "The user ID is not a valid UUID"))
		return
	}
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
	var req UpdateRedemptionRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request_body", "The request body is invalid: "+err.Error()))
		return
	}

	err = h.redeemablesService.UpdateRedemption(r.Context(), redeemableId, userId, req.amount)
	if err != nil {
		res.SendError(w, http.StatusInternalServerError, res.NewError("internal_server_error", "internal service error, failed to update redemption"))
		return
	}
	res.Send(w, http.StatusNoContent, nil)
}
