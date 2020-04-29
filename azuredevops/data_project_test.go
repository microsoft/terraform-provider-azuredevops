// +build all core resource_project

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
)

/**
 * Begin acceptance tests
 */

// Verifies that the following sequence of events occurrs without error:
//	(1) TF can create a project
//	(2) A data source is added to the configuration, and that data source can find the created project
func TestAccAzureDevOpsProject_DataSource(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfConfigStep1 := testhelper.TestAccProjectResource(projectName)
	tfConfigStep2 := fmt.Sprintf("%s\n%s", tfConfigStep1, testhelper.TestAccProjectDataSource(projectName))

	tfNode := "data.azuredevops_project.project"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testhelper.TestAccPreCheck(t, nil) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tfConfigStep1,
			}, {
				Config: tfConfigStep2,
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

func init() {
	InitProvider()
}
