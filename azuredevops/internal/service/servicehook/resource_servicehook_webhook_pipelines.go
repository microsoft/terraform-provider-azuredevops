package servicehook

import (
	"fmt"
	"log"
	"maps"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServicehookWebhookPipelines schemas a service hook subscription that
// uses the `pipelines` publisher and the `webHooks` (HTTP POST) consumer.
//
// It mirrors ResourceServicehookWebhookTfs but composes the pipelines publisher
// schema (`stage_state_changed_event`, `run_state_changed_event`) defined in
// pipelines_publisher.go instead of the TFS publisher schema.
func ResourceServicehookWebhookPipelines() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"project_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the project",
		},
		"url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  "The URL to send HTTP POST to",
		},
		"accept_untrusted_certs": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Accept untrusted SSL certificates",
		},
		"basic_auth_username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Basic authentication username",
		},
		"basic_auth_password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Basic authentication password",
		},
		"http_headers": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "HTTP headers as key-value pairs",
		},
		"resource_details_to_send": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "all",
			ValidateFunc: validation.StringInSlice([]string{"all", "minimal", "none"}, false),
			Description:  "Resource details to send - all, minimal, or none",
		},
		"messages_to_send": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "all",
			ValidateFunc: validation.StringInSlice([]string{"all", "text", "html", "markdown", "none"}, false),
			Description:  "Messages to send - all, text, html, markdown or none",
		},
		"detailed_messages_to_send": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "all",
			ValidateFunc: validation.StringInSlice([]string{"all", "text", "html", "markdown", "none"}, false),
			Description:  "Detailed messages to send - all, text, html, markdown or none",
		},
		"resource_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "5.1-preview.1",
			Description: "The resource version for the webhook subscription. The pipelines publisher events require a preview API version (default: 5.1-preview.1).",
		},
	}

	maps.Copy(resourceSchema, genPipelinesPublisherSchema())

	return &schema.Resource{
		Create: resourceServicehookWebhookPipelinesCreate,
		Read:   resourceServicehookWebhookPipelinesRead,
		Update: resourceServicehookWebhookPipelinesUpdate,
		Delete: resourceServicehookWebhookPipelinesDelete,
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

func resourceServicehookWebhookPipelinesCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookWebhookPipelines(d)
	createdSubscription, err := createSubscription(d, clients, subscription)
	if err != nil {
		return err
	}

	d.SetId(createdSubscription.Id.String())
	return resourceServicehookWebhookPipelinesRead(d, m)
}

func resourceServicehookWebhookPipelinesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionId := converter.UUID(d.Id())
	subscription, err := getSubscription(clients, subscriptionId)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			log.Printf("[INFO] Service hook subscription not found. ID: %s", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	flattenServicehookWebhookPipelines(d, subscription)
	return nil
}

func resourceServicehookWebhookPipelinesUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookWebhookPipelines(d)
	parsedID, err := uuid.Parse(d.Id())
	if err != nil {
		return err
	}
	subscription.Id = &parsedID

	if _, err := updateSubscription(clients, subscription); err != nil {
		return err
	}

	return resourceServicehookWebhookPipelinesRead(d, m)
}

func resourceServicehookWebhookPipelinesDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionID := converter.UUID(d.Id())
	return deleteSubscription(clients, subscriptionID)
}

func expandServicehookWebhookPipelines(d *schema.ResourceData) *servicehooks.Subscription {
	publisherInputs, eventType := expandPipelinesEventConfig(d)

	consumerInputs := map[string]string{
		"url":                    d.Get("url").(string),
		"acceptUntrustedCerts":   strconv.FormatBool(d.Get("accept_untrusted_certs").(bool)),
		"resourceDetailsToSend":  d.Get("resource_details_to_send").(string),
		"messagesToSend":         d.Get("messages_to_send").(string),
		"detailedMessagesToSend": d.Get("detailed_messages_to_send").(string),
	}

	if username := d.Get("basic_auth_username").(string); username != "" {
		consumerInputs["basicAuthUsername"] = username
	}
	if password := d.Get("basic_auth_password").(string); password != "" {
		consumerInputs["basicAuthPassword"] = password
	}

	if headersRaw, ok := d.GetOk("http_headers"); ok {
		headers := headersRaw.(map[string]interface{})
		keys := make([]string, 0, len(headers))
		for key := range headers {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		var headerString string
		for _, key := range keys {
			if headerString != "" {
				headerString += "\n"
			}
			headerString += fmt.Sprintf("%s:%s", key, headers[key].(string))
		}
		if headerString != "" {
			consumerInputs["httpHeaders"] = headerString
		}
	}

	return &servicehooks.Subscription{
		ConsumerActionId: converter.String("httpRequest"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerInputs:   &consumerInputs,
		EventType:        &eventType,
		PublisherId:      converter.String("pipelines"),
		PublisherInputs:  &publisherInputs,
		ResourceVersion:  converter.String(d.Get("resource_version").(string)),
	}
}

func flattenServicehookWebhookPipelines(d *schema.ResourceData, subscription *servicehooks.Subscription) {
	eventType, eventConfig := flattenPipelinesEventConfig(subscription)
	d.Set(eventType, eventConfig)
	d.Set("project_id", (*subscription.PublisherInputs)["projectId"])
	d.Set("url", (*subscription.ConsumerInputs)["url"])
	if subscription.ResourceVersion != nil {
		d.Set("resource_version", *subscription.ResourceVersion)
	}

	if acceptUntrustedCerts, exists := (*subscription.ConsumerInputs)["acceptUntrustedCerts"]; exists {
		if acceptBool, err := strconv.ParseBool(acceptUntrustedCerts); err == nil {
			d.Set("accept_untrusted_certs", acceptBool)
		}
	}

	if v, exists := (*subscription.ConsumerInputs)["resourceDetailsToSend"]; exists {
		d.Set("resource_details_to_send", v)
	}
	if v, exists := (*subscription.ConsumerInputs)["messagesToSend"]; exists {
		d.Set("messages_to_send", v)
	}
	if v, exists := (*subscription.ConsumerInputs)["detailedMessagesToSend"]; exists {
		d.Set("detailed_messages_to_send", v)
	}

	if headersString, exists := (*subscription.ConsumerInputs)["httpHeaders"]; exists && headersString != "" {
		headers := make(map[string]string)
		for _, line := range strings.Split(headersString, "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		if len(headers) > 0 {
			d.Set("http_headers", headers)
		}
	}
}
