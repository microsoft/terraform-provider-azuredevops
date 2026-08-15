package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccProjectPipelineRetentionSettings_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project_pipeline_retention_settings.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPipelineRetentionSettings(projectName, 30, 20, 15, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "run_retention", "30"),
					resource.TestCheckResourceAttr(tfNode, "artifact_retention", "20"),
					resource.TestCheckResourceAttr(tfNode, "pull_request_run_retention", "15"),
					resource.TestCheckResourceAttr(tfNode, "retain_runs_per_protected_branch", "10"),
				),
			},
			{
				Config: hclProjectPipelineRetentionSettings(projectName, 45, 25, 18, 12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "run_retention", "45"),
					resource.TestCheckResourceAttr(tfNode, "artifact_retention", "25"),
					resource.TestCheckResourceAttr(tfNode, "pull_request_run_retention", "18"),
					resource.TestCheckResourceAttr(tfNode, "retain_runs_per_protected_branch", "12"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclProjectPipelineRetentionSettings(projectName string, runRetention, artifactRetention, pullRequestRunRetention, retainRunsPerProtectedBranch int) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_project_pipeline_retention_settings" "test" {
  project_id                        = azuredevops_project.test.id
  run_retention                     = %d
  artifact_retention                = %d
  pull_request_run_retention        = %d
  retain_runs_per_protected_branch  = %d
}
`, projectName, runRetention, artifactRetention, pullRequestRunRetention, retainRunsPerProtectedBranch)
}
