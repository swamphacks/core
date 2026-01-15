package services

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
)

type RedeemablesService struct {
	redeemablesRepo *repository.RedeemablesRepository
	logger          zerolog.Logger
}

func NewRedeemablesService(
	redeemablesRepo *repository.RedeemablesRepository,
	logger zerolog.Logger) *RedeemablesService {
	return &RedeemablesService{
		redeemablesRepo: redeemablesRepo,
		logger:          logger,
	}
}
