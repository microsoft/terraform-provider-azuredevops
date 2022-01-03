//go:build (all || permissions || resource_tagging_permissions) && (!exclude_permissions || !exclude_resource_tagging_permissions)
// +build all permissions resource_tagging_permissions
// +build !exclude_permissions !exclude_resource_tagging_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func hclTaggingPermissions(projectName string, permissions map[string]map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions["root"], "=", "\n")

	return fmt.Sprintf(`
%s
data "azuredevops_group" "tf-project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}
resource "azuredevops_tagging_permissions" "acctest" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	permissions = {
		%s
	}
}
`, testutils.HclProjectResource(projectName), rootPermissions)
}

func TestAccTaggingPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclTaggingPermissions(projectName, map[string]map[string]string{
		"root": {
			"Enumerate": "Deny",
			"Create":    "NotSet",
			"Update":    "Deny",
			"Delete":    "Deny",
		},
	})
	tfNodeRoot := "azuredevops_tagging_permissions.acctest"

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

func TestAccTaggingPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclTaggingPermissions(projectName, map[string]map[string]string{
		"root": {
			"Enumerate": "Allow",
			"Create":    "NotSet",
			"Update":    "Deny",
			"Delete":    "Deny",
		},
	})
	config2 := hclTaggingPermissions(projectName, map[string]map[string]string{
		"root": {
			"Enumerate": "Deny",
			"Create":    "Deny",
			"Update":    "NotSet",
			"Delete":    "Allow",
		},
	})
	tfNodeRoot := "azuredevops_tagging_permissions.acctest"

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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Enumerate", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Create", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Update", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Delete", "deny"),
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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Enumerate", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Create", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Update", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Delete", "allow"),
				),
			},
		},
	})
}
