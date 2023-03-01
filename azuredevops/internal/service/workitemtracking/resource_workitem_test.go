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
		"System.State":    "To Do",
		"System.Title":    "TestTitle",
		"Custom.SomeName": "SomeValue",
		"Custom.foo":      "bar",
		"System.Tags":     "tag1; tag2=value",
	}
	flattenFields(d, &input)

	require.Equal(t, "TestTitle", d.Get("title").(string))
	require.Equal(t, "To Do", d.Get("state").(string))
	require.Equal(t, "tag1", d.Get("tags").([]interface{})[0].(string))
	require.Equal(t, "tag2=value", d.Get("tags").([]interface{})[1].(string))

	custom_fields := d.Get("custom_fields").(map[string]interface{})
	require.Equal(t, "SomeValue", custom_fields["SomeName"].(string))
	require.Equal(t, "bar", custom_fields["foo"].(string))

}
