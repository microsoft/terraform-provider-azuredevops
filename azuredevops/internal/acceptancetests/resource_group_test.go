// +build all core resource_group
// +build !exclude_resource_group

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func TestAccGroupResource_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccGroupResource_CreateAndUpdate: transient failures cause inconsistent results: https://github.com/microsoft/terraform-provider-azuredevops/issues/174")

	projectName := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGroupResource("mygroup", projectName, groupName),
				Check: resource.ComposeTestCheckFunc(
					checkGroupExists("mygroup", groupName),
					resource.TestCheckResourceAttrSet("azuredevops_group.mygroup", "scope"),
					resource.TestCheckResourceAttr("azuredevops_group.mygroup", "display_name", groupName),
				),
			},
			{
				ResourceName:      "azuredevops_group.mygroup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkGroupExists(resourceName, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		varGroup, ok := s.RootModule().Resources[fmt.Sprintf("azuredevops_group.%s", resourceName)]
		if !ok {
			return fmt.Errorf("Did not find a group resource with name %s in the TF state", resourceName)
		}

		getGroupArgs := graph.GetGroupArgs{
			GroupDescriptor: converter.String(varGroup.Primary.Attributes["display_name"]),
		}
		clients := testutils.GetProvider().Meta().(*config.AggregatedClient)
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
	clients := testutils.GetProvider().Meta().(*config.AggregatedClient)

	// verify that every project referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_group" {
			continue
		}

		id := resource.Primary.ID

		getGroupArgs := graph.GetGroupArgs{
			GroupDescriptor: converter.String(id),
		}
		group, err := clients.GraphClient.GetGroup(clients.Ctx, getGroupArgs)
		if err != nil {
			return err
		}
		if group.Descriptor != nil {
			return fmt.Errorf("Group with ID %s should not exist in scope %s", id, resource.Primary.Attributes["scope"])
		}
	}

	return nil
}
