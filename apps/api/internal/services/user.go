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
	ErrUserNotFound       = errors.New("user not found")
	ErrFailedToUpdateUser = errors.New("failed to update user")
	ErrFailedToGetUser    = errors.New("failed to get user")
)

type UserService struct {
	userRepo *repository.UserRepository
	logger   zerolog.Logger
}

func NewUserService(userRepo *repository.UserRepository, logger zerolog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger.With().Str("service", "UserService").Str("component", "user").Logger(),
	}
}

// GetUser retrieves a user by their ID
func (s *UserService) GetUser(ctx context.Context, userId uuid.UUID) (*sqlc.AuthUser, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		if err == repository.ErrUserNotFound {
			s.logger.Err(err).Msg(repository.ErrUserNotFound.Error())
			return nil, ErrUserNotFound
		} else {
			s.logger.Err(err).Msg("failed to get user by ID")
			return nil, ErrFailedToGetUser
		}
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userId uuid.UUID, params sqlc.UpdateUserParams) error {
	// Set the user ID in the params
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
	return s.userRepo.GetAllUsers(ctx, search, limit, offset)
}
