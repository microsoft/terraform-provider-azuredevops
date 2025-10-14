package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIdentityGroupDataSource(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	groupName := "Contributors"
	tfNode := "data.azuredevops_identity_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclIdentityGroupConfig(groupName, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfNode, "subject_descriptor"),
					resource.TestCheckResourceAttr(tfNode, "name", fmt.Sprintf("[%s]\\%s", projectName, groupName)),
				),
			},
		},
	})
}

func hclIdentityGroupConfig(groupName string, projectName string) string {
	combinedGroupName := fmt.Sprintf("[%s]\\\\%s", projectName, groupName)
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_identity_group" "test" {
  name       = "%s"
  project_id = azuredevops_project.project.id
}`, projectName, combinedGroupName)
}
