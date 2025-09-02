package ptr

import "github.com/google/uuid"

// Takes a UUID and returns a pointer to that UUID
func UUIDToPtr(id uuid.UUID) *uuid.UUID {
	return &id
}
