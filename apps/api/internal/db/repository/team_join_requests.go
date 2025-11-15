package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/logger"
	"github.com/swamphacks/core/apps/api/internal/ptr"
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

func (r *TeamJoinRequestRepository) GetTeamJoinRequestsByUser(ctx context.Context, userId uuid.UUID) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByUserID(ctx, userId)
}

func (r *TeamJoinRequestRepository) GetTeamJoinRequestsByTeam(ctx context.Context, teamId uuid.UUID, status sqlc.JoinRequestStatus) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByTeamIDAndStatus(ctx, sqlc.ListTeamJoinRequestsByTeamIDAndStatusParams{
		TeamID: teamId,
		Status: status,
	})
}

func (r *TeamJoinRequestRepository) GetTeamJoinRequestsByUserAndEvent(ctx context.Context, userId, eventId uuid.UUID, status sqlc.JoinRequestStatus) ([]sqlc.TeamJoinRequest, error) {
	return r.db.Query.ListTeamJoinRequestsByUserAndEventAndStatus(ctx, sqlc.ListTeamJoinRequestsByUserAndEventAndStatusParams{
		UserID:  userId,
		EventID: ptr.UUIDToPtr(eventId),
		Status:  status,
	})
}
