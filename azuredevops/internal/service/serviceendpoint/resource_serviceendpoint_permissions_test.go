package serviceendpoint

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Verifies that the following sequence of operations can be performed without error:
//	1. Creating a service endpoint with project permissions
//	2. Updating the service endpoint with different project permissions
func TestAccServiceEndpointProjectPermissions_CRUD(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	config := hclServiceEndpointProjectPermissionsResource(projectName, serviceEndpointName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_serviceendpoint_project_permissions.p", "project_id"),
					resource.TestCheckResourceAttrSet("azuredevops_serviceendpoint_project_permissions.p", "serviceendpoint_id"),
					resource.TestCheckResourceAttr("azuredevops_serviceendpoint_project_permissions.p", "project_reference.#", "2"),
					resource.TestCheckResourceAttr("azuredevops_serviceendpoint_project_permissions.p", "project_reference.0.project_id", "project-id-1"),
					resource.TestCheckResourceAttr("azuredevops_serviceendpoint_project_permissions.p", "project_reference.0.service_endpoint_name", "service-connection-shared"),
					resource.TestCheckResourceAttr("azuredevops_serviceendpoint_project_permissions.p", "project_reference.0.description", "Shared service connection"),
					resource.TestCheckResourceAttr("azuredevops_serviceendpoint_project_permissions.p", "project_reference.1.project_id", "project-id-2"),
					resource.TestCheckResourceAttr("azuredevops_serviceendpoint_project_permissions.p", "project_reference.1.service_endpoint_name", "service-connection-shared"),
					resource.TestCheckResourceAttr("azuredevops_serviceendpoint_project_permissions.p", "project_reference.1.description", "Shared service connection"),
				),
			},
		},
	})
}

func hclServiceEndpointProjectPermissionsResource(projectName string, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%s"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  credentials {
    serviceprincipalid  = "spn-id"
    serviceprincipalkey = "spn-key"
  }
  azurerm_spn_tenantid      = "tenant-id"
  azurerm_subscription_id   = "subscription-id"
  azurerm_subscription_name = "subscription-name"
}

resource "azuredevops_serviceendpoint_project_permissions" "p" {
  serviceendpoint_id = azuredevops_serviceendpoint_azurerm.example.id
  project_reference {
    project_id            = "project-id-1"
    service_endpoint_name = "service-connection-shared"
    description           = "Shared service connection"
  }
  project_reference {
    project_id            = "project-id-2"
    service_endpoint_name = "service-connection-shared"
    description           = "Shared service connection"
  }
}
`, projectName, serviceEndpointName)
}
