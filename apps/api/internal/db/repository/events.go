package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type EventRepository struct {
	db *db.DB
}

func NewEventRespository(db *db.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) NewTx(tx pgx.Tx) *EventRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &EventRepository{
		db: txDB,
	}
}

func (r *EventRepository) CreateEvent(ctx context.Context, params sqlc.CreateEventParams) (*sqlc.Event, error) {
	event, err := r.db.Query.CreateEvent(ctx, params)
	return &event, err
}
