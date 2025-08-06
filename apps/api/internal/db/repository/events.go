package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrEventNotFound         = errors.New("event not found")
	ErrEventRoleNotFound     = errors.New("event role not found")
	ErrDuplicateEvent        = errors.New("event already exists in database")
	ErrNoEventsDeleted       = errors.New("no events deleted")
	ErrMultipleEventsDeleted = errors.New("multiple events affected by delete query while only expecting one to delete one")
	ErrUnknown               = errors.New("an unkown error was caught")
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
	if db.IsUniqueViolation(err) {
		return nil, ErrDuplicateEvent
	}
	return &event, err
}

func (r *EventRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*sqlc.Event, error) {
	event, err := r.db.Query.GetEventByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return &event, err
}

func (r *EventRepository) UpdateEventById(ctx context.Context, params sqlc.UpdateEventByIdParams) error {
	err := r.db.Query.UpdateEventById(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrEventNotFound
		}
	}
	return err
}

func (r *EventRepository) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteEventById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrEventNotFound
		}
	}
	if affectedRows == 0 {
		return ErrNoEventsDeleted
	} else if affectedRows > 1 {
		return ErrMultipleEventsDeleted
	}

	return err
}

func (r *EventRepository) GetEventRoleByIds(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (*sqlc.EventRole, error) {
	params := sqlc.GetEventRoleByIdsParams{
		UserID:  userId,
		EventID: eventId,
	}

	eventRole, err := r.db.Query.GetEventRoleByIds(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEventRoleNotFound
		}
	}

	return &eventRole, err
}
