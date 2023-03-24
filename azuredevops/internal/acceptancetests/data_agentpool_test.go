//go:build (all || data_sources || data_agent_pool) && (!exclude_data_sources || !exclude_data_agent_pool)
// +build all data_sources data_agent_pool
// +build !exclude_data_sources !exclude_data_agent_pool

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccDataSourceAgentPool_basic(t *testing.T) {
	agentPoolName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_agent_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceAgentPoolBasic(agentPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", agentPoolName),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "auto_update", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
				),
			},
		},
	})
}

func hclDataSourceAgentPoolBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "test" {
  name           = "%s"
  auto_provision = false
  auto_update    = false
  pool_type      = "automation"
}

data "azuredevops_agent_pool" "test" {
  name = azuredevops_agent_pool.test.name
}

`, name)
}
