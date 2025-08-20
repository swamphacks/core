package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrEventNotFound          = errors.New("event not found")
	ErrEventRoleNotFound      = errors.New("event role not found")
	ErrDuplicateEvent         = errors.New("event already exists in database")
	ErrNoEventsDeleted        = errors.New("no events deleted")
	ErrMultipleEventsDeleted  = errors.New("multiple events affected by delete query while only expecting one to delete one")
	ErrCreateApplication      = errors.New("unable to create application")
	ErrSaveApplication        = errors.New("unable to save application")
	ErrSubmitApplication      = errors.New("unable to submit application")
	ErrInvalidApplicationData = errors.New("unable to parse application data")
	ErrGetApplication         = errors.New("unable to get application for user")
	ErrApplicationNotFound    = errors.New("can not find application for user")
	ErrUnknown                = errors.New("an unkown error was caught")
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
		if err == pgx.ErrNoRows {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return &event, err
}

func (r *EventRepository) UpdateEventById(ctx context.Context, params sqlc.UpdateEventByIdParams) error {
	err := r.db.Query.UpdateEventById(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrEventNotFound
		}
	}
	return err
}

func (r *EventRepository) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteEventById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
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

func (r *EventRepository) GetEventsWithRoles(ctx context.Context, userId *uuid.UUID, includeUnpublished bool) (*[]sqlc.GetEventsWithUserInfoRow, error) {
	events, err := r.db.Query.GetEventsWithUserInfo(ctx, sqlc.GetEventsWithUserInfoParams{
		UserID:             userId,
		IncludeUnpublished: includeUnpublished,
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
		if err == pgx.ErrNoRows {
			return nil, ErrEventRoleNotFound
		}
	}

	return &eventRole, err
}

func (r *EventRepository) GetEventStaff(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetEventStaffRow, error) {
	users, err := r.db.Query.GetEventStaff(ctx, eventId)
	return &users, err
}

func (r *EventRepository) AssignRole(ctx context.Context, params sqlc.AssignRoleParams) error {
	return r.db.Query.AssignRole(ctx, params)
}

func (r *EventRepository) CreateApplication(ctx context.Context, params sqlc.CreateApplicationParams) (*sqlc.Application, error) {
	application, err := r.db.Query.CreateApplication(ctx, params)

	if err != nil {
		return nil, ErrCreateApplication
	}

	return &application, nil
}

func (r *EventRepository) GetApplicationByUserAndEventID(ctx context.Context, params sqlc.GetApplicationByUserAndEventIDParams) (*sqlc.Application, error) {
	application, err := r.db.Query.GetApplicationByUserAndEventID(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}

		return nil, ErrGetApplication
	}

	return &application, nil
}

func (r *EventRepository) SubmitApplication(ctx context.Context, data any, userId, eventId uuid.UUID) error {
	resumeUrl := "TODO"

	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return ErrInvalidApplicationData
	}

	err = r.db.Query.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		StatusDoUpdate:      true,
		Status:              sqlc.ApplicationStatusStarted,
		ApplicationDoUpdate: true,
		Application:         jsonBytes,
		ResumeUrlDoUpdate:   true,
		ResumeUrl:           &resumeUrl,
		UserID:              userId,
		EventID:             eventId,
	})

	if err != nil {
		return ErrSubmitApplication
	}

	return nil
}

func (r *EventRepository) SaveApplication(ctx context.Context, data any, params sqlc.UpdateApplicationParams) error {
	resumeUrl := "TODO"

	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return ErrInvalidApplicationData
	}

	err = r.db.Query.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		StatusDoUpdate:      true,
		Status:              sqlc.ApplicationStatusStarted,
		ApplicationDoUpdate: true,
		Application:         jsonBytes,
		ResumeUrlDoUpdate:   true,
		ResumeUrl:           &resumeUrl,
		UserID:              params.UserID,
		EventID:             params.EventID,
	})

	if err != nil {
		return ErrSaveApplication
	}

	return nil
}
