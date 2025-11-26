package testutils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// CheckProcessDestroyed verifies that all processes referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func CheckProcessDestroyed(s *terraform.State) error {
	clients := GetProvider().Meta().(*client.AggregatedClient)

	// verify that every process referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_workitemtrackingprocess_process" {
			continue
		}

		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return err
		}

		if _, err := readProcess(clients, id); err == nil {
			return fmt.Errorf("process with ID %s should not exist", id.String())
		}
	}

	return nil
}

func readProcess(clients *client.AggregatedClient, identifier uuid.UUID) (*core.Process, error) {
	return clients.CoreClient.GetProcessById(clients.Ctx, core.GetProcessByIdArgs{
		ProcessId: &identifier,
	})
}
