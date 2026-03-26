package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type ApplicationRepository struct {
	db *database.DB
}

func NewApplicationRepository(db *database.DB) *ApplicationRepository {
	return &ApplicationRepository{
		db: db,
	}
}

func (r *ApplicationRepository) NewTx(tx pgx.Tx) *ApplicationRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &ApplicationRepository{
		db: txDB,
	}
}

func (r *ApplicationRepository) CreateApplication(ctx context.Context, userId uuid.UUID) (*sqlc.Application, error) {
	application, err := r.db.Query.CreateApplication(ctx, userId)

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (r *ApplicationRepository) GetApplicationByUserId(ctx context.Context, userId uuid.UUID) (*sqlc.Application, error) {
	application, err := r.db.Query.GetApplicationByUserId(ctx, userId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrApplicationNotFound
		}

		return nil, err
	}

	return &application, nil
}

func (r *ApplicationRepository) UpdateApplicationsStatuses(ctx context.Context, status sqlc.ApplicationStatus, userIds uuid.UUIDs) error {
	return r.db.Query.UpdateApplicationStatus(ctx, sqlc.UpdateApplicationStatusParams{
		Status:  status,
		UserIds: userIds,
	})
}

// List all candidates considered for admission.
// This queries for all applications who are 'under_review' and have their rating fields filled out.
// It also LEFT JOINs in their team id (if they have one) for further grouping based on teams.
func (r *ApplicationRepository) ListAdmissionCandidates(ctx context.Context) ([]sqlc.ListAdmissionCandidatesRow, error) {
	return r.db.Query.ListAdmissionCandidates(ctx)
}

func (r *ApplicationRepository) SubmitApplication(ctx context.Context, data any, userId uuid.UUID) error {
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
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepository) SaveApplication(ctx context.Context, data any, userId uuid.UUID) error {
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
		SavedAtDoUpdate:     true,
		SavedAt:             time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepository) UpdateApplication(ctx context.Context, params sqlc.UpdateApplicationParams) error {
	return r.db.Query.UpdateApplication(ctx, params)
}

func (r *ApplicationRepository) ListAvailableApplications(ctx context.Context) ([]uuid.UUID, error) {
	return r.db.Query.ListAvailableApplications(ctx)
}

func (r *ApplicationRepository) AssignApplicationToReview(ctx context.Context, reviewerId uuid.UUID, applicationIDs []uuid.UUID) error {
	return r.db.Query.AssignApplicationsToReviewer(ctx, sqlc.AssignApplicationsToReviewerParams{
		ReviewerID:     reviewerId,
		ApplicationIds: applicationIDs,
	})
}

func (r *ApplicationRepository) ListApplicationByReviewer(ctx context.Context, reviewerId uuid.UUID) ([]sqlc.ListApplicationByReviewerRow, error) {
	return r.db.Query.ListApplicationByReviewer(ctx, &reviewerId)
}

func (r *ApplicationRepository) ResetApplicationReviews(ctx context.Context) error {
	return r.db.Query.ResetApplicationReviews(ctx)
}

func (r *ApplicationRepository) GetSubmittedApplicationGenders(ctx context.Context) (sqlc.GetApplicationGenderSplitRow, error) {
	return r.db.Query.GetApplicationGenderSplit(ctx)
}

func (r *ApplicationRepository) GetSubmittedApplicationRaces(ctx context.Context) ([]sqlc.GetApplicationRaceSplitRow, error) {
	return r.db.Query.GetApplicationRaceSplit(ctx)
}

func (r *ApplicationRepository) GetSubmittedApplicationAges(ctx context.Context) (sqlc.GetApplicationAgeSplitRow, error) {
	return r.db.Query.GetApplicationAgeSplit(ctx)
}

func (r *ApplicationRepository) GetSubmittedApplicationMajors(ctx context.Context) ([]sqlc.GetApplicationMajorSplitRow, error) {
	return r.db.Query.GetApplicationMajorSplit(ctx)
}

func (r *ApplicationRepository) GetSubmittedApplicationSchools(ctx context.Context) ([]sqlc.GetApplicationSchoolSplitRow, error) {
	return r.db.Query.GetApplicationSchoolSplit(ctx)
}

func (r *ApplicationRepository) GetApplicationStatuses(ctx context.Context) (sqlc.GetApplicationStatusSplitRow, error) {
	return r.db.Query.GetApplicationStatusSplit(ctx)
}

func (r *ApplicationRepository) GetNonReviewedApplications(ctx context.Context) ([]uuid.UUID, error) {
	return r.db.Query.ListNonReviewedApplications(ctx)
}

func (r *ApplicationRepository) GetSubmissionTimes(ctx context.Context) ([]sqlc.GetSubmissionTimesRow, error) {
	return r.db.Query.GetSubmissionTimes(ctx)
}

func (r *ApplicationRepository) JoinWaitlist(ctx context.Context, userId uuid.UUID) error {
	return r.db.Query.JoinWaitlist(ctx, userId)
}

func (r *ApplicationRepository) TransitionAcceptedApplicationsToWaitlist(ctx context.Context) error {
	return r.db.Query.TransitionAcceptedApplicationsToWaitlist(ctx)
}

func (r *ApplicationRepository) TransitionWaitlistedApplicationsToAccepted(ctx context.Context, acceptanceCount int32) ([]uuid.UUID, error) {
	return r.db.Query.TransitionWaitlistedApplicationsToAccepted(ctx, acceptanceCount)
}

func (r *ApplicationRepository) GetAttendeeCount(ctx context.Context) (uint32, error) {
	amount, err := r.db.Query.GetAttendeeCount(ctx)

	return uint32(amount), err
}
