package repository

import (
	"context"


	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type WorkshopsRepository struct {
	db *database.DB
}

func NewWorkshopsRepository(db *database.DB) *WorkshopsRepository{
	return &WorkshopsRepository{
		db: db,
	}
}

func (r *WorkshopsRepository) GetWorkshop(ctx context.Context, workshopID uuid.UUID) (*sqlc.Workshop, error) {
	workshop, err := r.db.Query.GetWorkshop(ctx, workshopID)
	if err != nil {
		return nil, err
	}
	
	return &workshop, nil
}

func (r *WorkshopsRepository) UpdateWorkshop(ctx context.Context, params sqlc.UpdateWorkshopParams) (*sqlc.Workshop, error) {
	workshop, err := r.db.Query.UpdateWorkshop(ctx, params)
	if err != nil {
		return nil, err
	}
	
	return &workshop, nil
}

func (r *WorkshopsRepository) GetAllWorkshops(ctx context.Context) ([]sqlc.Workshop, error){
	workshop, err := r.db.Query.GetAllWorkshops(ctx)
	if err != nil {
		return nil, err
	}

	return workshop, nil
}

func (r *WorkshopsRepository) DeleteWorkshop(ctx context.Context, workshopID uuid.UUID) error {
	err := r.db.Query.DeleteWorkshop(ctx, workshopID)
	if err != nil {
		return err
	}

	return nil
}

func (r *WorkshopsRepository) DeleteWorkshopAll(ctx context.Context) error {
	err := r.db.Query.DeleteWorkshopAll(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *WorkshopsRepository) ViewAllWorkshops(ctx context.Context) ([]sqlc.Workshop, error){
	workshop, err := r.db.Query.ViewAllWorkshops(ctx)
	if err != nil {
		return nil, err
	}

	return workshop, nil
}

func (r *WorkshopsRepository) UnregisterUserForWorkshop(ctx context.Context, params sqlc.UnregisterUserForWorkshopParams) error {

	err := r.db.Query.UnregisterUserForWorkshop(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *WorkshopsRepository) RegisterUserForWorkshop(ctx context.Context, params sqlc.RegisterUserForWorkshopParams) error {

	err := r.db.Query.RegisterUserForWorkshop(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *WorkshopsRepository) DecrementWorkshopAttendees(ctx context.Context, workshopID uuid.UUID) error {

	err := r.db.Query.DecrementWorkshopAttendees(ctx, workshopID)
	if err != nil {
		return err
	}

	return nil
}

func (r *WorkshopsRepository) CreateWorkshop(ctx context.Context, params sqlc.CreateWorkshopParams) (*sqlc.Workshop, error){
	workshop, err := r.db.Query.CreateWorkshop(ctx, params)
	if err != nil {
		return nil, err
	}

	return &workshop, nil
}

func (r *WorkshopsRepository) IncrementWorkshopAttendees(ctx context.Context, workshopID uuid.UUID) error {
	err := r.db.Query.IncrementWorkshopAttendees(ctx, workshopID)
	if err != nil {
		return err
	}

	return nil
}
