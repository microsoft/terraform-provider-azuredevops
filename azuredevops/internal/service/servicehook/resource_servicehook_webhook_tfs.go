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

func ResourceServicehookWebhookTfs() *schema.Resource {
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
			Description:  "Resource details to send - all, text, html, markdown or none",
		},
		"detailed_messages_to_send": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "all",
			ValidateFunc: validation.StringInSlice([]string{"all", "text", "html", "markdown", "none"}, false),
			Description:  "Detailed messages to send - all, text, html, markdown or none",
		},
	}

	maps.Copy(resourceSchema, genTfsPublisherSchema())

	return &schema.Resource{
		Create: resourceServicehookWebhookTfsCreate,
		Read:   resourceServicehookWebhookTfsRead,
		Update: resourceServicehookWebhookTfsUpdate,
		Delete: resourceServicehookWebhookTfsDelete,
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

func resourceServicehookWebhookTfsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookWebhookTfs(d)
	createdSubscription, err := createSubscription(d, clients, subscription)
	if err != nil {
		return err
	}

	d.SetId(createdSubscription.Id.String())
	return resourceServicehookWebhookTfsRead(d, m)
}

func resourceServicehookWebhookTfsRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServicehookWebhookTfs(d, subscription)
	return nil
}

func resourceServicehookWebhookTfsUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookWebhookTfs(d)
	parsedID, err := uuid.Parse(d.Id())
	if err != nil {
		return err
	}
	subscription.Id = &parsedID

	_, err = updateSubscription(clients, subscription)
	if err != nil {
		return err
	}

	return resourceServicehookWebhookTfsRead(d, m)
}

func resourceServicehookWebhookTfsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionID := converter.UUID(d.Id())

	return deleteSubscription(clients, subscriptionID)
}

func expandServicehookWebhookTfs(d *schema.ResourceData) *servicehooks.Subscription {
	publisherInputs, eventType := expandTfsEventConfig(d)

	// Construct consumer inputs
	consumerInputs := map[string]string{
		"url":                    d.Get("url").(string),
		"acceptUntrustedCerts":   strconv.FormatBool(d.Get("accept_untrusted_certs").(bool)),
		"resourceDetailsToSend":  d.Get("resource_details_to_send").(string),
		"messagesToSend":         d.Get("messages_to_send").(string),
		"detailedMessagesToSend": d.Get("detailed_messages_to_send").(string),
	}

	// Add basic auth if provided
	if username := d.Get("basic_auth_username").(string); username != "" {
		consumerInputs["basicAuthUsername"] = username
	}
	if password := d.Get("basic_auth_password").(string); password != "" {
		consumerInputs["basicAuthPassword"] = password
	}

	// Add HTTP headers if provided
	if headersRaw, ok := d.GetOk("http_headers"); ok {
		headers := headersRaw.(map[string]interface{})
		var headerString string
		// Sort keys to ensure consistent ordering
		keys := make([]string, 0, len(headers))
		for key := range headers {
			keys = append(keys, key)
		}
		sort.Strings(keys)

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
		PublisherId:      converter.String("tfs"),
		PublisherInputs:  &publisherInputs,
		ResourceVersion:  converter.String("7.1"),
	}
}

func flattenServicehookWebhookTfs(d *schema.ResourceData, subscription *servicehooks.Subscription) {
	eventType, eventConfig := flattenTfsEventConfig(subscription)
	d.Set(eventType, eventConfig)
	d.Set("project_id", (*subscription.PublisherInputs)["projectId"])
	d.Set("url", (*subscription.ConsumerInputs)["url"])

	// Parse acceptUntrustedCerts
	if acceptUntrustedCerts, exists := (*subscription.ConsumerInputs)["acceptUntrustedCerts"]; exists {
		if acceptBool, err := strconv.ParseBool(acceptUntrustedCerts); err == nil {
			d.Set("accept_untrusted_certs", acceptBool)
		}
	}

	// Set resource details and messages
	if resourceDetails, exists := (*subscription.ConsumerInputs)["resourceDetailsToSend"]; exists {
		d.Set("resource_details_to_send", resourceDetails)
	}
	if messages, exists := (*subscription.ConsumerInputs)["messagesToSend"]; exists {
		d.Set("messages_to_send", messages)
	}
	if detailedMessages, exists := (*subscription.ConsumerInputs)["detailedMessagesToSend"]; exists {
		d.Set("detailed_messages_to_send", detailedMessages)
	}

	// Parse HTTP headers
	if headersString, exists := (*subscription.ConsumerInputs)["httpHeaders"]; exists && headersString != "" {
		headers := make(map[string]string)
		// Split by newlines and parse key:value pairs
		for _, line := range strings.Split(headersString, "\n") {
			if line != "" {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
		}
		if len(headers) > 0 {
			d.Set("http_headers", headers)
		}
	}
}
