package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type EventRolesRepository struct {
	db *database.DB
}

func NewEventRolesRepository(db *database.DB) *EventRolesRepository {
	return &EventRolesRepository{
		db: db,
	}
}

func (r *EventRolesRepository) NewTx(tx pgx.Tx) *EventRolesRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &EventRolesRepository{
		db: txDB,
	}
}

func (r *EventRolesRepository) AssignRole(ctx context.Context, params sqlc.AssignRoleParams) error {
	return r.db.Query.AssignRole(ctx, params)
}

func (r *EventRolesRepository) RemoveRole(ctx context.Context, userId uuid.UUID) error {
	return r.db.Query.RemoveRole(ctx, userId)
}

func (r *EventRolesRepository) GetRoleByUserId(ctx context.Context, userId uuid.UUID) (sqlc.EventRole, error) {
	return r.db.Query.GetRoleByUserId(ctx, userId)
}

// **Deprecated**. Use `UpdateEventRoleByIds` instead.
// Only kept for backwards compatibility. TODO: Refactor all current implementations to use new function.
func (r *EventRolesRepository) UpdateRole(ctx context.Context, userId uuid.UUID, role sqlc.EventRoleType) error {
	params := sqlc.UpdateRoleParams{
		UserID: userId,
		Role:   role,
	}
	return r.db.Query.UpdateRole(ctx, params)
}

func (r *EventRolesRepository) UpdateRoleByUserId(ctx context.Context, params sqlc.UpdateRoleByUserIdParams) error {
	return r.db.Query.UpdateRoleByUserId(ctx, params)
}
