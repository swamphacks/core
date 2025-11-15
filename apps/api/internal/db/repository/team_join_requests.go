package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/logger"
)

type TeamJoinRequestRepository struct {
	db *db.DB
}

func NewTeamJoinRequestRepository(db *db.DB) *TeamJoinRequestRepository {
	return &TeamJoinRequestRepository{
		db: db,
	}
}

func (r *TeamJoinRequestRepository) NewTx(tx pgx.Tx) *TeamJoinRequestRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamJoinRequestRepository{
		db: txDB,
	}
}

func (r *TeamJoinRequestRepository) Create(ctx context.Context, teamId, userId uuid.UUID, message *string) (*sqlc.TeamJoinRequest, error) {

	l := logger.New()
	l.Debug().Msg("Before creating...")
	request, err := r.db.Query.CreateTeamJoinRequest(ctx, sqlc.CreateTeamJoinRequestParams{
		TeamID:         teamId,
		UserID:         userId,
		RequestMessage: message,
	})
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *TeamJoinRequestRepository) GetById(ctx context.Context, requestId uuid.UUID) (*sqlc.TeamJoinRequest, error) {
	request, err := r.db.Query.GetTeamJoinRequestByID(ctx, requestId)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByUser(ctx context.Context, userId uuid.UUID) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByUserID(ctx, userId)
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByTeam(ctx context.Context, teamId uuid.UUID, status sqlc.JoinRequestStatus) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByTeamIDAndStatus(ctx, sqlc.ListTeamJoinRequestsByTeamIDAndStatusParams{
		TeamID: teamId,
		Status: status,
	})
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByTeamWithUser(ctx context.Context, teamId uuid.UUID, status sqlc.JoinRequestStatus) ([]sqlc.ListJoinRequestsByTeamAndStatusWithUserRow, error) {
	return r.db.Query.ListJoinRequestsByTeamAndStatusWithUser(ctx, sqlc.ListJoinRequestsByTeamAndStatusWithUserParams{
		TeamID: teamId,
		Status: status,
	})
}

func (r *TeamJoinRequestRepository) ListJoinRequestsByUserAndEvent(ctx context.Context, userId, eventId uuid.UUID, status sqlc.JoinRequestStatus) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByUserAndEventAndStatus(ctx, sqlc.ListTeamJoinRequestsByUserAndEventAndStatusParams{
		UserID:  userId,
		EventID: &eventId,
		Status:  status,
	})
}

func (r *TeamJoinRequestRepository) DeleteByUserAndEventAndStatus(ctx context.Context, userId, eventId uuid.UUID, status sqlc.JoinRequestStatus) error {
	return r.db.Query.DeleteJoinRequestsByUserAndEventAndStatus(ctx, sqlc.DeleteJoinRequestsByUserAndEventAndStatusParams{
		UserID:  userId,
		EventID: &eventId,
		Status:  status,
	})
}

func (r *TeamJoinRequestRepository) UpdateStatus(ctx context.Context, requestId uuid.UUID, status sqlc.JoinRequestStatus) (*sqlc.TeamJoinRequest, error) {
	request, err := r.db.Query.UpdateTeamJoinRequest(ctx, sqlc.UpdateTeamJoinRequestParams{
		ID:             requestId,
		StatusDoUpdate: true,
		Status:         status,
	})
	if err != nil {
		return nil, err
	}

	return &request, nil
}
