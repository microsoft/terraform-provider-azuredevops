package serviceendpoint

import (
	"context"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointDynamicsLifecycleServices() *schema.Resource {
	r := &schema.Resource{
		CreateContext: resourceServiceEndpointDynamicsLifecycleServicesCreate,
		ReadContext:   resourceServiceEndpointDynamicsLifecycleServicesRead,
		UpdateContext: resourceServiceEndpointDynamicsLifecycleServicesUpdate,
		DeleteContext: resourceServiceEndpointDynamicsLifecycleServicesDelete,
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
		"authorization_endpoint": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"lifecycle_services_api_endpoint": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"client_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"username": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	})
	return r
}

func resourceServiceEndpointDynamicsLifecycleServicesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointDynamicsLifecycleServices(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointDynamicsLifecycleServicesRead(ctx, d, m)
}

func resourceServiceEndpointDynamicsLifecycleServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return diag.FromErr(err)
		}
		return diag.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return diag.FromErr(err)
	}
	flattenServiceEndpointDynamicsLifecycleServices(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointDynamicsLifecycleServicesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointDynamicsLifecycleServices(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return diag.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointDynamicsLifecycleServicesRead(ctx, d, m)
}

func resourceServiceEndpointDynamicsLifecycleServicesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointDynamicsLifecycleServices(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	if err = deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func expandServiceEndpointDynamicsLifecycleServices(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"clientid": d.Get("client_id").(string),
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Data = &map[string]string{
		"apiurl": d.Get("lifecycle_services_api_endpoint").(string),
	}
	serviceEndpoint.Type = converter.String("lcsserviceendpoint")
	serviceEndpoint.Url = converter.String(d.Get("authorization_endpoint").(string))
	return serviceEndpoint, nil
}

func flattenServiceEndpointDynamicsLifecycleServices(d *schema.ResourceData, endpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, endpoint)
	if endpoint.Data != nil {
		if v, ok := (*endpoint.Data)["apiurl"]; ok {
			d.Set("lifecycle_services_api_endpoint", v)
		}
	}

	if endpoint.Url != nil {
		d.Set("authorization_endpoint", *endpoint.Url)
	}

	if endpoint.Authorization != nil && endpoint.Authorization.Parameters != nil {
		params := *endpoint.Authorization.Parameters
		if v, ok := params["clientid"]; ok {
			d.Set("client_id", v)
		}

		if v, ok := params["username"]; ok {
			d.Set("username", v)
		}
	}
}
