package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func TestAccServiceEndpointPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	config := hclServiceEndpointPermissions(projectName, serviceEndpointName, map[string]map[string]string{
		"root": {
			"Use":               "allow",
			"Administer":        "allow",
			"Create":            "allow",
			"ViewAuthorization": "allow",
			"ViewEndpoint":      "allow",
		},
		"service_endpoint": {
			"Use":               "allow",
			"Administer":        "deny",
			"Create":            "deny",
			"ViewAuthorization": "allow",
			"ViewEndpoint":      "allow",
		},
	})
	tfNodeRoot := "azuredevops_serviceendpoint_permissions.root-permissions"
	tfNodeServiceEndpoint := "azuredevops_serviceendpoint_permissions.serviceendpoint-permissions"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Administer", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Create", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewAuthorization", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewEndpoint", "allow"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Administer", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Create", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ViewAuthorization", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ViewEndpoint", "allow"),
				),
			},
		},
	})
}

func TestAccServiceEndpointPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	config1 := hclServiceEndpointPermissions(projectName, serviceEndpointName, map[string]map[string]string{
		"root": {
			"Use":               "allow",
			"Administer":        "allow",
			"Create":            "allow",
			"ViewAuthorization": "allow",
			"ViewEndpoint":      "allow",
		},
		"service_endpoint": {
			"Use":               "allow",
			"Administer":        "deny",
			"Create":            "deny",
			"ViewAuthorization": "allow",
			"ViewEndpoint":      "allow",
		},
	})
	config2 := hclServiceEndpointPermissions(projectName, serviceEndpointName, map[string]map[string]string{
		"root": {
			"Use":               "allow",
			"Administer":        "allow",
			"Create":            "allow",
			"ViewAuthorization": "notset",
			"ViewEndpoint":      "notset",
		},
		"service_endpoint": {
			"Use":               "allow",
			"Administer":        "allow",
			"Create":            "allow",
			"ViewAuthorization": "allow",
			"ViewEndpoint":      "allow",
		},
	})
	tfNodeRoot := "azuredevops_serviceendpoint_permissions.root-permissions"
	tfNodeServiceEndpoint := "azuredevops_serviceendpoint_permissions.serviceendpoint-permissions"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Administer", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Create", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewAuthorization", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewEndpoint", "allow"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Administer", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Create", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ViewAuthorization", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ViewEndpoint", "allow"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckNoResourceAttr(tfNodeRoot, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Administer", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.Create", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewAuthorization", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewEndpoint", "notset"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Use", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Administer", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.Create", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ViewAuthorization", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ViewEndpoint", "allow"),
				),
			},
		},
	})
}

func hclServiceEndpointPermissions(projectName string, serviceEndpointName string, permissions map[string]map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions["root"], "=", "\n")
	serviceEndpointPermissions := datahelper.JoinMap(permissions["service_endpoint"], "=", "\n")

	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_dockerregistry" "serviceendpoint" {
  docker_email          = "test@email.com"
  docker_username       = "testuser"
  docker_password       = "secret"
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
}

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_serviceendpoint_permissions" "root-permissions" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id
  permissions = {
		%s
  }
}

resource "azuredevops_serviceendpoint_permissions" "serviceendpoint-permissions" {
  project_id         = azuredevops_project.project.id
  principal          = data.azuredevops_group.tf-project-readers.id
  serviceendpoint_id = azuredevops_serviceendpoint_dockerregistry.serviceendpoint.id
  permissions = {
		%s
  }
}


`, testutils.HclProjectResource(projectName),
		serviceEndpointName,
		rootPermissions,
		serviceEndpointPermissions)
}
