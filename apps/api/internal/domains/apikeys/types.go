package apikeys

import (
	"strings"
	"time"

	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type ListApiKeysOutput struct {
	Body []sqlc.ListApiKeysRow
}

type GetApiKeyOutput struct {
	Body *sqlc.GetApiKeyByIdRow
}

type CreateApiKeyRequest struct {
	Name        string     `json:"name" minLength:"1"`
	Description *string    `json:"description,omitempty" minLength:"1"`
	Role        sqlc.Role  `json:"role"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

func (r *CreateApiKeyRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrApiKeyNameRequired
	}
	if r.Description != nil && strings.TrimSpace(*r.Description) == "" {
		return ErrApiKeyDescriptionRequired
	}
	return nil
}

type CreateApiKeyResponse struct {
	*sqlc.CreateApiKeyRow
	Secret string `json:"secret"`
}

type CreateApiKeyOutput struct {
	Body *CreateApiKeyResponse
}

type DeleteApiKeyOutput struct {
	Status int
}