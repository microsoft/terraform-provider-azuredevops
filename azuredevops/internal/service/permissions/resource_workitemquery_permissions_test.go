// +build all permissions resource_workitemquery_permissions
// +build !exclude_permissions !resource_workitemquery_permissions

package permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestWorkItemQueryPermissions_CreateWorkItemQueryToken(t *testing.T) {
}

func getWorkItemQueryPermissionsResource(t *testing.T, projectID string, path string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceWorkItemQueryPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	if path != "" {
		d.Set("path", path)
	}
	return d
}
