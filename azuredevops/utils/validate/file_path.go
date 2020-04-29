package validate

import (
	"fmt"
	"regexp"
)

// InvalidWindowsPathRegExp is a regex helper.
var InvalidWindowsPathRegExp = regexp.MustCompile("[<>|:$@\"/%+*?]")

// Path validates that the string does not contain characters (equal to [<>|:$@\"/%+*?])
func Path(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("path can not be empty"))
	}

	if len(v) >= 1 && v[:1] != `\` {
		errors = append(errors, fmt.Errorf("path must start with backslash"))
	}

	p := InvalidWindowsPathRegExp.MatchString(v)
	if p {
		errors = append(errors, fmt.Errorf("<>|:$@\"/%%+*? are not allowed in path"))
		return
	}

	return warnings, errors
}
