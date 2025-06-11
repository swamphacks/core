package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type SessionRepository struct {
	db *db.DB
}

func NewSessionRepository(db *db.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) NewTx(tx pgx.Tx) *SessionRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &SessionRepository{
		db: txDB,
	}
}

func (r *SessionRepository) Create(ctx context.Context, params sqlc.CreateSessionParams) (*sqlc.AuthSession, error) {
	session, err := r.db.Query.CreateSession(ctx, params)
	return &session, err
}

func (r *SessionRepository) Invalidate(ctx context.Context, sessionId uuid.UUID) error {
	return r.db.Query.InvalidateSessionByID(ctx, sessionId)
}
