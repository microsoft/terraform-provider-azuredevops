package servicehook

import (
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServicehookSubscription schema and implementation for service hook subscription resource
func ResourceServicehookSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceServicehookSubscriptionCreate,
		Read:   resourceServicehookSubscriptionRead,
		Update: resourceServicehookSubscriptionUpdate,
		Delete: resourceServicehookSubscriptionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the project. If not provided, the subscription will be created at the organization level",
			},
			"publisher_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The publisher ID (e.g., 'tfs' for Team Foundation Server, 'pipelines' for Azure Pipelines)",
			},
			"event_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The event type (e.g., 'git.push', 'build.complete', 'workitem.created')",
			},
			"consumer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The consumer ID (e.g., 'webHooks', 'azureServiceBus', 'azureStorageQueue')",
			},
			"consumer_action_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The consumer action ID (e.g., 'httpRequest', 'enqueue')",
			},
			"publisher_inputs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Publisher-specific inputs as key-value pairs",
			},
			"consumer_inputs": {
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Sensitive:   true,
				Description: "Consumer-specific inputs as key-value pairs (sensitive)",
			},
			"resource_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1.0",
				Description: "The resource version for the subscription",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enabled",
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled", "disabledByUser", "disabledBySystem", "onProbation"}, false),
				Description:  "The status of the subscription",
			},
		},
	}
}

func resourceServicehookSubscriptionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookSubscription(d)

	createdSubscription, err := createSubscription(clients, subscription)
	if err != nil {
		return err
	}

	d.SetId(createdSubscription.Id.String())
	return resourceServicehookSubscriptionRead(d, m)
}

func resourceServicehookSubscriptionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionId := converter.UUID(d.Id())

	subscription, err := getSubscription(clients, subscriptionId)
	if err != nil {
		return err
	}

	// Check if subscription was deleted
	if subscription == nil {
		d.SetId("")
		return nil
	}

	return flattenServicehookSubscription(d, subscription)
}

func resourceServicehookSubscriptionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookSubscription(d)

	parsedID, err := uuid.Parse(d.Id())
	if err != nil {
		return err
	}
	subscription.Id = &parsedID

	_, err = updateSubscription(clients, subscription)
	if err != nil {
		return err
	}

	return resourceServicehookSubscriptionRead(d, m)
}

func resourceServicehookSubscriptionDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	return clients.ServiceHooksClient.DeleteSubscription(clients.Ctx, servicehooks.DeleteSubscriptionArgs{
		SubscriptionId: converter.UUID(d.Id()),
	})
}

func expandServicehookSubscription(d *schema.ResourceData) *servicehooks.Subscription {
	publisherInputs := make(map[string]string)

	// Add project_id to publisher inputs if provided
	if projectID, ok := d.GetOk("project_id"); ok {
		publisherInputs["projectId"] = projectID.(string)
	}

	// Add any additional publisher inputs
	if inputs, ok := d.GetOk("publisher_inputs"); ok {
		for key, value := range inputs.(map[string]interface{}) {
			publisherInputs[key] = value.(string)
		}
	}

	consumerInputs := make(map[string]string)
	if inputs, ok := d.GetOk("consumer_inputs"); ok {
		for key, value := range inputs.(map[string]interface{}) {
			consumerInputs[key] = value.(string)
		}
	}

	status := convertStatus(d.Get("status").(string))

	return &servicehooks.Subscription{
		PublisherId:      converter.String(d.Get("publisher_id").(string)),
		EventType:        converter.String(d.Get("event_type").(string)),
		ConsumerId:       converter.String(d.Get("consumer_id").(string)),
		ConsumerActionId: converter.String(d.Get("consumer_action_id").(string)),
		PublisherInputs:  &publisherInputs,
		ConsumerInputs:   &consumerInputs,
		ResourceVersion:  converter.String(d.Get("resource_version").(string)),
		Status:           &status,
	}
}

func flattenServicehookSubscription(d *schema.ResourceData, subscription *servicehooks.Subscription) error {
	d.Set("publisher_id", subscription.PublisherId)
	d.Set("event_type", subscription.EventType)
	d.Set("consumer_id", subscription.ConsumerId)
	d.Set("consumer_action_id", subscription.ConsumerActionId)
	d.Set("resource_version", subscription.ResourceVersion)

	if subscription.Status != nil {
		d.Set("status", convertStatusFromAPI(*subscription.Status))
	}

	if subscription.PublisherInputs != nil {
		publisherInputs := make(map[string]interface{})
		for key, value := range *subscription.PublisherInputs {
			// Don't expose projectId as a publisher input if we manage it separately
			if key != "projectId" {
				publisherInputs[key] = value
			} else if key == "projectId" {
				// Set project_id if it exists in publisher inputs
				d.Set("project_id", value)
			}
		}
		d.Set("publisher_inputs", publisherInputs)
	}

	// Note: We don't flatten consumer_inputs as they are sensitive and may contain secrets.
	// The user needs to manage them in their configuration.

	return nil
}

func convertStatus(status string) servicehooks.SubscriptionStatus {
	switch status {
	case "enabled":
		return servicehooks.SubscriptionStatusValues.Enabled
	case "disabled":
		return servicehooks.SubscriptionStatusValues.DisabledByUser
	case "disabledByUser":
		return servicehooks.SubscriptionStatusValues.DisabledByUser
	case "disabledBySystem":
		return servicehooks.SubscriptionStatusValues.DisabledBySystem
	case "onProbation":
		return servicehooks.SubscriptionStatusValues.OnProbation
	default:
		return servicehooks.SubscriptionStatusValues.Enabled
	}
}

func convertStatusFromAPI(status servicehooks.SubscriptionStatus) string {
	switch status {
	case servicehooks.SubscriptionStatusValues.Enabled:
		return "enabled"
	case servicehooks.SubscriptionStatusValues.DisabledByUser:
		return "disabledByUser"
	case servicehooks.SubscriptionStatusValues.DisabledBySystem:
		return "disabledBySystem"
	case servicehooks.SubscriptionStatusValues.OnProbation:
		return "onProbation"
	default:
		return "enabled"
	}
}
