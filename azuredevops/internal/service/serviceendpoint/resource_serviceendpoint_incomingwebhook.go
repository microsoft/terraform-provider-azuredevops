package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointIncomingWebhook schema and implementation for incoming webhook service endpoint resource
func ResourceServiceEndpointIncomingWebhook() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointIncomingWebhookCreate,
		Read:   resourceServiceEndpointIncomingWebhookRead,
		Update: resourceServiceEndpointIncomingWebhookUpdate,
		Delete: resourceServiceEndpointIncomingWebhookDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}
	maps.Copy(r.Schema, map[string]*schema.Schema{
		"webhook_name": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_INCOMING_WEBHOOK_SERVICE_CONNECTION_WEBHOOK_NAME", nil),
			Description: "The name of the WebHook.",
		},

		"secret": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_INCOMING_WEBHOOK_SERVICE_CONNECTION_SECRET", nil),
			Description: "Optional secret for the webhook. WebHook service will use this secret to calculate the payload checksum.",
			Sensitive:   true,
		},

		"http_header": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_INCOMING_WEBHOOK_SERVICE_CONNECTION_HTTP_HEADER", nil),
			Description: "Optional http header name on which checksum will be sent.",
		},
	})
	return r
}

func resourceServiceEndpointIncomingWebhookCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointIncomingWebhook(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointIncomingWebhookRead(d, m)
}

func resourceServiceEndpointIncomingWebhookRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointIncomingWebhook(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointIncomingWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointIncomingWebhook(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointIncomingWebhookRead(d, m)
}

func resourceServiceEndpointIncomingWebhookDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointIncomingWebhook(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointIncomingWebhook(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Url = converter.String("https://dev.azure.com")
	serviceEndpoint.Type = converter.String("incomingwebhook")
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"webhookname": d.Get("webhook_name").(string),
			"secret":      d.Get("secret").(string),
			"header":      d.Get("http_header").(string),
		},
		Scheme: converter.String("None"),
	}
	return serviceEndpoint, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointIncomingWebhook(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("webhook_name", (*serviceEndpoint.Authorization.Parameters)["webhookname"])
	d.Set("http_header", (*serviceEndpoint.Authorization.Parameters)["header"])
}
