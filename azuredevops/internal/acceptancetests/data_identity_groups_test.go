package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIdentityGroupsDataSource(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_identity_groups.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclIdentityGroupsDataSourceConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "groups.#"),
					resource.TestCheckResourceAttrSet(tfNode, "groups.0.descriptor"),
					resource.TestCheckResourceAttrSet(tfNode, "groups.0.subject_descriptor"),
				),
			},
		},
	})
}

func hclIdentityGroupsDataSourceConfig(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_identity_groups" "test" {
  project_id = azuredevops_project.test.id
}
`, projectName)
}
