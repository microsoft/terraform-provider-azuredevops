package serviceendpoint

import (
	"fmt"
	"maps"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointOctopusDeploy() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointOctopusDeployCreate,
		Read:   resourceServiceEndpointOctopusDeployRead,
		Update: resourceServiceEndpointOctopusDeployUpdate,
		Delete: resourceServiceEndpointOctopusDeployDelete,
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
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"api_key": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"ignore_ssl_error": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	})

	return r
}

func resourceServiceEndpointOctopusDeployCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointOctopusDeploy(d)
	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointOctopusDeployRead(d, m)
}

func resourceServiceEndpointOctopusDeployRead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointOctopusDeploy(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointOctopusDeployUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointOctopusDeploy(d)
	if _, err := updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointOctopusDeployRead(d, m)
}

func resourceServiceEndpointOctopusDeployDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointOctopusDeploy(d)
	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointOctopusDeploy(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("api_key").(string),
		},
		Scheme: converter.String("Token"),
	}

	serviceEndpoint.Data = &map[string]string{
		"ignoreSslErrors": strconv.FormatBool(d.Get("ignore_ssl_error").(bool)),
	}
	serviceEndpoint.Type = converter.String("OctopusEndpoint")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint
}

func flattenServiceEndpointOctopusDeploy(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("url", *serviceEndpoint.Url)

	if v, ok := (*serviceEndpoint.Data)["ignoreSslErrors"]; ok && v != "" {
		ignoreSslErrors, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Errorf("Failed to parse OctopusDeploy.ignore_ssl_error.(Project: %s), (service endpoint:%s) ,Error: %+v",
				*serviceEndpoint.Name, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id, err))
		}
		d.Set("ignore_ssl_error", ignoreSslErrors)
	}
}
