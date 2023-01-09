//go:build (all || resource_environment) && !exclude_resource_environment
// +build all resource_environment
// +build !exclude_resource_environment

package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// Verifies that the following sequence of events occurrs without error:
//
//	(1) TF apply creates environment
//	(2) TF state values are set
//	(3) Environment can be queried by ID and has expected name
//	(4) TF apply updates environment with new name
//	(5) Environment can be queried by ID and has expected name
//	(6) TF destroy deletes environment
//	(7) Environment can no longer be queried by ID
func TestAccEnvironment_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	environmentNameFirst := testutils.GenerateResourceName()
	environmentNameSecond := testutils.GenerateResourceName()
	tfNode := "azuredevops_environment.environment"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkEnvironmentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclEnvironmentResource(projectName, environmentNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", environmentNameFirst),
					resource.TestCheckResourceAttr(tfNode, "description", ""),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					checkEnvironmentExists(environmentNameFirst),
				),
			},
			{
				Config: testutils.HclEnvironmentResource(projectName, environmentNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", environmentNameSecond),
					resource.TestCheckResourceAttr(tfNode, "description", ""),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					checkEnvironmentExists(environmentNameSecond),
				),
			},
			{
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Given the name of an environment, this will return a function that will check whether
// or not the environment (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkEnvironmentExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_environment.environment"]
		if !ok {
			return fmt.Errorf("Did not find an environment in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Parse ID error, ID:  %v !. Error= %v", resource.Primary.ID, err)
		}
		projectID := resource.Primary.Attributes["project_id"]

		environment, err := readEnvironment(clients, id, projectID)

		if err != nil {
			return fmt.Errorf("Environment with ID=%d cannot be found!. Error=%v", id, err)
		}

		if *environment.Name != expectedName {
			return fmt.Errorf("Environment with ID=%d has Name=%s, but expected Name=%s", id, *environment.Name, expectedName)
		}

		return nil
	}
}

// verifies that environment referenced in the state is destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkEnvironmentDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	// verify that every environment referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_environment" {
			continue
		}

		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Environment ID=%d cannot be parsed!. Error=%v", id, err)
		}
		projectID := resource.Primary.Attributes["project_id"]

		// indicates the environment still exists - this should fail the test
		if _, err := readEnvironment(clients, id, projectID); err == nil {
			return fmt.Errorf("Environment ID %d should not exist", id)
		}
	}

	return nil
}

// Lookup an Environment using the ID and the project ID.
func readEnvironment(clients *client.AggregatedClient, environmentID int, projectID string) (*taskagent.EnvironmentInstance, error) {
	return clients.TaskAgentClient.GetEnvironmentById(
		clients.Ctx,
		taskagent.GetEnvironmentByIdArgs{
			Project:       converter.String(projectID),
			EnvironmentId: &environmentID,
		},
	)
}
