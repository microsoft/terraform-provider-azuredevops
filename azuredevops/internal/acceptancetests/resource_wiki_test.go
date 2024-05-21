//go:build (all || wiki || resource_wiki) && !exclude_resource_wiki
// +build all wiki resource_wiki
// +build !exclude_resource_wiki

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccWikiResource_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	tf := "azuredevops_wiki.project_wiki"
	resourceType := "azuredevops_wiki"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWiki(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "name"),
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

func TestAccWikiResource_Complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	tf := "azuredevops_wiki.project_wiki"
	resourceType := "azuredevops_wiki"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWiki(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "name"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "name"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "repository_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "version"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "mapped_path"),
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

func TestAccWikiResource_CreateAndUpdate(t *testing.T) {

	projectName := testutils.GenerateResourceName()
	tf := "azuredevops_wiki.project_wiki"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed("azuredevops_wiki"),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWiki(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "name"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "name"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "repository_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "version"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "mapped_path"),
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

func TestAccWikiResource_ImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resourceType := "azuredevops_wiki"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkWikiDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWiki(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "type"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.project_wiki", "name"),
				),
			},
			{
				Config:      hclWikiResourceRequiresImport(projectName),
				ExpectError: wikiRequiresImportError(),
			},
		},
	})
}

func checkWikiDestroyed(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, resource := range s.RootModule().Resources {
			if resource.Type != resourceType {
				continue
			}

			// indicates the resource exists - this should fail the test
			clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
			_, err := clients.WikiClient.GetWiki(clients.Ctx, wiki.GetWikiArgs{WikiIdentifier: converter.String(resource.Primary.ID)})
			if err == nil {
				return fmt.Errorf("found wiki that should have been deleted")
			}
		}

		return nil
	}
}

func wikiRequiresImportError() *regexp.Regexp {
	message := "Error: Wiki already exists with name ('codeWikiRepo'|'projectWikiRepo')."
	return regexp.MustCompile(message)
}

func hclWikiResourceRequiresImport(projectName string) string {
	projectResource := testutils.HclProjectResource(projectName)
	projectFeatures := fmt.Sprintf(`
%s

resource "azuredevops_git_repository" "repository" {
	project_id = azuredevops_project.project.id
	name       = "Repo"
	initialization {
	  init_type = "Clean"
	}
}

resource "azuredevops_wiki" "code_wiki_test" {
	name = "codeWikiRepo"
	project_id = azuredevops_project.project.id
	repository_id = azuredevops_git_repository.repository.id
	version = "master"
	type = "codeWiki"
	mapped_path = "/"
}

resource "azuredevops_wiki" "project_wiki_test" {
	name = "projectWikiRepo"
	project_id = azuredevops_project.project.id
	type = "projectWiki"
}
`, projectResource)

	return fmt.Sprintf("%s", projectFeatures)
}
