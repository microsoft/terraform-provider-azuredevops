package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/google/uuid"
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
	serviceEndpoint, _, err := expandServiceEndpointIncomingWebhook(d)
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

	flattenServiceEndpointIncomingWebhook(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	return nil
}

func resourceServiceEndpointIncomingWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointIncomingWebhook(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointIncomingWebhook(d, updatedServiceEndpoint, projectID.String())
	return resourceServiceEndpointIncomingWebhookRead(d, m)
}

func resourceServiceEndpointIncomingWebhookDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointIncomingWebhook(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
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
func flattenServiceEndpointIncomingWebhook(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("webhook_name", (*serviceEndpoint.Authorization.Parameters)["webhookname"])
	d.Set("http_header", (*serviceEndpoint.Authorization.Parameters)["header"])
}
