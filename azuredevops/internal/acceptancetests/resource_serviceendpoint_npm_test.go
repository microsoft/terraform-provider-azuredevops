package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointNpm_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_npm"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointNpmResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "url"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointNpm_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_npm"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointNpmResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://url.com/"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointNpm_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_npm"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointNpmResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointNpmResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://url.com/2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointNpm_RequiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_npm"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointNpmResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointNpmResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointNpmResourceBasic(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_npm" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  access_token          = "redacted"
  url                   = "http://url.com/"
}`, serviceEndpointName)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointNpmResourceComplete(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_npm" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  description           = "%s"
  access_token          = "redacted"
  url                   = "https://url.com/"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointNpmResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_npm" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  description           = "%s"
  access_token          = "redacted2"
  url                   = "https://url.com/2"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointNpmResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointNpmResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_npm" "import" {
  project_id            = azuredevops_serviceendpoint_npm.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_npm.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_npm.test.description
  url                   = azuredevops_serviceendpoint_npm.test.url
  access_token          = "redacted"
}
`, template)
}
