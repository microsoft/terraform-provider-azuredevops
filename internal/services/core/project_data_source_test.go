package core_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance"
)

type ProjectDataSource struct{}

func TestAccDataSourceProject_withName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuredevops_project", "test")
	d := ProjectDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.withName(data),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "project_id"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "description"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "visibility", "private"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "version_control", "Git"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "work_item_template", "Basic"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "process_template_id"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "features.%", "5"),
			),
		},
	})
}

func TestAccDataSourceProject_withID(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuredevops_project", "test")
	d := ProjectDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.withID(data),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "project_id"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "description"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "visibility", "private"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "version_control", "Git"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "work_item_template", "Basic"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "process_template_id"),
				resource.TestCheckResourceAttr(data.ResourceAddr(), "features.%", "5"),
			),
		},
	})
}

func (d ProjectDataSource) withName(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  description        = "foo"
}

data "azuredevops_project" "test" {
  project_id = azuredevops_project.test.id
}
`, data.RandomString)
}

func (d ProjectDataSource) withID(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  description        = "foo"
}

data "azuredevops_project" "test" {
  name = azuredevops_project.test.name
}
`, data.RandomString)
}
