// +build all data_sources data_agent_pools
// +build !exclude_data_sources !exclude_data_agent_pools

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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestDataSourceAgentPools_Read_TestEmptyAgentPoolList(t *testing.T) {
	agentPoolListEmpty := []taskagent.TaskAgentPool{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &config.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	taskAgentClient.
		EXPECT().
		GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{}).
		Return(&agentPoolListEmpty, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataAzureAgentPools().Schema, nil)
	err := dataSourceAgentPoolsRead(resourceData, clients)
	require.Nil(t, err)
	agentPools := resourceData.Get("agent_pools").([]interface{})
	require.NotNil(t, agentPools)
	require.Equal(t, 0, len(agentPools))
}

var dataTestAgentPools = []taskagent.TaskAgentPool{
	{
		Id:            converter.Int(111),
		Name:          converter.String("AgentPool"),
		PoolType:      &taskagent.TaskAgentPoolTypeValues.Automation,
		AutoProvision: converter.Bool(false),
	},
	{
		Id:            converter.Int(65092),
		Name:          converter.String("AgentPool_AutoProvisioned"),
		PoolType:      &taskagent.TaskAgentPoolTypeValues.Automation,
		AutoProvision: converter.Bool(true),
	},
	{
		Id:            converter.Int(650792),
		Name:          converter.String("AgentPool_Deployment"),
		PoolType:      &taskagent.TaskAgentPoolTypeValues.Deployment,
		AutoProvision: converter.Bool(false),
	},
}

func TestDataSourceAgentPools_Read_TestFindAllAgentPools(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &config.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	taskAgentClient.
		EXPECT().
		GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{}).
		Return(&dataTestAgentPools, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataAzureAgentPools().Schema, nil)
	err := dataSourceAgentPoolsRead(resourceData, clients)
	require.Nil(t, err)
	agentPools := resourceData.Get("agent_pools").([]interface{})
	require.NotNil(t, agentPools)
	require.Equal(t, len(dataTestAgentPools), len(agentPools))
}
