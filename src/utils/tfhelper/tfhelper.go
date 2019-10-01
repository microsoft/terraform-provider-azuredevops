package tfhelper

import (
	"github.com/hashicorp/terraform/helper/schema"
)

import (
	"strings"
)

// DiffFuncSupressCaseSensitivity Suppress case sensitivity when comparing string values
func DiffFuncSupressCaseSensitivity(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}
