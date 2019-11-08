package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-uuid"
)

// UUIDRegExp is a regex helper.
var UUIDRegExp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

// UUID parses a UUID, returning warnings and errors.
func UUID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := uuid.ParseUUID(v); err != nil {
		errors = append(errors, fmt.Errorf("%q isn't a valid UUID (%q): %+v", k, v, err))
	}

	return warnings, errors
}

// UUIDOrEmpty parses a UUID, returning nil for warnings if i is empty.
func UUIDOrEmpty(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" {
		return
	}

	return UUID(i, k)
}
