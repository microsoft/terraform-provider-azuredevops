//go:build (all || policy) && !exclude_policy
// +build all policy
// +build !exclude_policy

package branch

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	projectID  = uuid.New().String()
	randomUUID = uuid.New()
	testPolicy = &policy.PolicyConfiguration{
		Id:         converter.Int(1),
		IsEnabled:  converter.Bool(true),
		IsBlocking: converter.Bool(true),
		Type: &policy.PolicyTypeRef{
			Id: &randomUUID,
		},
		Settings: map[string]interface{}{
			"scope": []map[string]interface{}{
				{
					"repositoryId": "test-repo-id",
					"refName":      "test-ref-name",
					"matchKind":    "test-match-kind",
				},
			},
		},
	}
)

var testResource = genBasePolicyResource(&policyCrudArgs{
	baseFlattenFunc,
	baseExpandFunc,
	randomUUID,
})

func getFlattenedResourceData(t *testing.T) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, testResource.Schema, nil)
	err := baseFlattenFunc(resourceData, testPolicy, &projectID)
	require.Nil(t, err)
	return resourceData
}

// verifies that the flatten/expand round trip path produces repeatable results
func TestBranchPolicyCRUD_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := getFlattenedResourceData(t)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))
	expandedPolicy, expandedProjectID, err := baseExpandFunc(resourceData, randomUUID)
	require.Nil(t, err)

	require.Equal(t, testPolicy, expandedPolicy)
	require.Equal(t, projectID, *expandedProjectID)
}

// verifies that CREATE failures are not swallowed
func TestBranchPolicyCRUD_CreateError_NotSwallowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := getFlattenedResourceData(t)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))

	policyClient := azdosdkmocks.NewMockPolicyClient(ctrl)
	clients := &client.AggregatedClient{PolicyClient: policyClient, Ctx: context.Background()}

	expectedArgs := policy.CreatePolicyConfigurationArgs{
		Configuration: testPolicy,
		Project:       &projectID,
	}

	policyClient.
		EXPECT().
		CreatePolicyConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreatePolicyConfiguration() Failed")).
		Times(1)

	err := testResource.Create(resourceData, clients)
	require.Regexp(t, ".*CreatePolicyConfiguration\\(\\) Failed$", err.Error())
}

// verifies that READ failures are not swallowed
func TestBranchPolicyCRUD_ReadError_NotSwallowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := getFlattenedResourceData(t)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))
	policyClient := azdosdkmocks.NewMockPolicyClient(ctrl)
	clients := &client.AggregatedClient{PolicyClient: policyClient, Ctx: context.Background()}

	expectedArgs := policy.GetPolicyConfigurationArgs{
		ConfigurationId: testPolicy.Id,
		Project:         &projectID,
	}

	policyClient.
		EXPECT().
		GetPolicyConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetPolicyConfiguration() Failed")).
		Times(1)

	err := testResource.Read(resourceData, clients)
	require.Regexp(t, ".*GetPolicyConfiguration\\(\\) Failed$", err.Error())
}

// verifies that UDPATE failures are not swallowed
func TestBranchPolicyCRUD_UpdateError_NotSwallowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := getFlattenedResourceData(t)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))

	policyClient := azdosdkmocks.NewMockPolicyClient(ctrl)
	clients := &client.AggregatedClient{PolicyClient: policyClient, Ctx: context.Background()}

	expectedArgs := policy.UpdatePolicyConfigurationArgs{
		ConfigurationId: testPolicy.Id,
		Configuration:   testPolicy,
		Project:         &projectID,
	}

	policyClient.
		EXPECT().
		UpdatePolicyConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdatePolicyConfiguration() Failed")).
		Times(1)

	err := testResource.Update(resourceData, clients)
	require.Regexp(t, ".*UpdatePolicyConfiguration\\(\\) Failed$", err.Error())
}

// verifies that DELETE failures are not swallowed
func TestBranchPolicyCRUD_DeleteError_NotSwallowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := getFlattenedResourceData(t)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))

	policyClient := azdosdkmocks.NewMockPolicyClient(ctrl)
	clients := &client.AggregatedClient{PolicyClient: policyClient, Ctx: context.Background()}

	expectedArgs := policy.DeletePolicyConfigurationArgs{
		ConfigurationId: testPolicy.Id,
		Project:         &projectID,
	}

	policyClient.
		EXPECT().
		DeletePolicyConfiguration(clients.Ctx, expectedArgs).
		Return(errors.New("DeletePolicyConfiguration() Failed")).
		Times(1)

	err := testResource.Delete(resourceData, clients)
	require.Regexp(t, ".*DeletePolicyConfiguration\\(\\) Failed$", err.Error())
}
