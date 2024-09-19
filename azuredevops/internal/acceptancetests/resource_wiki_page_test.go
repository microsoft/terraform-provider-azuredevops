//go:build (all || wiki || resource_wiki) && !exclude_resource_wiki
// +build all wiki resource_wiki
// +build !exclude_resource_wiki

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWikiPageResource_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	tf := "azuredevops_wiki.project_wiki"
	resourceType := "azuredevops_wiki"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: HclProjectWikiPage(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "name"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "wiki_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "path"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "content"),
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

func TestAccWikiPageResource_CreateAndUpdate(t *testing.T) {

	projectName := testutils.GenerateResourceName()
	tf := "azuredevops_wiki.project_wiki"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed("azuredevops_wiki"),
		Steps: []resource.TestStep{
			{
				Config: HclProjectWikiPage(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "name"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "wiki_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "path"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki_page.wiki_page", "content"),
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

func HclProjectWikiPage(projectName string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_wiki" "project_wiki" {
	name = "projectWikiRepo"
	project_id = azuredevops_project.project.id
	type = "projectWiki"
}

resource "azuredevops_wiki_page" "wiki_page" {
  project_id = azuredevops_project.project.id
  wiki_id = azuredevops_wiki.project_wiki.id
  path = "/path"
  content = "content"
}
`, projectResource)
}
