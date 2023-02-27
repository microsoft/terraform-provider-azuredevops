//go:build (all || wiki || resource_wiki) && !exclude_resource_wiki
// +build all wiki resource_wiki
// +build !exclude_resource_wiki

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccWikiResource_CreateAndUpdate(t *testing.T) {

	projectName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckWikiDestroyed("azuredevops_wiki"),
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
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "versions"),
					resource.TestCheckResourceAttrSet("azuredevops_wiki.code_wiki", "mapped_path"),
				),
			},
		},
	})
}

func CheckWikiDestroyed(resourceType string) resource.TestCheckFunc {
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
