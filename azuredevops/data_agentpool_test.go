// +build all core data_projects

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/stretchr/testify/require"
)

func TestDataSourceAgentPool_Read_TestAgentPoolNotFound(t *testing.T) {
	agentPoolListEmpty := []taskagent.TaskAgentPool{}
	name := "nonexistentAgentPool"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &config.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	taskAgentClient.
		EXPECT().
		GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{
			PoolName: &name,
		}).
		Return(&agentPoolListEmpty, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataAzureAgentPool().Schema, nil)
	resourceData.Set("name", &name)
	err := dataSourceAgentPoolRead(resourceData, clients)
	require.Contains(t, err.Error(), "Unable to find agent pool")
}

func TestDataSourceAgentPool_Read_TestMultipleAgentPoolsFound(t *testing.T) {
	agentPoolListEmpty := []taskagent.TaskAgentPool{{}, {}}
	name := "nonexistentAgentPool"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &config.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	taskAgentClient.
		EXPECT().
		GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{
			PoolName: &name,
		}).
		Return(&agentPoolListEmpty, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataAzureAgentPool().Schema, nil)
	resourceData.Set("name", &name)
	err := dataSourceAgentPoolRead(resourceData, clients)
	require.Contains(t, err.Error(), "Found multiple agent pools for name")
}
