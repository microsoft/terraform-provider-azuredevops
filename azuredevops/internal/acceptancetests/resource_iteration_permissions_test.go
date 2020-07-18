// +build all permissions resource_iteration_permissions
// +build !exclude_permissions !exclude_resource_iteration_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIterationPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_iteration_permissions" "root-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	permissions = {
	  CREATE_CHILDREN = "Deny"
	  GENERIC_READ    = "NotSet"
	  DELETE          = "Deny"
	}
}

resource "azuredevops_iteration_permissions" "iteration-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	path        = "Iteration 1"
	permissions = {
	  CREATE_CHILDREN = "Allow"
	  GENERIC_READ    = "NotSet"
	  DELETE          = "Allow"
	}
}

`, testutils.HclProjectResource(projectName))

	tfNode := "azuredevops_iteration_permissions.iteration-permissions"
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
					resource.TestCheckResourceAttr(tfNode, "path", "Iteration 1"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
				),
			},
		},
	})
}
