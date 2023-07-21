//go:build (all || permissions || resource_secure_file_permissions) && (!exclude_permissions || !resource_secure_file_permissions)
// +build all permissions resource_secure_file_permissions
// +build !exclude_permissions !resource_secure_file_permissions

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

var libraryToken = fmt.Sprintf("Library/%s", projectID)

func TestLibraryPermissions_CreateLibraryToken(t *testing.T) {
	var d *schema.ResourceData
	var token string
	var err error

	d = getLibraryPermissionsResource(t, projectID)
	token, err = createLibraryToken(d, nil)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, libraryToken, token)

	d = getLibraryPermissionsResource(t, "")
	token, err = createLibraryToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func getLibraryPermissionsResource(t *testing.T, projectID string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceLibraryPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	return d
}
