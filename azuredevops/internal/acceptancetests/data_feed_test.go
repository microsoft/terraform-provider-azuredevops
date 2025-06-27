//go:build (all || data_sources || data_feed) && (!exclude_data_sources || !exclude_data_feed)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccFeedDataSource_byName(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedDataSourceByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "feed_id"),
				),
			},
		},
	})
}

func TestAccFeedDataSource_byId(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedDataSourceByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "feed_id"),
				),
			},
		},
	})
}

func hclFeedDataSourceByName(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_feed" "test" {
  name = "%s"
}

data "azuredevops_feed" "test" {
  name = azuredevops_feed.test.name
}`, name)
}

func hclFeedDataSourceByID(feedID string) string {
	return fmt.Sprintf(`
resource "azuredevops_feed" "test" {
  name = "%s"
}

data "azuredevops_feed" "test" {
  feed_id = azuredevops_feed.test.id
}`, feedID)
}
