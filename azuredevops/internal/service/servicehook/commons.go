package servicehook

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// createSubscription creates a new service hook subscription in Azure DevOps
func createSubscription(d *schema.ResourceData, clients *client.AggregatedClient, subscription *servicehooks.Subscription) (*servicehooks.Subscription, error) {
	if subscription == nil {
		return nil, fmt.Errorf("subscription cannot be nil")
	}

	createdSubscription, err := clients.ServiceHooksClient.CreateSubscription(
		clients.Ctx,
		servicehooks.CreateSubscriptionArgs{
			Subscription: subscription,
		})
	if err != nil {
		return nil, fmt.Errorf("Error creating service hook subscription in Azure DevOps: %+v", err)
	}

	// Since service hooks are simpler and don't have the same asynchronous nature as service endpoints,
	// we can directly return the created subscription without complex state checking
	return createdSubscription, nil
}

// updateSubscription updates an existing service hook subscription in Azure DevOps
func updateSubscription(clients *client.AggregatedClient, subscription *servicehooks.Subscription) (*servicehooks.Subscription, error) {
	if subscription == nil || subscription.Id == nil {
		return nil, fmt.Errorf("subscription and subscription ID cannot be nil")
	}

	updatedSubscription, err := clients.ServiceHooksClient.ReplaceSubscription(
		clients.Ctx,
		servicehooks.ReplaceSubscriptionArgs{
			Subscription:   subscription,
			SubscriptionId: subscription.Id,
		})
	if err != nil {
		return nil, fmt.Errorf("Error updating service hook subscription in Azure DevOps: %+v", err)
	}

	return updatedSubscription, nil
}

// deleteSubscription deletes a service hook subscription from Azure DevOps
func deleteSubscription(clients *client.AggregatedClient, subscriptionID *uuid.UUID) error {
	if subscriptionID == nil {
		return fmt.Errorf("subscription ID cannot be nil")
	}

	err := clients.ServiceHooksClient.DeleteSubscription(
		clients.Ctx,
		servicehooks.DeleteSubscriptionArgs{
			SubscriptionId: subscriptionID,
		})
	if err != nil {
		return fmt.Errorf("Error deleting service hook subscription: %+v", err)
	}

	return nil
}

// getSubscription retrieves a service hook subscription from Azure DevOps
func getSubscription(clients *client.AggregatedClient, subscriptionID *uuid.UUID) (*servicehooks.Subscription, error) {
	if subscriptionID == nil {
		return nil, fmt.Errorf("subscription ID cannot be nil")
	}

	subscription, err := clients.ServiceHooksClient.GetSubscription(
		clients.Ctx,
		servicehooks.GetSubscriptionArgs{
			SubscriptionId: subscriptionID,
		})
	return subscription, err
}
