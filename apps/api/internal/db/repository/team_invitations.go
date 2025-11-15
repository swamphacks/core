package repository

import (
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type TeamInvitationRepository struct {
	db *db.DB
}

func NewTeamInvitationRespository(db *db.DB) *TeamInvitationRepository {
	return &TeamInvitationRepository{
		db: db,
	}
}

func (r *TeamInvitationRepository) NewTx(tx pgx.Tx) *TeamInvitationRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamInvitationRepository{
		db: txDB,
	}
}
