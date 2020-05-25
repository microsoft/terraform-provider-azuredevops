// +build all resource_branchpolicy_acceptance_test

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
)

/**
 * Begin acceptance tests
 */

// Verifies that the following sequence of events occurrs without error:
//	(1) Branch policies can be created with no errors
//	(2) Branch policies can be updated with no errors
//	(3) Branch policies can be deleted with no errors
func TestAccAzureDevOpsBranchPolicy_CreateAndUpdate(t *testing.T) {
	projName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	opts1 := hclOptions{
		projectName:            projName,
		repoName:               repoName,
		minReviewerOptions:     minReviewPolicyOpts{true, true, 1, false},
		buildValidationOptions: buildValidationPolicyOpts{true, true, "build validation", 0},
	}

	opts2 := hclOptions{
		projectName:            projName,
		repoName:               repoName,
		minReviewerOptions:     minReviewPolicyOpts{false, false, 2, true},
		buildValidationOptions: buildValidationPolicyOpts{false, false, "build validation rename", 720},
	}

	minReviewerTfNode := "azuredevops_branch_policy_min_reviewers.p"
	buildVlidationTfNode := "azuredevops_branch_policy_build_validation.p"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testhelper.TestAccPreCheck(t, nil) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getHCL(opts1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildVlidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildVlidationTfNode, "enabled", "true"),
				),
			}, {
				Config: getHCL(opts2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildVlidationTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildVlidationTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      minReviewerTfNode,
				ImportStateIdFunc: testAccImportStateIDFunc(minReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			}, {
				ResourceName:      buildVlidationTfNode,
				ImportStateIdFunc: testAccImportStateIDFunc(buildVlidationTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type minReviewPolicyOpts struct {
	enabled          bool
	blocking         bool
	reviewers        int
	submitterCanVote bool
}

type buildValidationPolicyOpts struct {
	enabled       bool
	blocking      bool
	displayName   string
	validDuration int
}

type hclOptions struct {
	projectName            string
	repoName               string
	minReviewerOptions     minReviewPolicyOpts
	buildValidationOptions buildValidationPolicyOpts
}

func getHCL(opts hclOptions) string {
	projectAndRepo := testhelper.TestAccAzureGitRepoResource(opts.projectName, opts.repoName, "Clean")
	buildDef := `
	resource "azuredevops_build_definition" "build" {
		project_id = azuredevops_project.project.id
		name       = "Sample Build Definition"

		repository {
			repo_type   = "TfsGit"
			repo_id     = azuredevops_git_repository.gitrepo.id
			yml_path    = "azure-pipelines.yml"
		}
	}
`
	minReviewCountPolicyFmt := `
	resource "azuredevops_branch_policy_min_reviewers" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			reviewer_count     = %d
			submitter_can_vote = %t
			scope {
				repository_id  = azuredevops_git_repository.gitrepo.id
				repository_ref = azuredevops_git_repository.gitrepo.default_branch
				match_type     = "exact"
			}
		}
	}
`

	minReviewCountPolicy := fmt.Sprintf(
		minReviewCountPolicyFmt,
		opts.minReviewerOptions.enabled,
		opts.minReviewerOptions.blocking,
		opts.minReviewerOptions.reviewers,
		opts.minReviewerOptions.submitterCanVote)

	buildValidationPolicyFmt := `
	resource "azuredevops_branch_policy_build_validation" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			display_name = "%s"
			valid_duration = %d
			build_definition_id = azuredevops_build_definition.build.id
			scope {
				repository_id  = azuredevops_git_repository.gitrepo.id
				repository_ref = azuredevops_git_repository.gitrepo.default_branch
				match_type     = "exact"
			}
		}
	}
`
	buildValidationPolicyFmt = fmt.Sprintf(
		buildValidationPolicyFmt,
		opts.buildValidationOptions.enabled,
		opts.buildValidationOptions.blocking,
		opts.buildValidationOptions.displayName,
		opts.buildValidationOptions.validDuration)

	return strings.Join(
		[]string{
			projectAndRepo,
			buildDef,
			minReviewCountPolicy,
			buildValidationPolicyFmt,
		},
		"\n",
	)

}
