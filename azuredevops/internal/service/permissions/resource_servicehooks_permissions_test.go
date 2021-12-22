//go:build (all || permissions || resource_servicehooks_permissions) && (!exclude_permissions || !exclude_resource_servicehooks_permissions)
// +build all permissions resource_servicehooks_permissions
// +build !exclude_permissions !exclude_resource_servicehooks_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func hclServiceHooksPermissions(projectName string, permissions map[string]map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions["root"], "=", "\n")

	return fmt.Sprintf(`
%s
data "azuredevops_group" "tf-project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}
resource "azuredevops_servicehooks_permissions" "acctest" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	permissions = {
		%s
	}
}
`, testutils.HclProjectResource(projectName), rootPermissions)
}

func TestAccServiceHooksPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclServiceHooksPermissions(projectName, map[string]map[string]string{
		"root": {
			"ViewSubscriptions":   "Deny",
			"EditSubscriptions":   "NotSet",
			"DeleteSubscriptions": "Deny",
			"PublishEvents":       "Deny",
		},
	})
	tfNodeRoot := "azuredevops_servicehooks_permissions.acctest"

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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewSubscriptions", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditSubscriptions", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteSubscriptions", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.PublishEvents", "deny"),
				),
			},
		},
	})
}

func TestAccServiceHooksPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclServiceHooksPermissions(projectName, map[string]map[string]string{
		"root": {
			"ViewSubscriptions":   "Allow",
			"EditSubscriptions":   "NotSet",
			"DeleteSubscriptions": "Deny",
			"PublishEvents":       "Deny",
		},
	})
	config2 := hclServiceHooksPermissions(projectName, map[string]map[string]string{
		"root": {
			"ViewSubscriptions":   "Deny",
			"EditSubscriptions":   "Deny",
			"DeleteSubscriptions": "NotSet",
			"PublishEvents":       "Allow",
		},
	})
	tfNodeRoot := "azuredevops_servicehooks_permissions.acctest"

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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewSubscriptions", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditSubscriptions", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteSubscriptions", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.PublishEvents", "deny"),
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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewSubscriptions", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditSubscriptions", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteSubscriptions", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.PublishEvents", "allow"),
				),
			},
		},
	})
}
