//go:build (all || permissions || resource_workitemtrackingprocess_process_permissions) && (!exclude_permissions || !resource_workitemtrackingprocess_process_permissions)
// +build all permissions resource_workitemtrackingprocess_process_permissions
// +build !exclude_permissions !resource_workitemtrackingprocess_process_permissions

package permissions

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

var (
	parentProcessID = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
	processID       = "0aa41603-5857-4155-bdfa-6a0db64d8045"
	processToken    = fmt.Sprintf("$PROCESS:%s:%s:", parentProcessID, processID)
)

func TestProcessPermissions_CreateProcessToken(t *testing.T) {
	var d *schema.ResourceData
	var token string
	var err error

	d = getProcessPermissionsResource(t, parentProcessID, processID)
	token, err = createProcessToken(d, nil)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, processToken, token)

	d = getProcessPermissionsResource(t, "", processID)
	token, err = createProcessToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)

	d = getProcessPermissionsResource(t, parentProcessID, "")
	token, err = createProcessToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func getProcessPermissionsResource(t *testing.T, parentProcessID string, processID string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceWorkItemTrackingProcessPermissions().Schema, nil)
	if parentProcessID != "" {
		d.Set("parent_process_id", parentProcessID)
	}
	if processID != "" {
		d.Set("process_id", processID)
	}
	return d
}
