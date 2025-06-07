package repository

import (
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
