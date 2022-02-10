//go:build (all || data_sources || data_variable_group) && (!exclude_data_sources || !exclude_data_variable_group)
// +build all data_sources data_variable_group
// +build !exclude_data_sources !exclude_data_variable_group

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccVariableGroup_DataSource(t *testing.T) {
	variableGroupName := testutils.GenerateResourceName()
	createVariableGroup := testutils.HclVariableGroupResource(variableGroupName, true)
	createAndGetVariableGroupData := fmt.Sprintf("%s\n%s", createVariableGroup, testutils.HclVariableGroupDataSource())

	tfNode := "data.azuredevops_variable_group.vg"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createAndGetVariableGroupData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", variableGroupName),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "variable.#"),
				),
			},
		},
	})
}
