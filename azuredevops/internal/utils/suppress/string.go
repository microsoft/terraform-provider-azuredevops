package suppress

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CaseDifference reports whether old and new, interpreted as UTF-8 strings,
// are equal under Unicode case-folding.
func CaseDifference(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
