package testhelper

import "github.com/google/uuid"

// CreateUUID creates a new UUID
func CreateUUID() *uuid.UUID {
	val := uuid.New()
	return &val
}

// ToUUID creates a UUID from a string value
func ToUUID(szUUID string) *uuid.UUID {
	val := uuid.MustParse(szUUID)
	return &val
}
