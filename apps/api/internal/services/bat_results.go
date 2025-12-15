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
	ErrFailedToCreateResult = errors.New("failed to result")
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

func (s *BatResultService) GetResultsByEventId(ctx context.Context, eventId uuid.UUID) (*[]sqlc.BatResult, error) {
	return s.batResultRepo.GetResultsByEvent(ctx, eventId)
}
