package users

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
		logger:   logger.With().Str("service", "UserService").Str("domain", "user").Logger(),
	}
}

func (s *UserService) GetUserById(ctx context.Context, userId uuid.UUID) (*sqlc.User, error) {
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

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
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

func (s *UserService) GetUserByRFID(ctx context.Context, rfid string) (*sqlc.User, error) {
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

// func (s *UserService) GetCheckedInStatusByUserId(ctx context.Context, userId uuid.UUID) (bool, error) {
// 	checkedIn, err := s.userRepo.GetCheckedInStatusByUserId(ctx, userId)

// 	if err != nil {
// 		s.logger.Err(err).Msg("check in status fail")
// 		return false, errors.New("Failed to get check in status for user")
// 	}

// 	return checkedIn, nil
// }

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

func (s *UserService) GetAllUsers(ctx context.Context, search *string, limit, offset int32) ([]sqlc.User, error) {
	users, err := s.userRepo.GetAllUsers(ctx, search, limit, offset)

	if err != nil {
		s.logger.Err(err).Msg("get all users fail")
		return []sqlc.User{}, errors.New("Failed to get users")
	}

	return users, nil
}

// func (s *UserService) GetRole(ctx context.Context, userId uuid.UUID) (*sqlc.RoleType, error) {
// 	role, err := s.eventRolesRepo.GetRoleByUserId(ctx, userId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return role, nil
// }

func (s *UserService) AssignRole(ctx context.Context, userId *uuid.UUID, email *string, role sqlc.UserRole) error {
	if userId == nil && email == nil {
		return errors.New("must provide either userId or email")
	}

	var selectedUser *sqlc.User
	var err error

	if userId != nil {
		selectedUser, err = s.userRepo.GetUserByID(ctx, *userId)
		// Do not return if user not found, the query needs to fallback to other optiosn
		if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Err(err).Msg("Something went wrong getting by id")
			return err
		}
	}

	if selectedUser == nil && email != nil {
		selectedUser, err = s.userRepo.GetUserByEmail(ctx, *email)
		if err != nil {
			s.logger.Err(err).Msg("Something went wrong getting by email")
			return err
		}
	}

	// Just a double safety check (should usually be caught by queries above)
	if selectedUser == nil {
		s.logger.Warn().Msg(("User not found from email OR id"))
		return repository.ErrUserNotFound
	}

	// Now assign the event role
	err = s.userRepo.UpdateRole(ctx, sqlc.UpdateRoleParams{
		UserID: selectedUser.ID,
		Role:   role,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) RevokeRole(ctx context.Context, userId uuid.UUID) error {
	return s.userRepo.RemoveRole(ctx, userId)
}

func (s *UserService) UpdateRole(ctx context.Context, userId uuid.UUID, role sqlc.UserRole) error {
	return s.userRepo.UpdateRole(ctx, sqlc.UpdateRoleParams{
		UserID: userId,
		Role:   role,
	})
}

// func (s *UserService) UpdateRoleById(ctx context.Context, userId uuid.UUID, role *sqlc.RoleType, checkedInAt *time.Time, RFID *string) error {
// 	if role == nil && checkedInAt == nil && RFID == nil {
// 		return errors.New("no fields provided to update")
// 	}

// 	return s.eventRolesRepo.UpdateRoleByUserId(ctx, sqlc.UpdateRoleByUserIdParams{
// 		UserID:       userId,
// 		Role:         *role,
// 		RoleDoUpdate: role != nil,

// 		CheckedInAt:         checkedInAt,
// 		CheckedInAtDoUpdate: checkedInAt != nil,

// 		Rfid:         RFID,
// 		RfidDoUpdate: RFID != nil,
// 	})
// }
