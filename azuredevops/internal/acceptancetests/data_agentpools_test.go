//go:build (all || data_sources || data_agent_pools) && (!exclude_data_sources || !exclude_data_agent_pools)
// +build all data_sources data_agent_pools
// +build !exclude_data_sources !exclude_data_agent_pools

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Verifies that the following sequence of events occurrs without error:
//	(1) TF can create a project
//	(2) A data source is added to the configuration, and that data source can find the created project

func TestAccAgentPools_DataSource(t *testing.T) {
	agentPoolName := testutils.GenerateResourceName()
	agentPool1Name := agentPoolName + "_1"
	agentPool2Name := agentPoolName + "_2"

	createAgent1 := testutils.HclAgentPoolResourceAppendPoolNameToResourceName(agentPool1Name)
	createAgent2 := testutils.HclAgentPoolResourceAppendPoolNameToResourceName(agentPool2Name)
	agentPoolsData := testutils.HclAgentPoolsDataSource()
	createAgentPools := fmt.Sprintf("%s\n%s", createAgent1, createAgent2)

	tfNode := "data.azuredevops_agent_pools.pools"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createAgentPools,
			},
			{
				Config: agentPoolsData,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckNestedKeyExistsWithValue(tfNode, "name", agentPool1Name),
					testutils.CheckNestedKeyExistsWithValue(tfNode, "name", agentPool2Name),
				),
			},
		},
	})
}
