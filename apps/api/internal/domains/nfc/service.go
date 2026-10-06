package nfc

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type NfcService struct {
	nfcRepo *repository.NFCRepository
	logger  zerolog.Logger
}

func NewService(nfcRepo *repository.NFCRepository, logger zerolog.Logger) *NfcService {
	return &NfcService{
		nfcRepo: nfcRepo,
		logger:  logger,
	}
}

func (s *NfcService) GetMeals(ctx context.Context) ([]sqlc.GetMealsRow, error) {
	meals, err := s.nfcRepo.GetMeals(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get meals")
		return nil, errors.New("Failed to get meals")
	}

	return meals, nil
}

func (s *NfcService) GetTshirts(ctx context.Context) ([]sqlc.GetTshirtsRow, error) {
	tshirts, err := s.nfcRepo.GetTshirts(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get tshirts")
		return nil, errors.New("Failed to get tshirts")
	}

	return tshirts, nil
}

func (s *NfcService) CheckinUser(ctx context.Context, tagID string, userID uuid.UUID) (*sqlc.NfcTagsUser, error) {
	user, err := s.nfcRepo.CheckinUser(ctx, sqlc.CheckinUserParams{
		TagID:  tagID,
		UserID: &userID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.logger.Err(err).Msg("Failed to check in user")
		}
		return nil, fmt.Errorf("check in user: %w", err)
	}
	return user, nil
}

func (s *NfcService) TagToWorkshop(ctx context.Context, tagID string, workshopID uuid.UUID) (int64, error) {
	rows, err := s.nfcRepo.TagToWorkshop(ctx, sqlc.TagToWorkshopParams{
		TagID:      tagID,
		WorkshopID: workshopID,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to tag to workshop")
		return 0, fmt.Errorf("tag to workshop: %w", err)
	}
	return rows, nil
}

func (s *NfcService) TagToRedeemable(ctx context.Context, tagID string, redeemableID uuid.UUID) error {
	err := s.nfcRepo.TagToRedeemable(ctx, sqlc.TagToRedeemableParams{
		TagID:        tagID,
		RedeemableID: redeemableID,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to tag to redeemable")
		return fmt.Errorf("tag to redeemable: %w", err)
	}
	return nil
}

func (s *NfcService) GetWorkshops(ctx context.Context) ([]sqlc.GetWorkshopsRow, error) {
	workshops, err := s.nfcRepo.GetWorkshops(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get workshops")
		return nil, errors.New("Failed to get workshops")
	}

	return workshops, nil
}

func (s *NfcService) GetSocials(ctx context.Context) ([]sqlc.GetSocialsRow, error) {
	socials, err := s.nfcRepo.GetSocials(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get socials")
		return nil, errors.New("Failed to get socials")
	}

	return socials, nil
}