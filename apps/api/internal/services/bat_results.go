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
	ErrResultConflict       = errors.New("result already exists for this event")
	ErrFailedToCreateResult = errors.New("failed to add result")
	ErrFailedToDeleteResult = errors.New("failed to delete result")
)

type BatResultService struct {
	batResultRepo *repository.BatResultsRepository
	logger        zerolog.Logger
}

func NewBatResultService(batResultRepo *repository.BatResultsRepository, logger zerolog.Logger) *BatResultService {
	return &BatResultService{
		batResultRepo: batResultRepo,
		logger:        logger.With().Str("service", "BatResultService").Str("component", "bat_result_service").Logger(),
	}
}

func (s *BatResultService) GetResultsByEventId(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetResultsByEventIdRow, error) {
	return s.batResultRepo.GetResultsByEvent(ctx, eventId)
}

func (s *BatResultService) DeleteResultById(ctx context.Context, id uuid.UUID) error {
	err := s.batResultRepo.DeleteEventById(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrResultNotFound:
			s.logger.Err(err).Msg(repository.ErrResultNotFound.Error())
		case repository.ErrNoResultsDeleted:
			s.logger.Err(err).Msg(repository.ErrNoResultsDeleted.Error())
		case repository.ErrMultipleResultsDeleted:
			s.logger.Err(err).Msg(repository.ErrMultipleResultsDeleted.Error())
		default:
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return ErrFailedToDeleteResult
	}

	return err
}
