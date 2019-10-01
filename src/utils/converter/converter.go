package converter

// String Get a pointer to a string
func String(value string) *string {
	return &value
}

// Bool Get a pointer to a boolean value
func Bool(value bool) *bool {
	return &value
}

// ToString Given a pointer return its value, or a default value of the poitner is nil
func ToString(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}

	return defaultValue
}
