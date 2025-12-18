package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessProcessPermissions_SetPermissions(t *testing.T) {
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_process_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclProcessPermissions(processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
				),
			},
		},
	})
}

func hclProcessPermissions(processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

data "azuredevops_group" "project-collection-administrators" {
  name = "Project Collection Administrators"
}

resource "azuredevops_workitemtrackingprocess_process_permissions" "test" {
  process_id = azuredevops_workitemtrackingprocess_process.test.id
  principal  = data.azuredevops_group.project-collection-administrators.id
  permissions = {
    "Edit"                         = "Allow"
    "Delete"                       = "Deny"
    "Create"                       = "NotSet"
    "AdministerProcessPermissions" = "Allow"
  }
}
`, processName, agileSystemProcessTypeId)
}
