// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthUserRole string

const (
	AuthUserRoleUser      AuthUserRole = "user"
	AuthUserRoleSuperuser AuthUserRole = "superuser"
)

func (e *AuthUserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AuthUserRole(s)
	case string:
		*e = AuthUserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for AuthUserRole: %T", src)
	}
	return nil
}

type NullAuthUserRole struct {
	AuthUserRole AuthUserRole `json:"auth_user_role"`
	Valid        bool         `json:"valid"` // Valid is true if AuthUserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAuthUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.AuthUserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AuthUserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAuthUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AuthUserRole), nil
}

type AuthAccount struct {
	ID                    uuid.UUID  `json:"id"`
	UserID                uuid.UUID  `json:"user_id"`
	ProviderID            string     `json:"provider_id"`
	AccountID             string     `json:"account_id"`
	HashedPassword        *string    `json:"hashed_password"`
	AccessToken           *string    `json:"access_token"`
	RefreshToken          *string    `json:"refresh_token"`
	IDToken               *string    `json:"id_token"`
	AccessTokenExpiresAt  *time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at"`
	Scope                 *string    `json:"scope"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type AuthSession struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	IpAddress  *string   `json:"ip_address"`
	UserAgent  *string   `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

type AuthUser struct {
	ID            uuid.UUID    `json:"id"`
	Name          string       `json:"name"`
	Email         *string      `json:"email"`
	EmailVerified bool         `json:"email_verified"`
	Onboarded     bool         `json:"onboarded"`
	Image         *string      `json:"image"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Role          AuthUserRole `json:"role"`
}
