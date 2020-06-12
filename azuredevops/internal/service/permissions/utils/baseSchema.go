package utils

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// CreatePermissionResourceSchema creates a resources schema for a Terraform permission resource
func CreatePermissionResourceSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	baseSchema := map[string]*schema.Schema{
		"principal": {
			Type:         schema.TypeString,
			ValidateFunc: validate.NoEmptyStrings,
			Required:     true,
			ForceNew:     true,
		},
		"replace": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true, // when set to false (merge mode), a permission of Allow or Deny CANNOT be replaced with NotSet
		},
		"permissions": {
			// Unable to define a validation function, because the
			// keys and values can only be validated with an initialized
			// security client as we must load the security namespace
			// definition and the available permission settings, and a validation
			// function in Terraform only receives the parameter name and the
			// current value as argument
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}

	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}
