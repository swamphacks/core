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

func (s *EventService) CreateEvent(ctx context.Context, params sqlc.CreateEventParams) (*sqlc.Event, error) {
	event, err := s.eventRepo.CreateEvent(ctx, params)
	if err != nil {
		if err == repository.ErrEventNotFound {
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToCreateEvent
	}

	return event, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*sqlc.Event, error) {
	event, err := s.eventRepo.GetEventByID(ctx, id)
	if err != nil {
		if err == repository.ErrEventNotFound {
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToGetEvent
	}

	return event, nil
}

func (s *EventService) UpdateEventById(ctx context.Context, params sqlc.UpdateEventByIdParams) (*sqlc.Event, error) {
	err := s.eventRepo.UpdateEventById(ctx, params)
	if err != nil {
		if err == repository.ErrEventNotFound {
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToUpdateEvent
	}

	event, err := s.eventRepo.GetEventByID(ctx, params.ID)

	return event, err
}

func (s *EventService) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	err := s.eventRepo.DeleteEventById(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrEventNotFound:
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		case repository.ErrNoEventsDeleted:
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		case repository.ErrMultipleEventsDeleted:
			s.logger.Err(err).Msg(repository.ErrMultipleEventsDeleted.Error())
		default:
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return ErrFailedToDeleteEvent
	}

	return err
}
