//go:build (all || permissions || resource_workitemtrackingprocess_process_permissions) && (!exclude_permissions || !resource_workitemtrackingprocess_process_permissions)
// +build all permissions resource_workitemtrackingprocess_process_permissions
// +build !exclude_permissions !resource_workitemtrackingprocess_process_permissions

package permissions

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	testParentProcessID      = uuid.MustParse("adcc42ab-9882-485e-a3ed-7678f01f66bc")
	testProcessID            = uuid.MustParse("0aa41603-5857-4155-bdfa-6a0db64d8045")
	inheritedProcessToken    = fmt.Sprintf("$PROCESS:%s:%s:", testParentProcessID.String(), testProcessID.String())
	nonInheritedProcessToken = fmt.Sprintf("$PROCESS:%s:", testProcessID.String())
)

func TestProcessPermissions_CreateProcessToken_InheritedProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	mockClient.EXPECT().GetProcessByItsId(clients.Ctx, workitemtrackingprocess.GetProcessByItsIdArgs{
		ProcessTypeId: &testProcessID,
	}).Return(&workitemtrackingprocess.ProcessInfo{
		ParentProcessTypeId: &testParentProcessID,
	}, nil).Times(1)

	d := getProcessPermissionsResource(t, testProcessID.String())
	token, err := createProcessToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, inheritedProcessToken, token)
}

func TestProcessPermissions_CreateProcessToken_NonInheritedProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

	mockClient.EXPECT().GetProcessByItsId(clients.Ctx, workitemtrackingprocess.GetProcessByItsIdArgs{
		ProcessTypeId: &testProcessID,
	}).Return(&workitemtrackingprocess.ProcessInfo{
		ParentProcessTypeId: nil,
	}, nil).Times(1)

	d := getProcessPermissionsResource(t, testProcessID.String())
	token, err := createProcessToken(d, clients)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	assert.Equal(t, nonInheritedProcessToken, token)
}

func TestProcessPermissions_CreateProcessToken_MissingProcessID(t *testing.T) {
	d := getProcessPermissionsResource(t, "")
	token, err := createProcessToken(d, nil)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func getProcessPermissionsResource(t *testing.T, processID string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceWorkItemTrackingProcessPermissions().Schema, nil)
	if processID != "" {
		d.Set("process_id", processID)
	}
	return d
}
