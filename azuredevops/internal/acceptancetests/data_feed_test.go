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

func TestAccAzureDevOps_DataSource_Feed_By_Name(t *testing.T) {
	name := testutils.GenerateResourceName()

	FeedData := fmt.Sprintf(`
		resource "azuredevops_feed" "feed" {
			name = "%s"
		}

		data "azuredevops_feed" "feed" {
			name = azuredevops_feed.feed.name
		}
		`, name)

	tfNode := "data.azuredevops_feed.feed"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "feed_id"),
				),
			},
		},
	})
}

func TestAccAzureDevOps_DataSource_Feed_By_Feed_Id(t *testing.T) {
	name := testutils.GenerateResourceName()

	FeedData := fmt.Sprintf(`
		resource "azuredevops_feed" "feed" {
			name = "%s"
		}

		data "azuredevops_feed" "feed" {
			feed_id = azuredevops_feed.feed.id
		}
		`, name)

	tfNode := "data.azuredevops_feed.feed"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "feed_id"),
				),
			},
		},
	})
}
