package parse

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

// Parses a string ptr to a UUID or nil if failed to parse
func ParseUUIDOrNil(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}

// Parses a string pointer to a pointer or nil
func ParseStrToPtr(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

func ParseGetEventScopeType(s string) (sqlc.GetEventScopeType, error) {
	enumType := sqlc.GetEventScopeType(s)

	// an empty string should default to published
	if s == "" {
		return sqlc.GetEventScopeTypePublished, nil
	}

	// sqlc generates constants for each enum value. We check if the input matches one of the valid, known constants.
	switch enumType {
	case sqlc.GetEventScopeTypePublished, sqlc.GetEventScopeTypeScoped, sqlc.GetEventScopeTypeAll:
		return enumType, nil
	default:
		// The input string is not a valid enum value.
		return "", fmt.Errorf("'%s' is not a valid GetEventScopeType", s)
	}
}
