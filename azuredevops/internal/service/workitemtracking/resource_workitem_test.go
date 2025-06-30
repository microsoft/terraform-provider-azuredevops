//go:build (all || resource_workitem) && !resource_workitem
// +build all resource_workitem
// +build !resource_workitem

package workitemtracking

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestWorkItem_GetWorkItem(t *testing.T) {
	r := ResourceWorkItem()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
	input := map[string]interface{}{
		"System.State":         "To Do",
		"System.Title":         "TestTitle",
		"Custom.SomeName":      "SomeValue",
		"Custom.foo":           "bar",
		"System.Tags":          "tag1; tag2=value",
		"System.AreaPath":      "SomePath",
		"System.IterationPath": "SomeIterationPath",
	}
	flattenFields(d, &input)

	require.Equal(t, "TestTitle", d.Get("title").(string))
	require.Equal(t, "To Do", d.Get("state").(string))
	require.Equal(t, "SomePath", d.Get("area_path").(string))
	require.Equal(t, "SomeIterationPath", d.Get("iteration_path").(string))
	require.Equal(t, 2, len(d.Get("tags").(*schema.Set).List()))

	custom_fields := d.Get("custom_fields").(map[string]interface{})
	require.Equal(t, "SomeValue", custom_fields["SomeName"].(string))
	require.Equal(t, "bar", custom_fields["foo"].(string))
}
