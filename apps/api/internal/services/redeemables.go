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
	params := sqlc.CreateRedeemableParams{
		EventID:       eventID,
		Name:          name,
		Amount:        int32(amount),
		MaxUserAmount: int32(maxUserAmount),
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
	// need to convert to int32, extra code because they are nullable pointers
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
	params := sqlc.UpdateRedeemableParams{
		ID:            redeemableID,
		Name:          name,
		Amount:        amount32,
		MaxUserAmount: maxUserAmount32,
	}
	redeemable, err := s.redeemablesRepo.UpdateRedeemable(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update redeemable")
		return nil, err
	}
	return redeemable, nil
}
