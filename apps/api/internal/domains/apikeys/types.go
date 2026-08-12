package apikeys

import (
	"time"

	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

// See https://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type omit *struct{}

type ApiKeyResponse struct {
	*sqlc.ApiKey
	SecretHash omit `json:"secret_hash,omitempty"`
}

type ListApiKeysOutput struct {
	Body []ApiKeyResponse
}

type CreateApiKeyRequest struct {
	Name        string     `json:"name" minLength:"1"`
	Description string     `json:"description" minLength:"1"`
	Role        *sqlc.Role `json:"role"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

type CreateApiKeyOutput struct {
	*ApiKeyResponse
	Secret string `json:"secret"`
}

type DeleteApiKeyOutput struct {
	Status int
}