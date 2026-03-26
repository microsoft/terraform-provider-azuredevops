package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointCheckMarxSAST_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sast.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sast"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSASTResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxSAST_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := "Managed by Terraform"

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sast.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sast"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSASTResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "team"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "team", "team"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "username"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxSAST_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName() + "update"

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sast.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sast"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSASTResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointCheckMarxSASTResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "team", "teamupdate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "preset", "presetupdate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxSAST_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_sast.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_sast"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxSASTResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointCheckMarxSASTResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointCheckMarxSASTResourceBasic(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_sast" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://server.com"
  username              = "username"
  password              = "password"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointCheckMarxSASTResourceComplete(projectName, serviceEndpointName, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_sast" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://server.com"
  username              = "username"
  password              = "password"
  preset                = "preset"
  team                  = "team"
  description           = "%s"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointCheckMarxSASTResourceUpdate(projectName, serviceEndpointName, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_sast" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://server.com/update"
  username              = "usernameupdate"
  password              = "passwordupdate"
  team                  = "teamupdate"
  preset                = "presetupdate"
  description           = "%s"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointCheckMarxSASTResourceRequiresImport(projectName, serviceEndpointName string) string {
	template := hclSvcEndpointCheckMarxSASTResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_checkmarx_sast" "import" {
  project_id            = azuredevops_serviceendpoint_checkmarx_sast.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_checkmarx_sast.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_checkmarx_sast.test.description
  server_url            = azuredevops_serviceendpoint_checkmarx_sast.test.server_url
  username              = azuredevops_serviceendpoint_checkmarx_sast.test.username
  password              = azuredevops_serviceendpoint_checkmarx_sast.test.password
}
`, template)
}
