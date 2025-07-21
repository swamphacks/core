package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrFailedToCreateEvent = errors.New("failed to create event")
	ErrFailedToGetEvent    = errors.New("failed to get event")
	ErrFailedToUpdateEvent = errors.New("failed to update event")
	ErrFailedToDeleteEvent = errors.New("failed to delete event")
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

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*sqlc.Event, error) {
	result, err := s.eventRepo.GetEventByID(ctx, id)
	if err != nil {
		s.logger.Err(err).Msg("An unkown error was caught!")
		return nil, ErrFailedToGetEvent
	}

	return result, nil
}

func (s *EventService) UpdateEventById(ctx context.Context, params sqlc.UpdateEventByIdParams) error {
	err := s.eventRepo.UpdateEventById(ctx, params)
	if err != nil {
		s.logger.Err(err).Msg("An unkown error was caught!")
		return ErrFailedToUpdateEvent
	}

	return nil
}

func (s *EventService) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	err := s.eventRepo.DeleteEventById(ctx, id)
	if err != nil {
		s.logger.Err(err).Msg("An unkown error was caught!")
		return ErrFailedToUpdateEvent
	}

	return nil
}
