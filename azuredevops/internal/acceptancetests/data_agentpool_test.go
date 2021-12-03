//go:build (all || data_sources || data_agent_pool) && (!exclude_data_sources || !exclude_data_agent_pool)
// +build all data_sources data_agent_pool
// +build !exclude_data_sources !exclude_data_agent_pool

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAgentPool_DataSource(t *testing.T) {
	agentPoolName := testutils.GenerateResourceName()
	createAgentPool := testutils.HclAgentPoolResource(agentPoolName)
	createAndGetAgentPoolData := fmt.Sprintf("%s\n%s", createAgentPool, testutils.HclAgentPoolDataSource())

	tfNode := "data.azuredevops_agent_pool.pool"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createAndGetAgentPoolData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", agentPoolName),
					resource.TestCheckResourceAttr(tfNode, "auto_provision", "false"),
					resource.TestCheckResourceAttr(tfNode, "pool_type", "automation"),
				),
			},
		},
	})
}
