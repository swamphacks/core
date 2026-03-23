package auth

import (
	"context"

	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type AuthService interface {
	GetMe(ctx context.Context) (*middleware.UserContext, error)
	AuthenticateWithOAuth(ctx context.Context, code, provider string, ipAddress, userAgent *string) (*sqlc.AuthSession, error)
	Logout(ctx context.Context) error
}
