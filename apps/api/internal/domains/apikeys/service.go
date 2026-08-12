package apikeys

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
)

type ApiKeysService struct {
	apiKeysRepo *repository.ApiKeysRepository
	logger      zerolog.Logger
}

func NewService(
	apiKeysRepo *repository.ApiKeysRepository,
	logger zerolog.Logger,
) *ApiKeysService {
	return &ApiKeysService{
		apiKeysRepo: apiKeysRepo,
		logger: logger.With().Str("service", "ApiKeysService").Str("domain", "apikeys").Logger(),
	}
}