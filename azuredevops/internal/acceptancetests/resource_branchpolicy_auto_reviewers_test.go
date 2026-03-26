package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyAutoReviewers_basic(t *testing.T) {
	if os.Getenv("AZDO_TEST_AAD_USER_EMAIL") == "" {
		t.Skip("Skip test due to AZDO_TEST_AAD_USER_EMAIL not set")
	}

	name := testutils.GenerateResourceName()
	autoReviewerTfNode := "azuredevops_branch_policy_auto_reviewers.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAutoReviewersBasic(name, true, true, false, "auto reviewer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "true"),
				),
			}, {
				Config: hclAutoReviewersBasic(name, false, false, true, "new auto reviewer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "false"),
				),
			}, {
				ResourceName:      autoReviewerTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(autoReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyAutoReviewers_minimumApproverCount(t *testing.T) {
	name := testutils.GenerateResourceName()
	autoReviewerTfNode := "azuredevops_branch_policy_auto_reviewers.test"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAutoReviewersMinimumApprover(name, true, true, true, "auto reviewer", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "settings.0.submitter_can_vote", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "settings.0.minimum_number_of_reviewers", "1"),
				),
			}, {
				Config: hclAutoReviewersMinimumApprover(name, true, true, true, "new auto reviewer", 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "settings.0.submitter_can_vote", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "settings.0.minimum_number_of_reviewers", "2"),
				),
			}, {
				ResourceName:      autoReviewerTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(autoReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclAutoReviewersBasic(name string, enabled, blocking, submitterCanVote bool, message string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name        = "%[1]s"
  description = "description"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%[2]s"
  account_license_type = "express"
}

resource "azuredevops_branch_policy_auto_reviewers" "test" {
  project_id = azuredevops_project.test.id
  enabled    = %[3]t
  blocking   = %[4]t
  settings {
    auto_reviewer_ids  = [azuredevops_user_entitlement.test.id]
    submitter_can_vote = %[5]t
    message            = "%[6]s"
    path_filters       = ["*/API*.cs", "README.md"]
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}
`, name, os.Getenv("AZDO_TEST_AAD_USER_EMAIL"), enabled, blocking, submitterCanVote, message)
}

func hclAutoReviewersMinimumApprover(name string, enabled, blocking, submitterCanVote bool, message string, numberOfApprovers int) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name        = "%[1]s"
  description = "description"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_group" "test" {
  scope        = azuredevops_project.test.id
  display_name = "test group"
}

resource "azuredevops_branch_policy_auto_reviewers" "test" {
  project_id = azuredevops_project.test.id
  enabled    = %[2]t
  blocking   = %[3]t
  settings {
    auto_reviewer_ids           = [azuredevops_group.test.origin_id]
    submitter_can_vote          = %[4]t
    message                     = "%[5]s"
    minimum_number_of_reviewers = %[6]d
    path_filters                = ["*/API*.cs", "README.md"]
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}
`, name, enabled, blocking, submitterCanVote, message, numberOfApprovers)
}
