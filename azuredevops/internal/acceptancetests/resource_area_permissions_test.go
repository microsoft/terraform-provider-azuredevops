//go:build (all || permissions || resource_area_permissions) && (!exclude_permissions || !exclude_resource_area_permissions)
// +build all permissions resource_area_permissions
// +build !exclude_permissions !exclude_resource_area_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func hclAreaPermissions(projectName string, permissions map[string]map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions["root"], "=", "\n")

	return fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_area_permissions" "root-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	permissions = {
		%s
	}
}

`, testutils.HclProjectResource(projectName), rootPermissions)
}

func TestAccAreaPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclAreaPermissions(projectName, map[string]map[string]string{
		"root": {
			"CREATE_CHILDREN": "Deny",
			"GENERIC_READ":    "NotSet",
			"DELETE":          "Deny",
			"WORK_ITEM_WRITE": "Deny",
		},
	})
	tfNodeRoot := "azuredevops_area_permissions.root-permissions"

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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE_CHILDREN", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.GENERIC_READ", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.WORK_ITEM_WRITE", "deny"),
				),
			},
		},
	})
}

func TestAccAreaPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclAreaPermissions(projectName, map[string]map[string]string{
		"root": {
			"CREATE_CHILDREN": "Deny",
			"GENERIC_READ":    "NotSet",
			"DELETE":          "Deny",
			"WORK_ITEM_WRITE": "Deny",
		},
	})
	config2 := hclAreaPermissions(projectName, map[string]map[string]string{
		"root": {
			"CREATE_CHILDREN": "Deny",
			"GENERIC_READ":    "Allow",
			"DELETE":          "Deny",
			"WORK_ITEM_WRITE": "Deny",
		},
	})
	tfNodeRoot := "azuredevops_area_permissions.root-permissions"

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
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE_CHILDREN", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.GENERIC_READ", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.WORK_ITEM_WRITE", "deny"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE_CHILDREN", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.WORK_ITEM_WRITE", "deny"),
				),
			},
		},
	})
}
