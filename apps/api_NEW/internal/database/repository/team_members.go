package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrTeamMemberNotFound = errors.New("team member not found")
)

type TeamMemberRepository struct {
	db *database.DB
}

func NewTeamMemberRespository(db *database.DB) *TeamMemberRepository {
	return &TeamMemberRepository{
		db: db,
	}
}

func (r *TeamMemberRepository) NewTx(tx pgx.Tx) *TeamMemberRepository {
	txDB := &database.DB{
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

func (r *TeamMemberRepository) GetTeamMemberByUser(ctx context.Context, userId uuid.UUID) (*sqlc.TeamMember, error) {
	member, err := r.db.Query.GetTeamMemberByUserId(ctx, userId)
	if err != nil && database.IsNotFound(err) {
		return nil, ErrTeamMemberNotFound
	}
	return &member, err
}

func (r *TeamMemberRepository) Create(ctx context.Context, teamId, userId uuid.UUID) (*sqlc.TeamMember, error) {
	member, err := r.db.Query.CreateTeamMember(ctx, sqlc.CreateTeamMemberParams{
		TeamID: teamId,
		UserID: userId,
	})

	return &member, err
}

func (r *TeamMemberRepository) Delete(ctx context.Context, teamId, userId uuid.UUID) error {
	return r.db.Query.RemoveTeamMember(ctx, sqlc.RemoveTeamMemberParams{
		UserID: userId,
		TeamID: teamId,
	})
}
