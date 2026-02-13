package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAws_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResource(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_key_id", "0000"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "secret_access_key", "secretkey"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "use_oidc", "false"),
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
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResourceComplete(projectName, serviceEndpointName, description, sessionToken, rta, rsn, externalId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_key_id", "0000"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "secret_access_key", "secretkey"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
					resource.TestCheckResourceAttr(tfSvcEpNode, "session_token", "foobar"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "role_to_assume", rta),
					resource.TestCheckResourceAttr(tfSvcEpNode, "role_session_name", rsn),
					resource.TestCheckResourceAttr(tfSvcEpNode, "external_id", externalId),
					resource.TestCheckResourceAttr(tfSvcEpNode, "use_oidc", "false"),
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
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResource(projectName, serviceEndpointNameFirst),
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

func TestAccServiceEndpointAws_oidc(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResourceOidc(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_key_id", "0000"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "secret_access_key", "secretkey"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "use_oidc", "true"),
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
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsResource(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointAwsResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointAwsResource(projectName string, serviceEndpointName string) string {
	return hclSvcEndpointAwsResourceUpdate(projectName, serviceEndpointName, "description")
}

func hclSvcEndpointAwsResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}
resource "azuredevops_serviceendpoint_aws" "test" {
  project_id            = azuredevops_project.project.id
  access_key_id         = "0000"
  secret_access_key     = "secretkey"
  service_endpoint_name = "%s"
  description           = "%s"
  use_oidc              = false
}`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointAwsResourceComplete(projectName string, serviceEndpointName string, description string, sessionToken string, rta string, rsn string, externalId string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}
resource "azuredevops_serviceendpoint_aws" "test" {
  project_id            = azuredevops_project.project.id
  access_key_id         = "0000"
  secret_access_key     = "secretkey"
  service_endpoint_name = "%s"
  description           = "%s"

  session_token     = "%s"
  role_to_assume    = "%s"
  role_session_name = "%s"
  external_id       = "%s"
  use_oidc          = false
}`, projectName, serviceEndpointName, description, sessionToken, rta, rsn, externalId)
}

func hclSvcEndpointAwsResourceOidc(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name = "%s"
}
resource "azuredevops_serviceendpoint_aws" "test" {
  project_id            = azuredevops_project.project.id
  access_key_id         = "0000"
  secret_access_key     = "secretkey"
  service_endpoint_name = "%s"
  use_oidc              = true
}`, projectName, serviceEndpointName)
}

func hclSvcEndpointAwsResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointAwsResource(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_aws" "import" {
  project_id            = azuredevops_serviceendpoint_aws.test.project_id
  access_key_id         = "0000"
  secret_access_key     = "secretkey"
  service_endpoint_name = azuredevops_serviceendpoint_aws.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_aws.test.description
  use_oidc              = azuredevops_serviceendpoint_aws.test.use_oidc
}
`, template)
}
