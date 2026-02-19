package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyMinReviewers_basic(t *testing.T) {
	name := testutils.GenerateResourceName()
	node := "azuredevops_branch_policy_min_reviewers.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPolicyMinReviewersBasic(1, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "id"),
					resource.TestCheckResourceAttr(node, "blocking", "true"),
					resource.TestCheckResourceAttr(node, "enabled", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.submitter_can_vote", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.allow_completion_with_rejects_or_waits", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.last_pusher_cannot_approve", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.on_last_iteration_require_vote", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.on_push_reset_approved_votes", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.on_each_iteration_require_vote", "false"),
				),
			}, {
				ResourceName:      node,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(node),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyMinReviewers_update(t *testing.T) {
	name := testutils.GenerateResourceName()
	node := "azuredevops_branch_policy_min_reviewers.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPolicyMinReviewersBasic(1, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "id"),
					resource.TestCheckResourceAttr(node, "blocking", "true"),
					resource.TestCheckResourceAttr(node, "enabled", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.submitter_can_vote", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.allow_completion_with_rejects_or_waits", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.last_pusher_cannot_approve", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.on_last_iteration_require_vote", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.on_push_reset_approved_votes", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.on_each_iteration_require_vote", "false"),
				),
			}, {
				Config: hclPolicyMinReviewersUpdate(2, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "id"),
					resource.TestCheckResourceAttr(node, "blocking", "false"),
					resource.TestCheckResourceAttr(node, "enabled", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.submitter_can_vote", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.allow_completion_with_rejects_or_waits", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.last_pusher_cannot_approve", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.on_last_iteration_require_vote", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.on_push_reset_all_votes", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.on_each_iteration_require_vote", "true"),
				),
			}, {
				ResourceName:      node,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(node),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyMinReviewers_resetAllVote(t *testing.T) {
	name := testutils.GenerateResourceName()
	node := "azuredevops_branch_policy_min_reviewers.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPolicyMinReviewersResetAllVote(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "id"),
					resource.TestCheckResourceAttr(node, "blocking", "true"),
					resource.TestCheckResourceAttr(node, "enabled", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.submitter_can_vote", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.on_push_reset_all_votes", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.on_push_reset_approved_votes", "true"),
				),
			}, {
				ResourceName:      node,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(node),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyMinReviewers_requiresImportError(t *testing.T) {
	name := testutils.GenerateResourceName()
	node := "azuredevops_branch_policy_min_reviewers.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPolicyMinReviewersResetAllVote(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "id"),
					resource.TestCheckResourceAttr(node, "blocking", "true"),
					resource.TestCheckResourceAttr(node, "enabled", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.submitter_can_vote", "false"),
					resource.TestCheckResourceAttr(node, "settings.0.on_push_reset_all_votes", "true"),
					resource.TestCheckResourceAttr(node, "settings.0.on_push_reset_approved_votes", "true"),
				),
			}, {
				Config:      hclPolicyMinReviewersResetRequireImportError(name),
				ExpectError: regexp.MustCompile(` creating policy in Azure DevOps: The update is rejected by policy`),
			},
		},
	})
}

func hclPolicyMinReviewersTemplate(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}
`, name)
}

func hclPolicyMinReviewersBasic(reviewers int, name string) string {
	template := hclPolicyMinReviewersTemplate(name)
	return fmt.Sprintf(`


%s

resource "azuredevops_branch_policy_min_reviewers" "test" {
  project_id = azuredevops_project.test.id
  enabled    = true
  blocking   = true
  settings {
    reviewer_count                         = %[2]d
    submitter_can_vote                     = false
    allow_completion_with_rejects_or_waits = false
    on_push_reset_approved_votes           = true
    on_each_iteration_require_vote         = false
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}
`, template, reviewers)
}

func hclPolicyMinReviewersUpdate(reviewers int, name string) string {
	template := hclPolicyMinReviewersTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_branch_policy_min_reviewers" "test" {
  project_id = azuredevops_project.test.id
  enabled    = false
  blocking   = false
  settings {
    reviewer_count                         = %[2]d
    submitter_can_vote                     = true
    allow_completion_with_rejects_or_waits = true
    last_pusher_cannot_approve             = true
    on_last_iteration_require_vote         = true
    on_each_iteration_require_vote         = true
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}
`, template, reviewers)
}

func hclPolicyMinReviewersResetAllVote(name string) string {
	template := hclPolicyMinReviewersTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_branch_policy_min_reviewers" "test" {
  project_id = azuredevops_project.test.id
  enabled    = true
  blocking   = true
  settings {
    reviewer_count               = 2
    submitter_can_vote           = false
    on_push_reset_all_votes      = true
    on_push_reset_approved_votes = true
    scope {
      repository_id  = data.azuredevops_git_repository.test.id
      repository_ref = "refs/heads/release"
      match_type     = "Exact"
    }
  }
}
`, template)
}

func hclPolicyMinReviewersResetRequireImportError(name string) string {
	template := hclPolicyMinReviewersResetAllVote(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_branch_policy_min_reviewers" "import" {
  project_id = azuredevops_branch_policy_min_reviewers.test.project_id
  enabled    = azuredevops_branch_policy_min_reviewers.test.enabled
  blocking   = azuredevops_branch_policy_min_reviewers.test.blocking
  settings {
    reviewer_count               = azuredevops_branch_policy_min_reviewers.test.settings.0.reviewer_count
    submitter_can_vote           = azuredevops_branch_policy_min_reviewers.test.settings.0.submitter_can_vote
    on_push_reset_all_votes      = azuredevops_branch_policy_min_reviewers.test.settings.0.on_push_reset_all_votes
    on_push_reset_approved_votes = azuredevops_branch_policy_min_reviewers.test.settings.0.on_push_reset_approved_votes
    scope {
      repository_id  = azuredevops_branch_policy_min_reviewers.test.settings.0.scope.0.repository_id
      repository_ref = azuredevops_branch_policy_min_reviewers.test.settings.0.scope.0.repository_ref
      match_type     = azuredevops_branch_policy_min_reviewers.test.settings.0.scope.0.match_type
    }
  }
}
`, template)
}
