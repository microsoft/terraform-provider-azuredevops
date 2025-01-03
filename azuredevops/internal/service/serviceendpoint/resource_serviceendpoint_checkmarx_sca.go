package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointCheckMarxSCA() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointCheckMarxSCACreate,
		Read:   resourceServiceEndpointCheckMarxSCARead,
		Update: resourceServiceEndpointCheckMarxSCAUpdate,
		Delete: resourceServiceEndpointCheckMarxSCADelete,
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
		"access_control_url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"server_url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"web_app_url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"account": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"username": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"team": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
	})
	return r
}

func resourceServiceEndpointCheckMarxSCACreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointCheckMarxSCA(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointCheckMarxSCARead(d, m)
}

func resourceServiceEndpointCheckMarxSCARead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointCheckMarxSCA(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointCheckMarxSCAUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	serviceEndpoint, err := expandServiceEndpointCheckMarxSCA(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}
	return resourceServiceEndpointCheckMarxSCARead(d, m)
}

func resourceServiceEndpointCheckMarxSCADelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointCheckMarxSCA(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointCheckMarxSCA(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("SCA-Endpoint")
	serviceEndpoint.Url = converter.String(d.Get("server_url").(string))
	serviceEndpoint.Data = &map[string]string{
		"dependencyAccessControlURL": d.Get("access_control_url").(string),
		"dependencyTenant":           d.Get("account").(string),
		"dependencyWebAppURL":        d.Get("web_app_url").(string),
		"teams":                      d.Get("team").(string),
	}
	return serviceEndpoint, nil
}

func flattenServiceEndpointCheckMarxSCA(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("server_url", *serviceEndpoint.Url)

	if serviceEndpoint.Data != nil {
		if v, ok := (*serviceEndpoint.Data)["dependencyAccessControlURL"]; ok {
			d.Set("access_control_url", v)

		}

		if v, ok := (*serviceEndpoint.Data)["dependencyWebAppURL"]; ok {
			d.Set("web_app_url", v)

		}

		if v, ok := (*serviceEndpoint.Data)["dependencyTenant"]; ok {
			d.Set("account", v)

		}

		if v, ok := (*serviceEndpoint.Data)["teams"]; ok {
			d.Set("team", v)
		}
	}

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		if v, ok := (*serviceEndpoint.Authorization.Parameters)["username"]; ok {
			d.Set("username", v)
		}
	}
}
