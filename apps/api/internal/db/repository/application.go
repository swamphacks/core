package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

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
		if errors.Is(err, pgx.ErrNoRows) {
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
		SubmittedAtDoUpdate: true,
		SubmittedAt:         time.Now(),
		SavedAtDoUpdate:     true,
		SavedAt:             time.Now(),
		UserID:              userId,
		EventID:             eventId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepository) SaveApplication(ctx context.Context, data any, userId, eventId uuid.UUID) error {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = r.db.Query.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		StatusDoUpdate:      true,
		Status:              sqlc.ApplicationStatusStarted,
		ApplicationDoUpdate: true,
		Application:         jsonBytes,
		UserID:              userId,
		EventID:             eventId,
		SavedAtDoUpdate:     true,
		SavedAt:             time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepository) ListAvailableApplicationForEvent(ctx context.Context, eventId uuid.UUID) ([]uuid.UUID, error) {
	return r.db.Query.ListAvailableApplicationsForEvent(ctx, eventId)
}

func (r *ApplicationRepository) AssignApplicationToReviewByEvent(ctx context.Context, reviewerId, eventId uuid.UUID, applicationIDs []uuid.UUID) error {
	return r.db.Query.AssignApplicationsToReviewer(ctx, sqlc.AssignApplicationsToReviewerParams{
		ReviewerID:     reviewerId,
		EventID:        eventId,
		ApplicationIds: applicationIDs,
	})
}

func (r *ApplicationRepository) ResetApplicationReviewsForEvent(ctx context.Context, eventId uuid.UUID) error {
	return r.db.Query.ResetApplicationReviews(ctx, eventId)
}

// Application statistics (Staff Dashboards)
func (r *ApplicationRepository) GetSubmittedApplicationGenders(ctx context.Context, eventId uuid.UUID) (sqlc.GetApplicationGenderSplitRow, error) {
	return r.db.Query.GetApplicationGenderSplit(ctx, eventId)
}

func (r *ApplicationRepository) GetSubmittedApplicationRaces(ctx context.Context, eventId uuid.UUID) ([]sqlc.GetApplicationRaceSplitRow, error) {
	return r.db.Query.GetApplicationRaceSplit(ctx, eventId)
}

func (r *ApplicationRepository) GetSubmittedApplicationAges(ctx context.Context, eventId uuid.UUID) (sqlc.GetApplicationAgeSplitRow, error) {
	return r.db.Query.GetApplicationAgeSplit(ctx, eventId)
}

func (r *ApplicationRepository) GetSubmittedApplicationMajors(ctx context.Context, eventId uuid.UUID) ([]sqlc.GetApplicationMajorSplitRow, error) {
	return r.db.Query.GetApplicationMajorSplit(ctx, eventId)
}

func (r *ApplicationRepository) GetSubmittedApplicationSchools(ctx context.Context, eventId uuid.UUID) ([]sqlc.GetApplicationSchoolSplitRow, error) {
	return r.db.Query.GetApplicationSchoolSplit(ctx, eventId)
}

func (r *ApplicationRepository) GetApplicationStatuses(ctx context.Context, eventId uuid.UUID) (sqlc.GetApplicationStatusSplitRow, error) {
	return r.db.Query.GetApplicationStatusSplit(ctx, eventId)
}

func (r *ApplicationRepository) GetSubmissionTimes(ctx context.Context, eventId uuid.UUID) ([]sqlc.GetSubmissionTimesRow, error) {
	return r.db.Query.GetSubmissionTimes(ctx, eventId)
}
