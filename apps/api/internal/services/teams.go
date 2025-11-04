package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type TeamService struct {
	teamRepo       *repository.TeamRepository
	teamMemberRepo *repository.TeamMemberRepository
	logger         zerolog.Logger
}

func NewTeamService(teamRepo *repository.TeamRepository, teamMemberRepo *repository.TeamMemberRepository, logger zerolog.Logger) *TeamService {
	return &TeamService{
		teamRepo:       teamRepo,
		teamMemberRepo: teamMemberRepo,
		logger:         logger.With().Str("service", "TeamService").Str("component", "team").Logger(),
	}
}

type TeamWithMembers struct {
	ID      uuid.UUID                `json:"id"`
	EventId *uuid.UUID               `json:"event_id"`
	OwnerId *uuid.UUID               `json:"owner_id"`
	Name    string                   `json:"name"`
	Members []sqlc.GetTeamMembersRow `json:"members"`
}

func (s *TeamService) GetUserTeamWithMembers(ctx context.Context, userId, eventId uuid.UUID) (*TeamWithMembers, error) {
	team, err := s.teamRepo.GetTeamByMemberAndEvent(ctx, userId, eventId)
	if err != nil {
		// If no team, just return nil
		if errors.Is(err, repository.ErrTeamNotFound) {
			return nil, nil
		}
		return nil, err
	}

	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	teamWithMembers := TeamWithMembers{
		ID:      team.ID,
		EventId: team.EventID,
		OwnerId: team.OwnerID,
		Name:    team.Name,
		Members: members,
	}

	return &teamWithMembers, nil
}

func (s *TeamService) GetTeamWithMembers(ctx context.Context, teamId uuid.UUID) (*TeamWithMembers, error) {
	team, err := s.teamRepo.GetByID(ctx, teamId)
	if err != nil {
		// If no team, just return nil
		if errors.Is(err, repository.ErrTeamNotFound) {
			return nil, nil
		}
		return nil, err
	}

	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	teamWithMembers := TeamWithMembers{
		ID:      team.ID,
		EventId: team.EventID,
		OwnerId: team.OwnerID,
		Name:    team.Name,
		Members: members,
	}

	return &teamWithMembers, nil
}
