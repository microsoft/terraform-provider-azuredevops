package testhelper

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Validates that the resource data conforms to the schema
func ValidateResourceData(t *testing.T, d *schema.ResourceData, r *schema.Resource) {
	for name, sch := range r.Schema {
		val, exists := d.GetOk(name)

		if sch.Required {
			if !exists {
				t.Fatalf("Missing required attribute: %s", name)
			}
		}

		if sch.ValidateDiagFunc != nil {
			diags := sch.ValidateDiagFunc(val, cty.Path{})
			if diags.HasError() {
				t.Fatalf("Validation error for attribute %q: %v", name, diags)
			}
		}
	}
}
