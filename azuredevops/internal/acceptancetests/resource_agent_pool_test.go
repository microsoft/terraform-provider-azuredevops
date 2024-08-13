//go:build (all || resource_agentpool) && !exclude_resource_agentpool
// +build all resource_agentpool
// +build !exclude_resource_agentpool

package acceptancetests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccAgentPool_basic(t *testing.T) {
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_agent_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAgentPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAgentPoolBasic(poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "auto_update", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAgentPool_update(t *testing.T) {
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_agent_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAgentPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAgentPoolBasic(poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "auto_update", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclAgentPoolUpdate(poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "true"),
					resource.TestCheckResourceAttr(tfNode, "auto_update", "true"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAgentPool_requiresImportErrorStep(t *testing.T) {
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_agent_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAgentPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAgentPoolBasic(poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "auto_update", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
				),
			},

			{
				Config:      hclAgentPoolResourceRequiresImport(poolName),
				ExpectError: requiresImportError(poolName),
			},
		},
	})
}

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

func requiresImportError(resourceName string) *regexp.Regexp {
	message := "creating agent pool in Azure DevOps: Agent pool %[1]s already exists."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}

func hclAgentPoolBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "test" {
  name           = "%s"
  auto_provision = false
  auto_update    = false
  pool_type      = "automation"
}`, name)
}

func hclAgentPoolUpdate(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "test" {
  name           = "%s"
  auto_provision = true
  auto_update    = true
  pool_type      = "automation"
}`, name)
}

func hclAgentPoolResourceRequiresImport(name string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_agent_pool" "import" {
  name           = azuredevops_agent_pool.test.name
  auto_provision = azuredevops_agent_pool.test.auto_provision
  auto_update    = azuredevops_agent_pool.test.auto_update
  pool_type      = azuredevops_agent_pool.test.pool_type
}`, hclAgentPoolBasic(name))
}
