//go:build (all || permissions || resource_iteration_permissions) && (!exclude_permissions || !exclude_resource_iteration_permissions)
// +build all permissions resource_iteration_permissions
// +build !exclude_permissions !exclude_resource_iteration_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func TestAccIterationPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclIterationPermissions(projectName, map[string]map[string]string{
		"root": {
			"CREATE_CHILDREN": "Deny",
			"GENERIC_READ":    "NotSet",
			"DELETE":          "Deny",
		},
		"iteration": {
			"CREATE_CHILDREN": "Allow",
			"GENERIC_READ":    "NotSet",
			"DELETE":          "Allow",
		},
	})
	tfNodeRoot := "azuredevops_iteration_permissions.root-permissions"
	tfNodeIteration := "azuredevops_iteration_permissions.iteration-permissions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE_CHILDREN", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.GENERIC_READ", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttrSet(tfNodeIteration, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeIteration, "principal"),
					resource.TestCheckResourceAttr(tfNodeIteration, "path", "Iteration 1"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.CREATE_CHILDREN", "allow"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.GENERIC_READ", "notset"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.DELETE", "allow"),
				),
			},
		},
	})
}

func TestAccIterationPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclIterationPermissions(projectName, map[string]map[string]string{
		"root": {
			"CREATE_CHILDREN": "Deny",
			"GENERIC_READ":    "NotSet",
			"DELETE":          "Deny",
		},
		"iteration": {
			"CREATE_CHILDREN": "Allow",
			"GENERIC_READ":    "NotSet",
			"DELETE":          "Allow",
		},
	})
	config2 := hclIterationPermissions(projectName, map[string]map[string]string{
		"root": {
			"CREATE_CHILDREN": "Allow",
			"GENERIC_READ":    "NotSet",
			"DELETE":          "Deny",
		},
		"iteration": {
			"CREATE_CHILDREN": "Deny",
			"GENERIC_READ":    "Allow",
			"DELETE":          "NotSet",
		},
	})
	tfNodeRoot := "azuredevops_iteration_permissions.root-permissions"
	tfNodeIteration := "azuredevops_iteration_permissions.iteration-permissions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE_CHILDREN", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.GENERIC_READ", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttrSet(tfNodeIteration, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeIteration, "principal"),
					resource.TestCheckResourceAttr(tfNodeIteration, "path", "Iteration 1"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.CREATE_CHILDREN", "allow"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.GENERIC_READ", "notset"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.DELETE", "allow"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE_CHILDREN", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.GENERIC_READ", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttrSet(tfNodeIteration, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeIteration, "principal"),
					resource.TestCheckResourceAttr(tfNodeIteration, "path", "Iteration 1"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.CREATE_CHILDREN", "deny"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNodeIteration, "permissions.DELETE", "notset"),
				),
			},
		},
	})
}

func hclIterationPermissions(projectName string, permissions map[string]map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions["root"], "=", "\n")
	iterationPermissions := datahelper.JoinMap(permissions["iteration"], "=", "\n")

	return fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_iteration_permissions" "root-permissions" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id
  permissions = {
		%s
  }
}

resource "azuredevops_iteration_permissions" "iteration-permissions" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id
  path       = "Iteration 1"
  permissions = {
		%s
  }
}


`, testutils.HclProjectResource(projectName), rootPermissions, iterationPermissions)
}
