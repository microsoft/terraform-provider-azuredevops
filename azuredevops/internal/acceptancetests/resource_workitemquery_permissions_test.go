// +build all permissions resource_workitemquery_permissions
// +build !exclude_permissions !resource_workitemquery_permissions

package acceptancetests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkItemQueryPermissions_SetProjectPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := testutils.HclWorkItemQueryPermissions(projectName, "")

	tfNode := "azuredevops_workitemquery_permissions.wiq-permissions"
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
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
				),
			},
		},
	})
}

func TestAccWorkItemQueryPermissions_SetSharedQueriesPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := testutils.HclWorkItemQueryPermissions(projectName, "/")

	tfNode := "azuredevops_workitemquery_permissions.wiq-permissions"
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
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
				),
			},
		},
	})
}

func TestAccWorkItemQueryPermissions_SetInvalidFolderPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := testutils.HclWorkItemQueryPermissions(projectName, "/invalid")

	tfNode := "azuredevops_workitemquery_permissions.wiq-permissions"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Unable to find query"),
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
