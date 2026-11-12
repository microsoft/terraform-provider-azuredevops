package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func TestAccAnalyticsPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclAnalyticsPermissions(projectName, map[string]string{
		"Read": "Allow",
	})
	tfNode := "azuredevops_analytics_permissions.acctest"

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
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "1"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
				),
			},
		},
	})
}

func TestAccAnalyticsPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclAnalyticsPermissions(projectName, map[string]string{
		"Read": "Allow",
	})
	config2 := hclAnalyticsPermissions(projectName, map[string]string{
		"Read": "Deny",
	})
	tfNode := "azuredevops_analytics_permissions.acctest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "1"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "1"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "deny"),
				),
			},
		},
	})
}

func TestAccAnalyticsViewsPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclAnalyticsViewsPermissions(projectName, map[string]string{
		"Read":   "Allow",
		"Edit":   "Deny",
		"Delete": "Deny",
	})
	tfNode := "azuredevops_analytics_views_permissions.acctest"

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
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Edit", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "deny"),
				),
			},
		},
	})
}

func TestAccAnalyticsViewsPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclAnalyticsViewsPermissions(projectName, map[string]string{
		"Read":   "Allow",
		"Edit":   "Deny",
		"Delete": "Deny",
	})
	config2 := hclAnalyticsViewsPermissions(projectName, map[string]string{
		"Read":   "Deny",
		"Edit":   "Allow",
		"Delete": "NotSet",
	})
	tfNode := "azuredevops_analytics_views_permissions.acctest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Edit", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "deny"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Edit", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Delete", "notset"),
				),
			},
		},
	})
}

func hclAnalyticsPermissions(projectName string, permissions map[string]string) string {
	return fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_analytics_permissions" "acctest" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id
  permissions = {
    %s
  }
}
`, testutils.HclProjectResource(projectName), datahelper.JoinMap(permissions, "=", "\n"))
}

func hclAnalyticsViewsPermissions(projectName string, permissions map[string]string) string {
	return fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_analytics_views_permissions" "acctest" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id
  permissions = {
    %s
  }
}
`, testutils.HclProjectResource(projectName), datahelper.JoinMap(permissions, "=", "\n"))
}
