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
    ErrNoEventRole = errors.New("user has no event role")
)

type DiscordService struct {
	eventRepo *repository.EventRepository
	logger zerolog.Logger
}

func NewDiscordService(
	eventRepo *repository.EventRepository,
	logger zerolog.Logger,
) *DiscordService {
	return &DiscordService{
		eventRepo: eventRepo,
		logger: logger.With().Str("service", "DiscordService").Str("component", "discord").Logger(),
	}
}

func (s *DiscordService) GetUserEventRoleByDiscordIDAndEventId(ctx context.Context, discordID string, eventID uuid.UUID) (*sqlc.EventRoleType, error) {
	eventRole, err := s.eventRepo.GetEventRoleByDiscordIDAndEventId(ctx, discordID, eventID)
	if err != nil {
		if err == repository.ErrEventRoleNotFound {
			return nil, ErrNoEventRole
		}
		s.logger.Err(err).Msg("failed to get event role by discord ID and event ID")
		return nil, err
	}

	return &eventRole.Role, nil
}

func (s *DiscordService) GetEventAttendeesWithDiscord(ctx context.Context, eventID uuid.UUID) (*[]sqlc.GetEventAttendeesWithDiscordRow, error) {
	attendees, err := s.eventRepo.GetEventAttendeesWithDiscord(ctx, eventID)
	if err != nil {
		s.logger.Err(err).Msg("failed to get event attendees with discord")
		return nil, err
	}
	return attendees, nil
}