//go:build (all || data_environment) && !exclude_data_environment

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccEnvironment_dataSource(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkEnvironmentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceEnvironmentBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", name),
				),
			},
		},
	})
}

func TestAccEnvironment_dataSource_by_name(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkEnvironmentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceEnvironmentBasicByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", name),
				),
			},
		},
	})
}

func hclDataSourceEnvironmentBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_environment" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%[1]s"
  description = "Managed by Terraform"
}

data "azuredevops_environment" "test" {
  project_id     = azuredevops_project.test.id
  environment_id = azuredevops_environment.test.id
}
`, name)
}

func hclDataSourceEnvironmentBasicByName(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_environment" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%[1]s"
  description = "Managed by Terraform"
}

data "azuredevops_environment" "test" {
  project_id = azuredevops_project.test.id
  name       = azuredevops_environment.test.name
}
`, name)
}
