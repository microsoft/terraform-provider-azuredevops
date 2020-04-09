package azuredevops

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
)

func TestAccAzureDevOpsServiceEndpointBitBucket_basic(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_bitbucket.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointBitBucketCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceEndpointBitBucketResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEndpointBitBucketResourceExists(serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
		},
	})
}

func TestAccAzureDevOpsServiceEndpointBitBucket_complete(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	description := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_bitbucket.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointBitBucketCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceEndpointBitBucketResourceComplete(projectName, serviceEndpointNameFirst, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEndpointBitBucketResourceExists(serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "username"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "password", ""),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccAzureDevOpsServiceEndpointBitBucket_update(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	description := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_bitbucket.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointBitBucketCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceEndpointBitBucketResourceBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEndpointBitBucketResourceExists(serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: testAccServiceEndpointBitBucketResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEndpointBitBucketResourceExists(serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "username"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "password", ""),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccAzureDevOpsServiceEndpointBitBucket_requiresImportErrorStep(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_bitbucket.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointBitBucketCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceEndpointBitBucketResourceBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEndpointBitBucketResourceExists(serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
			{
				Config:      testAccServiceEndpointBitBucketResourceRequiresImport(projectName, serviceEndpointName),
				ExpectError: requiresImportError(serviceEndpointName),
			},
		},
	})
}

func requiresImportError(resourceName string) *regexp.Regexp {
	message := "Error creating service endpoint in Azure DevOps: Service connection with name %[1]s already exists. Only a user having Administrator/User role permissions on service connection %[1]s can see it."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}

func testAccCheckServiceEndpointBitBucketResourceExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceEndpointDef, ok := s.RootModule().Resources["azuredevops_serviceendpoint_bitbucket.test"]
		if !ok {
			return fmt.Errorf("Did not find a service endpoint in the TF state")
		}

		serviceEndpoint, err := getServiceEndpointBitBucketFromResource(serviceEndpointDef)
		if err != nil {
			return err
		}

		if *serviceEndpoint.Name != expectedName {
			return fmt.Errorf("Service Endpoint has Name=%s, but expected Name=%s", *serviceEndpoint.Name, expectedName)
		}

		return nil
	}
}

func testAccServiceEndpointBitBucketCheckDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_serviceendpoint_bitbucket" {
			continue
		}

		if _, err := getServiceEndpointBitBucketFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a service endpoint that should be deleted")
		}
	}

	return nil
}

func getServiceEndpointBitBucketFromResource(resource *terraform.ResourceState) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpointDefID, err := uuid.Parse(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testAccProvider.Meta().(*config.AggregatedClient)
	return clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{
		Project:    &projectID,
		EndpointId: &serviceEndpointDefID,
	})
}

func preCheck(t *testing.T) {
	testhelper.TestAccPreCheck(t, &[]string{
		"AZDO_BITBUCKET_SERVICE_CONNECTION_USERNAME",
		"AZDO_BITBUCKET_SERVICE_CONNECTION_PASSWORD",
	})
}

func testAccServiceEndpointBitBucketResourceBasic(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_bitbucket" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
}`, serviceEndpointName)

	projectResource := testhelper.TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func testAccServiceEndpointBitBucketResourceComplete(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_bitbucket" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"
}`, serviceEndpointName, description)

	projectResource := testhelper.TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func testAccServiceEndpointBitBucketResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_bitbucket" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"
}`, serviceEndpointName, description)

	projectResource := testhelper.TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func testAccServiceEndpointBitBucketResourceRequiresImport(projectName string, serviceEndpointName string) string {
	template := testAccServiceEndpointBitBucketResourceBasic(projectName, serviceEndpointName)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_bitbucket" "import" {
  project_id                = azuredevops_serviceendpoint_bitbucket.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_bitbucket.test.service_endpoint_name
  description            = azuredevops_serviceendpoint_bitbucket.test.description
  authorization          = azuredevops_serviceendpoint_bitbucket.test.authorization
}
`, template)
}

func init() {
	InitProvider()
}
