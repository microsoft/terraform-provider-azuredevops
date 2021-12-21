//go:build all || resource_agent_queue
// +build all resource_agent_queue

package taskagent

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var agentQueueProject = "project"
var agentQueuePoolID = 100
var agentQueuePoolName = "foo-pool"
var agentQueueID = 200

// If the pool lookup fails, an error should be reported
func TestAgentQueue_DoesNotSwallowPoolLookupErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := generateResourceData(t, &agentQueueProject, &agentQueuePoolID, nil)
	agentClient, clients := generateMocks(ctrl)

	agentClient.
		EXPECT().
		GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{
			PoolId: &agentQueuePoolID,
		}).
		Return(nil, errors.New("GetAgentPool() Failed"))

	err := resourceAgentQueueCreate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "GetAgentPool() Failed")
}

// If the queue create fails, an error should be reported
func TestAgentQueue_DoesNotSwallowQueueCreateErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := generateResourceData(t, &agentQueueProject, &agentQueuePoolID, nil)
	agentClient, clients := generateMocks(ctrl)

	pool := &taskagent.TaskAgentPool{Id: &agentQueuePoolID, Name: &agentQueuePoolName}
	agentClient.
		EXPECT().
		GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{
			PoolId: &agentQueuePoolID,
		}).
		Return(pool, nil)

	agentClient.
		EXPECT().
		AddAgentQueue(clients.Ctx, taskagent.AddAgentQueueArgs{
			Queue: &taskagent.TaskAgentQueue{
				Name: &agentQueuePoolName,
				Pool: &taskagent.TaskAgentPoolReference{
					Id: &agentQueuePoolID,
				},
			},
			Project:            &agentQueueProject,
			AuthorizePipelines: converter.Bool(false),
		}).
		Return(nil, errors.New("AddAgentQueue() Failed"))

	err := resourceAgentQueueCreate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "AddAgentQueue() Failed")
}

// If a read fails, an error should be reported
func TestAgentQueue_DoesNotSwallowReadErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := generateResourceData(t, &agentQueueProject, &agentQueuePoolID, &agentQueueID)
	agentClient, clients := generateMocks(ctrl)

	agentClient.
		EXPECT().
		GetAgentQueue(clients.Ctx, taskagent.GetAgentQueueArgs{
			QueueId: &agentQueueID,
			Project: &agentQueueProject,
		}).
		Return(nil, errors.New("GetAgentQueue() Failed"))

	err := resourceAgentQueueRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "GetAgentQueue() Failed")
}

func TestAgentQueue_DoesNotSwallowDeleteErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := generateResourceData(t, &agentQueueProject, &agentQueuePoolID, &agentQueueID)
	agentClient, clients := generateMocks(ctrl)

	agentClient.
		EXPECT().
		DeleteAgentQueue(clients.Ctx, taskagent.DeleteAgentQueueArgs{
			QueueId: &agentQueueID,
			Project: &agentQueueProject,
		}).
		Return(errors.New("DeleteAgentQueue() Failed"))

	err := resourceAgentQueueDelete(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "DeleteAgentQueue() Failed")
}

func generateResourceData(t *testing.T, project *string, poolID *int, resourceID *int) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, ResourceAgentQueue().Schema, nil)
	if project != nil {
		resourceData.Set(projectID, *project)
	}

	if poolID != nil {
		resourceData.Set(agentPoolID, *poolID)
	}

	if resourceID != nil {
		resourceData.SetId(strconv.Itoa(*resourceID))
	}

	return resourceData
}

func generateMocks(ctrl *gomock.Controller) (*azdosdkmocks.MockTaskagentClient, *client.AggregatedClient) {
	agentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	return agentClient, &client.AggregatedClient{
		TaskAgentClient: agentClient,
		Ctx:             context.Background(),
	}
}
