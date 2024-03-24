//go:build (all || core || resource_securityrole_assignment) && !exclude_resource_securityrole_assignment
// +build all core resource_securityrole_assignment
// +build !exclude_resource_securityrole_assignment

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccSecurityroleAssignmentResource_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclSecurityroleAssignmentBasic(projectName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_securityrole_assignment.test", "scope"),
					resource.TestCheckResourceAttr("azuredevops_securityrole_assignment.test", "role_name", "Administrator"),
				),
			},
		},
	})
}

func hclSecurityroleAssignmentBasic(projectName, groupName string) string {

	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
name               = "%[1]s"
description        = "%[1]s-description"
visibility         = "private"
version_control    = "Git"
work_item_template = "Agile"
}

resource "azuredevops_group" "test" {
scope        = azuredevops_project.test.id
display_name = "%[2]s"
}

resource "azuredevops_environment" "test" {
project_id  = azuredevops_project.test.id
name        = "Example Environment"
description = "Example pipeline deployment environment"
}

resource "azuredevops_securityrole_assignment" "test" {
scope       = "distributedtask.environmentreferencerole"
resource_id = format("%s_%s", azuredevops_project.test.id, azuredevops_environment.test.id)
identity_id = azuredevops_group.test.id
role_name   = "Administrator"
}
`, projectName, groupName)
}
