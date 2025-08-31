package parse

import "github.com/google/uuid"

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