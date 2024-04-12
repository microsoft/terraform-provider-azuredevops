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

func TestAccAzureDevOps_Resource_FeedPermission(t *testing.T) {
	name := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()

	FeedResource := fmt.Sprintf(`
		%s

		resource "azuredevops_feed" "feed" {
			name = "%s"
			project_id = azuredevops_project.project.id
		}

		resource "azuredevops_feed_permission" "permission" {
			feed_id = azuredevops_feed.feed.id
			project_id = azuredevops_project.project.id
			role = "reader"
			identity_descriptor = azuredevops_group.group.descriptor
		}
	`, testutils.HclGroupResource("group", projectName, groupName), name)

	tfNode := "azuredevops_feed_permission.permission"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: FeedResource,
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
