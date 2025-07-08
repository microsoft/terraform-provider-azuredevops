//go:build (all || resource_serviceendpoint_checkmarx_sca) && !exclude_resource_serviceendpoint_checkmarx_sca

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointCheckMarxSCA_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sca.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sca"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSCAResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_control_url", "https://accesscontrol.com"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxSCA_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := "Managed by Terraform"

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sca.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sca"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSCAResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "team"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "team", "team"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_control_url", "https://accesscontrol.com"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxSCA_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName() + "update"

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sca.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sca"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSCAResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointCheckMarxSCAResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_control_url", "https://accesscontrol.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "web_app_url", "https://webapp.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "team", "teamupdate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxSCA_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sca.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sca"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSCAResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointCheckMarxSCAResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointCheckMarxSCAResourceBasic(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_sca" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  access_control_url    = "https://accesscontrol.com"
  server_url            = "https://server.com"
  web_app_url           = "https://webapp.com"
  account               = "account"
  username              = "username"
  password              = "password"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointCheckMarxSCAResourceComplete(projectName, serviceEndpointName, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_sca" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  access_control_url    = "https://accesscontrol.com"
  server_url            = "https://server.com"
  web_app_url           = "https://webapp.com"
  account               = "account"
  username              = "username"
  password              = "password"
  team                  = "team"
  description           = "%s"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointCheckMarxSCAResourceUpdate(projectName, serviceEndpointName, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_sca" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  access_control_url    = "https://accesscontrol.com/update"
  server_url            = "https://server.com/update"
  web_app_url           = "https://webapp.com/update"
  account               = "accountupdate"
  username              = "usernameupdate"
  password              = "passwordupdate"
  team                  = "teamupdate"
  description           = "%s"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointCheckMarxSCAResourceRequiresImport(projectName, serviceEndpointName string) string {
	template := hclSvcEndpointCheckMarxSCAResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_checkmarx_sca" "import" {
  project_id            = azuredevops_serviceendpoint_checkmarx_sca.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_checkmarx_sca.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_checkmarx_sca.test.description
  access_control_url    = azuredevops_serviceendpoint_checkmarx_sca.test.access_control_url
  server_url            = azuredevops_serviceendpoint_checkmarx_sca.test.server_url
  web_app_url           = azuredevops_serviceendpoint_checkmarx_sca.test.web_app_url
  account               = azuredevops_serviceendpoint_checkmarx_sca.test.account
  username              = azuredevops_serviceendpoint_checkmarx_sca.test.username
  password              = azuredevops_serviceendpoint_checkmarx_sca.test.password
}
`, template)
}
