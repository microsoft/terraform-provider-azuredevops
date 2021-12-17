package testutils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// CheckServiceHookWebhookExistsWithEventTypeAndUrl verifies that a service hook webhook exists in the state,
// and that it has the expected event type and url when compared against the data in Azure DevOps.
func CheckServiceHookWebhookExistsWithEventTypeAndUrl(tfNode string, expectedEventType string, expectedUrl string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[tfNode]
		if !ok {
			return fmt.Errorf("Did not find a service hook webhook in the state")
		}

		subscription, err := getSvcHookWebhookFromState(resourceState)
		if err != nil {
			return err
		}

		if *subscription.EventType != expectedEventType {
			return fmt.Errorf("Service Hook webhook has event type=%s, but expected event type=%s", *subscription.EventType, expectedEventType)
		}

		if (*subscription.ConsumerInputs)["url"] != expectedUrl {
			return fmt.Errorf("Service Hook webhook has url=%s, but expected url=%s", *subscription.Url, expectedUrl)
		}

		return nil
	}
}

// CheckServiceHookWebhookDestroyed verifies that all service hook webhoks the state are destroyed.
// This will be invoked *after* terraform destroys the resource but *before* the state is wiped clean.
func CheckServiceHookWebhookDestroyed(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, resource := range s.RootModule().Resources {
			if resource.Type != resourceType {
				continue
			}

			// indicates the resource exists - this should fail the test
			if _, err := getSvcHookWebhookFromState(resource); err == nil {
				return fmt.Errorf("Unexpectedly found a service hook webhook that should have been deleted")
			}
		}

		return nil
	}
}

// given a resource from the state, return a service hook webhook (and error)
func getSvcHookWebhookFromState(resource *terraform.ResourceState) (*servicehooks.Subscription, error) {
	serviceHookWebhookDefID, err := uuid.Parse(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	clients := GetProvider().Meta().(*client.AggregatedClient)
	return clients.ServiceHooksClient.GetSubscription(clients.Ctx, servicehooks.GetSubscriptionArgs{
		SubscriptionId: &serviceHookWebhookDefID,
	})
}
