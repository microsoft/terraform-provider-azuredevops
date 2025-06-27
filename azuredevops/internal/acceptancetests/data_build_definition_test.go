//go:build (all || data_sources || data_build_definition) && (!exclude_data_sources || !exclude_data_build_definition)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBuildDefinition_DataSource(t *testing.T) {
	buildDefinitionName := testutils.GenerateResourceName()
	createAndGetVariableGroupData := fmt.Sprintf("%s\n%s",
		testutils.HclBuildDefinitionWithVariables("foo1", "bar1", buildDefinitionName),
		testutils.HclBuildDefinitionDataSource(`\\`)) // `\\` is the default value for the path

	tfNode := "data.azuredevops_build_definition.build"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createAndGetVariableGroupData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", buildDefinitionName),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
		},
	})
}

func TestAccBuildDefinition_with_path_DataSource(t *testing.T) {
	name := testutils.GenerateResourceName()
	createAndGetVariableGroupData := fmt.Sprintf("%s\n%s",
		testutils.HclBuildDefinitionResourceGitHub(name, name, "\\some\\path"),
		testutils.HclBuildDefinitionDataSource(`\\some\\path`))

	tfNode := "data.azuredevops_build_definition.build"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createAndGetVariableGroupData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "path", "\\some\\path"),
				),
			},
		},
	})
}
