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
	ErrEmailConflict            = errors.New("email already exists in this mailing list")
	ErrFailedToCreateSubmission = errors.New("failed to create event interest submission")
)

type EventInterestService struct {
	eventInterestRepo *repository.EventInterestRepository
	logger            zerolog.Logger
}

func NewEventInterestService(eventInterestRepo *repository.EventInterestRepository, logger zerolog.Logger) *EventInterestService {
	return &EventInterestService{
		eventInterestRepo: eventInterestRepo,
		logger:            logger.With().Str("service", "EventInterestService").Str("component", "event_interest").Logger(),
	}
}

func (s *EventInterestService) CreateInterestSubmission(ctx context.Context, eventID uuid.UUID, email string, source *string) (*sqlc.EventInterestSubmission, error) {
	params := sqlc.AddEmailParams{
		EventID: eventID,
		Email:   email,
		Source:  source,
	}

	result, err := s.eventInterestRepo.AddEmail(ctx, params)
	if err != nil && err == repository.ErrDuplicateEmails {
		s.logger.Err(err).Msg("Could not insert email due to duplicate existing.")
		return nil, ErrEmailConflict
	} else if err != nil {
		s.logger.Err(err).Msg("An unknown error was caught!")
		return nil, ErrFailedToCreateSubmission
	}

	return result, nil
}
