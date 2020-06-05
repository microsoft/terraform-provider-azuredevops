package testhelper

import "github.com/google/uuid"

// CreateUUID creates a new UUID
func CreateUUID() *uuid.UUID {
	val := uuid.New()
	return &val
}
