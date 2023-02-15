//go:build (all || resource_serviceendpoint_gcp) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_gcp
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointGcp_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_gcp"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGcpResource(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "private_key", ""),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "private_key_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "token_uri", "0000"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "gcp_project_id", "project_id"),
				),
			},
		},
	})
}

func TestAccServiceEndpointGcp_Complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()
	scope := "scope"
	clientEmail := "client_email"
	tokenUri := "tokenUri"

	resourceType := "azuredevops_serviceendpoint_gcp"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGcpResourceComplete(projectName, serviceEndpointName, description, sessionToken, scope, clientEmail, tokenUri, projectId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "external_id", externalId),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "private_key", ""),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "private_key_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "token_uri", tokenUri),
					resource.TestCheckResourceAttr(tfSvcEpNode, "gcp_project_id", projectId),
					resource.TestCheckResourceAttr(tfSvcEpNode, "scope", scope),
					resource.TestCheckResourceAttr(tfSvcEpNode, "client_email", client_email),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointGcp_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_gcp"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGcpResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointGcpResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointGcp_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_gcp"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointGcpResource(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointGcpResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointGcpResource(projectName string, serviceEndpointName string) string {
	return hclSvcEndpointGcpResourceUpdate(projectName, serviceEndpointName, "description")
}

func hclSvcEndpointGcpResourceUpdate(projectName string, serviceEndpointName string, description string, tokenUri string) string {
	serviceEndpointResource := fmt.Sprintf(`
	resource "azuredevops_serviceendpoint_gcp" "test" {
		project_id             = azuredevops_project.project.id
	        private_key      = "secretkey"
		token_uri = "%s"
		service_endpoint_name  = "%s"
		description            = "%s"
		gcp_project_id = "project_id"
	}`, tokenUri, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointGcpResourceComplete(projectName string, serviceEndpointName string, description string, clientEmail string, scope string, tokenUri string, projectId string) string {
	serviceEndpointResource := fmt.Sprintf(`
	resource "azuredevops_serviceendpoint_gcp" "test" {
		project_id             = azuredevops_project.project.id
	        private_key      = "secretkey"
		token_uri = "%s"
		service_endpoint_name  = "%s"
		description            = "%s"
		client_email = "%s"
		scope = "%s"
		gcp_project_id = "%s"

	}`, tokenUri, serviceEndpointName, description, clientEmail, scope, projectId)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointGcpResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointGcpResource(projectName, serviceEndpointName)
	return fmt.Sprintf(`
	%s
	resource "azuredevops_serviceendpoint_gcp" "import" {
	project_id             = azuredevops_serviceendpoint_gcp.test.project_id
	private_key      = "secretkey"
	service_endpoint_name  = azuredevops_serviceendpoint_gcp.test.service_endpoint_name
	description            = azuredevops_serviceendpoint_gcp.test.description
	gcp_project_id            = azuredevops_serviceendpoint_gcp.test.gcp_project_id
	}
	`, template)
}
