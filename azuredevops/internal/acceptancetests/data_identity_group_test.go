//go:build (all || core || data_sources || data_group) && (!exclude_data_sources || !exclude_data_group)
// +build all core data_sources data_group
// +build !exclude_data_sources !exclude_data_group

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func testIdentityGroupDataSource(t *testing.T, groupName string) {
	tfNode := "data.azuredevops_identity_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createIdentityGroupConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", groupName),
				),
			},
		},
	})
}

func createIdentityGroupConfig(groupName string) string {
	return fmt.Sprintf(
		`
data "azuredevops_project" "project" {
	name = "default"
}

data "azuredevops_identity_group" "test" {
	name       = "%[1]s"
	project_id = data.azuredevops_project.project.id
}`, groupName)
}

func TestAccIdentityGroupDataSource(t *testing.T) {
	groupName := "[default]\\\\Contributors"
	testIdentityGroupDataSource(t, groupName)
}
