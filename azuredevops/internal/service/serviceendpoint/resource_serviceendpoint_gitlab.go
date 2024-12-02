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

func ResourceServiceEndpointGitLab() *schema.Resource {
	r := &schema.Resource{
		CreateContext: resourceServiceEndpointGitLabCreate,
		ReadContext:   resourceServiceEndpointGitLabRead,
		UpdateContext: resourceServiceEndpointGitLabUpdate,
		DeleteContext: resourceServiceEndpointGitLabDelete,
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
		"username": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"api_token": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	})
	return r
}

func resourceServiceEndpointGitLabCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGitLab(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointGitLabRead(ctx, d, m)
}

func resourceServiceEndpointGitLabRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return diag.FromErr(err)
	}
	flattenServiceEndpointGitLab(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointGitLabUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGitLab(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return diag.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointGitLabRead(ctx, d, m)
}

func resourceServiceEndpointGitLabDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGitLab(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	if err = deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func expandServiceEndpointGitLab(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("api_token").(string),
		},
		Scheme: converter.String("Token"),
	}
	serviceEndpoint.Type = converter.String("gitlab")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	serviceEndpoint.Data = &map[string]string{
		"username": d.Get("username").(string),
	}
	return serviceEndpoint, nil
}

func flattenServiceEndpointGitLab(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	if serviceEndpoint.Data != nil {
		if v, ok := (*serviceEndpoint.Data)["username"]; ok {
			d.Set("username", v)
		}
	}

	if serviceEndpoint.Url != nil {
		d.Set("url", *serviceEndpoint.Url)
	}
}
