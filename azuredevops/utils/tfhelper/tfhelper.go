package tfhelper

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// DiffFuncSupressCaseSensitivity Suppress case sensitivity when comparing string values
func DiffFuncSupressCaseSensitivity(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}
