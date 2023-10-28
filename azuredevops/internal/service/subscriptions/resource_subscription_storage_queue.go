package subscriptions

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceSubscriptionStorageQueue schema and implementation for storage queue subscription
func ResourceSubscriptionStorageQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionStorageQueueCreate,
		Read:   resourceSubscriptionStorageQueueRead,
		Update: resourceSubscriptionStorageQueueUpdate,
		Delete: resourceSubscriptionStorageQueueDelete,

		Schema: map[string]*schema.Schema{
			"consumer_action_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consumer_inputs": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"account_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"visi_timeout": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "0",
						},
						"ttl": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "604800",
						},
					},
				},
			},
			"event_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"publisher_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"publisher_inputs": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"stage_name_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"stage_state_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"stage_result_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"resource_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "5.1-preview.1",
			},
			"scope": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourceSubscriptionStorageQueueCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubscriptionStorageQueueRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubscriptionStorageQueueUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubscriptionStorageQueueDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func expandSubscriptionStorageQueue(d *schema.ResourceData) (*servicehooks.Subscription, error) {
	// Populate the Go structure from the Terraform schema
	return &servicehooks.Subscription{
		ConsumerActionId: converter.String(d.Get("consumer_action_id").(string)),
		ConsumerId:       converter.String(d.Get("consumer_id").(string)),
		ConsumerInputs:   d.Get("consumer_inputs").(map[string]interface{}),
		EventType:        converter.String(d.Get("event_type").(string)),
		PublisherId:      converter.String(d.Get("publisher_id").(string)),
		PublisherInputs:  d.Get("publisher_inputs").(map[string]interface{}),
		ResourceVersion:  converter.String(d.Get("resource_version").(string)),
	}, nil
}

func expandConsumerInputs(inputs map[string]interface{}) (*map[string]interface{}, error) {
	consumerInputs := make(map[string]interface{})
	consumerInputs["account_name"] = inputs["account_name"]
	consumerInputs["account_key"] = inputs["account_key"]
	consumerInputs["queue_name"] = inputs["queue_name"]
	consumerInputs["visi_timeout"] = inputs["visi_timeout"]
	consumerInputs["ttl"] = inputs["ttl"]

	return &consumerInputs, nil
}

func flattenSubscriptionStorageQueue(d *schema.ResourceData, subscription *servicehooks.Subscription) error {
	// Set the fields in the Terraform schema from the Go structure
	d.Set("consumer_action_id", subscription.ConsumerActionID)
	d.Set("consumer_id", subscription.ConsumerID)
	d.Set("consumer_inputs", subscription.ConsumerInputs)
	d.Set("event_type", subscription.EventType)
	d.Set("publisher_id", subscription.PublisherID)
	d.Set("publisher_inputs", subscription.PublisherInputs)
	d.Set("resource_version", subscription.ResourceVersion)
	d.Set("scope", subscription.Scope)

	return nil
}

func createSubscription(d *schema.ResourceData, clients *client.AggregatedClient, subscription *servicehooks.Subscription) (*servicehooks.Subscription, error) {
	createdSubscription, err := clients.ServiceHooksClient.CreateSubscription(
		clients.Ctx,
		servicehooks.CreateSubscriptionArgs{
			Subscription: subscription,
		})
	if err != nil {
		return nil, fmt.Errorf("Error creating subscription in Azure DevOps: %+v", err)
	}

	return createdSubscription, err
}

func updateSubscription(clients *client.AggregatedClient, subscription *servicehooks.Subscription) (*servicehooks.Subscription, error) {
	updatedSubscription, err := clients.ServiceHooksClient.ReplaceSubscription(
		clients.Ctx,
		servicehooks.ReplaceSubscriptionArgs{
			Subscription:   subscription,
			SubscriptionId: subscription.Id,
		})

	return updatedSubscription, err
}

func deleteSubscription(clients *client.AggregatedClient, subscriptionID *uuid.UUID) error {
	if err := clients.ServiceHooksClient.DeleteSubscription(
		clients.Ctx,
		servicehooks.DeleteSubscriptionArgs{
			SubscriptionId: subscriptionID,
		}); err != nil {
		return fmt.Errorf(" Delete subscription error %v", err)
	}

	return nil
}

func getSubscription(client *client.AggregatedClient, subscriptionID *uuid.UUID) (*servicehooks.Subscription, error) {
	return client.ServiceHooksClient.GetSubscription(
		client.Ctx,
		servicehooks.GetSubscriptionArgs{
			SubscriptionId: subscriptionID,
		})
}
