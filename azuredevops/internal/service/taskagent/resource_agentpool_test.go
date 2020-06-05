// +build all resource_agentpool
// +build !exclude_resource_agentpool

package taskagent

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var testAgentPoolID = rand.Intn(100)

var testAgentPool = taskagent.TaskAgentPool{
	Id:            &testAgentPoolID,
	Name:          converter.String("Name"),
	PoolType:      &taskagent.TaskAgentPoolTypeValues.Automation,
	AutoProvision: converter.Bool(false),
}

// verifies that the flatten/expand round trip yields the same agent pool definition
func TestAgentPool_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceAgentPool().Schema, nil)
	flattenAzureAgentPool(resourceData, &testAgentPool)

	agentPoolAfterRoundTrip, err := expandAgentPool(resourceData, true)
	require.Nil(t, err)
	require.Equal(t, testAgentPool, *agentPoolAfterRoundTrip)
}

// verifies that the create operation is considered failed if the API call fails.
func TestAgentPool_CreateAgentPool_DoesNotSwallowErrorFromFailedAddAgentCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedProjectCreateArgs := taskagent.AddAgentPoolArgs{
		Pool: &testAgentPool,
	}

	taskAgentClient.
		EXPECT().
		AddAgentPool(clients.Ctx, expectedProjectCreateArgs).
		Return(nil, errors.New("AddAgentPool() Failed")).
		Times(1)

	newTaskAgentPool, err := createAzureAgentPool(clients, &testAgentPool)
	require.Nil(t, newTaskAgentPool)
	require.Equal(t, "AddAgentPool() Failed", err.Error())
}

func TestAgentPool_DeleteAgentPool_ReturnsErrorIfIdReadFails(t *testing.T) {
	client := &client.AggregatedClient{}

	resourceData := schema.TestResourceDataRaw(t, ResourceAgentPool().Schema, nil)
	flattenAzureAgentPool(resourceData, &testAgentPool)
	resourceData.SetId("")

	err := resourceAzureAgentPoolDelete(resourceData, client)
	require.Equal(t, "Error getting agent pool Id: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
}

func TestAgentPool_UpdateAgentPool_ReturnsErrorIfIdReadFails(t *testing.T) {
	client := &client.AggregatedClient{}

	resourceData := schema.TestResourceDataRaw(t, ResourceAgentPool().Schema, nil)
	flattenAzureAgentPool(resourceData, &testAgentPool)
	resourceData.SetId("")

	err := resourceAzureAgentPoolUpdate(resourceData, client)
	require.Equal(t, "Error converting terraform data model to AzDO agent pool reference: Error getting agent pool Id: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
}

func TestAgentPool_UpdateAgentPool_UpdateAndRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	agentToUpdate := taskagent.TaskAgentPool{
		Id:            &testAgentPoolID,
		Name:          converter.String("Foo"),
		PoolType:      &taskagent.TaskAgentPoolTypeValues.Deployment,
		AutoProvision: converter.Bool(true),
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceAgentPool().Schema, nil)
	flattenAzureAgentPool(resourceData, &agentToUpdate)

	taskAgentClient.
		EXPECT().
		UpdateAgentPool(clients.Ctx, taskagent.UpdateAgentPoolArgs{
			PoolId: &testAgentPoolID,
			Pool: &taskagent.TaskAgentPool{
				Name:          agentToUpdate.Name,
				PoolType:      agentToUpdate.PoolType,
				AutoProvision: agentToUpdate.AutoProvision,
			},
		}).
		Return(&agentToUpdate, nil).
		Times(1)

	taskAgentClient.
		EXPECT().
		GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{
			PoolId: &testAgentPoolID,
		}).
		Return(&agentToUpdate, nil).
		Times(1)

	err := resourceAzureAgentPoolUpdate(resourceData, clients)
	require.Nil(t, err)

	updatedTaskAgent, _ := expandAgentPool(resourceData, false)
	require.Equal(t, agentToUpdate.Id, updatedTaskAgent.Id)
	require.Equal(t, agentToUpdate.Name, updatedTaskAgent.Name)
	require.Equal(t, agentToUpdate.PoolType, updatedTaskAgent.PoolType)
	require.Equal(t, agentToUpdate.AutoProvision, updatedTaskAgent.AutoProvision)
}

// validates supported pool types are allowed by the schema
func TestAgentPoolDefinition_PoolTypeIsCorrect(t *testing.T) {
	validPoolTypes := []string{
		string(taskagent.TaskAgentPoolTypeValues.Automation),
		string(taskagent.TaskAgentPoolTypeValues.Deployment),
	}
	poolTypeSchema := ResourceAgentPool().Schema["pool_type"]

	for _, repoType := range validPoolTypes {
		_, errors := poolTypeSchema.ValidateFunc(repoType, "")
		require.Equal(t, 0, len(errors), "Agent pool type unexpectedly did not pass validation")
	}
}

// validates invalid pool types are rejected by the schema
func TestAgentPoolDefinition_WhenPoolTypeIsNotCorrect_ReturnsError(t *testing.T) {
	invalidPoolTypes := []string{"", "unknown"}
	poolTypeSchema := ResourceAgentPool().Schema["pool_type"]

	for _, poolType := range invalidPoolTypes {
		_, errors := poolTypeSchema.ValidateFunc(poolType, "pool_type")
		expectedError := fmt.Sprintf("expected pool_type to be one of [automation deployment], got %s", poolType)
		require.Equal(t, 1, len(errors), "Agent pool type %v unexpectedly passed validation", poolType)
		require.Equal(t, expectedError, errors[0].Error())
	}
}
