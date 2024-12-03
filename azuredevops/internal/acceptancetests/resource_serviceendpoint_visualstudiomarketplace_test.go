package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointMarketplace_basicToken(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_visualstudiomarketplace"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointMarketplaceResourceBasicToken(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://marketplace.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_token.#", "1"),
				),
			},
		},
	})
}

func TestAccServiceEndpointMarketplace_basicUsernamePassword(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_visualstudiomarketplace"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointMarketplaceResourceBasicUsernamePasword(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://marketplace.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_basic.#", "1"),
				),
			},
		},
	})
}

func TestAccServiceEndpointMarketplace_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_visualstudiomarketplace"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointMarketplaceResourceBasicToken(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://marketplace.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_token.#", "1"),
				),
			},
			{
				Config: hclSvcEndpointMarketplaceResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://marketplace.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_token.#", "1"),
				),
			},
		},
	})
}

func TestAccServiceEndpointMarketplace_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_visualstudiomarketplace"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointMarketplaceResourceBasicToken(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointMarketplaceResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointMarketplaceResourceBasicToken(projectName string, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_visualstudiomarketplace" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  url                   = "https://marketplace.com"
  authentication_token {
    token = "token"
  }
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointMarketplaceResourceBasicUsernamePasword(projectName string, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_visualstudiomarketplace" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  url                   = "https://marketplace.com"
  authentication_basic {
    username = "uname"
    password = "pwd"
  }
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointMarketplaceResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_visualstudiomarketplace" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  description           = "%s"
  url                   = "https://marketplace.com/update"
  authentication_token {
    token = "tokenupdate"
  }
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointMarketplaceResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointMarketplaceResourceBasicToken(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_visualstudiomarketplace" "import" {
  project_id            = azuredevops_serviceendpoint_visualstudiomarketplace.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_visualstudiomarketplace.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_visualstudiomarketplace.test.description
  url                   = azuredevops_serviceendpoint_visualstudiomarketplace.test.url
  authentication_token {
    token = azuredevops_serviceendpoint_visualstudiomarketplace.test.authentication_token.0.token
  }
}
`, template)
}
