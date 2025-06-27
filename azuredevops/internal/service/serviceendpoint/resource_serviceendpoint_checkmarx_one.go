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

func ResourceServiceEndpointCheckMarxOneService() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointCheckMarxOneServiceCreate,
		Read:   resourceServiceEndpointCheckMarxOneServiceRead,
		Update: resourceServiceEndpointCheckMarxOneServiceUpdate,
		Delete: resourceServiceEndpointCheckMarxOneServiceDelete,
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

		"api_key": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotWhiteSpace,
			ConflictsWith: []string{"client_id", "client_secret", "authorization_url"},
			AtLeastOneOf:  []string{"client_id", "api_key"},
		},

		"client_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotWhiteSpace,
			ConflictsWith: []string{"api_key"},
			RequiredWith:  []string{"client_secret"},
		},

		"client_secret": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotWhiteSpace,
			ConflictsWith: []string{"api_key"},
			RequiredWith:  []string{"client_id"},
		},

		"authorization_url": {
			Type:          schema.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsURLWithHTTPorHTTPS,
			ConflictsWith: []string{"api_key"},
		},
	})
	return r
}

func resourceServiceEndpointCheckMarxOneServiceCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointCheckMarxOneService(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointCheckMarxOneServiceRead(d, m)
}

func resourceServiceEndpointCheckMarxOneServiceRead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointCheckMarxOneService(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointCheckMarxOneServiceUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	serviceEndpoint, err := expandServiceEndpointCheckMarxOneService(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}
	return resourceServiceEndpointCheckMarxOneServiceRead(d, m)
}

func resourceServiceEndpointCheckMarxOneServiceDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointCheckMarxOneService(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointCheckMarxOneService(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)

	if v, ok := d.GetOk("api_key"); ok {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"apitoken": v.(string),
			},
			Scheme: converter.String("Token"),
		}
	} else if _, ok := d.GetOk("client_id"); ok {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"authURL":  d.Get("authorization_url").(string),
				"username": d.Get("client_id").(string),
				"password": d.Get("client_secret").(string),
			},
			Scheme: converter.String("UsernamePassword"),
		}
	}
	serviceEndpoint.Type = converter.String("CheckmarxASTService")
	serviceEndpoint.Url = converter.String(d.Get("server_url").(string))
	return serviceEndpoint, nil
}

func flattenServiceEndpointCheckMarxOneService(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("server_url", *serviceEndpoint.Url)

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		if *serviceEndpoint.Authorization.Scheme == "" {
			// TOKEN will not return
		} else if *serviceEndpoint.Authorization.Scheme == "UsernamePassword" {
			if v, ok := (*serviceEndpoint.Authorization.Parameters)["authURL"]; ok {
				d.Set("authorization_url", v)
			}

			if v, ok := (*serviceEndpoint.Authorization.Parameters)["username"]; ok {
				d.Set("client_id", v)
			}
		}
	}
}
