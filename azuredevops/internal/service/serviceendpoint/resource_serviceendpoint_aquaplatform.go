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

// ResourceServiceEndpointAquaPlatform schema and implementation for Aqua Platform service endpoint resource
func ResourceServiceEndpointAquaPlatform() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointAquaPlatformCreate,
		Read:   resourceServiceEndpointAquaPlatformRead,
		Update: resourceServiceEndpointAquaPlatformUpdate,
		Delete: resourceServiceEndpointAquaPlatformDelete,
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
		"aqua_platform_url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  "The URL of the Aqua Platform.",
		},
		"aqua_auth_url": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "https://api.cloudsploit.com",
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  "The URL used for authentication.",
		},
		"aqua_key": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
			Description:  "The API key for the Aqua Platform.",
		},
		"aqua_secret": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
			Description:  "The API secret for the Aqua Platform.",
		},
	})

	return r
}

func resourceServiceEndpointAquaPlatformCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointAquaPlatform(d)
	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointAquaPlatformRead(d, m)
}

func resourceServiceEndpointAquaPlatformRead(d *schema.ResourceData, m interface{}) error {
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

	if serviceEndpoint.Id == nil {
		d.SetId("")
		return nil
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointAquaPlatform(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointAquaPlatformUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointAquaPlatform(d)
	if _, err := updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointAquaPlatformRead(d, m)
}

func resourceServiceEndpointAquaPlatformDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointAquaPlatform(d)
	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAquaPlatform(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("None"),
		Parameters: &map[string]string{
			"aquaKey":    d.Get("aqua_key").(string),
			"aquaSecret": d.Get("aqua_secret").(string),
		},
	}
	serviceEndpoint.Data = &map[string]string{
		"aquaPlatformUrl": d.Get("aqua_platform_url").(string),
		"aquaAuthUrl":     d.Get("aqua_auth_url").(string),
	}
	serviceEndpoint.Type = converter.String("aquaplatform")
	serviceEndpoint.Url = converter.String(d.Get("aqua_platform_url").(string))
	return serviceEndpoint
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAquaPlatform(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)

	if serviceEndpoint.Data != nil {
		d.Set("aqua_platform_url", (*serviceEndpoint.Data)["aquaPlatformUrl"])
		d.Set("aqua_auth_url", (*serviceEndpoint.Data)["aquaAuthUrl"])
	}
}
