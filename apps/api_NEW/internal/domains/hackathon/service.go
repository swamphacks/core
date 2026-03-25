package hackathon

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type HackathonService struct {
	hackathonRepo  *repository.HackathonRepository
	eventRolesRepo *repository.EventRolesRepository
	logger         zerolog.Logger
}

func NewService(hackathonRepo *repository.HackathonRepository, eventRolesRepo *repository.EventRolesRepository, logger zerolog.Logger) *HackathonService {
	return &HackathonService{
		hackathonRepo:  hackathonRepo,
		eventRolesRepo: eventRolesRepo,
		logger:         logger.With().Str("service", "HackathonService").Str("domain", "hackathon").Logger(),
	}
}

func (s *HackathonService) CreateHackathon(ctx context.Context, params sqlc.CreateHackathonParams) (*sqlc.Hackathon, error) {
	hackathon, err := s.hackathonRepo.CreateHackathon(ctx, params)

	if err != nil {
		return nil, errors.New("Failed to create hackathon")
	}

	return hackathon, nil
}

func (s *HackathonService) GetHackathon(ctx context.Context) (*sqlc.Hackathon, error) {
	hackathon, err := s.hackathonRepo.GetHackathon(ctx)

	if err != nil {
		if errors.Is(err, repository.ErrEntityNotFound) {
			return nil, repository.ErrEntityNotFound
		}
		s.logger.Err(err).Msg("")
		return nil, errors.New("Failed to get hackathon")
	}

	return hackathon, nil
}

func (s *HackathonService) UpdateHackathon(ctx context.Context, params sqlc.UpdateHackathonParams) error {
	err := s.hackathonRepo.UpdateHackathon(ctx, params)

	if err != nil {
		s.logger.Err(err).Msg("")
		return errors.New("Failed to update hackathon")
	}

	return nil
}

func (s *HackathonService) GetStaff(ctx context.Context) ([]sqlc.GetStaffRow, error) {
	staff, err := s.hackathonRepo.GetStaff(ctx)

	if err != nil {
		return nil, errors.New("Failed to get staff")
	}

	if staff == nil {
		return []sqlc.GetStaffRow{}, nil
	}

	return *staff, nil
}

func (s *HackathonService) GetAttendeesWithDiscord(ctx context.Context) ([]sqlc.GetAttendeesWithDiscordRow, error) {
	attendees, err := s.hackathonRepo.GetAttendeesWithDiscord(ctx)

	if err != nil {
		return nil, errors.New("Failed to get attendees with Discord")
	}

	if attendees == nil {
		return []sqlc.GetAttendeesWithDiscordRow{}, nil
	}

	return *attendees, nil
}

func (s *HackathonService) GetAttendeeUserIds(ctx context.Context) ([]uuid.UUID, error) {
	userIds, err := s.hackathonRepo.GetAttendeeUserIds(ctx)

	if err != nil {
		return nil, errors.New("Failed to get attendee user ids")
	}

	return userIds, nil
}

func (s *HackathonService) GetAttendeeCount(ctx context.Context) (int64, error) {
	count, err := s.hackathonRepo.GetAttendeeCount(ctx)

	if err != nil {
		return -1, errors.New("Failed to get attendee count")
	}

	return count, nil
}

var (
	ErrRolesNotFound   = errors.New("roles not found")
	ErrUserNotAttendee = errors.New("user is not an attendee")
	ErrUserCheckedIn   = errors.New("user already checked in")
)

func (s *HackathonService) CheckInAttendee(ctx context.Context, userId uuid.UUID, RFID *string) error {
	// Retrieve user with their current event role
	role, err := s.eventRolesRepo.GetRoleByUserId(ctx, userId)
	if err != nil {
		return ErrRolesNotFound
	}

	if role.Role != sqlc.EventRoleTypeAttendee {
		return ErrUserNotAttendee
	}

	if role.CheckedInAt != nil {
		return ErrUserCheckedIn
	}

	now := time.Now()
	// Update user role checked in AND rfid
	return s.eventRolesRepo.UpdateRoleByUserId(ctx, sqlc.UpdateRoleByUserIdParams{
		UserID: userId,

		Role:         sqlc.EventRoleTypeAttendee,
		RoleDoUpdate: false,

		CheckedInAt:         &now,
		CheckedInAtDoUpdate: true,

		Rfid:         RFID,
		RfidDoUpdate: RFID != nil,
	})
}
