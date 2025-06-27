//go:build (all || data_sources || data_serviceendpoint_github) && (!exclude_data_sources || !exclude_data_serviceendpoint_github)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointBitbucketDataSource_withID(t *testing.T) {
	serviceEndpointName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_serviceendpoint_bitbucket.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceServiceConnectionBitbucketWithID(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
				),
			},
		},
	})
}

func TestAccServiceEndpointBitbucketDataSource_withName(t *testing.T) {
	serviceEndpointName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_serviceendpoint_bitbucket.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceServiceConnectionBitbucketWithName(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
				),
			},
		},
	})
}

func hclDataSourceServiceConnectionBitbucketWithID(projectName, serviceConnectionName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_bitbucket" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  username              = "username"
  password              = "password"
}

data "azuredevops_serviceendpoint_bitbucket" "test" {
  project_id          = azuredevops_project.test.id
  service_endpoint_id = azuredevops_serviceendpoint_bitbucket.test.id
}
`, projectName, serviceConnectionName)
}

func hclDataSourceServiceConnectionBitbucketWithName(projectName, serviceConnectionName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_bitbucket" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  username              = "username"
  password              = "password"
}

data "azuredevops_serviceendpoint_bitbucket" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = azuredevops_serviceendpoint_bitbucket.test.service_endpoint_name
}
`, projectName, serviceConnectionName)
}
