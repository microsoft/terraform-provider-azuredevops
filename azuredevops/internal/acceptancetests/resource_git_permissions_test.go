//go:build (all || permissions || resource_git_permissions) && (!exclude_permissions || !exclude_resource_git_permissions)
// +build all permissions resource_git_permissions
// +build !exclude_permissions !exclude_resource_git_permissions

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := testutils.HclGitPermissions(projectName)

	tfNode := "azuredevops_git_permissions.git-permissions"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
				),
			},
		},
	})
}
