package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointGeneric schema and implementation for generic service endpoint resource
func ResourceServiceEndpointGeneric() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointGenericCreate,
		Read:   resourceServiceEndpointGenericRead,
		Update: resourceServiceEndpointGenericUpdate,
		Delete: resourceServiceEndpointGenericDelete,
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
		"server_url": {
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			Description:  "The server URL of the generic service connection.",
		},

		"username": {
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_SERVICE_CONNECTION_USERNAME", nil),
			Description: "The username to use for the generic service connection.",
			Optional:    true,
		},

		"password": {
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_SERVICE_CONNECTION_PASSWORD", nil),
			Description: "The password or token key to use for the generic service connection.",
			Sensitive:   true,
			Optional:    true,
		},
	})
	return r
}

func resourceServiceEndpointGenericCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointGeneric(d)
	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointGenericRead(d, m)
}

func resourceServiceEndpointGenericRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if isServiceEndpointDeleted(d, err, serviceEndpoint, getArgs) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointGeneric(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointGenericUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointGeneric(d)
	if _, err := updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointGenericRead(d, m)
}

func resourceServiceEndpointGenericDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointGeneric(d)
	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointGeneric(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("generic")
	serviceEndpoint.Url = converter.String(d.Get("server_url").(string))
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	return serviceEndpoint
}

func flattenServiceEndpointGeneric(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("server_url", *serviceEndpoint.Url)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
}
