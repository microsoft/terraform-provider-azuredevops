// +build all resource_release_definition
// +build !exclude_resource_release_definition

package acceptancetests

import (
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccReleaseDefinition_Create_Update_Import(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	releaseDefinitionPathEmpty := `\`
	releaseDefinitionNameFirst := testutils.GenerateResourceName()
	releaseDefinitionNameSecond := testutils.GenerateResourceName()
	releaseDefinitionNameThird := testutils.GenerateResourceName()

	releaseDefinitionPathFirst := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	releaseDefinitionPathSecond := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	releaseDefinitionPathThird := `\` + releaseDefinitionNameFirst + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	releaseDefinitionPathFourth := `\` + releaseDefinitionNameSecond + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfReleaseDefNode := "azuredevops_release_definition.release"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkReleaseDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclReleaseDefinitionAgentless(projectName, releaseDefinitionNameFirst, releaseDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(releaseDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathEmpty),
				),
			}, {
				Config: testutils.HclReleaseDefinitionAgentless(projectName, releaseDefinitionNameSecond, releaseDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(releaseDefinitionNameSecond),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameSecond),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathEmpty),
				),
			}, {
				Config: testutils.HclReleaseDefinitionAgentless(projectName, releaseDefinitionNameFirst, releaseDefinitionPathFirst),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(releaseDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathFirst),
				),
			}, {
				Config: testutils.HclReleaseDefinitionAgentless(projectName, releaseDefinitionNameFirst,
					releaseDefinitionPathSecond),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(releaseDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathSecond),
				),
			}, {
				Config: testutils.HclReleaseDefinitionAgentless(projectName, releaseDefinitionNameFirst, releaseDefinitionPathThird),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(releaseDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathThird),
				),
			}, {
				Config: testutils.HclReleaseDefinitionAgentless(projectName, releaseDefinitionNameFirst, releaseDefinitionPathFourth),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(releaseDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathFourth),
				),
			}, {
				Config: testutils.HclReleaseDefinitionAgentless(projectName, gitRepoName, releaseDefinitionNameThird, releaseDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(releaseDefinitionNameThird),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "repository.0.repo_id"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameThird),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathEmpty),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfReleaseDefNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfReleaseDefNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Given the name of an AzDO release definition, this will return a function that will check whether
// or not the definition (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkReleaseDefinitionExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		releaseDef, ok := s.RootModule().Resources["azuredevops_release_definition.release"]
		if !ok {
			return fmt.Errorf("Did not find a release definition in the TF state")
		}

		releaseDefinition, err := getReleaseDefinitionFromResource(releaseDef)
		if err != nil {
			return err
		}

		if *releaseDefinition.Name != expectedName {
			return fmt.Errorf("Release Definition has Name=%s, but expected Name=%s", *releaseDefinition.Name, expectedName)
		}

		return nil
	}
}

// verifies that all release definitions referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkReleaseDefinitionDestroyed(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_release_definition" {
			continue
		}

		// indicates the release definition still exists - this should fail the test
		if _, err := getReleaseDefinitionFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a release definition that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a release definition (and error)
func getReleaseDefinitionFromResource(resource *terraform.ResourceState) (*release.ReleaseDefinition, error) {
	releaseDefID, err := strconv.Atoi(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	return clients.ReleaseClient.GetReleaseDefinition(clients.Ctx, release.GetReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefID,
	})
}
