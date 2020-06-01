// +build all core data_projects

package azuredevops

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
)

// Verifies that the following sequence of events occurrs without error:
//	(1) TF can create a project
//	(2) A data source is added to the configuration, and that data source can find the created project

func TestAccAgentPools_DataSource(t *testing.T) {
	agentPoolName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	agentPool1Name := agentPoolName + "_1"
	agentPool2Name := agentPoolName + "_2"

	createAgent1 := testhelper.TestAccAgentPoolResourceAppendPoolNameToResourceName(agentPool1Name)
	createAgent2 := testhelper.TestAccAgentPoolResourceAppendPoolNameToResourceName(agentPool2Name)
	agentPoolsData := testhelper.TestAccAgentPoolsDataSource()
	createAgentPools := fmt.Sprintf("%s\n%s", createAgent1, createAgent2)

	tfNode := "data.azuredevops_agent_pools.pools"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testhelper.TestAccPreCheck(t, nil) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createAgentPools,
			},
			{
				Config: agentPoolsData,
				Check: resource.ComposeTestCheckFunc(
					testhelper.CheckNestedKeyExistsWithValue(tfNode, "name", agentPool1Name),
					testhelper.CheckNestedKeyExistsWithValue(tfNode, "name", agentPool2Name),
				),
			},
		},
	})
}
