//go:build (all || permissions || resource_project_permissions) && (!exclude_permissions || !resource_project_permissions)
// +build all permissions resource_project_permissions
// +build !exclude_permissions !resource_project_permissions

package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

/**
 * Begin unit tests
 */

var (
	projectID    = "9083e944-8e9e-405e-960a-c80180aa71e6"
	projectToken = fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID)
)

func TestProjectPermissions_CreateProjectToken(t *testing.T) {
	var d *schema.ResourceData
	var token string
	var err error

	d = getProjecPermissionsResource(t, projectID)
	token, err = createProjectToken(d, nil)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, projectToken, token)

	d = getProjecPermissionsResource(t, "")
	token, err = createProjectToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func getProjecPermissionsResource(t *testing.T, projectID string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceProjectPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	return d
}
