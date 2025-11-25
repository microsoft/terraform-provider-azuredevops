package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccFeedsDataSource_List(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	feedName1 := testutils.GenerateResourceName()
	feedName2 := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_feeds.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedsDataSourceList(projectName, feedName1, feedName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "feeds.#", "2"),
					resource.TestCheckResourceAttrSet(tfNode, "feeds.0.name"),
					resource.TestCheckResourceAttrSet(tfNode, "feeds.0.feed_id"),
					resource.TestCheckResourceAttrSet(tfNode, "feeds.1.name"),
				),
			},
		},
	})
}

func hclFeedsDataSourceList(projectName, feedName1, feedName2 string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

resource "azuredevops_feed" "feed1" {
  name       = "%[2]s"
  project_id = azuredevops_project.test.id
}

resource "azuredevops_feed" "feed2" {
  name       = "%[3]s"
  project_id = azuredevops_project.test.id
}

data "azuredevops_feeds" "test" {
  project_id = azuredevops_project.test.id
  depends_on = [
    azuredevops_feed.feed1,
    azuredevops_feed.feed2
  ]
}`, projectName, feedName1, feedName2)
}
