package testutils

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

// GenerateFieldName generates a valid field name without hyphens or other invalid characters
func GenerateFieldName() string {
	return strings.ReplaceAll(GenerateResourceName(), "-", "")
}

// CheckFieldDestroyed verifies that all fields referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func CheckFieldDestroyed(s *terraform.State) error {
	clients := GetProvider().Meta().(*client.AggregatedClient)
	timeout := 10 * time.Second

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_workitemtracking_field" {
			continue
		}

		referenceName := resource.Primary.ID
		var project *string
		if projectID, ok := resource.Primary.Attributes["project_id"]; ok && projectID != "" {
			project = &projectID
		}

		err := retry.RetryContext(clients.Ctx, timeout, func() *retry.RetryError {
			_, err := readField(clients, referenceName, project)
			if err == nil {
				return retry.RetryableError(fmt.Errorf("field with reference name %s should not exist", referenceName))
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

func readField(clients *client.AggregatedClient, referenceName string, project *string) (*workitemtracking.WorkItemField2, error) {
	return clients.WorkItemTrackingClient.GetWorkItemField(clients.Ctx, workitemtracking.GetWorkItemFieldArgs{
		FieldNameOrRefName: &referenceName,
		Project:            project,
	})
}