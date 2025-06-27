//go:build (all || wiki || resource_wiki) && !exclude_resource_wiki

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWikiPageResource_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tf := "azuredevops_wiki.test"
	resourceType := "azuredevops_wiki"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclProjectWikiPageBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "wiki_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "path"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "content"),
				),
			},
			{
				ResourceName:      tf,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWikiPageResource_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tf := "azuredevops_wiki.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed("azuredevops_wiki"),
		Steps: []resource.TestStep{
			{
				Config: hclProjectWikiPageBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "wiki_id"),
					resource.TestCheckResourceAttr("azuredevops_wiki_page.test", "path", "/path"),
					resource.TestCheckResourceAttr("azuredevops_wiki_page.test", "content", "content"),
				),
			},
			{
				ResourceName:      tf,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclProjectWikiPageUpdate(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "wiki_id"),
					resource.TestCheckResourceAttr("azuredevops_wiki_page.test", "path", "/path"),
					resource.TestCheckResourceAttr("azuredevops_wiki_page.test", "content", "contentupdate"),
				),
			},
			{
				ResourceName:      tf,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWikiPageResource_requireImportError(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tf := "azuredevops_wiki.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed("azuredevops_wiki"),
		Steps: []resource.TestStep{
			{
				Config: hclProjectWikiPageBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.test", "wiki_id"),
					resource.TestCheckResourceAttr("azuredevops_wiki_page.test", "path", "/path"),
					resource.TestCheckResourceAttr("azuredevops_wiki_page.test", "content", "content"),
				),
			},
			{
				ResourceName:      tf,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      hclProjectWikiPageImport(projectName),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The page '/path' specified in the add operation already exists in the wiki. Please specify a new page path.`)),
			},
		},
	})
}

func hclProjectWikiPageBasic(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_wiki" "test" {
  project_id = azuredevops_project.test.id
  name       = "projectWikiRepo"
  type       = "projectWiki"
}

resource "azuredevops_wiki_page" "test" {
  project_id = azuredevops_project.test.id
  wiki_id    = azuredevops_wiki.test.id
  path       = "/path"
  content    = "content"
}
`, projectName)
}

func hclProjectWikiPageUpdate(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_wiki" "test" {
  project_id = azuredevops_project.test.id
  name       = "projectWikiRepo"
  type       = "projectWiki"
}

resource "azuredevops_wiki_page" "test" {
  project_id = azuredevops_project.test.id
  wiki_id    = azuredevops_wiki.test.id
  path       = "/path"
  content    = "contentupdate"
}
`, projectName)
}

func hclProjectWikiPageImport(projectName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_wiki_page" "import" {
  project_id = azuredevops_wiki_page.test.project_id
  wiki_id    = azuredevops_wiki_page.test.wiki_id
  path       = azuredevops_wiki_page.test.path
  content    = azuredevops_wiki_page.test.content
}
`, hclProjectWikiPageBasic(projectName))
}
