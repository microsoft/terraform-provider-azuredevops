package acceptance

import (
	"regexp"
)

func ResourceExistError() *regexp.Regexp {
	return regexp.MustCompile(`Resource already exists`)
}
