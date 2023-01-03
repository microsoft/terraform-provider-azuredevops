//go:build (all || core || data_sources || resource_project || data_project) && (!exclude_data_sources || !exclude_data_project)
// +build all core data_sources resource_project data_project
// +build !exclude_data_sources !exclude_data_project

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Verifies that the following sequence of events occurrs without error:
//
//	(1) TF can create a project
//	(2) A data source is added to the configuration, and that data source can find the created project
func TestAccProject_DataSource(t *testing.T) {
	var tests = []struct {
		Name                string
		Identifier          string
		IdentifierOnProject string
	}{
		{
			Name:                "Get project with id",
			Identifier:          "project_id",
			IdentifierOnProject: "id",
		},
		{
			Name:                "Get project with name",
			Identifier:          "name",
			IdentifierOnProject: "name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			projectName := testutils.GenerateResourceName()
			projectData := fmt.Sprintf(`
			%s

			data "azuredevops_project" "project" {
				%s = azuredevops_project.project.%s
			}`, testutils.HclProjectResource(projectName), tc.Identifier, tc.IdentifierOnProject)

			tfNode := "data.azuredevops_project.project"
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                  func() { testutils.PreCheck(t, nil) },
				ProviderFactories:         testutils.GetProviderFactories(),
				PreventPostDestroyRefresh: true,
				Steps: []resource.TestStep{
					{
						Config: projectData,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
							resource.TestCheckResourceAttr(tfNode, "name", projectName),
							resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
							resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
							resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
						),
					},
				},
			})

		})
	}
}

func TestAccProject_DataSource_ErrorWhenBothNameAndIdSet(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      errorWithNameAndIdSet(),
				ExpectError: regexp.MustCompile(`Either project_id or name must be set`),
			},
		},
	})
}

func TestAccProject_DataSource_ErrorWhenDescriptionSet(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      errorWhenDescriptionSet(testutils.GenerateResourceName()),
				ExpectError: regexp.MustCompile(`Value for unconfigurable attribute`),
			},
		},
	})
}

func errorWithNameAndIdSet() string {
	return `data "azuredevops_project" "project" {}`
}

func errorWhenDescriptionSet(projectName string) string {
	return fmt.Sprintf(`
	data "azuredevops_project" "project" {
		name = "%s"
		description = "A project description"
	}`, projectName)
}
