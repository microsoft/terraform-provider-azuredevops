package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointIncomingWebhook schema and implementation for incoming webhook service endpoint resource
func ResourceServiceEndpointIncomingWebhook() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointIncomingWebhook, expandServiceEndpointIncomingWebhook)
	r.Schema["webhook_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_INCOMING_WEBHOOK_SERVICE_CONNECTION_WEBHOOK_NAME", nil),
		Description: "The name of the WebHook.",
	}
	r.Schema["secret"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_INCOMING_WEBHOOK_SERVICE_CONNECTION_SECRET", nil),
		Description: "Optional secret for the webhook. WebHook service will use this secret to calculate the payload checksum.",
		Sensitive:   true,
	}
	r.Schema["http_header"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_INCOMING_WEBHOOK_SERVICE_CONNECTION_HTTP_HEADER", nil),
		Description: "Optional http header name on which checksum will be sent.",
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointIncomingWebhook(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
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
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointIncomingWebhook(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("webhook_name", (*serviceEndpoint.Authorization.Parameters)["webhookname"])
	d.Set("http_header", (*serviceEndpoint.Authorization.Parameters)["header"])
}
