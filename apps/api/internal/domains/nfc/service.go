package nfc

import (
	"context"
	"errors"

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
func (s *NfcService) GetMeals(ctx context.Context) ([]sqlc.GetMealsRow, error) {
	meals, err := s.nfcRepo.GetMeals(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get meals")
		return nil, errors.New("Failed to get meals")
	}

	return meals, nil
}

func (s *NfcService) GetTshirts(ctx context.Context) ([]sqlc.GetTshirtsRow, error) {
func (s *NfcService) GetTshirts(ctx context.Context) ([]sqlc.GetTshirtsRow, error) {
	tshirts, err := s.nfcRepo.GetTshirts(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get tshirts")
		return nil, errors.New("Failed to get tshirts")
	}

	return tshirts, nil
}

func (s *NfcService) CheckinUser(ctx context.Context, TagId string, UserID uuid.UUID) (*sqlc.NfcTagsUser, error) {
	params := sqlc.CheckinUserParams{
		TagID: TagId,
		UserID: &UserID,
	}
	user, err := s.nfcRepo.CheckinUser(ctx, params)
	if err != nil {
		s.logger.Err(err).Msg("Failed to check in user")
		return nil, errors.New("Failed to check in user")
	}

	return user, nil
}

func (s *NfcService) TagToWorkshop(ctx context.Context, TagId string, WorkshopId uuid.UUID) (int64, error) {
	params := sqlc.TagToWorkshopParams{
		TagID: TagId,
		WorkshopID: WorkshopId,
	}
	workshopID, err := s.nfcRepo.TagToWorkshop(ctx, params)
	if err != nil {
		s.logger.Err(err).Msg("Failed to tag to workshop")
		return 0, errors.New("Failed to tag to workshop")
	}

	return workshopID, nil
}

func (s *NfcService) TagToRedeemable(ctx context.Context, TagId string, RedeemableId uuid.UUID) error {
	params := sqlc.TagToRedeemableParams{
		TagID: TagId,
		RedeemableID: RedeemableId,
	}
	err := s.nfcRepo.TagToRedeemable(ctx, params)
	if err != nil {
		s.logger.Err(err).Msg("Failed to tag to redeemable")
		return errors.New("Failed to tag to redeemable")
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

func (s *NfcService) GetWorkshops(ctx context.Context) ([]sqlc.GetWorkshopsRow, error) {
	workshops, err := s.nfcRepo.GetWorkshops(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get workshops")
		return nil, errors.New("Failed to get workshops")
	}

	return workshops, nil
}