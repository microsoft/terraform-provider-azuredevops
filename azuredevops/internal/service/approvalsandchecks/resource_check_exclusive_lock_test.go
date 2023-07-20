//go:build (all || resource_check_exclusive_lock) && !exclude_approvalsandchecks
// +build all resource_check_exclusive_lock
// +build !exclude_approvalsandchecks

package approvalsandchecks

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
	"github.com/stretchr/testify/require"
)

var CheckExclusiveLockID = 123456789
var CheckExclusiveLockProjectID = uuid.New().String()
var CheckExclusiveLockTestProjectID = &CheckExclusiveLockProjectID

var CheckExclusiveLockInputs = map[string]interface{}{}

var CheckExclusiveLockSettings = map[string]interface{}{}

var CheckExclusiveLockTest = pipelineschecksextras.CheckConfiguration{
	Id:       &CheckExclusiveLockID,
	Type:     approvalAndCheckType.ExclusiveLock,
	Settings: CheckExclusiveLockSettings,
	Timeout:  converter.ToPtr(20000),
	Resource: &endpointResource,
}

// verifies that the flatten/expand round trip yields the same exclusive lock check
func TestCheckExclusiveLock_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceCheckApproval().Schema, nil)
	flattenExclusiveLock(resourceData, &CheckExclusiveLockTest, CheckExclusiveLockProjectID)

	CheckExclusiveLockAfterRoundTrip, projectID, err := expandExclusiveLock(resourceData)

	require.Equal(t, CheckExclusiveLockTest, *CheckExclusiveLockAfterRoundTrip)
	require.Equal(t, CheckExclusiveLockProjectID, projectID)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestCheckExclusiveLock_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckExclusiveLock()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenExclusiveLock(resourceData, &CheckExclusiveLockTest, CheckExclusiveLockProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.AddCheckConfigurationArgs{Configuration: &CheckExclusiveLockTest, Project: &CheckExclusiveLockProjectID}
	pipelinesCheckClient.
		EXPECT().
		AddCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("AddCheckConfiguration() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "AddCheckConfiguration() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestCheckExclusiveLock_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckExclusiveLock()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenExclusiveLock(resourceData, &CheckExclusiveLockTest, CheckExclusiveLockProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.GetCheckConfigurationArgs{
		Id:      CheckExclusiveLockTest.Id,
		Project: &CheckExclusiveLockProjectID,
		Expand:  converter.ToPtr(pipelineschecksextras.CheckConfigurationExpandParameterValues.Settings),
	}

	pipelinesCheckClient.
		EXPECT().
		GetCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestCheckExclusiveLock_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckExclusiveLock()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenExclusiveLock(resourceData, &CheckExclusiveLockTest, CheckExclusiveLockProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.DeleteCheckConfigurationArgs{
		Id:      CheckExclusiveLockTest.Id,
		Project: &CheckExclusiveLockProjectID,
	}

	pipelinesCheckClient.
		EXPECT().
		DeleteCheckConfiguration(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestCheckExclusiveLock_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckExclusiveLock()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenExclusiveLock(resourceData, &CheckExclusiveLockTest, CheckExclusiveLockProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.UpdateCheckConfigurationArgs{
		Project:       &CheckExclusiveLockProjectID,
		Configuration: &CheckExclusiveLockTest,
		Id:            &CheckExclusiveLockID,
	}

	pipelinesCheckClient.
		EXPECT().
		UpdateCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}
