package converter

import "strings"

// String Get a pointer to a string
func String(value string) *string {
	if strings.EqualFold(value, "") {
		return nil
	}
	return &value
}

// Bool Get a pointer to a boolean value
func Bool(value bool) *bool {
	return &value
}

// Int Get a pointer to an integer value
func Int(value int) *int {
	return &value
}

// UInt64 Get a pointer to an uint64 value
func UInt64(value uint64) *uint64 {
	return &value
}

// ToString Given a pointer return its value, or a default value of the poitner is nil
func ToString(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}

	return defaultValue
}

// ToBool Given a pointer return its value, or a default value of the pointer is nil
func ToBool(value *bool, defaultValue bool) bool {
	if value != nil {
		return *value
	}

	return defaultValue
}
