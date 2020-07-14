// +build all resource_serviceendpoint_aws
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

var requiredAwsTestEnvVars = &[]string{
	"AZDO_AWS_SERVICE_CONNECTION_ACCESS_KEY_ID",
	"AZDO_AWS_SERVICE_CONNECTION_SECRET_ACCESS_KEY",
}

func TestAccServiceEndpointAws_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, requiredAwsTestEnvVars) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "access_key_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "secret_access_key", ""),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "secret_access_key_hash"),
				),
			},
		},
	})
}

func TestAccServiceEndpointAws_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()
	sessionToken := "foobar"
	rta := "rta"
	rsn := "rsn"
	externalId := "external_id"

	resourceType := "azuredevops_serviceendpoint_aws"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, requiredAwsTestEnvVars) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResourceComplete(projectName, serviceEndpointName, description, sessionToken, rta, rsn, externalId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "access_key_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "secret_access_key", ""),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "secret_access_key_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
					resource.TestCheckResourceAttr(tfSvcEpNode, "session_token", ""),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "session_token_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "role_to_assume", rta),
					resource.TestCheckResourceAttr(tfSvcEpNode, "role_session_name", rsn),
					resource.TestCheckResourceAttr(tfSvcEpNode, "external_id", externalId),
				),
			},
		},
	})
}

func TestAccServiceEndpointAws_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, requiredAwsTestEnvVars) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointAwsResourceUpdate(projectName, serviceEndpointNameSecond, description),
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

func TestAccServiceEndpointAws_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, requiredAwsTestEnvVars) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointAwsResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: requiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointAwsResourceBasic(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_aws" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
}`, serviceEndpointName)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointAwsResourceComplete(projectName string, serviceEndpointName string, description string, sessionToken string, rta string, rsn string, externalId string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_aws" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"

	session_token = "%s"
	role_to_assume = "%s"
	role_session_name = "%s"
	external_id = "%s"
}`, serviceEndpointName, description, sessionToken, rta, rsn, externalId)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointAwsResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_aws" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointAwsResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointAwsResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_aws" "import" {
  project_id                = azuredevops_serviceendpoint_aws.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_aws.test.service_endpoint_name
  description            = azuredevops_serviceendpoint_aws.test.description
  authorization          = azuredevops_serviceendpoint_aws.test.authorization
}
`, template)
}
