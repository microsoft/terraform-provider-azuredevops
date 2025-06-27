//go:build (all || resource_feed) && !exclude_resource_feed

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

func TestAccFeed_basic(t *testing.T) {
	feedName := testutils.GenerateResourceName()
	tfNode := "azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFeedDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclFeedBasic(feedName),
				Check: resource.ComposeTestCheckFunc(
					CheckFeedExist(feedName),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
		},
	})
}

func TestAccFeed_project(t *testing.T) {
	feedName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFeedDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclFeedWithProject(projectName, feedName),
				Check: resource.ComposeTestCheckFunc(
					CheckFeedExist(feedName),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
				),
			},
		},
	})
}

func TestAccFeed_softDeleteRecovery(t *testing.T) {
	feedName := testutils.GenerateResourceName()
	tfNode := "azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFeedDestroyed,
		Steps: []resource.TestStep{
			{
				Config:  hclFeedBasic(feedName),
				Destroy: true,
			},
			{
				Config: hclFeedRestore(feedName),
				Check: resource.ComposeTestCheckFunc(
					CheckFeedExist(feedName),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
		},
	})
}

func TestAccFeed_requiresImportErrorOrg(t *testing.T) {
	feedName := testutils.GenerateResourceName()
	tfNode := "azuredevops_feed.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFeedDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclFeedBasic(feedName),
				Check: resource.ComposeTestCheckFunc(
					CheckFeedExist(feedName),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
			{
				Config:      hclFeedImportOrg(feedName),
				ExpectError: requiresFeedImportError(feedName),
			},
		},
	})
}

func TestAccFeed_requiresImportErrorProject(t *testing.T) {
	feedName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()

	tfNode := "azuredevops_feed.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedWithProject(projectName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
			{
				Config:      hclFeedImportProject(projectName, feedName),
				ExpectError: requiresFeedImportError(feedName),
			},
		},
	})
}

func checkFeedDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_feed" {
			continue
		}
		id := res.Primary.ID
		projectID := res.Primary.Attributes["project_id"]

		_, err := clients.FeedClient.GetFeed(clients.Ctx, feed.GetFeedArgs{
			FeedId:  &id,
			Project: &projectID,
		})
		if err == nil {
			return fmt.Errorf(" Feed (Feed ID: %s) should not exist", id)
		}
	}
	return nil
}

func CheckFeedExist(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["azuredevops_feed.test"]
		if !ok {
			return fmt.Errorf(" Did not find `azuredevops_feed` in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id := res.Primary.ID
		projectID := res.Primary.Attributes["project_id"]

		feeds, err := clients.FeedClient.GetFeed(clients.Ctx, feed.GetFeedArgs{
			FeedId:  &id,
			Project: &projectID,
		})

		if err != nil {
			return fmt.Errorf(" Feed with ID=%s cannot be found!. Error=%v", id, err)
		}

		if *feeds.Name != expectedName {
			return fmt.Errorf(" Feed with ID=%s has Name=%s, but expected Name=%s", id, *feeds.Name, expectedName)
		}
		return nil
	}
}

func requiresFeedImportError(resourceName string) *regexp.Regexp {
	message := "creating new feed. Name: %[1]s, Error: A feed named '%[1]s' already exists"
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}

func hclFeedBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_feed" "test" {
  name = "%s"
}`, name)
}

func hclFeedWithProject(projectName, feedName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_feed" "test" {
  name       = "%[2]s"
  project_id = azuredevops_project.test.id
}`, projectName, feedName)
}

func hclFeedRestore(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_feed" "test" {
  name = "%s"
  features {
    restore = true
  }
}
`, name)
}

func hclFeedImportOrg(name string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_feed" "import" {
  name = azuredevops_feed.test.name
}
`, hclFeedBasic(name))
}

func hclFeedImportProject(projectName, feedName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_feed" "import" {
  name       = azuredevops_feed.test.name
  project_id = azuredevops_project.test.id
}
`, hclFeedWithProject(projectName, feedName))
}
