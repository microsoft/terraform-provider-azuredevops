// +build all permissions resource_serviceendpoint_permissions
// +build !exclude_permissions !exclude_resource_serviceendpoint_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

func hclServiceEndpointPermissions(projectName string, serviceEndpointName string, permissions map[string]map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions["root"], "=", "\n")
	serviceEndpointPermissions := datahelper.JoinMap(permissions["service_endpoint"], "=", "\n")

	return fmt.Sprintf(`
%s

%s

data "azuredevops_group" "tf-project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_serviceendpoint_permissions" "root-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	permissions = {
		%s
	}
}

resource "azuredevops_serviceendpoint_permissions" "serviceendpoint-permissions" {
	project_id  			 = azuredevops_project.project.id
	principal   			 = data.azuredevops_group.tf-project-readers.id
	serviceendpoint_id = azuredevops_serviceendpoint_dockerregistry.serviceendpoint.id
	permissions = {
		%s
	}
}

`,
		testutils.HclProjectResource(projectName),
		testutils.HclServiceEndpointDockerRegistryResource(projectName, serviceEndpointName),
		rootPermissions,
		serviceEndpointPermissions)
}

func TestAccServiceEndpointPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	config := hclServiceEndpointPermissions(projectName, serviceEndpointName, map[string]map[string]string{
		"root": {
			"USE":               "allow",
			"ADMINISTER":        "allow",
			"CREATE":            "allow",
			"VIEWAUTHORIZATION": "allow",
			"VIEWENDPOINT":      "allow",
		},
		"service_endpoint": {
			"USE":               "allow",
			"ADMINISTER":        "deny",
			"CREATE":            "deny",
			"VIEWAUTHORIZATION": "allow",
			"VIEWENDPOINT":      "allow",
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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.USE", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ADMINISTER", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.VIEWAUTHORIZATION", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.VIEWENDPOINT", "allow"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.USE", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ADMINISTER", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.CREATE", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.VIEWAUTHORIZATION", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.VIEWENDPOINT", "allow"),
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
			"USE":               "allow",
			"ADMINISTER":        "allow",
			"CREATE":            "allow",
			"VIEWAUTHORIZATION": "allow",
			"VIEWENDPOINT":      "allow",
		},
		"service_endpoint": {
			"USE":               "allow",
			"ADMINISTER":        "deny",
			"CREATE":            "deny",
			"VIEWAUTHORIZATION": "allow",
			"VIEWENDPOINT":      "allow",
		},
	})
	config2 := hclServiceEndpointPermissions(projectName, serviceEndpointName, map[string]map[string]string{
		"root": {
			"USE":               "allow",
			"ADMINISTER":        "allow",
			"CREATE":            "allow",
			"VIEWAUTHORIZATION": "notset",
			"VIEWENDPOINT":      "notset",
		},
		"service_endpoint": {
			"USE":               "allow",
			"ADMINISTER":        "allow",
			"CREATE":            "allow",
			"VIEWAUTHORIZATION": "allow",
			"VIEWENDPOINT":      "allow",
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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.USE", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ADMINISTER", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.VIEWAUTHORIZATION", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.VIEWENDPOINT", "allow"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.USE", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ADMINISTER", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.CREATE", "deny"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.VIEWAUTHORIZATION", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.VIEWENDPOINT", "allow"),
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
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.USE", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ADMINISTER", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.CREATE", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.VIEWAUTHORIZATION", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.VIEWENDPOINT", "notset"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeServiceEndpoint, "serviceendpoint_id"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.%", "5"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.USE", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.ADMINISTER", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.CREATE", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.VIEWAUTHORIZATION", "allow"),
					resource.TestCheckResourceAttr(tfNodeServiceEndpoint, "permissions.VIEWENDPOINT", "allow"),
				),
			},
		},
	})
}
