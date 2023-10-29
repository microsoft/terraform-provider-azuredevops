package subscription

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

var (
	consumerActionId = "enqueue"
	consumerId       = "azureStorageQueue"
)

func ResourceSubscriptionStorageQueue() *schema.Resource {
	resourceSchema := genPublisherSchema()
	resourceSchema["project_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.IsUUID,
	}
	resourceSchema["account_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	resourceSchema["account_key"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Sensitive:    true,
		ValidateFunc: validation.StringLenBetween(64, 100),
	}
	resourceSchema["queue_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	resourceSchema["visi_timeout"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  "0",
	}
	resourceSchema["ttl"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  "604800",
	}

	return &schema.Resource{
		Create: resourceSubscriptionStorageQueueCreate,
		Read:   resourceSubscriptionStorageQueueRead,
		Update: resourceSubscriptionStorageQueueUpdate,
		Delete: resourceSubscriptionStorageQueueDelete,

		Schema: resourceSchema,
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
	flattenSubscriptionStorageQueue(d, subscription, d.Get("account_key").(string))
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

	flattenSubscriptionStorageQueue(d, newSubscription, d.Get("account_key").(string))
	return resourceSubscriptionStorageQueueRead(d, m)
}

func resourceSubscriptionStorageQueueDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	return clients.ServiceHooksClient.DeleteSubscription(clients.Ctx, servicehooks.DeleteSubscriptionArgs{
		SubscriptionId: converter.UUID(d.Id()),
	})
}

func expandSubscriptionStorageQueue(d *schema.ResourceData) (*servicehooks.Subscription, error) {
	var subscriptionId *uuid.UUID
	parsedID, err := uuid.Parse(d.Id())
	if err == nil {
		subscriptionId = &parsedID
	}
	visiTimeout := strconv.Itoa(d.Get("visi_timeout").(int))
	ttl := strconv.Itoa(d.Get("ttl").(int))
	publisherInputs := expandPublisherInputs(d.Get("project_id").(string), d.Get("publisher").([]interface{}))
	resourceVersion := publisherResourceVersionMap[d.Get("publisher.0.name").(string)]
	return &servicehooks.Subscription{
		Id:               subscriptionId,
		ConsumerActionId: &consumerActionId,
		ConsumerId:       &consumerId,
		ConsumerInputs: &map[string]string{
			"accountName": d.Get("account_name").(string),
			"accountKey":  d.Get("account_key").(string),
			"queueName":   d.Get("queue_name").(string),
			"visiTimeout": visiTimeout,
			"ttl":         ttl,
		},
		EventType:       getEventType(d.Get("publisher").([]interface{})[0].(map[string]interface{})),
		PublisherId:     converter.String(d.Get("publisher.0.name").(string)),
		PublisherInputs: publisherInputs,
		ResourceVersion: &resourceVersion,
	}, nil
}

func flattenSubscriptionStorageQueue(d *schema.ResourceData, subscription *servicehooks.Subscription, accountKey string) {
	d.SetId(subscription.Id.String())
	visiTimeout, err := strconv.Atoi((*subscription.ConsumerInputs)["visiTimeout"])
	if err != nil {
		visiTimeout = 0
	}
	ttl, err := strconv.Atoi((*subscription.ConsumerInputs)["ttl"])
	if err != nil {
		ttl = 604800
	}
	publisher := flattenPublisherInputs(*subscription.PublisherId, *subscription.PublisherInputs)
	d.Set("project_id", (*subscription.PublisherInputs)["projectId"])
	d.Set("account_name", (*subscription.ConsumerInputs)["accountName"])
	d.Set("account_key", accountKey)
	d.Set("queue_name", (*subscription.ConsumerInputs)["queueName"])
	d.Set("visi_timeout", visiTimeout)
	d.Set("ttl", ttl)
	d.Set("publisher", publisher)
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
