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
	ErrTeamNotFound = errors.New("team was not found")
)

type TeamRepository struct {
	db *database.DB
}

func NewTeamRespository(db *database.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) NewTx(tx pgx.Tx) *TeamRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamRepository{
		db: txDB,
	}
}

func (r *TeamRepository) Create(ctx context.Context, arg sqlc.CreateTeamParams) (*sqlc.Team, error) {
	team, err := r.db.Query.CreateTeam(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &team, err
}

func (r *TeamRepository) GetByID(ctx context.Context, teamID uuid.UUID) (*sqlc.Team, error) {
	team, err := r.db.Query.GetTeamById(ctx, teamID)
	if err != nil && database.IsNotFound(err) {
		return nil, ErrTeamNotFound
	}

	return &team, err
}

func (r *TeamRepository) GetTeamByMember(ctx context.Context, userID uuid.UUID) (*sqlc.GetUserTeamRow, error) {
	team, err := r.db.Query.GetUserTeam(ctx, userID)

	if err != nil && database.IsNotFound(err) {
		return nil, ErrTeamNotFound
	}

	return &team, err
}

func (r *TeamRepository) GetTeamsWithMembers(ctx context.Context, params sqlc.ListTeamsWithMembersParams) ([]sqlc.ListTeamsWithMembersRow, error) {
	return r.db.Query.ListTeamsWithMembers(ctx, params)
}

func (r *TeamRepository) Delete(ctx context.Context, teamID uuid.UUID) error {
	return r.db.Query.DeleteTeam(ctx, teamID)
}

func (r *TeamRepository) Update(ctx context.Context, params sqlc.UpdateTeamByIdParams) (*sqlc.Team, error) {
	// params := sqlc.UpdateTeamByIdParams{
	// 	ID:              teamId,
	// 	OwnerIDDoUpdate: ownerId != nil && *ownerId != uuid.Nil,
	// 	NameDoUpdate:    name != nil && *name != "",
	// }

	// if ownerId != nil {
	// 	params.OwnerID = ownerId
	// }
	// if name != nil {
	// 	params.Name = *name
	// }

	team, err := r.db.Query.UpdateTeamById(ctx, params)
	return &team, err
}
