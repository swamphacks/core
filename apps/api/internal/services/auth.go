package services

import (
	"context"

	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) AuthenticateWithOAuth(ctx context.Context, code, provider string) (string, error) {
	// Authenticate logic here

	return "token", nil
}

func (s *AuthService) GetMe(ctx context.Context) (*sqlc.AuthUser, error) {
	// Pull the userId from the context
	return nil, nil
}
