// +build all resource_agent_queue

package azuredevops

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"github.com/stretchr/testify/require"
)

/**
 * Begin unit tests
 */

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
	resourceData := schema.TestResourceDataRaw(t, resourceAgentQueue().Schema, nil)
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

func generateMocks(ctrl *gomock.Controller) (*azdosdkmocks.MockTaskagentClient, *config.AggregatedClient) {
	agentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	return agentClient, &config.AggregatedClient{
		TaskAgentClient: agentClient,
		Ctx:             context.Background(),
	}
}

/**
 * Begin acceptance tests
 */
func TestAccResourceAgentQueue_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	poolName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfNode := "azuredevops_agent_queue.q"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testhelper.TestAccPreCheck(t, nil) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccAgentQueueResource(projectName, poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "agent_pool_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			}, {
				ResourceName:      tfNode,
				ImportStateIdFunc: testAccImportStateIDFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func init() {
	InitProvider()
}
