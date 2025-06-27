package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// ResourceServiceEndpointNexus schema and implementation for Nexus service endpoint resource
func ResourceServiceEndpointNexus() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointNexusCreate,
		Read:   resourceServiceEndpointNexusRead,
		Update: resourceServiceEndpointNexusUpdate,
		Delete: resourceServiceEndpointNexusDelete,
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
		"url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.Url,
			Description:  "Url for the Nexus Repository",
		},

		"username": {
			Description: "The Nexus user name.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"password": {
			Description: "The Nexus password.",
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
		},
	})

	return r
}

func resourceServiceEndpointNexusCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointNexus(d)
	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointNexusRead(d, m)
}

func resourceServiceEndpointNexusRead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointNexus(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointNexusUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointNexus(d)
	if _, err := updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointNexusRead(d, m)
}

func resourceServiceEndpointNexusDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointNexus(d)
	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointNexus(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("NexusIqServiceConnection")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}

	return serviceEndpoint
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointNexus(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("url", *serviceEndpoint.Url)
}
