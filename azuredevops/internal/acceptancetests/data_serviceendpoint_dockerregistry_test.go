//go:build (all || data_sources || data_serviceendpoint_dockerregistry) && (!exclude_data_sources || !exclude_data_serviceendpoint_dockerregistry)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointDockerRegistry_data_withName(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "data.azuredevops_serviceendpoint_dockerregistry"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclDataServiceConnectionDockerRegistryWithName(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointDockerRegistry_data_withID(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "data.azuredevops_serviceendpoint_dockerregistry"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclDataServiceConnectionDockerRegistryWithID(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_endpoint_id"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
		},
	})
}

func hclDataServiceConnectionDockerRegistryWithName(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_dockerregistry" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  docker_email          = "test@email.com"
  docker_username       = "testuser"
  docker_password       = "secret"
}

data "azuredevops_serviceendpoint_dockerregistry" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = azuredevops_serviceendpoint_dockerregistry.test.service_endpoint_name
}
`, projectName, serviceEndpointName)
}

func hclDataServiceConnectionDockerRegistryWithID(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_dockerregistry" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  docker_email          = "test@email.com"
  docker_username       = "testuser"
  docker_password       = "secret"
}

data "azuredevops_serviceendpoint_dockerregistry" "test" {
  project_id          = azuredevops_project.test.id
  service_endpoint_id = azuredevops_serviceendpoint_dockerregistry.test.id
}
`, projectName, serviceEndpointName)
}
