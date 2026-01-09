package testutils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

// CheckListDestroyed verifies that all lists referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func CheckListDestroyed(s *terraform.State) error {
	clients := GetProvider().Meta().(*client.AggregatedClient)
	timeout := 10 * time.Second

	// verify that every list referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_workitemtrackingprocess_list" {
			continue
		}

		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return err
		}

		err = retry.RetryContext(clients.Ctx, timeout, func() *retry.RetryError {
			_, err := readList(clients, id)
			if err == nil {
				return retry.RetryableError(fmt.Errorf("list with ID %s should not exist", id.String()))
			}
			if utils.ResponseWasNotFound(err) {
				return nil
			}

			return retry.NonRetryableError(err)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func readList(clients *client.AggregatedClient, id uuid.UUID) (*workitemtrackingprocess.PickList, error) {
	return clients.WorkItemTrackingProcessClient.GetList(clients.Ctx, workitemtrackingprocess.GetListArgs{
		ListId: &id,
	})
}
