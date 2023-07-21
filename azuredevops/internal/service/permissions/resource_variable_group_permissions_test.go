//go:build (all || permissions || resource_variable_group_permissions) && (!exclude_permissions || !resource_variable_group_permissions)
// +build all permissions resource_variable_group_permissions
// +build !exclude_permissions !resource_variable_group_permissions

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

var variableGroupID = "5"
var variableGroupToken = fmt.Sprintf("Library/%s/VariableGroup/%s", projectID, variableGroupID)

func TestVariableGroupsPermissions_CreateVariableGroupToken(t *testing.T) {
	var d *schema.ResourceData
	var token string
	var err error

	d = getVariableGroupPermissionsResource(t, projectID, variableGroupID)
	token, err = createVariableGroupToken(d, nil)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, variableGroupToken, token)

	d = getVariableGroupPermissionsResource(t, "", "")
	token, err = createVariableGroupToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func getVariableGroupPermissionsResource(t *testing.T, projectID string, variableGroupID string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceVariableGroupPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	if variableGroupID != "" {
		d.Set("variable_group_id", variableGroupID)
	}
	return d
}
