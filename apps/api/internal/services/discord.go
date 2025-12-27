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
    ErrDiscordAccountNotFound = errors.New("discord account not found")
    ErrNoEventRole = errors.New("user has no event role")
)

type DiscordService struct {
	accountRepo *repository.AccountRepository
	eventRepo *repository.EventRepository
	logger zerolog.Logger
}

func NewDiscordService(
	accountRepo *repository.AccountRepository,
	eventRepo *repository.EventRepository,
	logger zerolog.Logger,
) *DiscordService {
	return &DiscordService{
		accountRepo: accountRepo,
		eventRepo: eventRepo,
		logger: logger.With().Str("service", "DiscordService").Str("component", "discord").Logger(),
	}
}

func (s *DiscordService) GetUserEventRoleByDiscordID(ctx context.Context, discordID string) (*uuid.UUID, *sqlc.EventRoleType, error) {
	userID, err := s.accountRepo.GetUserIDByDiscordAccountID(ctx, discordID)
	if err != nil {
		if err == repository.ErrAccountNotFound {
			return nil, nil, ErrDiscordAccountNotFound
		}
		s.logger.Err(err).Msg("failed to get user ID from discord account")
		return nil, nil, err
	}

	eventRole, err := s.eventRepo.GetEventRoleByUserID(ctx, *userID)
	if err != nil {
		if err == repository.ErrEventRoleNotFound {
			return nil, nil, ErrNoEventRole
		}
		s.logger.Err(err).Msg("failed to get event role")
		return nil, nil, err
	}

	return &eventRole.EventID, &eventRole.Role, nil
}

func (s *DiscordService) GetEventAttendeesWithDiscord(ctx context.Context, eventID uuid.UUID) (*[]sqlc.GetEventAttendeesWithDiscordRow, error) {
	attendees, err := s.eventRepo.GetEventAttendeesWithDiscord(ctx, eventID)
	if err != nil {
		s.logger.Err(err).Msg("failed to get event attendees with discord")
		return nil, err
	}
	return attendees, nil
}