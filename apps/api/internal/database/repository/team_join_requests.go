package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type TeamJoinRequestRepository struct {
	db *database.DB
}

func NewTeamJoinRequestRepository(db *database.DB) *TeamJoinRequestRepository {
	return &TeamJoinRequestRepository{
		db: db,
	}
}

func (r *TeamJoinRequestRepository) NewTx(tx pgx.Tx) *TeamJoinRequestRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamJoinRequestRepository{
		db: txDB,
	}
}

func (r *TeamJoinRequestRepository) Create(ctx context.Context, params sqlc.CreateTeamJoinRequestParams) (*sqlc.TeamJoinRequest, error) {
	request, err := r.db.Query.CreateTeamJoinRequest(ctx, params)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *TeamJoinRequestRepository) GetById(ctx context.Context, requestID uuid.UUID) (*sqlc.TeamJoinRequest, error) {
	request, err := r.db.Query.GetTeamJoinRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByUser(ctx context.Context, userID uuid.UUID) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByUserID(ctx, userID)
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByTeam(ctx context.Context, params sqlc.ListTeamJoinRequestsByTeamIDAndStatusParams) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByTeamIDAndStatus(ctx, params)
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByTeamWithUser(ctx context.Context, params sqlc.ListJoinRequestsByTeamAndStatusWithUserParams) ([]sqlc.ListJoinRequestsByTeamAndStatusWithUserRow, error) {
	return r.db.Query.ListJoinRequestsByTeamAndStatusWithUser(ctx, params)
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByUserAndStatus(ctx context.Context, params sqlc.ListTeamJoinRequestsByUserAndStatusParams) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByUserAndStatus(ctx, params)
}

func (r *TeamJoinRequestRepository) DeleteByUserAndStatus(ctx context.Context, params sqlc.DeleteJoinRequestsByUserAndStatusParams) error {
	return r.db.Query.DeleteJoinRequestsByUserAndStatus(ctx, params)
}

func (r *TeamJoinRequestRepository) UpdateStatus(ctx context.Context, params sqlc.UpdateTeamJoinRequestParams) (*sqlc.TeamJoinRequest, error) {
	request, err := r.db.Query.UpdateTeamJoinRequest(ctx, params)

	if err != nil {
		return nil, err
	}

	return &request, nil
}
