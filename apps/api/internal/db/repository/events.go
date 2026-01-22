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
	ErrDuplicateRFID         = errors.New("RFID has already been scanned")
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

func (r *EventRepository) GetUserByRFID(ctx context.Context, eventId uuid.UUID, rfid string) (*sqlc.AuthUser, error) {
	params := sqlc.GetUserByRFIDParams{
		EventID: eventId,
		Rfid:    &rfid,
	}
	user, err := r.db.Query.GetUserByRFID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventRoleNotFound
		}
		return nil, ErrUnknown
	}
	return &user, nil
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

// **Deprecated**. Use `UpdateEventRoleByIds` instead.
// Only kept for backwards compatibility. TODO: Refactor all current implementations to use new function.
func (r *EventRepository) UpdateRole(ctx context.Context, userId uuid.UUID, eventId uuid.UUID, role sqlc.EventRoleType) error {
	params := sqlc.UpdateRoleParams{
		UserID:  userId,
		EventID: eventId,
		Role:    role,
	}
	return r.db.Query.UpdateRole(ctx, params)
}

func (r *EventRepository) UpdateEventRoleByIds(ctx context.Context, params sqlc.UpdateEventRoleByIdsParams) error {
	err := r.db.Query.UpdateEventRoleByIds(ctx, params)
	if err != nil {
		// Check if the error is a unique constraint violation on RFID
		if db.IsUniqueViolation(err) {
			return ErrDuplicateRFID
		}
	}
	return err
}

func (r *EventRepository) GetApplicationStatuses(ctx context.Context, eventId uuid.UUID) (sqlc.GetApplicationStatusSplitRow, error) {
	return r.db.Query.GetApplicationStatusSplit(ctx, eventId)
}

func (r *EventRepository) GetSubmissionTimes(ctx context.Context, eventId uuid.UUID) ([]sqlc.GetSubmissionTimesRow, error) {
	return r.db.Query.GetSubmissionTimes(ctx, eventId)
}

func (r *EventRepository) GetEventAttendeesWithDiscord(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetEventAttendeesWithDiscordRow, error) {
	attendees, err := r.db.Query.GetEventAttendeesWithDiscord(ctx, eventId)
	if err != nil {
		return nil, err
	}
	return &attendees, nil
}

func (r *EventRepository) GetEventRoleByDiscordIDAndEventId(ctx context.Context, discordID string, eventID uuid.UUID) (*sqlc.GetEventRoleByDiscordIDAndEventIdRow, error) {
	params := sqlc.GetEventRoleByDiscordIDAndEventIdParams{
		AccountID: discordID,
		EventID:   eventID,
	}

	eventRole, err := r.db.Query.GetEventRoleByDiscordIDAndEventId(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventRoleNotFound
		}
		return nil, err
	}

	return &eventRole, nil
}

func (r *EventRepository) GetAttendeeUserIdsByEventId(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error) {
	return r.db.Query.GetAttendeeUserIdsByEventId(ctx, eventID)
}
