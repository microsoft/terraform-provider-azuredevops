// +build all core data_sources resource_project data_project
// +build !exclude_data_sources !exclude_data_project

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Verifies that the following sequence of events occurrs without error:
//	(1) TF can create a project
//	(2) A data source is added to the configuration, and that data source can find the created project
func TestAccProject_DataSource(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	projectData := testutils.HclProjectDataSource(projectName)

	tfNode := "data.azuredevops_project.project"
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: projectData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "project_name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
				),
			},
		},
	})
}

func TestAccProject_DataSource_IncorrectParameters(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	dataProject := fmt.Sprintf(`
	data "azuredevops_project" "project" {
		project_name = "%s"
		description = "A project description"
	}`, projectName)

	errorRegex, _ := regexp.Compile("config is invalid: \"description\": this field cannot be set")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config:      dataProject,
				ExpectError: errorRegex,
			},
		},
	})
}
