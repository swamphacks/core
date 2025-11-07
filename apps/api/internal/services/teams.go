package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrTeamExists          = errors.New("team already exists")
	ErrTeamNotFound        = errors.New("team does not exist")
	ErrNoEligibleNextOwner = errors.New("no eligible next owner for team")
)

type TeamService struct {
	teamRepo       *repository.TeamRepository
	teamMemberRepo *repository.TeamMemberRepository
	txm            *db.TransactionManager
	logger         zerolog.Logger
}

func NewTeamService(teamRepo *repository.TeamRepository, teamMemberRepo *repository.TeamMemberRepository, txm *db.TransactionManager, logger zerolog.Logger) *TeamService {
	return &TeamService{
		teamRepo:       teamRepo,
		teamMemberRepo: teamMemberRepo,
		txm:            txm,
		logger:         logger.With().Str("service", "TeamService").Str("component", "team").Logger(),
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string, eventId, userId uuid.UUID) (*sqlc.Team, error) {
	// Check if user already has a team for this event.
	member, err := s.teamMemberRepo.GetTeamMemberByUserAndEvent(ctx, userId, eventId)
	if err == nil && member != nil {
		// User already has a team
		return nil, ErrTeamExists
	}
	if err != nil && !errors.Is(err, repository.ErrTeamMemberNotFound) {
		return nil, err
	}

	var newTeam sqlc.Team

	// Transactionally create a new team and assign the user as the owner.
	if err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txTeamRepo := s.teamRepo.NewTx(tx)
		txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)

		team, err := txTeamRepo.Create(ctx, name, userId, eventId)
		if err != nil {
			return err
		}

		if _, err = txTeamMemberRepo.Create(ctx, team.ID, userId); err != nil {
			return err
		}

		newTeam = *team
		return nil
	}); err != nil {
		return nil, err
	}

	return &newTeam, nil
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

func (s *TeamService) JoinTeam(ctx context.Context, userId, teamId uuid.UUID) error {
	_, err := s.teamMemberRepo.Create(ctx, teamId, userId)
	return err
}

func (s *TeamService) LeaveTeam(ctx context.Context, userId, teamId uuid.UUID) error {
	team, err := s.teamRepo.GetByID(ctx, teamId)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			return ErrTeamNotFound
		}
		return err
	}

	// Check if user is NOT the owner
	if team.OwnerID == nil || *team.OwnerID != userId {
		return s.teamMemberRepo.Delete(ctx, team.ID, userId)
	}

	// User IS the owner
	members, err := s.teamMemberRepo.GetTeamMembers(ctx, team.ID)
	if err != nil {
		return err
	}

	// Confirm that team has members. This should not run unless there are inconsistencies.
	if len(members) == 0 {
		s.logger.Warn().Msg("Team has no members but owner exists — inconsistent state.")
		return s.teamRepo.Delete(ctx, team.ID)
	}

	// If current user is the last member → delete both team + membership
	if len(members) == 1 {
		return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
			txTeamRepo := s.teamRepo.NewTx(tx)
			txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)

			if err := txTeamMemberRepo.Delete(ctx, team.ID, userId); err != nil {
				return err
			}
			return txTeamRepo.Delete(ctx, team.ID)
		})
	}

	// Choose the next owner deterministically
	var nextOwner sqlc.GetTeamMembersRow
	for _, m := range members {
		if m.UserID != userId {
			nextOwner = m
			break
		}
	}

	if nextOwner.UserID == uuid.Nil {
		// Should not happen unless data is inconsistent
		s.logger.Error().Msg("No eligible next owner found.")
		return ErrNoEligibleNextOwner
	}

	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txTeamRepo := s.teamRepo.NewTx(tx)
		txTeamMemberRepo := s.teamMemberRepo.NewTx(tx)

		if err := txTeamMemberRepo.Delete(ctx, team.ID, userId); err != nil {
			return err
		}

		_, err := txTeamRepo.Update(ctx, team.ID, nil, &nextOwner.UserID)
		return err
	})
}
