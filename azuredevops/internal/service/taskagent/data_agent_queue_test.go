//go:build all || core || data_projects
// +build all core data_projects

package taskagent

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/require"
)

func TestDataSourceAgentQueue_Read_TestAgentQueueNotFound(t *testing.T) {
	agentQueueListEmpty := []taskagent.TaskAgentQueue{}
	name := "nonexistentAgentQueue"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	taskAgentClient.
		EXPECT().
		GetAgentQueues(clients.Ctx, taskagent.GetAgentQueuesArgs{
			Project:   &agentQueueProject,
			QueueName: &name,
		}).
		Return(&agentQueueListEmpty, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataAgentQueue().Schema, nil)
	resourceData.Set("name", &name)
	resourceData.Set("project_id", &agentQueueProject)
	err := dataSourceAgentQueueRead(resourceData, clients)
	require.Contains(t, err.Error(), "Unable to find agent queue")
}

func TestDataSourceAgentQueue_Read_TestMultipleAgentQueuesFound(t *testing.T) {
	agentQueueList := []taskagent.TaskAgentQueue{{}, {}}
	name := "multipleQueues"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	taskAgentClient.
		EXPECT().
		GetAgentQueues(clients.Ctx, taskagent.GetAgentQueuesArgs{
			Project:   &agentQueueProject,
			QueueName: &name,
		}).
		Return(&agentQueueList, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataAgentQueue().Schema, nil)
	resourceData.Set("name", &name)
	resourceData.Set("project_id", &agentQueueProject)
	err := dataSourceAgentQueueRead(resourceData, clients)
	require.Contains(t, err.Error(), "Found multiple agent queues for name")
}
