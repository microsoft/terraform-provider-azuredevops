package testutils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// CheckServiceEndpointExistsWithName verifies that a service endpoint of a particular type exists in the state,
// and that it has the expected name when compared against the data in Azure DevOps.
func CheckServiceEndpointExistsWithName(tfNode string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[tfNode]
		if !ok {
			return fmt.Errorf("Did not find a service endpoint in the state")
		}

		serviceEndpoint, err := getSvcEndpointFromState(resourceState)
		if err != nil {
			return err
		}

		if *serviceEndpoint.Name != expectedName {
			return fmt.Errorf("Service Endpoint has Name=%s, but expected Name=%s", *serviceEndpoint.Name, expectedName)
		}

		return nil
	}
}

// CheckServiceEndpointDestroyed verifies that all service endpoints of the given type in the state are destroyed.
// This will be invoked *after* terraform destroys the resource but *before* the state is wiped clean.
func CheckServiceEndpointDestroyed(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, resource := range s.RootModule().Resources {
			if resource.Type != resourceType {
				continue
			}

			// indicates the resource exists - this should fail the test
			if _, err := getSvcEndpointFromState(resource); err == nil {
				return fmt.Errorf("Unexpectedly found a service endpoint that should have been deleted")
			}
		}

		return nil
	}
}

// given a resource from the state, return a service endpoint (and error)
func getSvcEndpointFromState(resource *terraform.ResourceState) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpointDefID, err := uuid.Parse(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := GetProvider().Meta().(*client.AggregatedClient)
	return clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{
		Project:    &projectID,
		EndpointId: &serviceEndpointDefID,
	})
}
