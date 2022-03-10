//go:build (all || core || resource_project || resource_project_features) && !exclude_resource_project_features
// +build all core resource_project resource_project_features
// +build !exclude_resource_project_features

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccProjectPipelineSettings_Enabled(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project_pipeline_settings.this"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclProjectPipelineSettings(projectName, false, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "enforce_job_scope", "false"),
					resource.TestCheckResourceAttr(tfNode, "enforce_referenced_repo_scoped_token", "false"),
					resource.TestCheckResourceAttr(tfNode, "enforce_settable_var", "false"),
					resource.TestCheckResourceAttr(tfNode, "publish_pipeline_metadata", "false"),
					resource.TestCheckResourceAttr(tfNode, "status_badges_are_private", "false"),
				),
				Destroy: false,
			},
			{
				Config: testutils.HclProjectPipelineSettings(projectName, true, true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "enforce_job_scope", "true"),
					resource.TestCheckResourceAttr(tfNode, "enforce_referenced_repo_scoped_token", "true"),
					resource.TestCheckResourceAttr(tfNode, "enforce_settable_var", "true"),
					resource.TestCheckResourceAttr(tfNode, "publish_pipeline_metadata", "true"),
					resource.TestCheckResourceAttr(tfNode, "status_badges_are_private", "true"),
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
