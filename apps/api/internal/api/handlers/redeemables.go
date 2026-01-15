package handlers

import (
	"net/http"

	"github.com/rs/zerolog"
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
	return
}

func (h *RedeemablesHandler) CreateRedeemable(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *RedeemablesHandler) UpdateRedeemable(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *RedeemablesHandler) DeleteRedeemable(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *RedeemablesHandler) RedeemRedeemable(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *RedeemablesHandler) UpdateRedemption(w http.ResponseWriter, r *http.Request) {
	return
}
