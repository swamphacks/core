package apikeys

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type handler struct {
	apiKeysService *ApiKeysService
	logger         zerolog.Logger
}

func NewHandler(apiKeysService *ApiKeysService, logger zerolog.Logger) *handler {
	return &handler{
		apiKeysService: apiKeysService,
		logger:               logger.With().Str("handler", "ApiKeysHandler").Str("domain", "apikeys").Logger(),
	}
}

func (h *handler) getAllApiKeys(
	ctx context.Context,
	input *struct{},
) (*ListApiKeysOutput, error) {
	apiKeys, err := h.apiKeysService.ListApiKeys(ctx)
	if err != nil {
		return nil, apiKeyHTTPError(err, "Failed to list API Keys")
	}
	return &ListApiKeysOutput{Body: apiKeys}, nil
}

func (h *handler) getApiKey(
	ctx context.Context,
	input *struct{
		ApiKeyID string `path:"apiKeyId"`
	},
) (*GetApiKeyOutput, error) {
	apiKeyID, err := uuid.Parse(input.ApiKeyID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid API Key ID")
	}

	apiKey, err := h.apiKeysService.GetApiKeyByID(ctx, apiKeyID)
	if err != nil {
		return nil, apiKeyHTTPError(err, "Failed to get API Key")
	}

	return &GetApiKeyOutput{Body: apiKey}, nil
}

func (h *handler) createApiKey(
	ctx context.Context,
	input *struct{
		Body *CreateApiKeyRequest
	},
) (*CreateApiKeyOutput, error) {
	apiKey, err := h.apiKeysService.CreateApiKey(ctx, input.Body)
	if err != nil {
		return nil, apiKeyHTTPError(err, "Failed to create API Key")
	}

	return &CreateApiKeyOutput{Body: apiKey}, nil
}

func (h *handler) deleteApiKey(
	ctx context.Context,
	input *struct{
		ApiKeyID string `path:"apiKeyId"`
	},
) (*DeleteApiKeyOutput, error) {
	apiKeyID, err := uuid.Parse(input.ApiKeyID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid API Key ID")
	}

	if err := h.apiKeysService.DeleteApiKeyByID(ctx, apiKeyID); err != nil {
		return nil, apiKeyHTTPError(err, "Failed to delete API Key")
	}

	return &DeleteApiKeyOutput{Status: http.StatusOK}, nil
}

func apiKeyHTTPError(err error, fallback string) error {
	if errors.Is(err, ErrApiKeyNotFound) {
		return huma.Error404NotFound("API Key not found")
	}

	return huma.Error500InternalServerError(fallback)
}