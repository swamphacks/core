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
	ErrCampaignNotFound         = errors.New("campaign not found")
	ErrNoCampaignsDeleted       = errors.New("no campsigns deleted")
	ErrMultipleCampaignsDeleted = errors.New("multiple campaigns affected by delete query while only expecting one to delete one")
)

type CampaignRepository struct {
	db *db.DB
}

func NewCampaignRepository(db *db.DB) *CampaignRepository {
	return &CampaignRepository{
		db: db,
	}
}

func (r *CampaignRepository) CreateCampaign(ctx context.Context, params sqlc.CreateCampaignParams) (*sqlc.Campaign, error) {
	campaign, err := r.db.Query.CreateCampaign(ctx, params)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (r *CampaignRepository) GetCampaignByID(ctx context.Context, id uuid.UUID) (*sqlc.Campaign, error) {
	campaign, err := r.db.Query.GetCampaignById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCampaignNotFound
	} else if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (r *CampaignRepository) UpdateCampaignById(ctx context.Context, params sqlc.UpdateCampaignByIdParams) error {
	err := r.db.Query.UpdateCampaignById(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCampaignNotFound
	}
	return err
}

func (r *CampaignRepository) DeleteCampaignById(ctx context.Context, id uuid.UUID) error {
	affectedRows, err := r.db.Query.DeleteCampaignById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCampaignNotFound
	}
	if affectedRows == 0 {
		return ErrNoCampaignsDeleted
	} else if affectedRows > 1 {
		return ErrMultipleCampaignsDeleted
	}
	return err
}
