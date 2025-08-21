package services

import (
	"context"

	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/storage"
)

type ApplicationService struct {
	appRepo *repository.ApplicationRepository
	storage storage.Storage
}

func NewApplicationService(appRepo *repository.ApplicationRepository, storage storage.Storage) *ApplicationService {
	return &ApplicationService{
		appRepo: appRepo,
		storage: storage,
	}
}

// THIS IS AN EXAMPLE FUNCTION
func (s *ApplicationService) CreateApplication() error {
	s.storage.Store(context.Background(), "core-application-resumes-dev", "application_data", []byte("application data"))
	return nil
}
