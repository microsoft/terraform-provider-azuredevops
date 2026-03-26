package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointExternalTFS_PersonalTokenBasic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_externaltfs"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointExternalTFSResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_url", "https://dev.azure.com/myorganization"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointExternalTFS_PersonalTokenUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	description := "Managed by Terraform"
	resourceType := "azuredevops_serviceendpoint_externaltfs"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointExternalTFSResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_url", "https://dev.azure.com/myorganization"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointExternalTFSResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_url", "https://dev.azure.com/myorganization"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}

func TestAccServiceEndpointExternalTFS_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_externaltfs"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointExternalTFSResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_url", "https://dev.azure.com/myorganization"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointExternalTFSResourceBasic(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "connection_url", "https://dev.azure.com/myorganization"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_personal"},
			},
		},
	})
}

func TestAccServiceEndpointExternalTFS_RequiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_externaltfs"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointExternalTFSResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointExternalTFSResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointExternalTFSResourceBasic(projectName string, serviceEndpointName string) string {
	projectResource := testutils.HclProjectResource(projectName)
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_externaltfs" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%[1]s"
  connection_url        = "https://dev.azure.com/myorganization"
  auth_personal {
    personal_access_token = "test_token_basic"
  }
}`, serviceEndpointName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointExternalTFSResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	projectResource := testutils.HclProjectResource(projectName)
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_externaltfs" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%[1]s"
  connection_url        = "https://dev.azure.com/myorganization"
  auth_personal {
    personal_access_token = "test_token_update"
  }
  description = "%[2]s"
}`, serviceEndpointName, description)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointExternalTFSResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointExternalTFSResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_externaltfs" "import" {
  project_id            = azuredevops_serviceendpoint_externaltfs.serviceendpoint.project_id
  service_endpoint_name = azuredevops_serviceendpoint_externaltfs.serviceendpoint.service_endpoint_name
  connection_url        = azuredevops_serviceendpoint_externaltfs.serviceendpoint.connection_url
  auth_personal {
    personal_access_token = "test_token_basic"
  }
}
`, template)
}
