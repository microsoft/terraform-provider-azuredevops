//go:build (all || core || resource_project || resource_project_features) && !exclude_resource_project_features
// +build all core resource_project resource_project_features
// +build !exclude_resource_project_features

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccProjectFeatures_EnableUpdateFeature(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project_features.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectFeatureBasic(projectName, "disabled", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "features.testplans", "disabled"),
					resource.TestCheckResourceAttr(tfNode, "features.artifacts", "disabled"),
				),
			},
			{
				Config: hclProjectFeatureBasic(projectName, "enabled", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "features.testplans", "enabled"),
					resource.TestCheckResourceAttr(tfNode, "features.artifacts", "disabled"),
				),
			},
		},
	})
}

func hclProjectFeatureBasic(name, testPlanState, artifactState string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_project_features" "test" {
  project_id = azuredevops_project.test.id
  features = {
    "testplans" = "%[2]s"
    "artifacts" = "%[3]s"
  }
}`, name, testPlanState, artifactState)
}
