package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/ptr"
)

type TeamRepository struct {
	db *db.DB
}

func NewTeamRespository(db *db.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) NewTx(tx pgx.Tx) *TeamRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamRepository{
		db: txDB,
	}
}

func (r *TeamRepository) Create(ctx context.Context, name string, owner_id, event_id uuid.UUID) (*sqlc.Team, error) {
	team, err := r.db.Query.CreateTeam(ctx, sqlc.CreateTeamParams{
		Name:    name,
		OwnerID: ptr.UUIDToPtr(owner_id),
		EventID: ptr.UUIDToPtr(event_id),
	})
	if err != nil {
		return nil, err
	}

	return &team, err
}
