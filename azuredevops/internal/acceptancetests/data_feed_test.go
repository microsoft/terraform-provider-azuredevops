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

func TestAccFeedDataSource_WithUpstream(t *testing.T) {
	name := testutils.GenerateResourceName()
	upstreamName := "npmjs"
	upstreamLocation := "https://registry.npmjs.org"

	tfNode := "data.azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedDataSourceWithUpstream(name, upstreamName, upstreamLocation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttr(tfNode, "upstream_sources.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "upstream_sources.0.name", upstreamName),
					resource.TestCheckResourceAttr(tfNode, "upstream_sources.0.location", upstreamLocation),
					resource.TestCheckResourceAttr(tfNode, "upstream_sources.0.protocol", "npm"),
					resource.TestCheckResourceAttr(tfNode, "upstream_sources.0.upstream_source_type", "public"),
					resource.TestCheckResourceAttr(tfNode, "upstream_enabled", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
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

func hclFeedDataSourceWithUpstream(name, upstreamName, upstreamLocation string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

resource "azuredevops_feed" "test" {
  name       = "%[1]s"
  project_id = azuredevops_project.test.id
  
  upstream_sources {
    name                 = "%[2]s"
    protocol             = "npm"
    location             = "%[3]s"
    upstream_source_type = "public"
  }
}

data "azuredevops_feed" "test" {
  name       = azuredevops_feed.test.name
  project_id = azuredevops_project.test.id
}`, name, upstreamName, upstreamLocation)
}
