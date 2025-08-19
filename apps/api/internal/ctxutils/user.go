package ctxutils

import (
	"context"

	"github.com/google/uuid"
	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

// Takes in a context and returns whether the user has the superuser role
func IsSuperuser(ctx context.Context) bool {
	userCtx, ok := ctx.Value(mw.UserContextKey).(*mw.UserContext)
	if !ok {
		return false
	}

	return userCtx.Role == sqlc.AuthUserRoleSuperuser
}

// Takes in a context and returns whether the user has the user role
func IsUser(ctx context.Context) bool {
	userCtx, ok := ctx.Value(mw.UserContextKey).(*mw.UserContext)
	if !ok {
		return false
	}

	return userCtx.Role == sqlc.AuthUserRoleUser
}

// Takes in the request context and returns the userID or nil if not retrievable
func GetUserIdFromCtx(ctx context.Context) *uuid.UUID {
	userCtx, ok := ctx.Value(mw.UserContextKey).(*mw.UserContext)
	if !ok {
		return nil
	}

	return &userCtx.UserID
}
