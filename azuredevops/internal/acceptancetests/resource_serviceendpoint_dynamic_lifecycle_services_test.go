//go:build (all || resource_serviceendpoint_dynamic_lifecycle_services) && !exclude_resource_serviceendpoint_dynamic_lifecycle_services

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointDynamicLifecycleServices_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_dynamics_lifecycle_services"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointDynamicLifecycleServicesResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_endpoint", "https://login.microsoftonline.com/organization"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "lifecycle_services_api_endpoint", "https://lcsapi.lcs.dynamics.com"),
				),
			},
		},
	})
}

func TestAccServiceEndpointDynamicLifecycleServices_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_dynamics_lifecycle_services"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointDynamicLifecycleServicesResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_endpoint", "https://login.microsoftonline.com/organization"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "lifecycle_services_api_endpoint", "https://lcsapi.lcs.dynamics.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointDynamicLifecycleServices_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_dynamics_lifecycle_services"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointDynamicLifecycleServicesResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointDynamicLifecycleServicesResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_endpoint", "https://login.microsoftonline.com/organization/update/"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "lifecycle_services_api_endpoint", "https://lcsapi.lcs.dynamics.com/update/"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "client_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "usernameupdate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointDynamicLifecycleServices_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_dynamics_lifecycle_services"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointDynamicLifecycleServicesResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointDynamicLifecycleServicesResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointDynamicLifecycleServicesResourceBasic(projectName string, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_dynamics_lifecycle_services" "test" {
  project_id                      = azuredevops_project.test.id
  service_endpoint_name           = "%s"
  authorization_endpoint          = "https://login.microsoftonline.com/organization"
  lifecycle_services_api_endpoint = "https://lcsapi.lcs.dynamics.com"
  client_id                       = "00000000-0000-0000-0000-000000000000"
  username                        = "username"
  password                        = "password"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointDynamicLifecycleServicesResourceComplete(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_dynamics_lifecycle_services" "test" {
  project_id                      = azuredevops_project.test.id
  service_endpoint_name           = "%s"
  description                     = "%s"
  authorization_endpoint          = "https://login.microsoftonline.com/organization"
  lifecycle_services_api_endpoint = "https://lcsapi.lcs.dynamics.com"
  client_id                       = "00000000-0000-0000-0000-000000000000"
  username                        = "username"
  password                        = "password"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointDynamicLifecycleServicesResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_dynamics_lifecycle_services" "test" {
  project_id                      = azuredevops_project.test.id
  service_endpoint_name           = "%s"
  description                     = "%s"
  authorization_endpoint          = "https://login.microsoftonline.com/organization/update/"
  lifecycle_services_api_endpoint = "https://lcsapi.lcs.dynamics.com/update/"
  client_id                       = "00000000-0000-0000-0000-000000000002"
  username                        = "usernameupdate"
  password                        = "passwordupdate"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointDynamicLifecycleServicesResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointDynamicLifecycleServicesResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_dynamics_lifecycle_services" "import" {
  project_id                      = azuredevops_serviceendpoint_dynamics_lifecycle_services.test.project_id
  service_endpoint_name           = azuredevops_serviceendpoint_dynamics_lifecycle_services.test.service_endpoint_name
  authorization_endpoint          = azuredevops_serviceendpoint_dynamics_lifecycle_services.test.authorization_endpoint
  lifecycle_services_api_endpoint = azuredevops_serviceendpoint_dynamics_lifecycle_services.test.lifecycle_services_api_endpoint
  client_id                       = azuredevops_serviceendpoint_dynamics_lifecycle_services.test.client_id
  username                        = azuredevops_serviceendpoint_dynamics_lifecycle_services.test.username
  password                        = azuredevops_serviceendpoint_dynamics_lifecycle_services.test.password
}
`, template)
}
