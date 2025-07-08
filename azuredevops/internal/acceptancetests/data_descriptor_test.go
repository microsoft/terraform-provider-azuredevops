//go:build (all || data_sources || data_descriptor) && (!exclude_data_sources || !exclude_data_descriptor)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccDescriptorDatasource_user(t *testing.T) {
	name := testutils.GenerateResourceName() + "@contoso.com"
	tfNode := "data.azuredevops_descriptor.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		ProviderFactories:         testutils.GetProviderFactories(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hclDescriptorDataSourceUser(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func TestAccDescriptorDatasource_project(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_descriptor.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		ProviderFactories:         testutils.GetProviderFactories(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hclDescriptorDataSourceProject(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func TestAccDescriptorDatasource_group(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_descriptor.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		ProviderFactories:         testutils.GetProviderFactories(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hclDescriptorDataSourceGroup(projectName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func hclDescriptorDataSourceUser(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%s"
  account_license_type = "express"
}

data "azuredevops_descriptor" "test" {
  storage_key = azuredevops_user_entitlement.test.id
}`, name)
}

func hclDescriptorDataSourceProject(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_descriptor" "test" {
  storage_key = azuredevops_project.test.id
}`, projectName)
}

func hclDescriptorDataSourceGroup(projectName, groupName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_group" "test" {
  scope        = azuredevops_project.test.id
  display_name = "%s"
}

data "azuredevops_descriptor" "test" {
  storage_key = azuredevops_project.test.id
}`, projectName, groupName)
}
