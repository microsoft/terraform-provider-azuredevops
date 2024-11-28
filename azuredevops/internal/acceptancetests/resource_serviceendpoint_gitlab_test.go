package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointGitLab_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_gitlab"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGitLabResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://gitlab.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "username"),
				),
			},
		},
	})
}

func TestAccServiceEndpointGitLab_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_gitlab"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGitLabResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://gitlab.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "username"),
				),
			},
		},
	})
}

func TestAccServiceEndpointGitLab_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_gitlab"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGitLabResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://gitlab.com"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "username"),
				),
			},
			{
				Config: hclSvcEndpointGitLabResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://gitlab.com/update"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "usernameupdate"),
				),
			},
		},
	})
}

func TestAccServiceEndpointGitLab_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_gitlab"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGitLabResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointGitLabResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointGitLabResourceBasic(projectName string, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_gitlab" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  url                   = "https://gitlab.com"
  username              = "username"
  api_token             = "token"
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointGitLabResourceComplete(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_gitlab" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  description           = "%s"
  url                   = "https://gitlab.com"
  username              = "username"
  api_token             = "token"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointGitLabResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_gitlab" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  description           = "%s"
  url                   = "https://gitlab.com/update"
  username              = "usernameupdate"
  api_token             = "tokenupdate"
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointGitLabResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointGitLabResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_gitlab" "import" {
  project_id            = azuredevops_serviceendpoint_gitlab.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_gitlab.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_gitlab.test.description
  url                   = azuredevops_serviceendpoint_gitlab.test.url
  username              = azuredevops_serviceendpoint_gitlab.test.username
  api_token             = azuredevops_serviceendpoint_gitlab.test.api_token
}
`, template)
}
