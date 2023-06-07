//go:build (all || resource_check_branch_control) && !exclude_approvalsandchecks
// +build all resource_check_branch_control
// +build !exclude_approvalsandchecks

package approvalsandchecks

import (
	"context"
	"errors"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
	"github.com/stretchr/testify/require"
)

var branchControlCheckID = 123456789
var branchControlEndpointID = uuid.New().String()
var branchControlCheckProjectID = uuid.New().String()
var branchControlCheckTestProjectID = &branchControlCheckProjectID

var endpointType = "endpoint"
var endpointResource = pipelineschecksextras.Resource{
	Id:   &branchControlEndpointID,
	Type: &endpointType,
}

var branchControlInputs = map[string]interface{}{
	"allowedBranches":          "refs/heads/releases/*",
	"ensureProtectionOfBranch": strconv.FormatBool(false),
	"allowUnknownStatusBranch": strconv.FormatBool(true),
}

var branchControlCheckSettings = map[string]interface{}{
	"definitionRef": evaluateBranchProtectionDef,
	"displayName":   "Test Branch Control",
	"inputs":        branchControlInputs,
}

var branchControlCheckTest = pipelineschecksextras.CheckConfiguration{
	Id:       &branchControlCheckID,
	Type:     approvalAndCheckType.BranchProtection,
	Settings: branchControlCheckSettings,
	Resource: &endpointResource,
}

// verifies that the flatten/expand round trip yields the same branch control
func TestCheckBranchControl_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceCheckBranchControl().Schema, nil)
	flattenBranchControlCheck(resourceData, &branchControlCheckTest, branchControlCheckProjectID)

	branchControlCheckAfterRoundTrip, projectID, err := expandBranchControlCheck(resourceData)

	require.Equal(t, branchControlCheckTest, *branchControlCheckAfterRoundTrip)
	require.Equal(t, branchControlCheckProjectID, projectID)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestCheckBranchControl_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBranchControl()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenBranchControlCheck(resourceData, &branchControlCheckTest, branchControlCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.AddCheckConfigurationArgs{Configuration: &branchControlCheckTest, Project: &branchControlCheckProjectID}
	pipelinesChecksClient.
		EXPECT().
		AddCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("AddCheckConfiguration() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "AddCheckConfiguration() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestCheckBranchControl_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBranchControl()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenBranchControlCheck(resourceData, &branchControlCheckTest, branchControlCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.GetCheckConfigurationArgs{
		Id:      branchControlCheckTest.Id,
		Project: &branchControlCheckProjectID,
		Expand:  converter.ToPtr(pipelineschecksextras.CheckConfigurationExpandParameterValues.Settings),
	}

	pipelinesChecksClient.
		EXPECT().
		GetCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestCheckBranchControl_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBranchControl()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenBranchControlCheck(resourceData, &branchControlCheckTest, branchControlCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.DeleteCheckConfigurationArgs{
		Id:      branchControlCheckTest.Id,
		Project: &branchControlCheckProjectID,
	}

	pipelinesChecksClient.
		EXPECT().
		DeleteCheckConfiguration(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestCheckBranchControl_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBranchControl()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenBranchControlCheck(resourceData, &branchControlCheckTest, branchControlCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.UpdateCheckConfigurationArgs{
		Project:       &branchControlCheckProjectID,
		Configuration: &branchControlCheckTest,
		Id:            &branchControlCheckID,
	}

	pipelinesChecksClient.
		EXPECT().
		UpdateCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}
