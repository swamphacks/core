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
	nfcRepo       *repository.NFCRepository
	logger        zerolog.Logger
	checkinLogger zerolog.Logger
}

func NewService(nfcRepo *repository.NFCRepository, logger zerolog.Logger, checkinLogger zerolog.Logger) *NfcService {
	return &NfcService{
		nfcRepo:       nfcRepo,
		logger:        logger,
		checkinLogger: checkinLogger,
	}
}

func (s *NfcService) GetMeals(ctx context.Context) ([]sqlc.GetMealsRow, error) {
	meals, err := s.nfcRepo.GetMeals(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get meals")
		return nil, errors.New("failed to get meals")
	}

	return meals, nil
}

func (s *NfcService) GetTshirts(ctx context.Context) ([]sqlc.GetTshirtsRow, error) {
	tshirts, err := s.nfcRepo.GetTshirts(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get tshirts")
		return nil, errors.New("failed to get tshirts")
	}

	return tshirts, nil
}

func (s *NfcService) GetRedeemablesByType(ctx context.Context, redeemableType string) ([]sqlc.GetRedeemablesByTypeRow, error) {
	redeemables, err := s.nfcRepo.GetRedeemablesByType(ctx, redeemableType)
	if err != nil {
		s.logger.Err(err).Str("type", redeemableType).Msg("Failed to get redeemables by type")
		return nil, errors.New("failed to get redeemables")
	}

	return redeemables, nil
}

func (s *NfcService) CheckinUser(ctx context.Context, tagID string, userID uuid.UUID) (*sqlc.CheckinUserRow, error) {
	user, err := s.nfcRepo.CheckinUser(ctx, sqlc.CheckinUserParams{
		TagID:  tagID,
		UserID: &userID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.logger.Err(err).Msg("Failed to check in user")
		}
		if user == nil {
			s.checkinLogger.Info().Str("type", "checkin").Msgf("❌ User can't check in - UserId: %s - Error: %s", userID, err.Error())
		} else {
			s.checkinLogger.Info().Str("type", "checkin").Msgf("❌ User can't check in - Name: %s - Status: %s - Error: %s", user.Name, user.ApplicationStatus, err.Error())
		}
		return nil, fmt.Errorf("check in user: %w", err)
	}

	s.checkinLogger.Info().Str("type", "checkin").Msgf("✅ User checked in - Name: %s - Status: %s", user.Name, user.ApplicationStatus)

	return user, nil
}

func (s *NfcService) TagToWorkshop(ctx context.Context, tagID string, workshopID uuid.UUID) (*sqlc.TagToWorkshopRow, error) {
	user, err := s.getUserIDFromTagId(ctx, tagID)

	if err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("Failed to tag to workshop, can't get user id, tagID: %s", tagID))
		return nil, fmt.Errorf("tag to workshop: %w", err)
	}

	if user.UserID == nil {
		s.logger.Err(err).Msg(fmt.Sprintf("userID is null for tag, %s", tagID))
		return nil, fmt.Errorf("tag to redeemable: %w", err)
	}

	workshop, err := s.nfcRepo.TagToWorkshop(ctx, sqlc.TagToWorkshopParams{
		TagID:      &tagID,
		UserID:     *user.UserID,
		WorkshopID: workshopID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.logger.Err(err).Msg("Failed to tag to workshop")
		}
		return nil, fmt.Errorf("tag to workshop: %w", err)
	}

	s.checkinLogger.Info().Str("type", "workshop").Msgf(
		"✅ User checked in successfully - Name: %s - Workshop: %s",
		workshop.Name,
		workshop.WorkshopName,
	)

	return workshop, nil
}

func (s *NfcService) TagToRedeemable(ctx context.Context, tagID string, redeemableID uuid.UUID) error {
	user, err := s.getUserIDFromTagId(ctx, tagID)

	if err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("Failed to tag to redeemable, can't get user id, tagID: %s", tagID))
		return fmt.Errorf("tag to redeemable: %w", err)
	}

	if user.UserID == nil {
		s.logger.Err(err).Msg(fmt.Sprintf("userID is null for tag, %s", tagID))
		return fmt.Errorf("tag to redeemable: %w", err)
	}

	redeemable, err := s.nfcRepo.TagToRedeemable(ctx, sqlc.TagToRedeemableParams{
		UserID:       *user.UserID,
		TagID:        &tagID,
		RedeemableID: redeemableID,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to tag to redeemable")
		return fmt.Errorf("tag to redeemable: %w", err)
	}

	s.checkinLogger.Info().Str("type", "redeemable").Msgf(
		"✅ User checked in successfully - Name: %s - Redeemable: %s - Items remaining: %d",
		redeemable.Name,
		redeemable.RedeemableName,
		redeemable.ItemsRemaining,
	)

	return nil
}

func (s *NfcService) GetWorkshops(ctx context.Context) ([]sqlc.GetWorkshopsRow, error) {
	workshops, err := s.nfcRepo.GetWorkshops(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get workshops")
		return nil, errors.New("failed to get workshops")
	}

	return workshops, nil
}

func (s *NfcService) GetSocials(ctx context.Context) ([]sqlc.GetSocialsRow, error) {
	socials, err := s.nfcRepo.GetSocials(ctx)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get socials")
		return nil, errors.New("failed to get socials")
	}

	return socials, nil
}

func (s *NfcService) UnlinkTag(ctx context.Context, userID uuid.UUID) (bool, error) {
	rowsAffected, err := s.nfcRepo.DeleteTagByUserID(ctx, userID)
	if err != nil {
		s.logger.Err(err).Str("user_id", userID.String()).Msg("Failed to unlink NFC tag")
		return false, errors.New("failed to unlink NFC tag")
	}

	return rowsAffected > 0, nil
}

func (s *NfcService) getUserIDFromTagId(ctx context.Context, tagID string) (*sqlc.NfcTagsUser, error) {
	user, err := s.nfcRepo.GetTag(ctx, tagID)

	if err != nil {
		s.logger.Err(err).Msg("Failed to get userid from tagid")
		return nil, errors.New("failed to get userid from tagid")
	}

	if user == nil {
		s.logger.Err(err).Msg("user is null from GetTag")
		return nil, errors.New("failed to get user from tagid")
	}

	return user, nil
}
