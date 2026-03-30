package repository

import (
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type TeamInvitationRepository struct {
	db *database.DB
}

func NewTeamInvitationRespository(db *database.DB) *TeamInvitationRepository {
	return &TeamInvitationRepository{
		db: db,
	}
}

func (r *TeamInvitationRepository) NewTx(tx pgx.Tx) *TeamInvitationRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamInvitationRepository{
		db: txDB,
	}
}
