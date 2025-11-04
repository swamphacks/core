package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type TeamMemberRepository struct {
	db *db.DB
}

func NewTeamMemberRespository(db *db.DB) *TeamMemberRepository {
	return &TeamMemberRepository{
		db: db,
	}
}

func (r *TeamMemberRepository) NewTx(tx pgx.Tx) *TeamMemberRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamMemberRepository{
		db: txDB,
	}
}

func (r *TeamMemberRepository) GetTeamMembers(ctx context.Context, teamId uuid.UUID) ([]sqlc.GetTeamMembersRow, error) {
	return r.db.Query.GetTeamMembers(ctx, teamId)
}
