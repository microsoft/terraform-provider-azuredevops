//go:build (all || resource_agent_queue) && !exclude_resource_agent_queue
// +build all resource_agent_queue
// +build !exclude_resource_agent_queue

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccResourceAgentQueue_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_agent_queue.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAgentQueueBasic(projectName, poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			}, {
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAgentQueue_basedOnPool(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_agent_queue.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAgentQueueBasedObnPool(projectName, poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "agent_pool_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			}, {
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclAgentQueueBasic(projectName, queueName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_agent_queue" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
}
`, projectName, queueName)
}

func hclAgentQueueBasedObnPool(projectName, poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_agent_pool" "test" {
  name           = "%s"
  auto_provision = false
  auto_update    = false
  pool_type      = "automation"
}

resource "azuredevops_agent_queue" "test" {
  project_id    = azuredevops_project.test.id
  agent_pool_id = azuredevops_agent_pool.test.id
}
`, projectName, poolName)
}
