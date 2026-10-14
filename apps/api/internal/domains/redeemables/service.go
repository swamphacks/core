package redeemables

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

// ErrNoActiveHackathon is re-exported so the handler checks the service-level
// name, matching ErrEmailCampaignNotFound in the email package.
var ErrNoActiveHackathon = repository.ErrNoActiveHackathon

type RedeemablesService struct {
	redeemablesRepo *repository.RedeemablesRepository
	logger          zerolog.Logger
}

func NewService(redeemablesRepo *repository.RedeemablesRepository, logger zerolog.Logger) *RedeemablesService {
	return &RedeemablesService{
		redeemablesRepo: redeemablesRepo,
		logger:          logger,
	}
}

func (s *RedeemablesService) GetRedeemables(ctx context.Context) (*[]sqlc.GetRedeemablesRow, error) {
	redeemables, err := s.redeemablesRepo.GetRedeemables(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get redeemables by event ID")
		return nil, err
	}
	return redeemables, nil
}

func (s *RedeemablesService) CreateRedeemable(ctx context.Context, name string, amount int, maxUserAmount int, redeemableType string) (*sqlc.Redeemable, error) {
	params := sqlc.CreateRedeemableParams{
		Name:          name,
		Amount:        int32(amount),
		MaxUserAmount: int32(maxUserAmount),
		Type:          redeemableType,
	}
	redeemable, err := s.redeemablesRepo.CreateRedeemable(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create redeemable")
		return nil, err
	}
	return redeemable, nil
}

func (s *RedeemablesService) DeleteRedeemable(ctx context.Context, redeemableID uuid.UUID) error {
	err := s.redeemablesRepo.DeleteRedeemable(ctx, redeemableID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to delete redeemable")
		return err
	}
	return nil
}

func (s *RedeemablesService) UpdateRedeemable(ctx context.Context, redeemableID uuid.UUID, name *string, amount *int, maxUserAmount *int) (*sqlc.Redeemable, error) {

	var amount32 *int32
	if amount != nil {
		v := int32(*amount)
		amount32 = &v
	}
	var maxUserAmount32 *int32
	if maxUserAmount != nil {
		v := int32(*maxUserAmount)
		maxUserAmount32 = &v
	}

	redeemable, err := s.redeemablesRepo.UpdateRedeemable(ctx, sqlc.UpdateRedeemableParams{
		ID:            redeemableID,
		Name:          name,
		Amount:        amount32,
		MaxUserAmount: maxUserAmount32,
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update redeemable")
		return nil, err
	}
	return redeemable, nil
}
