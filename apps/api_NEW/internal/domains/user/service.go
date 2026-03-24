package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrFailedToUpdateUser = errors.New("failed to update user")
	ErrFailedToGetUser    = errors.New("failed to get user")
)

type UserService struct {
	userRepo *repository.UserRepository
	logger   zerolog.Logger
}

func NewService(userRepo *repository.UserRepository, logger zerolog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger.With().Str("service", "UserService").Str("component", "user").Logger(),
	}
}

func (s *UserService) GetUser(ctx context.Context, userId uuid.UUID) (*sqlc.AuthUser, error) {
	user, err := s.userRepo.GetUserByID(ctx, userId)

	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		} else {
			s.logger.Err(err).Msg("failed to get user by ID")
			return nil, ErrFailedToGetUser
		}
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*sqlc.AuthUser, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		} else {
			s.logger.Err(err).Msg("get user by email fail")
			return nil, ErrFailedToGetUser
		}
	}

	return user, nil
}

func (s *UserService) GetUserEmailInfoById(ctx context.Context, userId uuid.UUID) (*sqlc.GetUserEmailInfoByIdRow, error) {
	emailInfo, err := s.userRepo.GetUserEmailInfoById(ctx, userId)

	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		} else {
			s.logger.Err(err).Msg("get user by email fail")
			return nil, ErrFailedToGetUser
		}
	}

	return emailInfo, nil
}

func (s *UserService) GetUserByRFID(ctx context.Context, rfid string) (*sqlc.AuthUser, error) {
	user, err := s.userRepo.GetUserByRFID(ctx, rfid)

	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		} else {
			s.logger.Err(err).Msg("get user by rfid fail")
			return nil, ErrFailedToGetUser
		}
	}

	return user, nil
}

func (s *UserService) GetCheckedInStatusByUserId(ctx context.Context, userId uuid.UUID) (bool, error) {
	checkedIn, err := s.userRepo.GetCheckedInStatusByUserId(ctx, userId)

	if err != nil {
		s.logger.Err(err).Msg("check in status fail")
		return false, errors.New("Failed to get check in status for user")
	}

	return checkedIn, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userId uuid.UUID, params sqlc.UpdateUserParams) error {
	params.ID = userId

	err := s.userRepo.UpdateUser(ctx, params)
	if err != nil {
		if err == repository.ErrUserNotFound {
			s.logger.Err(err).Msg(repository.ErrUserNotFound.Error())
			return ErrUserNotFound
		} else {
			s.logger.Err(err).Msg("failed to update user")
			return ErrFailedToUpdateUser
		}
	}

	return nil
}

func (s *UserService) CompleteOnboarding(ctx context.Context, userId uuid.UUID, name, email string) error {
	params := sqlc.UpdateUserParams{
		ID:                     userId,
		NameDoUpdate:           true,
		Name:                   name,
		PreferredEmailDoUpdate: true,
		PreferredEmail:         &email,
		OnboardedDoUpdate:      true,
		Onboarded:              true,
	}

	return s.UpdateUser(ctx, userId, params)
}

func (s *UserService) GetAllUsers(ctx context.Context, search *string, limit, offset int32) ([]sqlc.AuthUser, error) {
	users, err := s.userRepo.GetAllUsers(ctx, search, limit, offset)

	if err != nil {
		s.logger.Err(err).Msg("get all users fail")
		return []sqlc.AuthUser{}, errors.New("Failed to get users")
	}

	return users, nil
}
