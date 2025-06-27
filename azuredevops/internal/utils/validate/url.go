package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Url(i interface{}, key string) (_ []string, errors []error) {
	url, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
		return
	}
	if strings.HasSuffix(url, "/") {
		errors = append(errors, fmt.Errorf("%q should not end with slash, got %q.", key, url))
		return
	}
	return validation.IsURLWithHTTPorHTTPS(url, key)
}
