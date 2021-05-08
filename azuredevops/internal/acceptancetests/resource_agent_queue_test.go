// +build all resource_agent_queue
// +build !exclude_resource_agent_queue

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccResourceAgentQueue_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_agent_queue.q"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAgentQueueResource(projectName, poolName),
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
