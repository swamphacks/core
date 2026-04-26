package workshops

import (
	"context"
	"errors"


	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type WorkshopService struct {
	workshopsRepo *repository.WorkshopsRepository
	logger 	  zerolog.Logger
}

func NewService(workshopsRepo *repository.WorkshopsRepository, logger zerolog.Logger) *WorkshopService {
	return &WorkshopService{
		workshopsRepo: workshopsRepo,
		logger: logger,
	}
}


func (s *WorkshopService) GetWorkshop(ctx context.Context, workshopID uuid.UUID) (*sqlc.Workshop, error) {
	

	workshop, err := s.workshopsRepo.GetWorkshop(ctx, workshopID)

	if err != nil {
		s.logger.Err(err).Msg("Failed to get workshop (1)")
		return nil, errors.New("Failed to get workshop")
	}

	return workshop, nil
}

func (s *WorkshopService) GetAllWorkshops(ctx context.Context) ([]sqlc.Workshop, error){
	workshops, err := s.workshopsRepo.GetAllWorkshops(ctx)

	if err != nil {
		s.logger.Err(err).Msg("Failed to get all workshops")
		return nil, errors.New("Failed to get all workshop")
	}
	return workshops, nil
}

func (s *WorkshopService) ViewAllWorkshops(ctx context.Context) ([]sqlc.Workshop, error){
	workshops, err := s.workshopsRepo.ViewAllWorkshops(ctx)

	if err != nil {
		s.logger.Err(err).Msg("Failed to get all workshops")
		return nil, errors.New("Failed to get all workshop")
	}
	return workshops, nil
}



func (s *WorkshopService) DeleteWorkshop(ctx context.Context, workshopID uuid.UUID) error {
	err := s.workshopsRepo.DeleteWorkshop(ctx, workshopID)

	if err != nil {
		s.logger.Err(err).Msg("")
		return errors.New("Failed to delete workshop")
	}
	return nil
}

func (s *WorkshopService) DeleteAllWorkshops(ctx context.Context) error {
	err := s.workshopsRepo.DeleteWorkshopAll(ctx)

	if err != nil {
		s.logger.Err(err).Msg("")
		return errors.New("Failed to delete all workshops")
	}
	return nil
}

func (s *WorkshopService) CreateWorkshop(ctx context.Context, params sqlc.CreateWorkshopParams) (*sqlc.Workshop, error) {
	workshop, err := s.workshopsRepo.CreateWorkshop(ctx, params)

	if err != nil {
		s.logger.Err(err).Msg("")
		return nil, errors.New("Failed to create workshop")
	}

	return workshop, nil
}

func (s *WorkshopService) RegisterWorkshop(ctx context.Context, userID uuid.UUID, workshopID uuid.UUID) (*sqlc.Workshop, error){
	err := s.workshopsRepo.RegisterUserForWorkshop(ctx, sqlc.RegisterUserForWorkshopParams{
		UserID:     userID,
		WorkshopID: workshopID,
	})

	if err != nil {
		s.logger.Err(err).Msg("failed to register user for workshop")
		return nil, errors.New("failed to register for workshop")
	}

	err = s.workshopsRepo.IncrementWorkshopAttendees(ctx, workshopID)
	if err != nil {
		s.logger.Err(err).Msg("failed to increment attendees")
		return nil, errors.New("failed to update attendees")
	}

	workshop, err := s.workshopsRepo.GetWorkshop(ctx, workshopID)
	if err != nil {
		s.logger.Err(err).Msg("failed to fetch updated workshop")
		return nil, errors.New("failed to fetch workshop")
	}

	return workshop, nil
}

func (s *WorkshopService) UnregisterWorkshop(ctx context.Context,userID uuid.UUID, workshopID uuid.UUID)   error {

	err := s.workshopsRepo.UnregisterUserForWorkshop(ctx, sqlc.UnregisterUserForWorkshopParams{
		UserID:     userID,
		WorkshopID: workshopID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Err(err).Msg("registration not found")
			return errors.New("registration not found")
		}
		s.logger.Err(err).Msg("failed to unregister user from workshop")
		return errors.New("failed to unregister from workshop")
	}

	err = s.workshopsRepo.DecrementWorkshopAttendees(ctx, workshopID)
	if err != nil {
		s.logger.Err(err).Msg("failed to decrement attendees")
		return errors.New("failed to update attendees")
	}
	

	return nil
}

func (s *WorkshopService) UpdateWorkshop(ctx context.Context, params sqlc.UpdateWorkshopParams) (*sqlc.Workshop, error) {
	workshop, err := s.workshopsRepo.UpdateWorkshop(ctx, sqlc.UpdateWorkshopParams{
		Title: params.Title,
		Description: params.Description,
		StartTime: params.StartTime,
		EndTime: params.EndTime,
		Location: params.Location,
		Presenter: params.Presenter,
		WorkshopID: params.WorkshopID,
	})

	if err != nil {
		s.logger.Err(err).Msg("")
		return nil, errors.New("Failed to update workshop")
	}

	return workshop, nil
}