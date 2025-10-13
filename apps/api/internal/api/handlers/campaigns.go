package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	res "github.com/swamphacks/core/apps/api/internal/api/response"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/services"
)

type CampaignHandler struct {
	campaignService *services.CampaignService
	cfg             *config.Config
	logger          zerolog.Logger
}

func NewCampaignHandler(campaignService *services.CampaignService, logger zerolog.Logger) *CampaignHandler {
	return &CampaignHandler{
		campaignService: campaignService,
		logger:          logger.With().Str("handler", "CampaignHandler").Str("component", "campaigns").Logger(),
	}
}

type CreateCampaignFields struct {
	EventID        uuid.UUID `json:"event_id" validate:"required"`
	Title          string    `json:"title" validate:"required"`
	Description    *string   `json:"description"`
	RecipientRoles *[]string `json:"recipient_roles"`
	CreatedBy      uuid.UUID `json:"created_by" validate:"required"`
}

func (h *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {

	// Parse JSON body
	var req CreateCampaignFields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", "Could not parse request body"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		res.SendError(w, http.StatusBadRequest, res.NewError("invalid_request", err.Error()))
	}

	params := sqlc.CreateCampaignParams{
		EventID:        req.EventID,
		Title:          req.Title,
		Description:    req.Description,
		RecipientRoles: req.RecipientRoles,
		CreatedBy:      req.CreatedBy,
	}

	campaign, err := h.campaignService.CreateCampaign(r.Context(), params)
	if err != nil {
		if errors.Is(err, services.ErrFailedToCreateCampaign) {
			res.SendError(w, http.StatusInternalServerError, res.NewError("creation_error", "Failed to create campaign"))
		} else {
			res.SendError(w, http.StatusInternalServerError, res.NewError("internal_err", "Something went wrong"))
		}
	}

	res.Send(w, http.StatusCreated, campaign)
}
