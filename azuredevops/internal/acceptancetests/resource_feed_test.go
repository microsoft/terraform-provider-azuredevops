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

func TestAccAzureDevOps_Resource_Feed(t *testing.T) {
	name := testutils.GenerateResourceName()

	FeedResource := fmt.Sprintf(`
		resource "azuredevops_feed" "feed" {
			name = "%s"
		}
	`, name)

	tfNode := "azuredevops_feed.feed"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
		},
	})
}

func TestAccAzureDevOps_Resource_Feed_with_Project(t *testing.T) {
	name := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()

	ProjectResource := testutils.HclProjectResource(projectName)
	FeedResource := fmt.Sprintf(`
	%s

	resource "azuredevops_feed" "feed" {
		name       = "%s"
		project_id    = azuredevops_project.project.id
	}
	
	`, ProjectResource, name)

	tfNode := "azuredevops_feed.feed"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
				),
			},
		},
	})
}

func TestAccAzureDevOps_Resource_Feed_Soft_Delete(t *testing.T) {
	name := testutils.GenerateResourceName()

	FeedResource := fmt.Sprintf(`
		resource "azuredevops_feed" "feed" {
			name = "%s"
			permanent_delete = false
		}
	`, name)

	tfNode := "azuredevops_feed.feed"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
		},
	})
}
