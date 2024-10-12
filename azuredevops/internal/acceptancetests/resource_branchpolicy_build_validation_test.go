package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyBuildValidation_basic(t *testing.T) {
	name := testutils.GenerateResourceName()
	buildValidationTfNode := "azuredevops_branch_policy_build_validation.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclBuildValidationBasic(name, true, true, "build validation", 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.filename_patterns.#", "3"),
				),
			}, {
				Config: hclBuildValidationBasic(name, false, false, "build validation rename", 720),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.filename_patterns.#", "3"),
				),
			}, {
				ResourceName:      buildValidationTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(buildValidationTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclBuildValidationBasic(name string, enabled, blocking bool, displayName string, validDuration int) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name        = "%[1]s"
  description = "description"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id      = azuredevops_project.test.id
  name            = "Example Build Definition"
  agent_pool_name = "Azure Pipelines"
  path            = "\\"

  repository {
    repo_type   = "TfsGit"
    repo_id     = data.azuredevops_git_repository.test.id
    branch_name = "main"
    yml_path    = "path/to/yaml"
  }
}

resource "azuredevops_branch_policy_build_validation" "test" {
  project_id = azuredevops_project.test.id
  enabled    = %[2]t
  blocking   = %[3]t
  settings {
    display_name        = "%[4]s"
    valid_duration      = %[5]d
    build_definition_id = azuredevops_build_definition.test.id
    filename_patterns = [
      "/WebApp/*",
      "!/WebApp/Tests/*",
      "*.cs"
    ]
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}`, name, enabled, blocking, displayName, validDuration)
}
