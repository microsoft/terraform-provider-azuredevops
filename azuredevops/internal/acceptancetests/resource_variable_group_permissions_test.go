package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func TestAccVariableGroupPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	variableGroupName := testutils.GenerateResourceName()
	config := hclVariableGroupPermissions(projectName, variableGroupName, map[string]string{
		"View":        "allow",
		"Administer":  "allow",
		"Create":      "allow",
		"ViewSecrets": "notset",
		"Use":         "allow",
		"Owner":       "allow",
	})
	tfNode := "azuredevops_variable_group_permissions.permissions"

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
					resource.TestCheckResourceAttrSet(tfNode, "variable_group_id"),
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

func TestAccVariableGroupPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	variableGroupName := testutils.GenerateResourceName()
	config1 := hclVariableGroupPermissions(projectName, variableGroupName, map[string]string{
		"View":        "allow",
		"Administer":  "allow",
		"Create":      "allow",
		"ViewSecrets": "notset",
		"Use":         "allow",
		"Owner":       "allow",
	})
	config2 := hclVariableGroupPermissions(projectName, variableGroupName, map[string]string{
		"View":        "allow",
		"Administer":  "notset",
		"Create":      "notset",
		"ViewSecrets": "notset",
		"Use":         "notset",
		"Owner":       "notset",
	})
	tfNode := "azuredevops_variable_group_permissions.permissions"

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
					resource.TestCheckResourceAttrSet(tfNode, "variable_group_id"),
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
					resource.TestCheckResourceAttrSet(tfNode, "variable_group_id"),
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

func hclVariableGroupPermissions(projectName string, variableGroupName string, permissions map[string]string) string {
	variableGroupPermissions := datahelper.JoinMap(permissions, "=", "\n")

	return fmt.Sprintf(`
%s

resource "azuredevops_variable_group" "example" {
  project_id   = azuredevops_project.project.id
  name         = "%s"
  description  = "Test Description"
  allow_access = true

  variable {
    name  = "key1"
    value = "val1"
  }
}

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_variable_group_permissions" "permissions" {
  project_id        = azuredevops_project.project.id
  variable_group_id = azuredevops_variable_group.example.id
  principal         = data.azuredevops_group.tf-project-readers.id
  permissions = {
		%s
  }
}


`, testutils.HclProjectResource(projectName),
		variableGroupName,
		variableGroupPermissions,
	)
}
