package servicehooks

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceHookWebhook() *schema.Resource {
	return &schema.Resource{
		Create:   resourceWebhookCreate,
		Read:     createResourceWebhookRead(false),
		Update:   resourceWebhookUpdate,
		Delete:   resourceWebhookDelete,
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"event_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"basic_auth": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      nil,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
							Sensitive:    true,
						},
					},
				},
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"http_headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"filters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceWebhookCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	subscriptionData := getSubscription(d)

	subscription, err := clients.ServiceHooksClient.CreateSubscription(clients.Ctx, servicehooks.CreateSubscriptionArgs{
		Subscription: &subscriptionData,
	})

	if err != nil {
		return err
	}

	d.SetId(subscription.Id.String())

	return createResourceWebhookRead(true)(d, m)
}

func createResourceWebhookRead(afterCreateOrUpdate bool) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)

		subscriptionId := d.Id()

		subscription, err := clients.ServiceHooksClient.GetSubscription(clients.Ctx, servicehooks.GetSubscriptionArgs{
			SubscriptionId: converter.UUID(subscriptionId),
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil
			}
			return err
		}

		d.Set("project_id", (*subscription.PublisherInputs)["projectId"])
		d.Set("url", (*subscription.ConsumerInputs)["url"])
		d.Set("event_type", *subscription.EventType)

		oldUpdatedAt := d.Get("updated_at")
		newUpdatedAt := subscription.ModifiedDate.String()
		d.Set("updated_at", newUpdatedAt)

		if basicAuthList, ok := d.GetOk("basic_auth"); ok {
			basicAuth := basicAuthList.([]interface{})[0].(map[string]interface{})

			if username, ok := (*subscription.ConsumerInputs)["basicAuthUsername"]; ok {
				basicAuth["username"] = username
			}

			if password, ok := basicAuth["password"]; ok && !afterCreateOrUpdate && oldUpdatedAt != newUpdatedAt {
				// note: condition above means someone modified webhook subscription externally and since we can't
				// know whether they've changed password (API returns mask ********) we'll force a password change
				basicAuth["password"] = password.(string) + " - this suffix will cause a diff and therefore change"
			}

			d.Set("basic_auth", basicAuthList)
		}

		// http headers are returned as string -> we need to parse them
		httpHeadersString := (*subscription.ConsumerInputs)["httpHeaders"]
		reader := bufio.NewReader(strings.NewReader("GET / HTTP/1.1\r\n" + (*subscription.ConsumerInputs)["httpHeaders"] + "\r\n\n"))
		req, err := http.ReadRequest(reader)
		if err != nil {
			return fmt.Errorf("could not parse subscription http headers: %s", httpHeadersString)
		}

		httpHeaders := map[string]string{}
		for header, values := range req.Header {
			httpHeaders[header] = strings.Join(values, ", ")
		}
		d.Set("http_headers", httpHeaders)

		filters := map[string]string{}
		for key, value := range *subscription.PublisherInputs {
			if key == "projectId" || key == "tfsSubscriptionId" {
				continue
			}
			filters[key] = value
		}
		d.Set("filters", filters)

		return nil
	}
}

func resourceWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	subscriptionData := getSubscription(d)

	if _, err := clients.ServiceHooksClient.ReplaceSubscription(clients.Ctx, servicehooks.ReplaceSubscriptionArgs{
		SubscriptionId: converter.UUID(d.Id()),
		Subscription:   &subscriptionData,
	}); err != nil {
		return err
	}

	return createResourceWebhookRead(true)(d, m)
}

func resourceWebhookDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	err := clients.ServiceHooksClient.DeleteSubscription(clients.Ctx, servicehooks.DeleteSubscriptionArgs{
		SubscriptionId: converter.UUID(d.Id()),
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func getSubscription(d *schema.ResourceData) servicehooks.Subscription {
	publisherId := "tfs"
	eventType := d.Get("event_type").(string)
	url := d.Get("url").(string)
	consumerId := "webHooks"
	consumerActionId := "httpRequest"
	httpHeaders := d.Get("http_headers").(map[string]interface{})
	filters := d.Get("filters").(map[string]interface{})

	consumerInputs := map[string]string{
		"url": url,
	}
	if basicAuthList, ok := d.GetOk("basic_auth"); ok {
		basicAuth := basicAuthList.([]interface{})[0].(map[string]interface{})
		if username, ok := basicAuth["username"]; ok && username != "" {
			consumerInputs["basicAuthUsername"] = username.(string)
		}
		if password, ok := basicAuth["password"]; ok && password != "" {
			consumerInputs["basicAuthPassword"] = password.(string)
		}
	}

	httpHeadersSlice := []string{}
	for header, value := range httpHeaders {
		httpHeadersSlice = append(httpHeadersSlice, fmt.Sprintf("%s: %s", header, value.(string)))
	}
	httpHeadersStr := strings.Join(httpHeadersSlice, "\n")
	if httpHeadersStr != "" {
		consumerInputs["httpHeaders"] = httpHeadersStr
	}

	publisherInputs := map[string]string{}
	for key, value := range filters {
		publisherInputs[key] = value.(string)
	}
	publisherInputs["projectId"] = d.Get("project_id").(string)

	subscriptionData := servicehooks.Subscription{
		PublisherId:      &publisherId,
		EventType:        &eventType,
		ConsumerId:       &consumerId,
		ConsumerActionId: &consumerActionId,
		PublisherInputs:  &publisherInputs,
		ConsumerInputs:   &consumerInputs,
	}

	return subscriptionData
}
