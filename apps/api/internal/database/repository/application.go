package repository

import (
	"context"
	"errors"

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

func (r *ApplicationRepository) CreateApplication(ctx context.Context, params sqlc.CreateApplicationParams) (*sqlc.Application, error) {
	application, err := r.db.Query.CreateApplication(ctx, params)

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (r *ApplicationRepository) UpdateApplication(ctx context.Context, params sqlc.UpdateApplicationParams) error {
	return r.db.Query.UpdateApplication(ctx, params)
}

func (r *ApplicationRepository) GetApplicationByUserId(ctx context.Context, userID uuid.UUID) (*sqlc.Application, error) {
	application, err := r.db.Query.GetApplicationByUserId(ctx, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrApplicationNotFound
		}

		return nil, err
	}

	return &application, nil
}

func (r *ApplicationRepository) UpdateApplicationsStatuses(ctx context.Context, params sqlc.UpdateApplicationStatusParams) error {
	return r.db.Query.UpdateApplicationStatus(ctx, params)
}

// List all candidates considered for admission.
// This queries for all applications who are 'under_review' and have their rating fields filled out.
// It also LEFT JOINs in their team id (if they have one) for further grouping based on teams.
func (r *ApplicationRepository) ListAdmissionCandidates(ctx context.Context) ([]sqlc.ListAdmissionCandidatesRow, error) {
	return r.db.Query.ListAdmissionCandidates(ctx)
}

func (r *ApplicationRepository) ListAvailableApplications(ctx context.Context) ([]uuid.UUID, error) {
	return r.db.Query.ListAvailableApplications(ctx)
}

func (r *ApplicationRepository) AssignApplicationToReview(ctx context.Context, params sqlc.AssignApplicationsToReviewerParams) error {
	return r.db.Query.AssignApplicationsToReviewer(ctx, params)
}

func (r *ApplicationRepository) ListApplicationByReviewer(ctx context.Context, reviewerID uuid.UUID) ([]sqlc.ListApplicationByReviewerRow, error) {
	return r.db.Query.ListApplicationByReviewer(ctx, &reviewerID)
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

func (r *ApplicationRepository) JoinWaitlist(ctx context.Context, userID uuid.UUID) error {
	return r.db.Query.JoinWaitlist(ctx, userID)
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
