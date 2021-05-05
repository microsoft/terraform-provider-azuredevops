// +build all permissions resource_release_permission_permissions
// +build !exclude_permissions !resource_release_permission_permissions

package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

/**
 * Begin unit tests
 */

var releaseID = "9083e944-8e9e-405e-960a-c80180aa71e6"

func TestReleasePermissions_createReleaseToken(t *testing.T) {
	var d *schema.ResourceData
	var token string
	var err error

	d = getReleasePermissionsResource(t, releaseID)
	token, err = createReleaseToken(d, nil)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, token)

	d = getReleasePermissionsResource(t, "")
	token, err = createReleaseToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func getReleasePermissionsResource(t *testing.T, releaseID string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceReleasePermissions().Schema, nil)
	if releaseID != "" {
		d.Set("project_id", releaseID)
	}
	return d
}
