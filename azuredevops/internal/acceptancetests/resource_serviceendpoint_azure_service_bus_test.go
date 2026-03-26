package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAzureServiceBus_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azure_service_bus"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAzureServiceBusResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_string", "connectionstring"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "queue_name", "testqueue"),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureServiceBus_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azure_service_bus"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAzureServiceBusResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_string", "connectionstring"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "queue_name", "testqueue"),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureServiceBus_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azure_service_bus"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAzureServiceBusResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_string", "connectionstring"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "queue_name", "testqueue"),
				),
			},
			{
				Config: hclSvcEndpointAzureServiceBusResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_string", "connectionstringupdate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "queue_name", "testqueueupdate"),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureServiceBus_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_azure_service_bus"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAzureServiceBusResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointAzureServiceBusResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointAzureServiceBusResourceBasic(projectName string, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_azure_service_bus" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  connection_string     = "connectionstring"
  queue_name            = "testqueue"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointAzureServiceBusResourceComplete(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_azure_service_bus" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  description           = "%s"
  connection_string     = "connectionstring"
  queue_name            = "testqueue"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointAzureServiceBusResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_azure_service_bus" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  description           = "%s"
  connection_string     = "connectionstringupdate"
  queue_name            = "testqueueupdate"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointAzureServiceBusResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointAzureServiceBusResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_azure_service_bus" "import" {
  project_id            = azuredevops_serviceendpoint_azure_service_bus.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_azure_service_bus.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_azure_service_bus.test.description
  connection_string     = azuredevops_serviceendpoint_azure_service_bus.test.connection_string
  queue_name            = azuredevops_serviceendpoint_azure_service_bus.test.queue_name
}
`, template)
}
