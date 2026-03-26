package ctxutils

import (
	"context"

	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
)

func GetUserFromCtx(ctx context.Context) *mw.UserContext {
	userCtx, ok := ctx.Value(mw.UserContextKey).(*mw.UserContext)

	if !ok {
		return nil
	}

	return userCtx
}
