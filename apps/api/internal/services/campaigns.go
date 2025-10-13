package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrFailedToCreateCampaign = errors.New("failed to create campaign")
	ErrCampaignNotFound       = errors.New("campaign not found")
	ErrFailedToGetCampaign    = errors.New("failed to get campaign")
	ErrFailedToUpdateCampaign = errors.New("failed to update campaign")
	ErrFailedToDeleteCampaign = errors.New("failed to get campaign")
)

type CampaignService struct {
	campaignRepo *repository.CampaignRepository
	logger       zerolog.Logger
}

func NewCampaignService(campaignRepo *repository.CampaignRepository, logger zerolog.Logger) *CampaignService {
	return &CampaignService{
		campaignRepo: campaignRepo,
		logger:       logger.With().Str("service", "CampaignService").Str("component", "campaign").Logger(),
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, params sqlc.CreateCampaignParams) (*sqlc.Campaign, error) {
	event, err := s.campaignRepo.CreateCampaign(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrCampaignNotFound) {
			s.logger.Err(err).Msg(repository.ErrCampaignNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToCreateCampaign
	}

	return event, nil
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, campaignId uuid.UUID) (*sqlc.Campaign, error) {
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, campaignId)
	if err != nil {
		if err == repository.ErrCampaignNotFound {
			s.logger.Err(err).Msg(repository.ErrCampaignNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToGetCampaign
	}

	return campaign, nil
}

func (s *CampaignService) UpdateCampaignById(ctx context.Context, params sqlc.UpdateCampaignByIdParams) (*sqlc.Campaign, error) {
	err := s.campaignRepo.UpdateCampaignById(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrCampaignNotFound) {
			s.logger.Err(err).Msg(repository.ErrCampaignNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToUpdateCampaign
	}

	event, err := s.campaignRepo.GetCampaignByID(ctx, params.ID)

	return event, err
}

func (s *CampaignService) DeleteCampaignById(ctx context.Context, id uuid.UUID) error {
	err := s.campaignRepo.DeleteCampaignById(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrCampaignNotFound:
			s.logger.Err(err).Msg(repository.ErrCampaignNotFound.Error())
		case repository.ErrNoCampaignsDeleted:
			s.logger.Err(err).Msg(repository.ErrCampaignNotFound.Error())
		case repository.ErrMultipleCampaignsDeleted:
			s.logger.Err(err).Msg(repository.ErrMultipleCampaignsDeleted.Error())
		default:
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return ErrFailedToDeleteCampaign
	}

	return err
}
