// +build all resource_agentpool

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"github.com/stretchr/testify/require"
)

var testAgentPoolID = rand.Intn(100)

var testAgentPool = taskagent.TaskAgentPool{
	Id:            &testAgentPoolID,
	Name:          converter.String("Name"),
	PoolType:      &taskagent.TaskAgentPoolTypeValues.Automation,
	AutoProvision: converter.Bool(false),
}

/**
 * Begin unit tests
 */

// verifies that the flatten/expand round trip yields the same agent pool definition
func TestAzureDevOpsAgentPool_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceAzureAgentPool().Schema, nil)
	flattenAzureAgentPool(resourceData, &testAgentPool)

	agentPoolAfterRoundTrip, err := expandAgentPool(resourceData, true)
	require.Nil(t, err)
	require.Equal(t, testAgentPool, *agentPoolAfterRoundTrip)
}

// verifies that the create operation is considered failed if the API call fails.
func TestAzureDevOpsAgentPool_CreateAgentPool_DoesNotSwallowErrorFromFailedAddAgentCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &config.AggregatedClient{
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

func TestAzureDevOpsAgentPool_DeleteAgentPool_ReturnsErrorIfIdReadFails(t *testing.T) {
	client := &config.AggregatedClient{}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureAgentPool().Schema, nil)
	flattenAzureAgentPool(resourceData, &testAgentPool)
	resourceData.SetId("")

	err := resourceAzureAgentPoolDelete(resourceData, client)
	require.Equal(t, "Error getting agent pool Id: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
}

func TestAzureDevOpsAgentPool_UpdateAgentPool_ReturnsErrorIfIdReadFails(t *testing.T) {
	client := &config.AggregatedClient{}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureAgentPool().Schema, nil)
	flattenAzureAgentPool(resourceData, &testAgentPool)
	resourceData.SetId("")

	err := resourceAzureAgentPoolUpdate(resourceData, client)
	require.Equal(t, "Error converting terraform data model to AzDO agent pool reference: Error getting agent pool Id: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
}

func TestAzureDevOpsAgentPool_UpdateAgentPool_UpdateAndRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &config.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	agentToUpdate := taskagent.TaskAgentPool{
		Id:            &testAgentPoolID,
		Name:          converter.String("Foo"),
		PoolType:      &taskagent.TaskAgentPoolTypeValues.Deployment,
		AutoProvision: converter.Bool(true),
	}

	resourceData := schema.TestResourceDataRaw(t, resourceAzureAgentPool().Schema, nil)
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
func TestAzureDevOpsAgentPoolDefinition_PoolTypeIsCorrect(t *testing.T) {
	validPoolTypes := []string{
		string(taskagent.TaskAgentPoolTypeValues.Automation),
		string(taskagent.TaskAgentPoolTypeValues.Deployment),
	}
	poolTypeSchema := resourceAzureAgentPool().Schema["pool_type"]

	for _, repoType := range validPoolTypes {
		_, errors := poolTypeSchema.ValidateFunc(repoType, "")
		require.Equal(t, 0, len(errors), "Agent pool type unexpectedly did not pass validation")
	}
}

// validates invalid pool types are rejected by the schema
func TestAzureDevOpsAgentPoolDefinition_WhenPoolTypeIsNotCorrect_ReturnsError(t *testing.T) {
	invalidPoolTypes := []string{"", "unknown"}
	poolTypeSchema := resourceAzureAgentPool().Schema["pool_type"]

	for _, poolType := range invalidPoolTypes {
		_, errors := poolTypeSchema.ValidateFunc(poolType, "pool_type")
		expectedError := fmt.Sprintf("expected pool_type to be one of [automation deployment], got %s", poolType)
		require.Equal(t, 1, len(errors), "Agent pool type %v unexpectedly passed validation", poolType)
		require.Equal(t, expectedError, errors[0].Error())
	}
}

/**
 * Begin acceptance tests
 */

// Verifies that the following sequence of events occurrs without error:
//	(1) TF apply creates agent pool
//	(2) TF state values are set
//	(3) Agent pool can be queried by ID and has expected name
//  (4) TF apply updates agent pool with new name
//  (5) Agent pool can be queried by ID and has expected name
// 	(6) TF destroy deletes agent pool
//	(7) Agent pool can no longer be queried by ID
func TestAccAzureDevOpsAgentPool_CreateAndUpdate(t *testing.T) {
	poolNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	poolNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfNode := "azuredevops_agent_pool.pool"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAgentPoolCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccAgentPoolResource(poolNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolNameFirst),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
					testAccCheckAgentPoolResourceExists(poolNameFirst),
				),
			},
			{
				Config: testhelper.TestAccAgentPoolResource(poolNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolNameSecond),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
					testAccCheckAgentPoolResourceExists(poolNameSecond),
				),
			},
			{
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Given the name of an AzDO project, this will return a function that will check whether
// or not the project (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func testAccCheckAgentPoolResourceExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_agent_pool.pool"]
		if !ok {
			return fmt.Errorf("Did not find a agent pool in the TF state")
		}

		clients := testAccProvider.Meta().(*config.AggregatedClient)
		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Parse ID error, ID:  %v !. Error= %v", resource.Primary.ID, err)
		}

		project, agentPoolErr := azureAgentPoolRead(clients, id)

		if agentPoolErr != nil {
			return fmt.Errorf("Agent Pool with ID=%d cannot be found!. Error=%v", id, err)
		}

		if *project.Name != expectedName {
			return fmt.Errorf("Agent Pool with ID=%d has Name=%s, but expected Name=%s", id, *project.Name, expectedName)
		}

		return nil
	}
}

// verifies that agent pool referenced in the state is destroyed. This will be invoked
// *after* terrafform destroys the resource but *before* the state is wiped clean.
func testAccAgentPoolCheckDestroy(s *terraform.State) error {
	clients := testAccProvider.Meta().(*config.AggregatedClient)

	// verify that every agent pool referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_agent_pool" {
			continue
		}

		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Agent Pool ID=%d cannot be parsed!. Error=%v", id, err)
		}

		// indicates the agent pool still exists - this should fail the test
		if _, err := azureAgentPoolRead(clients, id); err == nil {
			return fmt.Errorf("Agent Pool ID %d should not exist", id)
		}
	}

	return nil
}

func init() {
	InitProvider()
}
