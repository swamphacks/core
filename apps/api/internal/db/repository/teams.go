package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/ptr"
)

var (
	ErrTeamNotFound = errors.New("team was not found")
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

func (r *TeamRepository) GetByID(ctx context.Context, teamId uuid.UUID) (*sqlc.Team, error) {
	team, err := r.db.Query.GetTeamById(ctx, teamId)
	if err != nil && db.IsNotFound(err) {
		return nil, ErrTeamNotFound
	}

	return &team, err
}

func (r *TeamRepository) GetTeamByMemberAndEvent(ctx context.Context, userId, eventId uuid.UUID) (*sqlc.GetUserEventTeamRow, error) {
	team, err := r.db.Query.GetUserEventTeam(ctx, sqlc.GetUserEventTeamParams{
		UserID:  userId,
		EventID: ptr.UUIDToPtr(eventId),
	})

	if err != nil && db.IsNotFound(err) {
		return nil, ErrTeamNotFound
	}

	return &team, err
}

func (r *TeamRepository) GetTeamsWithMembersByEvent(ctx context.Context, eventId uuid.UUID, limit, offset int32) ([]sqlc.ListTeamsWithMembersByEventRow, error) {
	return r.db.Query.ListTeamsWithMembersByEvent(ctx, sqlc.ListTeamsWithMembersByEventParams{
		EventID: ptr.UUIDToPtr(eventId),
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *TeamRepository) Delete(ctx context.Context, teamId uuid.UUID) error {
	return r.db.Query.DeleteTeam(ctx, teamId)
}

func (r *TeamRepository) Update(ctx context.Context, teamId uuid.UUID, name *string, ownerId *uuid.UUID) (*sqlc.Team, error) {
	params := sqlc.UpdateTeamByIdParams{
		ID:              teamId,
		OwnerIDDoUpdate: ownerId != nil && *ownerId != uuid.Nil,
		NameDoUpdate:    name != nil && *name != "",
	}

	if ownerId != nil {
		params.OwnerID = ownerId
	}
	if name != nil {
		params.Name = *name
	}

	team, err := r.db.Query.UpdateTeamById(ctx, params)
	return &team, err
}
