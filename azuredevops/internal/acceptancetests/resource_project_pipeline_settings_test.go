package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccProjectPipelineSettings_Enabled(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project_pipeline_settings.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPipelineSettings(projectName, false, false, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "enforce_job_scope", "false"),
					resource.TestCheckResourceAttr(tfNode, "enforce_referenced_repo_scoped_token", "false"),
					resource.TestCheckResourceAttr(tfNode, "enforce_settable_var", "false"),
					resource.TestCheckResourceAttr(tfNode, "publish_pipeline_metadata", "false"),
					resource.TestCheckResourceAttr(tfNode, "status_badges_are_private", "false"),
					resource.TestCheckResourceAttr(tfNode, "enforce_job_scope_for_release", "false"),
				),
			},
			{
				Config: hclProjectPipelineSettings(projectName, true, true, true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "enforce_job_scope", "true"),
					resource.TestCheckResourceAttr(tfNode, "enforce_referenced_repo_scoped_token", "true"),
					resource.TestCheckResourceAttr(tfNode, "enforce_settable_var", "true"),
					resource.TestCheckResourceAttr(tfNode, "publish_pipeline_metadata", "true"),
					resource.TestCheckResourceAttr(tfNode, "status_badges_are_private", "true"),
					resource.TestCheckResourceAttr(tfNode, "enforce_job_scope_for_release", "true"),
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

func hclProjectPipelineSettings(projectName string, enforceJobAuthScope, enforceReferencedRepoScopedToken, enforceSettableVar, publishPipelineMetadata, statusBadgesArePrivate, enforceJobAuthScopeForReleases bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_project_pipeline_settings" "test" {
  project_id                           = azuredevops_project.test.id
  enforce_job_scope                    = %t
  enforce_referenced_repo_scoped_token = %t
  enforce_settable_var                 = %t
  publish_pipeline_metadata            = %t
  status_badges_are_private            = %t
  enforce_job_scope_for_release        = %t
}
`, projectName, enforceJobAuthScope, enforceReferencedRepoScopedToken, enforceSettableVar, publishPipelineMetadata, statusBadgesArePrivate, enforceJobAuthScopeForReleases)
}
