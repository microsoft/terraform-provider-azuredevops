//go:build (all || resource_workitemtracking_field) && !resource_workitemtracking_field
// +build all resource_workitemtracking_field
// +build !resource_workitemtracking_field

package workitemtracking

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestResourceField_UpdateIsLockedAndIsDeleted_Error(t *testing.T) {
	r := ResourceField()

	// Create state with initial values
	state := &terraform.InstanceState{
		ID: "Custom.TestField",
		Attributes: map[string]string{
			"name":           "TestField",
			"reference_name": "Custom.TestField",
			"type":           "string",
			"is_locked":      "false",
			"is_deleted":     "false",
		},
	}

	// Create config with changed values for both is_locked and is_deleted
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":           "TestField",
		"reference_name": "Custom.TestField",
		"type":           "string",
		"is_locked":      true,
		"is_deleted":     true,
	})

	sm := schema.InternalMap(r.Schema)
	diff, err := sm.Diff(context.Background(), state, config, nil, nil, true)
	require.NoError(t, err)

	d, err := sm.Data(state, diff)
	require.NoError(t, err)
	d.SetId("Custom.TestField")

	diags := resourceFieldUpdate(context.Background(), d, nil)

	require.Len(t, diags, 1)
	require.Contains(t, diags[0].Summary, "cannot update is_locked and is_deleted at the same time")
}
