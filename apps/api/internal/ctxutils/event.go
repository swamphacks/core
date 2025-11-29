package ctxutils

import (
	"context"

	mw "github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

// Takes in the request context and returns the eventRole or nil if not retrievable
func GetEventRoleFromCtx(ctx context.Context) *sqlc.EventRole {
	eventCtx, ok := ctx.Value(mw.EventContextKey).(*mw.EventContext)
	if !ok {
		return nil
	}

	return eventCtx.EventRole
}
