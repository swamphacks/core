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

// Should we have the helper functions,
// UpdateUser updates user information with conditional field updates
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

// UpdateUserOnboarded sets the user's onboarded status to true
func (s *UserService) UpdateUserOnboarded(ctx context.Context, userId uuid.UUID) error {
	err := s.userRepo.UpdateUserOnboarded(ctx, userId)
	if err != nil {
		if err == repository.ErrUserNotFound {
			s.logger.Err(err).Msg(repository.ErrUserNotFound.Error())
			return ErrUserNotFound
		} else {
			s.logger.Err(err).Msg("failed to update user onboarded status")
			return ErrFailedToUpdateUser
		}
	}

	return nil
}

// UpdateUserProfile is a convenience method for updating name and email, calls above update function
func (s *UserService) UpdateUserProfile(ctx context.Context, userId uuid.UUID, name *string, email *string) error {
	params := sqlc.UpdateUserParams{
		ID: userId,
	}

	// Set update flags and values based on what's provided
	if name != nil {
		params.NameDoUpdate = true
		params.Name = *name
	}

	if email != nil {
		params.EmailDoUpdate = true
		params.Email = email
	}

	return s.UpdateUser(ctx, userId, params)
}
