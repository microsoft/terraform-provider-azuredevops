package validate

import (
	"fmt"
	"regexp"
)

func EnvironmentName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile(`^[^,^"^/^\\^\[^\]^:^|^<^>^+^=^;^?^*]{0,127}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("test: %s, Environment name '%s' is not valid. "+
			"A valid name is less than 128 characters in length and does not contain the following "+
			"characters: ',', '\"', '/', '\\', '[', ']', ':', '|', '<', '>', '+', '=', ';', '?', '*', '.'", k, v))
	}
	return warnings, errors
}
