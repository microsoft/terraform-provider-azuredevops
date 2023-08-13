package testutils

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

// CheckServiceEndpointExistsWithName verifies that a service endpoint of a particular type exists in the state,
// and that it has the expected name when compared against the data in Azure DevOps.
func CheckPipelineCheckExistsWithName(tfNode string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[tfNode]
		if !ok {
			return fmt.Errorf("Did not find a check in the state")
		}

		check, err := getCheckFromState(resourceState)
		if err != nil {
			return err
		}

		if DisplayName, found := check.Settings.(map[string]interface{})["displayName"]; found {
			if DisplayName != expectedName {
				return fmt.Errorf("Check has Name=%s, but expected Name=%s", DisplayName, expectedName)
			}
		} else {
			return fmt.Errorf("displayName setting not found")
		}

		return nil
	}
}

// CheckPipelineCheckDestroyed verifies that all checks of the given type in the state are destroyed.
// This will be invoked *after* terraform destroys the resource but *before* the state is wiped clean.
func CheckPipelineCheckDestroyed(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, resource := range s.RootModule().Resources {
			if resource.Type != resourceType {
				continue
			}

			// indicates the resource exists - this should fail the test
			if _, err := getSvcEndpointFromState(resource); err == nil {
				return fmt.Errorf("Unexpectedly found a check that should have been deleted")
			}
		}

		return nil
	}
}

// given a resource from the state, return a check (and error)
func getCheckFromState(resource *terraform.ResourceState) (*pipelineschecksextras.CheckConfiguration, error) {
	branchControlCheckID, err := strconv.Atoi(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := GetProvider().Meta().(*client.AggregatedClient)
	return clients.PipelinesChecksClientExtras.GetCheckConfiguration(clients.Ctx, pipelineschecksextras.GetCheckConfigurationArgs{
		Project: &projectID,
		Id:      &branchControlCheckID,
		Expand:  converter.ToPtr(pipelineschecksextras.CheckConfigurationExpandParameterValues.Settings),
	})
}
