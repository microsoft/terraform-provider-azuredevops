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

func TestAccAzureDevOps_DataSource_Feed(t *testing.T) {
	name := testutils.GenerateResourceName()

	FeedResource := testutils.HclFeedResource(name)
	FeedData := fmt.Sprintf(`
		%s

		data "azuredevops_feed" "feed" {
			name = azuredevops_feed.feed.name
		}
		`, FeedResource)

	tfNode := "data.azuredevops_feed.feed"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
				),
			},
		},
	})
}
