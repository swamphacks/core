package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrEmailCampaignNotFound     = errors.New("email campaign not found")
	ErrEmailCampaignNotClaimable = errors.New("email campaign is no longer claimable for sending")
)

type EmailCampaignRepository struct {
	db *database.DB
}

func (r *EmailCampaignRepository) NewTx(tx pgx.Tx) *EmailCampaignRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}
	return &EmailCampaignRepository{db: txDB}
}

func NewEmailCampaignRepository(db *database.DB) *EmailCampaignRepository {
	return &EmailCampaignRepository{db: db}
}

func (r *EmailCampaignRepository) CreateEmailCampaign(
	ctx context.Context,
	params sqlc.CreateEmailCampaignParams,
) (*sqlc.EmailCampaign, error) {
	campaign, err := r.db.Query.CreateEmailCampaign(ctx, params)
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *EmailCampaignRepository) GetEmailCampaignByID(
	ctx context.Context,
	params sqlc.GetEmailCampaignByIDParams,
) (*sqlc.EmailCampaign, error) {
	campaign, err := r.db.Query.GetEmailCampaignByID(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailCampaignNotFound
	}
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *EmailCampaignRepository) ListEmailCampaigns(
	ctx context.Context,
	hackathonID string,
) ([]sqlc.EmailCampaign, error) {
	return r.db.Query.ListEmailCampaigns(ctx, hackathonID)
}

func (r *EmailCampaignRepository) UpdateEmailCampaign(
	ctx context.Context,
	params sqlc.UpdateEmailCampaignParams,
) (*sqlc.EmailCampaign, error) {
	campaign, err := r.db.Query.UpdateEmailCampaign(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailCampaignNotFound
	}
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *EmailCampaignRepository) UpdateEmailCampaignStatus(
	ctx context.Context,
	params sqlc.UpdateEmailCampaignStatusParams,
) (*sqlc.EmailCampaign, error) {
	campaign, err := r.db.Query.UpdateEmailCampaignStatus(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailCampaignNotFound
	}
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *EmailCampaignRepository) GetApplicantContactEmailsByStatus(
	ctx context.Context,
	params sqlc.GetApplicantContactEmailsByStatusParams,
) ([]string, error) {
	return r.db.Query.GetApplicantContactEmailsByStatus(ctx, params)
}

func (r *EmailCampaignRepository) GetUserContactEmailsByRoles(
	ctx context.Context,
	roles []string,
) ([]string, error) {
	return r.db.Query.GetUserContactEmailsByRoles(ctx, roles)
}

func (r *EmailCampaignRepository) DeleteEmailCampaign(
	ctx context.Context,
	params sqlc.DeleteEmailCampaignParams,
) error {
	rows, err := r.db.Query.DeleteEmailCampaign(ctx, params)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrEmailCampaignNotFound
	}
	return nil
}

// ClaimCampaignForSending atomically moves a campaign into "sending".
// A campaign that is no longer in a claimable state matches no row, which we
// surface as ErrEmailCampaignNotClaimable rather than a generic not-found.
func (r *EmailCampaignRepository) ClaimCampaignForSending(
	ctx context.Context,
	params sqlc.ClaimCampaignForSendingParams,
) (*sqlc.EmailCampaign, error) {
	campaign, err := r.db.Query.ClaimCampaignForSending(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailCampaignNotClaimable
	}
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *EmailCampaignRepository) ListDueScheduledCampaigns(
	ctx context.Context,
) ([]sqlc.EmailCampaign, error) {
	return r.db.Query.ListDueScheduledCampaigns(ctx)
}

func (r *EmailCampaignRepository) GetInterestSubscriberEmails(
	ctx context.Context,
	hackathonID string,
) ([]string, error) {
	return r.db.Query.GetInterestSubscriberEmails(ctx, hackathonID)
}
