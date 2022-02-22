//go:build all || utils || securitynamespaces
// +build all utils securitynamespaces

package utils

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/stretchr/testify/assert"
)

func TestCreatePermissionResourceSchema(t *testing.T) {
	schema := CreatePermissionResourceSchema(map[string]*schema.Schema{
		"project_id": {
			Type:         schema.TypeString,
			ValidateFunc: validation.IsUUID,
			Required:     true,
			ForceNew:     true,
		},
		"repository_id": {
			Type:         schema.TypeString,
			ValidateFunc: validation.IsUUID,
			Optional:     true,
			ForceNew:     true,
		},
		"branch_name": {
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotWhiteSpace,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"repository_id"},
		},
	})

	requiredFields := []string{
		"principal",
		"replace",
		"permissions",
		"project_id",
		"repository_id",
		"branch_name",
	}

	for _, field := range requiredFields {
		_, ok := schema[field]
		assert.True(t, ok, fmt.Sprintf("Schema should contain a field [%s]", field))
	}
}
