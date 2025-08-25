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
	ErrCreateApplication      = errors.New("unable to create application")
	ErrSaveApplication        = errors.New("unable to save application")
	ErrSubmitApplication      = errors.New("unable to submit application")
	ErrInvalidApplicationData = errors.New("unable to parse application data")
	ErrGetApplication         = errors.New("unable to get application for user")
	ErrApplicationNotFound    = errors.New("can not find application for user")
)

type ApplicationRepository struct {
	db *db.DB
}

func NewApplicationRepository(db *db.DB) *ApplicationRepository {
	return &ApplicationRepository{
		db: db,
	}
}

func (r *ApplicationRepository) NewTx(tx pgx.Tx) *ApplicationRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &ApplicationRepository{
		db: txDB,
	}
}

func (r *ApplicationRepository) CreateApplication(ctx context.Context, params sqlc.CreateApplicationParams) (*sqlc.Application, error) {
	application, err := r.db.Query.CreateApplication(ctx, params)

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (r *ApplicationRepository) GetApplicationByUserAndEventID(ctx context.Context, params sqlc.GetApplicationByUserAndEventIDParams) (*sqlc.Application, error) {
	application, err := r.db.Query.GetApplicationByUserAndEventID(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}

		return nil, err
	}

	return &application, nil
}

func (r *ApplicationRepository) SubmitApplication(ctx context.Context, data any, userId, eventId uuid.UUID) error {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = r.db.Query.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		StatusDoUpdate:      true,
		Status:              sqlc.ApplicationStatusSubmitted,
		ApplicationDoUpdate: true,
		Application:         jsonBytes,
		UserID:              userId,
		EventID:             eventId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepository) SaveApplication(ctx context.Context, data any, params sqlc.UpdateApplicationParams) error {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = r.db.Query.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		StatusDoUpdate:      true,
		Status:              sqlc.ApplicationStatusStarted,
		ApplicationDoUpdate: true,
		Application:         jsonBytes,
		UserID:              params.UserID,
		EventID:             params.EventID,
	})

	if err != nil {
		return err
	}

	return nil
}
