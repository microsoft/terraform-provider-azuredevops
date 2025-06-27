//go:build (all || permissions || resource_build_definition_permissions) && (!exclude_permissions || !exclude_resource_build_definition_permissions)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func TestAccBuildDefinitionPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclBuildDefinitionPermissions(projectName, `\`, map[string]string{
		"ViewBuilds":         "Allow",
		"EditBuildQuality":   "NotSet",
		"RetainIndefinitely": "Deny",
		"DeleteBuilds":       "Deny",
	})
	tfNodeRoot := "azuredevops_build_definition_permissions.permissions"

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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewBuilds", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditBuildQuality", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.RetainIndefinitely", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteBuilds", "deny"),
				),
			},
		},
	})
}

func TestAccBuildDefinitionPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclBuildDefinitionPermissions(projectName, `\`, map[string]string{
		"ViewBuilds":         "Deny",
		"EditBuildQuality":   "NotSet",
		"RetainIndefinitely": "Deny",
		"DeleteBuilds":       "Deny",
	})
	config2 := hclBuildDefinitionPermissions(projectName, `\dir1\dir2`, map[string]string{
		"ViewBuilds":         "Deny",
		"EditBuildQuality":   "Allow",
		"RetainIndefinitely": "Deny",
		"DeleteBuilds":       "Deny",
	})
	tfNodeRoot := "azuredevops_build_definition_permissions.permissions"

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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewBuilds", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditBuildQuality", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.RetainIndefinitely", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteBuilds", "deny"),
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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewBuilds", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditBuildQuality", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.RetainIndefinitely", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteBuilds", "deny"),
				),
			},
		},
	})
}

func hclBuildDefinitionPermissions(projectName string, path string, permissions map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions, "=", "\n")
	buildDefinitionNameFirst := testutils.GenerateResourceName()

	return fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_build_definition_permissions" "permissions" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id

  build_definition_id = azuredevops_build_definition.build.id

  permissions = {
		%s
  }
}
`, testutils.HclBuildDefinitionResourceGitHub(projectName, buildDefinitionNameFirst, path),
		rootPermissions,
	)
}
