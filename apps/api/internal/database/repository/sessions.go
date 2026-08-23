package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type SessionRepository struct {
	db *database.DB
}

func NewSessionRepository(db *database.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) NewTx(tx pgx.Tx) *SessionRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &SessionRepository{
		db: txDB,
	}
}

func (r *SessionRepository) CreateForUser(ctx context.Context, params sqlc.CreateSessionForUserParams) (*sqlc.Session, error) {
	session, err := r.db.Query.CreateSessionForUser(ctx, params)
	return &session, err
}

func (r *SessionRepository) CreateForAPIKey(ctx context.Context, params sqlc.CreateSessionForAPIKeyParams) (*sqlc.Session, error) {
	session, err := r.db.Query.CreateSessionForAPIKey(ctx, params)
	return &session, err
}

func (r *SessionRepository) Invalidate(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.Query.InvalidateSessionByID(ctx, sessionID)
}
