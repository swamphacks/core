package apikeys

import (
	"context"

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
	return nil, nil
}

func (h *handler) getApiKey(
	ctx context.Context,
	input *struct{
		ApiKeyId string `path:"apiKeyId"`
	},
) (*ApiKeyResponse, error) {
	return nil, nil
}

func (h *handler) createApiKey(
	ctx context.Context,
	input *struct{
		Body *CreateApiKeyRequest
	},
) (*CreateApiKeyOutput, error) {
	return nil, nil
}

func (h *handler) deleteApiKey(
	ctx context.Context,
	input *struct{
		ApiKeyId string `path:"apiKeyId"`
	},
) (*DeleteApiKeyOutput, error) {
	return nil, nil
}