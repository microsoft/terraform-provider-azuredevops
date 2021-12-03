//go:build (all || core || data_sources || resource_project || data_projects) && (!data_sources || !exclude_data_projects)
// +build all core data_sources resource_project data_projects
// +build !data_sources !exclude_data_projects

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAzureDevOpsProjects_DataSource_SingleProject(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	projectData := testutils.HclProjectsDataSource(projectName)

	tfNode := "data.azuredevops_projects.project-list"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: projectData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "projects.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureDevOpsProjects_DataSource_EmptyResult(t *testing.T) {
	projectData := testutils.HclProjectsDataSourceWithStateAndInvalidName()

	tfNode := "data.azuredevops_projects.project-list"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: projectData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "projects.#", "0"),
				),
			},
		},
	})
}
