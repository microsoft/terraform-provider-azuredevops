package acceptancetests

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// Tests basic functionality of the Generic Service Endpoint V2 resource
func TestAccServiceEndpointGenericV2_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "github"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2TokenBasic(projectName, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://github.com"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "id"),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportState:             true,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization_parameters"},
			},
		},
	})
}

// Tests if the Generic Service Endpoint V2 can be updated with a different server_url
func TestAccServiceEndpointGenericV2_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "github"
	serverUrlInitial := "https://github.com"
	serverUrlUpdated := "https://api.github.com"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2TokenCustomUrl(projectName, serviceEndpointName, serviceEndpointType, serverUrlInitial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", serverUrlInitial),
				),
			},
			{
				Config: hclServiceEndpointGenericV2TokenCustomUrl(projectName, serviceEndpointName, serviceEndpointType, serverUrlUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", serverUrlUpdated),
				),
			},
		},
	})
}

// Tests username/password authentication for Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_UsernamePassword(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "generic"
	username := "testuser"
	password := "testpass"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2UsernamePassword(projectName, serviceEndpointName, serviceEndpointType, username, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportState:             true,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization_parameters"},
			},
		},
	})
}

// Tests shared_project_ids functionality for Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_SharedProjects(t *testing.T) {
	projectName1 := testutils.GenerateResourceName()
	projectName2 := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "generic"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2WithSharedProjects(projectName1, projectName2, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
					resource.TestCheckResourceAttr(tfSvcEpNode, "shared_project_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(tfSvcEpNode, "shared_project_ids.*", "azuredevops_project.project2", "id"),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportState:             true,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization_parameters"},
			},
		},
	})
}

// Tests updating shared_project_ids for Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_SharedProjectsUpdate(t *testing.T) {
	projectName1 := testutils.GenerateResourceName()
	projectName2 := testutils.GenerateResourceName()
	projectName3 := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "generic"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2WithSharedProjects(projectName1, projectName2, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "shared_project_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(tfSvcEpNode, "shared_project_ids.*", "azuredevops_project.project2", "id"),
				),
			},
			{
				Config: hclServiceEndpointGenericV2WithMultipleSharedProjects(projectName1, projectName2, projectName3, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "shared_project_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(tfSvcEpNode, "shared_project_ids.*", "azuredevops_project.project2", "id"),
					resource.TestCheckTypeSetElemAttrPair(tfSvcEpNode, "shared_project_ids.*", "azuredevops_project.project3", "id"),
				),
			},
			{
				Config: hclServiceEndpointGenericV2WithSharedProjects(projectName1, projectName3, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "shared_project_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(tfSvcEpNode, "shared_project_ids.*", "azuredevops_project.project2", "id"),
				),
			},
		},
	})
}

// Tests removing all shared_project_ids for Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_SharedProjectsRemoveAll(t *testing.T) {
	projectName1 := testutils.GenerateResourceName()
	projectName2 := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "generic"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2WithSharedProjects(projectName1, projectName2, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "shared_project_ids.#", "1"),
				),
			},
			{
				Config: hclServiceEndpointGenericV2WithoutSharedProjects(projectName1, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "shared_project_ids.#", "0"),
				),
			},
		},
	})
}

// Tests resource validation for type in Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_ValidateServiceEndpointType(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "invalidtype" // This should fail validation

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config:      hclServiceEndpointGenericV2TokenBasic(projectName, serviceEndpointName, serviceEndpointType),
				ExpectError: validateServiceEndpointTypeError(serviceEndpointType),
			},
		},
	})
}

// Helper function to validate service endpoint type error
func validateServiceEndpointTypeError(serviceEndpointType string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("service endpoint type '%s' is not available", serviceEndpointType))
}

// Helper function to generate HCL for a generic service endpoint with token auth
func hclServiceEndpointGenericV2TokenBasic(projectName string, serviceEndpointName string, serviceEndpointType string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  description = "Managed by Terraform"
  type        = "%s"
  server_url  = "https://github.com"

  authorization_scheme = "Token"
  authorization_parameters = {
    AccessToken = "test-token"
  }
}`, projectResource, serviceEndpointName, serviceEndpointType)
}

// Helper function to generate HCL for a generic service endpoint with token auth and custom URL
func hclServiceEndpointGenericV2TokenCustomUrl(projectName string, serviceEndpointName string, serviceEndpointType string, serverUrl string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  description = "Managed by Terraform"
  type        = "%s"
  server_url  = "%s"

  authorization_scheme = "Token"
  authorization_parameters = {
    AccessToken = "test-token"
  }
}`, projectResource, serviceEndpointName, serviceEndpointType, serverUrl)
}

