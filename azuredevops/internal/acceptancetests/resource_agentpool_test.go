//go:build (all || resource_agentpool) && !exclude_resource_agentpool
// +build all resource_agentpool
// +build !exclude_resource_agentpool

package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// Verifies that the following sequence of events occurrs without error:
//
//	(1) TF apply creates agent pool
//	(2) TF state values are set
//	(3) Agent pool can be queried by ID and has expected name
//	(4) TF apply updates agent pool with new name
//	(5) Agent pool can be queried by ID and has expected name
//	(6) TF destroy deletes agent pool
//	(7) Agent pool can no longer be queried by ID
func TestAccAgentPool_CreateAndUpdate(t *testing.T) {
	poolNameFirst := testutils.GenerateResourceName()
	poolNameSecond := testutils.GenerateResourceName()
	tfNode := "azuredevops_agent_pool.pool"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAgentPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAgentPoolResource(poolNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolNameFirst),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
					checkAgentPoolExists(poolNameFirst),
				),
			},
			{
				Config: testutils.HclAgentPoolResource(poolNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolNameSecond),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
					checkAgentPoolExists(poolNameSecond),
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

// Given the name of an agent pool, this will return a function that will check whether
// or not the pool (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkAgentPoolExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_agent_pool.pool"]
		if !ok {
			return fmt.Errorf("Did not find a agent pool in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Parse ID error, ID:  %v !. Error= %v", resource.Primary.ID, err)
		}

		pool, agentPoolErr := clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{PoolId: &id})

		if agentPoolErr != nil {
			return fmt.Errorf("Agent Pool with ID=%d cannot be found!. Error=%v", id, err)
		}

		if *pool.Name != expectedName {
			return fmt.Errorf("Agent Pool with ID=%d has Name=%s, but expected Name=%s", id, *pool.Name, expectedName)
		}

		return nil
	}
}

// verifies that agent pool referenced in the state is destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkAgentPoolDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

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
		if _, err := clients.TaskAgentClient.GetAgentPool(clients.Ctx, taskagent.GetAgentPoolArgs{PoolId: &id}); err == nil {
			return fmt.Errorf("Agent Pool ID %d should not exist", id)
		}
	}

	return nil
}
