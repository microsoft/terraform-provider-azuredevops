//go:build (all || core || resource_group) && !exclude_resource_group
// +build all core resource_group
// +build !exclude_resource_group

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccGroupResource_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGroupBasic(projectName, groupName),
				Check: resource.ComposeTestCheckFunc(
					checkGroupExists(groupName),
					resource.TestCheckResourceAttrSet("azuredevops_group.test", "scope"),
					resource.TestCheckResourceAttrSet("azuredevops_group.test", "group_id"),
					resource.TestCheckResourceAttr("azuredevops_group.test", "display_name", groupName),
				),
			},
			{
				ResourceName:      "azuredevops_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkGroupExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		varGroup, ok := s.RootModule().Resources["azuredevops_group.test"]
		if !ok {
			return fmt.Errorf(" Did not find a group resource in the TF state")
		}

		getGroupArgs := graph.GetGroupArgs{
			GroupDescriptor: converter.String(varGroup.Primary.Attributes["descriptor"]),
		}
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		group, err := clients.GraphClient.GetGroup(clients.Ctx, getGroupArgs)
		if err != nil {
			return err
		}
		if group == nil {
			return fmt.Errorf("Group with Name=%s does not exit", varGroup.Primary.Attributes["display_name"])
		}
		if *group.DisplayName != expectedName {
			return fmt.Errorf("Group has Name=%s, but expected %s", *group.DisplayName, expectedName)
		}

		return nil
	}
}

func checkGroupDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	// verify that every project referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_group" {
			continue
		}

		// The group will be returned even if it has been deleted from the account or has had all its memberships deleted.
		id := resource.Primary.ID
		err := clients.GraphClient.DeleteGroup(clients.Ctx, graph.DeleteGroupArgs{
			GroupDescriptor: converter.String(id),
		})
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				return nil
			}
			return fmt.Errorf(" Group with ID %s should not exist in scope %s", id, resource.Primary.Attributes["scope"])
		}

	}

	return nil
}

func hclGroupBasic(projectName, groupName string) string {

	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_group" "test" {
  scope        = azuredevops_project.test.id
  display_name = "%[2]s"
}
`, projectName, groupName)

}
