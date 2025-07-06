//go:build (all || resource_serviceendpoint_checkmarx_one) && !exclude_resource_serviceendpoint_checkmarx_one

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointCheckMarxOne_apiKey(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_one.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_one"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxOneServiceResourceApiKey(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxOne_apiKeyUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_one.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_one"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxOneServiceResourceApiKeyUpdate(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "api_key"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com/update"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxOne_clientIdSecret(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_one.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_one"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxOneServiceResourceClientIdSecret(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxOne_clientIdSecretUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_one.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_one"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxOneServiceResourceClientIdSecret(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com"),
				),
			},
			{
				Config: hclSvcEndpointCheckMarxOneServiceResourceClientIdSecretUpdate(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://server.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "client_id", "clientidupdate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_url", "https://authurl.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "descriptionupdate"),
				),
			},
		},
	})
}

func TestAccServiceEndpointCheckMarxOne_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	tfSvcEpNode := "azuredevops_serviceendpoint_checkmarx_one.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_checkmarx_one"),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointCheckMarxOneServiceResourceApiKey(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointCheckMarxOneServiceResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointCheckMarxOneServiceResourceApiKey(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_one" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://server.com"
  api_key               = "apikey"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointCheckMarxOneServiceResourceApiKeyUpdate(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_one" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://server.com/update"
  api_key               = "apikeyupdate"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointCheckMarxOneServiceResourceClientIdSecret(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_one" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://server.com"
  client_id             = "clientid"
  client_secret         = "secret"
  authorization_url     = "https://authurl.com"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointCheckMarxOneServiceResourceClientIdSecretUpdate(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_checkmarx_one" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://server.com/update"
  client_id             = "clientidupdate"
  client_secret         = "secretupdate"
  authorization_url     = "https://authurl.com/update"
  description           = "descriptionupdate"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointCheckMarxOneServiceResourceRequiresImport(projectName, serviceEndpointName string) string {
	template := hclSvcEndpointCheckMarxOneServiceResourceApiKey(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_checkmarx_one" "import" {
  project_id            = azuredevops_serviceendpoint_checkmarx_one.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_checkmarx_one.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_checkmarx_one.test.description
  server_url            = azuredevops_serviceendpoint_checkmarx_one.test.server_url
  api_key               = azuredevops_serviceendpoint_checkmarx_one.test.api_key
}
`, template)
}