// Helper function to generate HCL for a generic service endpoint with username/password auth
func hclServiceEndpointGenericV2UsernamePassword(projectName string, serviceEndpointName string, serviceEndpointType string, username string, password string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  description = "Managed by Terraform"
  type        = "%s"
  server_url  = "https://example.com"

  authorization_scheme = "UsernamePassword"
  authorization_parameters = {
    username = "%s"
    password = "%s"
  }
}`, projectResource, serviceEndpointName, serviceEndpointType, username, password)
}

// Helper function to generate HCL for a generic service endpoint with shared projects
func hclServiceEndpointGenericV2WithSharedProjects(projectName1 string, projectName2 string, serviceEndpointName string, serviceEndpointType string) string {
	projectResource1 := testutils.HclProjectResource(projectName1)
	projectResource2 := testutils.HclProjectResource(projectName2)
	// Replace "project" with "project2" in the second project
	projectResource2 = strings.Replace(projectResource2, "azuredevops_project\" \"project", "azuredevops_project\" \"project2", 1)
	return fmt.Sprintf(`
%s

%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  description = "Managed by Terraform"
  type        = "%s"
  server_url  = "https://example.com"

  shared_project_ids = [
    azuredevops_project.project2.id
  ]

  authorization_scheme = "UsernamePassword"
  authorization_parameters = {
    username = "test-token"
    password = "test-password"
  }
}`, projectResource1, projectResource2, serviceEndpointName, serviceEndpointType)
}

// Helper function to generate HCL for a generic service endpoint without shared projects
func hclServiceEndpointGenericV2WithoutSharedProjects(projectName string, serviceEndpointName string, serviceEndpointType string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  description = "Managed by Terraform"
  type        = "%s"
  server_url  = "https://example.com"

  authorization_scheme = "UsernamePassword"
  authorization_parameters = {
    username = "test-token"
    password = "test-password"
  }
}`, projectResource, serviceEndpointName, serviceEndpointType)
}

// Helper function to generate HCL for a generic service endpoint with multiple shared projects
func hclServiceEndpointGenericV2WithMultipleSharedProjects(projectName1 string, projectName2 string, projectName3 string, serviceEndpointName string, serviceEndpointType string) string {
	projectResource1 := testutils.HclProjectResource(projectName1)
	projectResource2 := testutils.HclProjectResource(projectName2)
	projectResource3 := testutils.HclProjectResource(projectName3)
	// Replace "project" with "project2" and "project3" in the additional projects
	projectResource2 = strings.Replace(projectResource2, "azuredevops_project\" \"project", "azuredevops_project\" \"project2", 1)
	projectResource3 = strings.Replace(projectResource3, "azuredevops_project\" \"project", "azuredevops_project\" \"project3", 1)
	return fmt.Sprintf(`
%s

%s

%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  description = "Managed by Terraform"
  type        = "%s"
  server_url  = "https://example.com"

  shared_project_ids = [
    azuredevops_project.project2.id,
    azuredevops_project.project3.id
  ]

  authorization_scheme = "UsernamePassword"
  authorization_parameters = {
    username = "test-user"
    password = "test-password"
  }
}`, projectResource1, projectResource2, projectResource3, serviceEndpointName, serviceEndpointType)
}

// checkServiceEndpointDestroyed verifies that all service endpoints with the specified type have been destroyed
func checkServiceEndpointDestroyed(resourceType string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		// verify that every service endpoint referenced in the state does not exist in AzDO
		for _, res := range s.RootModule().Resources {
			if res.Type != resourceType {
				continue
			}

			endpointIDStr := res.Primary.ID
			endpointID, err := uuid.Parse(endpointIDStr)
			if err != nil {
				return fmt.Errorf("Service Endpoint ID %s is not a valid UUID: %v", endpointIDStr, err)
			}
			projectIDStr := res.Primary.Attributes["project_id"]

			// Ensure the service endpoint does not exist in the main project
			_, err = clients.ServiceEndpointClient.GetServiceEndpointDetails(
				clients.Ctx,
				serviceendpoint.GetServiceEndpointDetailsArgs{
					EndpointId: &endpointID,
					Project:    &projectIDStr,
				},
			)

			if err == nil {
				return fmt.Errorf("Service Endpoint ID %s still exists", endpointID)
			}
		}

		return nil
	}
}
