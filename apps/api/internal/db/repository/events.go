package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return &event, err
}

func (r *EventRepository) UpdateEventById(ctx context.Context, params sqlc.UpdateEventByIdParams) error {
	err := r.db.Query.UpdateEventById(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrEventNotFound
		}
	}
	return err
}

func (r *EventRepository) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteEventById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

func (r *EventRepository) GetAllEvents(ctx context.Context, userId uuid.UUID) (*[]sqlc.GetAllEventsRow, error) {
	events, err := r.db.Query.GetAllEvents(ctx, userId)
	return &events, err
}

func (r *EventRepository) GetPublishedEvents(ctx context.Context, userId uuid.UUID) (*[]sqlc.GetPublishedEventsRow, error) {
	events, err := r.db.Query.GetPublishedEvents(ctx, userId)
	return &events, err
}

// find a struct for valid roles to enforce type safety
func (r *EventRepository) GetEventsWithRoles(ctx context.Context, userId *uuid.UUID, scope sqlc.GetEventScopeType) (*[]sqlc.GetEventsWithUserInfoRow, error) {
	events, err := r.db.Query.GetEventsWithUserInfo(ctx, sqlc.GetEventsWithUserInfoParams{
		UserID: userId,
		Scope:  scope,
	})
	return &events, err
}

func (r *EventRepository) GetEventRoleByIds(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (*sqlc.EventRole, error) {
	params := sqlc.GetEventRoleByIdsParams{
		UserID:  userId,
		EventID: eventId,
	}

	eventRole, err := r.db.Query.GetEventRoleByIds(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventRoleNotFound
		}
	}

	return &eventRole, err
}

func (r *EventRepository) GetEventStaff(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetEventStaffRow, error) {
	users, err := r.db.Query.GetEventStaff(ctx, eventId)
	return &users, err
}

func (r *EventRepository) GetEventUsers(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetEventUsersRow, error) {
	users, err := r.db.Query.GetEventUsers(ctx, eventId)
	return &users, err
}

func (r *EventRepository) AssignRole(ctx context.Context, params sqlc.AssignRoleParams) error {
	return r.db.Query.AssignRole(ctx, params)
}

func (r *EventRepository) RevokeRole(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) error {
	params := sqlc.RemoveRoleParams{
		UserID:  userId,
		EventID: eventId,
	}

	return r.db.Query.RemoveRole(ctx, params)
}
