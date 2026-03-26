package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyCommentResolution_basic(t *testing.T) {
	name := testutils.GenerateResourceName()
	resourceNode := "azuredevops_branch_policy_comment_resolution.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclCommentResolutionBasic(name, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNode, "enabled", "true"),
				),
			}, {
				Config: hclCommentResolutionBasic(name, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNode, "enabled", "false"),
				),
			}, {
				ResourceName:      resourceNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(resourceNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclCommentResolutionBasic(name string, enabled, blocking bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name        = "%[1]s"
  description = "description"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_branch_policy_comment_resolution" "test" {
  project_id = azuredevops_project.test.id
  enabled    = %[2]t
  blocking   = %[3]t
  settings {
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}`, name, enabled, blocking)
}
