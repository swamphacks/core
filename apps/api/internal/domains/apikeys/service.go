package apikeys

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrApiKeyNotFound = repository.ErrApiKeyNotFound

	ErrApiKeyNameRequired        = errors.New("API Key name cannot be empty")
	ErrApiKeyDescriptionRequired = errors.New("API Key description cannot be empty")
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

func (s *ApiKeysService) ListApiKeys(
	ctx context.Context,
) ([]sqlc.ListApiKeysRow, error) {
	return s.apiKeysRepo.ListApiKeys(ctx)
}

func (s *ApiKeysService) GetApiKeyByID(
	ctx context.Context,
	id uuid.UUID,
) (*sqlc.GetApiKeyByIdRow, error) {
	return s.apiKeysRepo.GetApiKeyByID(ctx, id)
}

// Generates secret of length 64
func generateAPIKeySecret() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func hashAPIKeySecret(secret string) string {
    hash := sha256.Sum256([]byte(secret))
    return hex.EncodeToString(hash[:])
}

func (s *ApiKeysService) CreateApiKey(
	ctx context.Context,
	request *CreateApiKeyRequest,
) (*CreateApiKeyResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	secret, err := generateAPIKeySecret()
	if err != nil {
		return nil, err
	}
	secretHash := hashAPIKeySecret(secret)

	params := sqlc.CreateApiKeyParams{
		Name:        request.Name,
		Description: request.Description,
		Role:        request.Role,
		ExpiresAt:   request.ExpiresAt,
		SecretHash:  secretHash,
	}

	apiKey, err := s.apiKeysRepo.CreateApiKey(ctx, params)
	if err != nil {
		return nil, err
	}

	return &CreateApiKeyResponse{
		CreateApiKeyRow: apiKey,
		Secret: secret,
	}, nil
}

func (s *ApiKeysService) DeleteApiKeyByID(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.apiKeysRepo.DeleteApiKeyByID(ctx, id)
}