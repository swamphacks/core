package redeemables

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

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

func (s *RedeemablesService) CreateRedeemable(ctx context.Context, name string, amount int, maxUserAmount int) (*sqlc.Redeemable, error) {
	redeemable, err := s.redeemablesRepo.CreateRedeemable(ctx, name, amount, maxUserAmount)
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
	redeemable, err := s.redeemablesRepo.UpdateRedeemable(ctx, redeemableID, name, amount, maxUserAmount)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update redeemable")
		return nil, err
	}
	return redeemable, nil
}

func (s *RedeemablesService) RedeemRedeemable(ctx context.Context, redeemableID uuid.UUID, userID uuid.UUID) error {
	// Need to do a check to see if the user is checked in
	// Probably need event service

	// CREATE NEW SQL function for getting checked in status
	_, err := s.redeemablesRepo.RedeemRedeemable(ctx, redeemableID, userID)

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to redeem redeemable")
		return err
	}
	return nil
}

func (s *RedeemablesService) UpdateRedemption(ctx context.Context, redeemableID uuid.UUID, userID uuid.UUID, amount int) error {
	err := s.redeemablesRepo.UpdateRedemption(ctx, redeemableID, userID, amount)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update redemption")
		return err
	}
	return nil
}
