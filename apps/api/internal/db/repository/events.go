package repository

import (
	"context"

	"github.com/google/uuid"
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

func (r *EventRepository) CreateEvent(ctx context.Context, params sqlc.CreateEventParams) (*sqlc.Event, error) {
	event, err := r.db.Query.CreateEvent(ctx, params)
	return &event, err
}

func (r *EventRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*sqlc.Event, error) {
	event, err := r.db.Query.GetEventByID(ctx, id)
	return &event, err
}

func (r *EventRepository) UpdateEventById(ctx context.Context, params sqlc.UpdateEventByIdParams) error {
	err := r.db.Query.UpdateEventById(ctx, params)
	return err
}
