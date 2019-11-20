package validate

import (
	"fmt"
	"regexp"
)

// InvalidWindowsFilePathRegExp is a regex helper.
var InvalidWindowsFilePathRegExp = regexp.MustCompile("[<>|:$@\"/%+*?]")

// FilePath validates that the string does not contain characters (equal to [<>|:$@\"/%+*?])
func FilePath(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	p := InvalidWindowsFilePathRegExp.MatchString(v)
	if p {
		errors = append(errors, fmt.Errorf("<>|:$@\"/%%+*? are not allowed"))
		return
	}

	return warnings, errors
}

// FilePathOrEmpty allows empty string otherwise validates that the string does not contain characters (equal to [<>|:$@\"/%+*?])
func FilePathOrEmpty(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" {
		return
	}

	return FilePath(i, k)
}
