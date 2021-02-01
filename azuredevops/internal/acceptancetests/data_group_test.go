// +build all core data_sources data_group
// +build !exclude_data_sources !exclude_data_group

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Validates that a configuration containing a project group lookup is able to read the resource correctly.
// Because this is a data source, there are no resources to inspect in AzDO
func TestAccGroupDataSource_Read_HappyPath(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	group := "Build Administrators"
	tfBuildDefNode := "data.azuredevops_group.group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGroupDataSource(projectName, group),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "name"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin_id"),
				),
			},
		},
	})
}

func TestAccGroupDataSource_Read_ProjectCollectionAdministrators(t *testing.T) {
	group := "Project Collection Administrators"
	tfBuildDefNode := "data.azuredevops_group.group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGroupDataSource("", group),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "name"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin_id"),
				),
			},
		},
	})
}
