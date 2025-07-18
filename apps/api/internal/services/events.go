package services

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrFailedToCreateEvent = errors.New("failed to create event")
)

type EventService struct {
	eventRepo *repository.EventRepository
	logger    zerolog.Logger
}

func NewEventService(eventRepo *repository.EventRepository, logger zerolog.Logger) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		logger:    logger.With().Str("service", "EventService").Str("component", "events").Logger(),
	}
}

func (s *EventService) CreateEvent(ctx context.Context, name string, applicationOpen time.Time, applicationClose time.Time, startTime time.Time, endTime time.Time) (*sqlc.Event, error) {
	params := sqlc.CreateEventParams{
		Name:             name,
		ApplicationOpen:  applicationOpen,
		ApplicationClose: applicationClose,
		StartTime:        startTime,
		EndTime:          endTime,
	}

	result, err := s.eventRepo.CreateEvent(ctx, params)
	if err != nil {
		s.logger.Err(err).Msg("An unkown error was caught!")
		return nil, ErrFailedToCreateEvent
	}

	return result, nil
}
