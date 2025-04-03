//go:build (all || core || data_sources || data_feed) && (!data_sources || !exclude_feed)
// +build all core data_sources data_feed
// +build !data_sources !exclude_feed

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
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

func TestAccFeedPermission_importErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	tfNode := "azuredevops_feed_permission.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFeedPermissionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclFeedPermissionBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					CheckFeedPermissionExist(),
					resource.TestCheckResourceAttrSet(tfNode, "feed_id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "role"),
					resource.TestCheckResourceAttrSet(tfNode, "identity_descriptor"),
				),
			},
			{
				Config: hclFeedPermissionImport(projectName),
				ExpectError: feedPermissionRequiresImportError(),
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

func hclFeedPermissionImport(name string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_feed_permission" "import" {
  feed_id             = azuredevops_feed.test.id
  project_id          = azuredevops_project.test.id
  role                = "reader"
  identity_descriptor = azuredevops_group.test.descriptor
}
`, hclFeedPermissionBasic(name))
}

func checkFeedPermissionDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_feed_permission" {
			continue
		}
		id := res.Primary.Attributes["feed_id"]
		projectID := res.Primary.Attributes["project_id"]
		permissions, err := clients.FeedClient.GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
			FeedId:  &id,
			Project: &projectID,
		})

		if err == nil {
			if permissions != nil && len(*permissions) > 0 {
				return fmt.Errorf(" Feed permissions (Feed ID: %s) should not exist", id)
			}
		}
	}
	return nil
}

func CheckFeedPermissionExist() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["azuredevops_feed_permission.test"]
		if !ok {
			return fmt.Errorf(" Did not find a `azuredevops_feed_permission` in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id := res.Primary.Attributes["feed_id"]
		projectID := res.Primary.Attributes["project_id"]

		_, err := clients.FeedClient.GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
			FeedId:  &id,
			Project: &projectID,
		})

		if err != nil {
			return fmt.Errorf(" Feed permissions with Feed ( Feed ID=%s ) cannot be found!. Error=%v", id, err)
		}

		return nil
	}
}

func feedPermissionRequiresImportError() *regexp.Regexp {
	return regexp.MustCompile(`Error: feed Permission for Feed : .* and Identity : .* already exists`)
}
