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

func testIdentityGroupDataSource(t *testing.T, groupName string, projectName string) {
	tfNode := "data.azuredevops_identity_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createIdentityGroupConfig(groupName, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", fmt.Sprintf("[%s]\\%s", projectName, groupName)),
				),
			},
		},
	})
}

func createIdentityGroupConfig(groupName string, projectName string) string {
	combinedgroupName := fmt.Sprintf("[%s]\\\\%s", projectName, groupName)
	dataSource := fmt.Sprintf(
		`
data "azuredevops_identity_group" "test" {
	name       = "%s"
	project_id = azuredevops_project.project.id
}`, combinedgroupName)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, dataSource)
}

func TestAccIdentityGroupDataSource(t *testing.T) {
	groupName := "Contributors"
	projectName := testutils.GenerateResourceName()
	testIdentityGroupDataSource(t, groupName, projectName)
}
