//go:build (all || core || data_sources || data_feed) && (!data_sources || !exclude_feed)
// +build all core data_sources data_feed
// +build !data_sources !exclude_feed

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccFeedPermission_basic(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_feed_permission.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedPermissionBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "feed_id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "role"),
					resource.TestCheckResourceAttrSet(tfNode, "identity_descriptor"),
				),
			},
		},
	})
}

func hclFeedPermissionBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_feed" "test" {
  name       = "%[1]s"
  project_id = azuredevops_project.test.id
}

resource "azuredevops_group" "test" {
  scope        = azuredevops_project.test.id
  display_name = "%[1]s"
}

resource "azuredevops_feed_permission" "test" {
  feed_id             = azuredevops_feed.test.id
  project_id          = azuredevops_project.test.id
  role                = "reader"
  identity_descriptor = azuredevops_group.test.descriptor
}
`, name)
}
