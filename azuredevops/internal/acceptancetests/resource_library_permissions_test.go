//go:build (all || permissions || resource_library_permissions) && (!exclude_permissions || !exclude_resource_library_permissions)
// +build all permissions resource_library_permissions
// +build !exclude_permissions !exclude_resource_library_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func TestAccLibraryPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclLibraryPermissions(projectName, map[string]string{
		"View":        "allow",
		"Administer":  "allow",
		"Create":      "allow",
		"ViewSecrets": "notset",
		"Use":         "allow",
		"Owner":       "allow",
	})
	tfNode := "azuredevops_library_permissions.permissions"

	resource.Test(t, resource.TestCase{
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
					resource.TestCheckNoResourceAttr(tfNode, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "6"),
					resource.TestCheckResourceAttr(tfNode, "permissions.View", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Administer", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Create", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ViewSecrets", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Owner", "allow"),
				),
			},
		},
	})
}

func TestAccLibraryPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclLibraryPermissions(projectName, map[string]string{
		"View":        "allow",
		"Administer":  "allow",
		"Create":      "allow",
		"ViewSecrets": "notset",
		"Use":         "allow",
		"Owner":       "allow",
	})
	config2 := hclLibraryPermissions(projectName, map[string]string{
		"View":        "allow",
		"Administer":  "notset",
		"Create":      "notset",
		"ViewSecrets": "notset",
		"Use":         "notset",
		"Owner":       "notset",
	})
	tfNode := "azuredevops_library_permissions.permissions"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "6"),
					resource.TestCheckResourceAttr(tfNode, "permissions.View", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Administer", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Create", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ViewSecrets", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Owner", "allow"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "6"),
					resource.TestCheckResourceAttr(tfNode, "permissions.View", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Administer", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Create", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ViewSecrets", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Use", "notset"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Owner", "notset"),
				),
			},
		},
	})
}

func hclLibraryPermissions(projectName string, permissions map[string]string) string {
	LibraryPermissions := datahelper.JoinMap(permissions, "=", "\n")

	return fmt.Sprintf(`
%s
data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_library_permissions" "permissions" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id
  permissions = {
		%s
  }
}


`, testutils.HclProjectResource(projectName),
		LibraryPermissions,
	)
}
