package repository

import (
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
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
