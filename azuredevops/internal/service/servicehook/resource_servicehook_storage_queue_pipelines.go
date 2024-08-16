package servicehook

import (
	"fmt"
	"maps"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceServicehookStorageQueuePipelines() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"project_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the project",
		},
		"account_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The queue's storage account name",
		},
		"account_key": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(64, 100),
			Description:  "A valid account key from the queue's storage account",
		},
		"queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the queue that will store the events",
		},
		"visi_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     "0",
			Description: "event visibility timout - how long a message is invisible to other consumers after it's been dequeued",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     "604800",
			Description: "event time-to-live - the duration a message can remain in the queue before it's automatically removed",
		},
	}

	maps.Copy(resourceSchema, genPipelinesPublisherSchema())

	return &schema.Resource{
		Create: resourceServicehookStorageQueuePipelinesCreate,
		Read:   resourceServicehookStorageQueuePipelinesRead,
		Update: resourceServicehookStorageQueuePipelinesUpdate,
		Delete: resourceServicehookStorageQueuePipelinesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

func resourceServicehookStorageQueuePipelinesCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription, err := expandServicehookStorageQueuePipelines(d)
	if err != nil {
		return err
	}

	createdSubscription, err := createSubscription(clients, subscription)
	if err != nil {
		return err
	}

	d.SetId(createdSubscription.Id.String())
	return resourceServicehookStorageQueuePipelinesRead(d, m)
}

func resourceServicehookStorageQueuePipelinesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionId := converter.UUID(d.Id())
	subscription, err := getSubscription(clients, subscriptionId)
	if err != nil {
		return err
	}
	flattenServicehookStorageQueuePipelines(d, subscription, d.Get("account_key").(string))
	return nil
}

func resourceServicehookStorageQueuePipelinesUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription, err := expandServicehookStorageQueuePipelines(d)
	if err != nil {
		return err
	}

	parsedID, err := uuid.Parse(d.Id())
	if err != nil {
		return err
	}
	subscription.Id = &parsedID

	_, err = updateSubscription(clients, subscription)
	if err != nil {
		return err
	}

	return resourceServicehookStorageQueuePipelinesRead(d, m)
}

func resourceServicehookStorageQueuePipelinesDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	return clients.ServiceHooksClient.DeleteSubscription(clients.Ctx, servicehooks.DeleteSubscriptionArgs{
		SubscriptionId: converter.UUID(d.Id()),
	})
}

func expandServicehookStorageQueuePipelines(d *schema.ResourceData) (*servicehooks.Subscription, error) {
	visiTimeout := strconv.Itoa(d.Get("visi_timeout").(int))
	ttl := strconv.Itoa(d.Get("ttl").(int))
	publisherInputs, eventType := expandPipelinesEventConfig(d)
	return &servicehooks.Subscription{
		ConsumerActionId: converter.String("enqueue"),
		ConsumerId:       converter.String("azureStorageQueue"),
		ConsumerInputs: &map[string]string{
			"accountName": d.Get("account_name").(string),
			"accountKey":  d.Get("account_key").(string),
			"queueName":   d.Get("queue_name").(string),
			"visiTimeout": visiTimeout,
			"ttl":         ttl,
		},
		EventType:       &eventType,
		PublisherId:     converter.String("pipelines"),
		PublisherInputs: &publisherInputs,
		ResourceVersion: converter.String("5.1-preview.1"),
	}, nil
}

func flattenServicehookStorageQueuePipelines(d *schema.ResourceData, subscription *servicehooks.Subscription, accountKey string) {
	visiTimeout, err := strconv.Atoi((*subscription.ConsumerInputs)["visiTimeout"])
	if err != nil {
		visiTimeout = 0
	}
	ttl, err := strconv.Atoi((*subscription.ConsumerInputs)["ttl"])
	if err != nil {
		ttl = 604800
	}

	eventType, eventConfig := flattenPipelinesEventConfig(subscription)
	d.Set(eventType, eventConfig)
	d.Set("project_id", (*subscription.PublisherInputs)["projectId"])
	d.Set("account_name", (*subscription.ConsumerInputs)["accountName"])
	d.Set("account_key", accountKey)
	d.Set("queue_name", (*subscription.ConsumerInputs)["queueName"])
	d.Set("visi_timeout", visiTimeout)
	d.Set("ttl", ttl)
}

func createSubscription(clients *client.AggregatedClient, subscription *servicehooks.Subscription) (*servicehooks.Subscription, error) {
	createdSubscription, err := clients.ServiceHooksClient.CreateSubscription(
		clients.Ctx,
		servicehooks.CreateSubscriptionArgs{
			Subscription: subscription,
		})
	if err != nil {
		return nil, fmt.Errorf(" creating subscription in Azure DevOps: %+v", err)
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

func getSubscription(client *client.AggregatedClient, subscriptionID *uuid.UUID) (*servicehooks.Subscription, error) {
	return client.ServiceHooksClient.GetSubscription(
		client.Ctx,
		servicehooks.GetSubscriptionArgs{
			SubscriptionId: subscriptionID,
		})
}
