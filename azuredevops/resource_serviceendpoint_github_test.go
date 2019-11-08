package azuredevops

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

var ghTestServiceEndpointID = uuid.New()
var ghRandomServiceEndpointProjectID = uuid.New().String()
var ghTestServiceEndpointProjectID = &ghRandomServiceEndpointProjectID

var ghTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"accessToken": "UNIT_TEST_ACCESS_TOKEN",
		},
		Scheme: converter.String("PersonalAccessToken"),
	},
	Id:    &ghTestServiceEndpointID,
	Name:  converter.String("UNIT_TEST_NAME"),
	Owner: converter.String("library"),
	Type:  converter.String("github"),
	Url:   converter.String("http://github.com"),
}

/**
 * Begin unit tests
 */

// verifies that the flatten/expand round trip yields the same service endpoint
func TestAzureDevOpsServiceEndpointGitHub_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpointGitHub().Schema, nil)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint, ghTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID := expandServiceEndpointGitHub(resourceData)

	require.Equal(t, ghTestServiceEndpoint, *serviceEndpointAfterRoundTrip)
	require.Equal(t, ghTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAzureDevOpsServiceEndpointGitHub_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint, ghTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &ghTestServiceEndpoint, Project: ghTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointGitHub_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint, ghTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: ghTestServiceEndpoint.Id, Project: ghTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpointGitHub_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint, ghTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: ghTestServiceEndpoint.Id, Project: ghTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpointGitHub_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointGitHub()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointGitHub(resourceData, &ghTestServiceEndpoint, ghTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &ghTestServiceEndpoint,
		EndpointId: ghTestServiceEndpoint.Id,
		Project:    ghTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

/**
 * Begin acceptance tests
 */

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsServiceEndpointGitHub_CreateAndUpdate(t *testing.T) {
	projectName := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameSecond := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_github.serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointGitHubCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceEndpointGitHubResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "github_service_endpoint_pat", ""),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "github_service_endpoint_pat_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testAccCheckServiceEndpointGitHubResourceExists(serviceEndpointNameFirst),
				),
			}, {
				Config: testAccServiceEndpointGitHubResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "github_service_endpoint_pat", ""),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "github_service_endpoint_pat_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testAccCheckServiceEndpointGitHubResourceExists(serviceEndpointNameSecond),
				),
			},
		},
	})
}

// HCL describing an AzDO service endpoint
func testAccServiceEndpointGitHubResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_github" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	
}`, serviceEndpointName)

	projectResource := testAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// Given the name of an AzDO service endpoint, this will return a function that will check whether
// or not the resource (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func testAccCheckServiceEndpointGitHubResourceExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceEndpointDef, ok := s.RootModule().Resources["azuredevops_serviceendpoint_github.serviceendpoint"]
		if !ok {
			return fmt.Errorf("Did not find a service endpoint in the TF state")
		}

		serviceEndpoint, err := getServiceEndpointGitHubFromResource(serviceEndpointDef)
		if err != nil {
			return err
		}

		if *serviceEndpoint.Name != expectedName {
			return fmt.Errorf("Service Endpoint has Name=%s, but expected Name=%s", *serviceEndpoint.Name, expectedName)
		}

		return nil
	}
}

// verifies that all service endpoints referenced in the state are destroyed. This will be invoked
// *after* terrafform destroys the resource but *before* the state is wiped clean.
func testAccServiceEndpointGitHubCheckDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_serviceendpoint_github" {
			continue
		}

		// indicates the service endpoint still exists - this should fail the test
		if _, err := getServiceEndpointGitHubFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a service endpoint that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a service endpoint (and error)
func getServiceEndpointGitHubFromResource(resource *terraform.ResourceState) (*serviceendpoint.ServiceEndpoint, error) {
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
