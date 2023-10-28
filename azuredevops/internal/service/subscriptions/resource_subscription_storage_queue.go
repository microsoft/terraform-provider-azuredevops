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
		},
	}
}

func resourceSubscriptionStorageQueueCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription, err := expandSubscriptionStorageQueue(d)
	if err != nil {
		return err
	}

	createdSubscription, err := createSubscription(d, clients, subscription)
	if err != nil {
		return err
	}

	d.SetId(createdSubscription.Id.String())
	return resourceSubscriptionStorageQueueRead(d, m)
}

func resourceSubscriptionStorageQueueRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionId := converter.UUID(d.Id())
	subscription, err := getSubscription(clients, subscriptionId)
	if err != nil {
		return err
	}
	flattenSubscriptionStorageQueue(d, subscription)
	return nil
}

func resourceSubscriptionStorageQueueUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription, err := expandSubscriptionStorageQueue(d)
	if err != nil {
		return err
	}

	newSubscription, err := updateSubscription(clients, subscription)
	if err != nil {
		return err
	}

	flattenSubscriptionStorageQueue(d, newSubscription)
	return resourceSubscriptionStorageQueueRead(d, m)
}

func resourceSubscriptionStorageQueueDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	return clients.ServiceHooksClient.DeleteSubscription(clients.Ctx, servicehooks.DeleteSubscriptionArgs{
		SubscriptionId: converter.UUID(d.Id()),
	})
}

func expandSubscriptionStorageQueue(d *schema.ResourceData) (*servicehooks.Subscription, error) {
	consumerInputs := expandConsumerInputs(d.Get("consumer_inputs").(map[string]interface{}))
	publisherInputs := expandPublisherInputs(d.Get("publisher_inputs").(map[string]interface{}))
	return &servicehooks.Subscription{
		Id:               converter.UUID(d.Id()),
		ConsumerActionId: converter.String(d.Get("consumer_action_id").(string)),
		ConsumerId:       converter.String(d.Get("consumer_id").(string)),
		ConsumerInputs:   consumerInputs,
		EventType:        converter.String(d.Get("event_type").(string)),
		PublisherId:      converter.String(d.Get("publisher_id").(string)),
		PublisherInputs:  publisherInputs,
		ResourceVersion:  converter.String(d.Get("resource_version").(string)),
	}, nil
}

func expandConsumerInputs(inputs map[string]interface{}) *map[string]string {
	consumerInputs := make(map[string]string)
	consumerInputs["accountName"] = inputs["account_name"].(string)
	consumerInputs["accountKey"] = inputs["account_key"].(string)
	consumerInputs["queueName"] = inputs["queue_name"].(string)
	consumerInputs["visiTimeout"] = inputs["visi_timeout"].(string)
	consumerInputs["ttl"] = inputs["ttl"].(string)

	return &consumerInputs
}

func expandPublisherInputs(inputs map[string]interface{}) *map[string]string {
	publisherInputs := make(map[string]string)
	publisherInputs["pipelineId"] = inputs["pipeline_id"].(string)
	publisherInputs["stageNameId"] = inputs["stage_name_id"].(string)
	publisherInputs["stageStateId"] = inputs["stage_state_id"].(string)
	publisherInputs["stageResultId"] = inputs["stage_result_id"].(string)
	publisherInputs["projectId"] = inputs["project_id"].(string)

	return &publisherInputs
}

func flattenSubscriptionStorageQueue(d *schema.ResourceData, subscription *servicehooks.Subscription) {
	d.SetId(subscription.Id.String())
	consumerInputs := flattenConsumerInputs(subscription.ConsumerInputs)
	publisherInputs := flattenPublisherInputs(subscription.PublisherInputs)
	d.Set("consumer_action_id", subscription.ConsumerActionId)
	d.Set("consumer_id", subscription.ConsumerId)
	d.Set("consumer_inputs", consumerInputs)
	d.Set("event_type", subscription.EventType)
	d.Set("publisher_id", subscription.PublisherId)
	d.Set("publisher_inputs", publisherInputs)
	d.Set("resource_version", subscription.ResourceVersion)
}

func flattenConsumerInputs(inputs *map[string]string) map[string]interface{} {
	consumerInputs := make(map[string]interface{})
	consumerInputs["account_name"] = (*inputs)["accountName"]
	consumerInputs["account_key"] = (*inputs)["accountKey"]
	consumerInputs["queue_name"] = (*inputs)["queueName"]
	consumerInputs["visi_timeout"] = (*inputs)["visiTimeout"]
	consumerInputs["ttl"] = (*inputs)["ttl"]

	return consumerInputs
}

func flattenPublisherInputs(inputs *map[string]string) map[string]interface{} {
	publisherInputs := make(map[string]interface{})
	publisherInputs["pipeline_id"] = (*inputs)["pipelineId"]
	publisherInputs["stage_name_id"] = (*inputs)["stageNameId"]
	publisherInputs["stage_state_id"] = (*inputs)["stageStateId"]
	publisherInputs["stage_result_id"] = (*inputs)["stageResultId"]
	publisherInputs["project_id"] = (*inputs)["projectId"]

	return publisherInputs
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
