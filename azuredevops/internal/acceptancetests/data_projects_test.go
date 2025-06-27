//go:build (all || core || data_sources || resource_project || data_projects) && (!data_sources || !exclude_data_projects)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccProjects_DataSource_SingleProject(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_projects.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceProjectsSingle(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "projects.#", "1"),
				),
			},
		},
	})
}

func TestAccProjects_DataSource_EmptyResult(t *testing.T) {
	tfNode := "data.azuredevops_projects.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceProjectsEmptyResult(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "projects.#", "0"),
				),
			},
		},
	})
}

func hclDataSourceProjectsSingle(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_projects" "test" {
  name = azuredevops_project.test.name
}
`, name)
}

func hclDataSourceProjectsEmptyResult() string {
	return fmt.Sprintf(`
data "azuredevops_projects" "test" {
  name  = "invalid_name"
  state = "wellFormed"
}
`)
}
