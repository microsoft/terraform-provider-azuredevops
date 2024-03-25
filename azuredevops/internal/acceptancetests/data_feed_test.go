//go:build (all || core || data_sources || data_feed) && (!data_sources || !exclude_feed)
// +build all core data_sources data_feed
// +build !data_sources !exclude_feed

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAzureDevOps_DataSource_Feed(t *testing.T) {
	FeedData := testutils.HclFeedDataSource()

	tfNode := "data.azuredevops_feed.releases-feed"
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
