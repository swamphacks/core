package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
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

func (s *RedeemablesService) GetRedeemablesByEventID(ctx context.Context, eventID uuid.UUID) (*[]sqlc.GetRedeemablesByEventIDRow, error) {
	redeemables, err := s.redeemablesRepo.GetRedeemablesByEventID(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get redeemables by event ID")
		return nil, err
	}
	return redeemables, nil
}

func (s *RedeemablesService) CreateRedeemable(ctx context.Context, eventID uuid.UUID, name string, amount int, maxUserAmount int) (*sqlc.Redeemable, error) {
	redeemable, err := s.redeemablesRepo.CreateRedeemable(ctx, eventID, name, amount, maxUserAmount)
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
	_, err := s.redeemablesRepo.RedeemRedeemable(ctx, redeemableID, userID)

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to redeem redeemable")
		return err
	}
	return nil
}

func (s *RedeemablesService) UpdateRedemption(ctx context.Context, redeemableID uuid.UUID, userID uuid.UUID) error {
	err := s.redeemablesRepo.UpdateRedemption(ctx, redeemableID, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update redemption")
		return err
	}
	return nil
}
