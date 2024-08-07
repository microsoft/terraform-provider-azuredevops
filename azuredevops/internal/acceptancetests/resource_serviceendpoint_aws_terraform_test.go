package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAwsTerraform_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws_terraform"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsTerraformResource(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_key_id", "0000"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "secret_access_key", "secretkey"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "region", "us-east-1"),
				),
			},
		},
	})
}

func TestAccServiceEndpointAwsTerraform_Complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()
	sessionToken := "foobar"
	rta := "rta"
	rsn := "rsn"
	externalId := "external_id"

	resourceType := "azuredevops_serviceendpoint_aws_terraform"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsTerraformResourceComplete(projectName, serviceEndpointName, description, sessionToken, rta, rsn, externalId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "access_key_id", "0000"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "secret_access_key", "secretkey"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
					resource.TestCheckResourceAttr(tfSvcEpNode, "region", "us-east-1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "external_id", externalId),
				),
			},
		},
	})
}

func TestAccServiceEndpointAwsTerraform_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws_terraform"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsTerraformResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointAwsTerraformResourceUpdate(projectName, serviceEndpointNameSecond, description),
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

func TestAccServiceEndpointAwsTerraform_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_aws_terraform"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointAwsTerraformResource(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointAwsTerraformResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointAwsTerraformResource(projectName string, serviceEndpointName string) string {
	return hclSvcEndpointAwsTerraformResourceUpdate(projectName, serviceEndpointName, "description")
}

func hclSvcEndpointAwsTerraformResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
	resource "azuredevops_serviceendpoint_aws_terraform" "test" {
		project_id             = azuredevops_project.project.id
		access_key_id          = "0000"
		secret_access_key      = "secretkey"
		region                 = "us-east-1"
		service_endpoint_name  = "%s"
		description            = "%s"
	}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointAwsTerraformResourceComplete(projectName string, serviceEndpointName string, description string, sessionToken string, rta string, rsn string, externalId string) string {
	serviceEndpointResource := fmt.Sprintf(`
	resource "azuredevops_serviceendpoint_aws_terraform" "test" {
		project_id             = azuredevops_project.project.id
		access_key_id          = "0000"
		secret_access_key      = "secretkey"
		region				   = "us-east-1"
		service_endpoint_name  = "%s"
		description            = "%s"
		external_id = "%s"
	}`, serviceEndpointName, description, externalId)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointAwsTerraformResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := hclSvcEndpointAwsTerraformResource(projectName, serviceEndpointName)
	return fmt.Sprintf(`
	%s
	resource "azuredevops_serviceendpoint_aws_terraform" "import" {
	project_id             = azuredevops_serviceendpoint_aws_terraform.test.project_id
	access_key_id          = "0000"
	secret_access_key      = "secretkey"
	region                 = "us-east-1"
	service_endpoint_name  = azuredevops_serviceendpoint_aws_terraform.test.service_endpoint_name
	description            = azuredevops_serviceendpoint_aws_terraform.test.description
	}
	`, template)
}
