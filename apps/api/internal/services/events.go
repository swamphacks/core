package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrFailedToCreateEvent = errors.New("failed to create event")
	ErrFailedToGetEvent    = errors.New("failed to get event")
	ErrFailedToUpdateEvent = errors.New("failed to update event")
	ErrFailedToDeleteEvent = errors.New("failed to delete event")
	ErrFailedToParseUUID   = errors.New("failed to parse uuid")
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

func (s *EventService) GetAllEvents(ctx context.Context) (*[]sqlc.Event, error) {
	// Check role, if role is none or user, return published events only
	userCtx, ok := ctx.Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		s.logger.Warn().Msg("Couldn't get user context")
		return s.eventRepo.GetPublishedEvents(ctx)
	}

	//TODO: Replace with switch later
	if userCtx.Role == sqlc.AuthUserRoleUser {
		return s.eventRepo.GetPublishedEvents(ctx)
	}

	return s.eventRepo.GetAllEvents(ctx)

}

func (s *EventService) GetEventRoleByIds(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (*sqlc.EventRole, error) {
	eventRole, err := s.eventRepo.GetEventRoleByIds(ctx, userId, eventId)
	if err != nil {
		if err == repository.ErrEventRoleNotFound {
			s.logger.Err(err).Msg(repository.ErrEventRoleNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, err
	}

	return eventRole, err
}

func (s *EventService) GetEventStaffUsers(ctx context.Context, eventId uuid.UUID) (*[]sqlc.AuthUser, error) {
	return s.eventRepo.GetEventStaff(ctx, eventId)
}
