package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyMergeTypes_basic(t *testing.T) {
	name := testutils.GenerateResourceName()
	buildValidationTfNode := "azuredevops_branch_policy_merge_types.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclMergeTypesBasic(name, true, true, true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_squash", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_and_fast_forward", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_basic_no_fast_forward", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_with_merge", "true"),
				),
			}, {
				Config: hclMergeTypesBasic(name, false, false, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_squash", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_and_fast_forward", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_basic_no_fast_forward", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_with_merge", "false"),
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

func hclMergeTypesBasic(name string, enabled, blocking, allowSquash, allowRebase, allowNoFastForward, allowRebaseMerge bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name        = "%[1]s"
  description = "description"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_branch_policy_merge_types" "test" {
  project_id = azuredevops_project.test.id
  enabled    = %[2]t
  blocking   = %[3]t
  settings {
    allow_squash                  = %[4]t
    allow_rebase_and_fast_forward = %[5]t
    allow_basic_no_fast_forward   = %[6]t
    allow_rebase_with_merge       = %[7]t
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}`, name, enabled, blocking, allowSquash, allowRebase, allowNoFastForward, allowRebaseMerge)
}
