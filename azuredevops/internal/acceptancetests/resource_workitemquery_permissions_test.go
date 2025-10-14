package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func TestAccWorkItemQueryPermissions_SetProjectPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	permissions := map[string]string{
		"Contribute":        "Allow",
		"Delete":            "Deny",
		"ManagePermissions": "NotSet",
	}
	config := hclWorkItemQueryPermissions(projectName, "", permissions)

	tfNode := "azuredevops_workitemquery_permissions.wiq-permissions"
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
					resource.TestCheckResourceAttr(tfNode, "permissions.Contribute", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManagePermissions", "notset"),
				),
			},
		},
	})
}

func TestAccWorkItemQueryPermissions_UpdateProjectPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclWorkItemQueryPermissions(projectName, "", map[string]string{
		"Contribute":        "Allow",
		"Delete":            "Deny",
		"ManagePermissions": "NotSet",
	})
	config2 := hclWorkItemQueryPermissions(projectName, "", map[string]string{
		"Contribute":        "Deny",
		"Delete":            "Allow",
		"ManagePermissions": "Deny",
	})

	tfNode := "azuredevops_workitemquery_permissions.wiq-permissions"
	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Contribute", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManagePermissions", "notset"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Contribute", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManagePermissions", "deny"),
				),
			},
		},
	})
}

func TestAccWorkItemQueryPermissions_SetSharedQueriesPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	permissions := map[string]string{
		"Contribute":        "Allow",
		"Delete":            "Deny",
		"ManagePermissions": "NotSet",
	}
	config := hclWorkItemQueryPermissions(projectName, "/", permissions)

	tfNode := "azuredevops_workitemquery_permissions.wiq-permissions"
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
					resource.TestCheckResourceAttr(tfNode, "permissions.Contribute", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManagePermissions", "notset"),
				),
			},
		},
	})
}

func TestAccWorkItemQueryPermissions_SetInvalidFolderPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	permissions := map[string]string{
		"Contribute":        "Allow",
		"Delete":            "Deny",
		"ManagePermissions": "NotSet",
	}
	config := hclWorkItemQueryPermissions(projectName, "/invalid", permissions)

	tfNode := "azuredevops_workitemquery_permissions.wiq-permissions"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Unable to find query"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckNoResourceAttr(tfNode, "project_id"),
					resource.TestCheckNoResourceAttr(tfNode, "principal"),
					resource.TestCheckNoResourceAttr(tfNode, "permissions.Contribute"),
					resource.TestCheckNoResourceAttr(tfNode, "permissions.Delete"),
					resource.TestCheckNoResourceAttr(tfNode, "permissions.ManagePermissions"),
				),
			},
		},
	})
}

func hclWorkItemQueryPermissions(projectName string, path string, permissions map[string]string) string {
	projectResource := testutils.HclProjectResource(projectName)
	szPermissions := datahelper.JoinMap(permissions, "=", "\n")
	szPath := ""
	if path != "" {
		szPath = fmt.Sprintf("path = \"%s\"", path)
	}

	return fmt.Sprintf(`
%s

data "azuredevops_group" "project-administrators" {
  project_id = azuredevops_project.project.id
  name       = "Project administrators"
}

resource "azuredevops_workitemquery_permissions" "wiq-permissions" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.project-administrators.id
	%s
  permissions = {
		%s
  }
}
`, projectResource, szPath, szPermissions)
}
