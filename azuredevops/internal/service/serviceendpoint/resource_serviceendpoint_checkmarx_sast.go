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

func ResourceServiceEndpointCheckMarxSAST() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointCheckMarxSASTCreate,
		Read:   resourceServiceEndpointCheckMarxSASTRead,
		Update: resourceServiceEndpointCheckMarxSASTUpdate,
		Delete: resourceServiceEndpointCheckMarxSASTDelete,
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
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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

		"preset": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
	})
	return r
}

func resourceServiceEndpointCheckMarxSASTCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointCheckMarxSAST(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointCheckMarxSASTRead(d, m)
}

func resourceServiceEndpointCheckMarxSASTRead(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf(" looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointCheckMarxSAST(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointCheckMarxSASTUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	serviceEndpoint, err := expandServiceEndpointCheckMarxSAST(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}
	return resourceServiceEndpointCheckMarxSASTRead(d, m)
}

func resourceServiceEndpointCheckMarxSASTDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointCheckMarxSAST(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointCheckMarxSAST(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
			"preset":   d.Get("preset").(string),
			"teams":    d.Get("team").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("Checkmarx-Endpoint")
	serviceEndpoint.Url = converter.String(d.Get("server_url").(string))
	return serviceEndpoint, nil
}

func flattenServiceEndpointCheckMarxSAST(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("server_url", *serviceEndpoint.Url)

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		if v, ok := (*serviceEndpoint.Authorization.Parameters)["username"]; ok {
			d.Set("username", v)
		}

		if v, ok := (*serviceEndpoint.Authorization.Parameters)["preset"]; ok {
			d.Set("preset", v)
		}

		if v, ok := (*serviceEndpoint.Authorization.Parameters)["teams"]; ok {
			d.Set("team", v)
		}
	}
}
