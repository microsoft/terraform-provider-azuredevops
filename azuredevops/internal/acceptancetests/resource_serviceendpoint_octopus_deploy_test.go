package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointOctopusDeploy_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_octopusdeploy"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointOctopusDeployResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointOctopusDeploy_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_octopusdeploy"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointOctopusDeployResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "url"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "api_key"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointOctopusDeploy_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_octopusdeploy"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointOctopusDeployResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointOctopusDeployResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "url"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "api_key"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointOctopusDeploy_RequiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_octopusdeploy"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointOctopusDeployResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointOctopusDeployResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointOctopusDeployResourceBasic(projectName string, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "Managed by Terraform"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_octopusdeploy" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%[2]s"
  url                   = "https://dev.azure.com"
  api_key               = "000000000000000000000000000000000000"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointOctopusDeployResourceComplete(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "Managed by Terraform"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_octopusdeploy" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%[2]s"
  url                   = "https://dev.azure.com"
  api_key               = "000000000000000000000000000000000000"
  ignore_ssl_error      = true
  description           = "%[3]s"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointOctopusDeployResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "Managed by Terraform"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_octopusdeploy" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%[2]s"
  url                   = "https://dev.azure.com"
  api_key               = "000000000000000000000000000000000000"
  ignore_ssl_error      = true
  description           = "%[3]s"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointOctopusDeployResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointOctopusDeployResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_octopusdeploy" "import" {
  project_id            = azuredevops_serviceendpoint_octopusdeploy.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_octopusdeploy.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_octopusdeploy.test.description
  url                   = azuredevops_serviceendpoint_octopusdeploy.test.url
  api_key               = "000000000000000000000000000000000000"
  ignore_ssl_error      = azuredevops_serviceendpoint_octopusdeploy.test.ignore_ssl_error
}
`, template)
}
