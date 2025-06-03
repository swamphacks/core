package services

import (
	"context"
	"errors"

	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrProviderUnsupported = errors.New("this provider is unsupported")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) AuthenticateWithOAuth(ctx context.Context, code, provider string) (*sqlc.AuthSession, error) {
	// Authenticate logic here
	switch provider {
	case "discord":
		// Do discord auth here
		break
	default:
		return nil, ErrProviderUnsupported
	}

	return nil, nil
}

func (s *AuthService) GetMe(ctx context.Context) (*sqlc.AuthUser, error) {
	// Pull the userId from the context
	return nil, nil
}
