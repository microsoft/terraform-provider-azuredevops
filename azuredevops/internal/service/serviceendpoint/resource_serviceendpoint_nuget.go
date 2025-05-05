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

func ResourceServiceEndpointNuGet() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointNuGetCreate,
		Read:   resourceServiceEndpointNuGetRead,
		Update: resourceServiceEndpointNuGetUpdate,
		Delete: resourceServiceEndpointNuGetDelete,
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
		"feed_url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"api_key": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"personal_access_token", "username", "password"},
			AtLeastOneOf:  []string{"api_key", "personal_access_token", "username", "password"},
		},

		"personal_access_token": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"api_key", "username", "password"},
		},

		"username": {
			Type:          schema.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"personal_access_token", "api_key"},
			RequiredWith:  []string{"password"},
		},

		"password": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"personal_access_token", "api_key"},
			RequiredWith:  []string{"username"},
		},
	})

	return r
}

func resourceServiceEndpointNuGetCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointNuGet(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointNuGetRead(d, m)
}

func resourceServiceEndpointNuGetRead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointNuGet(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointNuGetUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointNuGet(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointNuGetRead(d, m)
}

func resourceServiceEndpointNuGetDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointNuGet(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointNuGet(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("externalnugetfeed")
	serviceEndpoint.Url = converter.String(d.Get("feed_url").(string))
	if apiKey := d.Get("api_key"); apiKey != "" {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"nugetkey": apiKey.(string),
			},
			Scheme: converter.String("None"),
		}
	}

	if pat := d.Get("personal_access_token"); pat != "" {
		serviceEndpoint.Type = converter.String("externalnugetfeed")
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"apitoken": pat.(string),
			},
			Scheme: converter.String("Token"),
		}
	}

	if uname := d.Get("username"); uname != "" {
		serviceEndpoint.Type = converter.String("externalnugetfeed")
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"username": uname.(string),
				"password": d.Get("password").(string),
			},
			Scheme: converter.String("UsernamePassword"),
		}
	}
	return serviceEndpoint, nil
}

func flattenServiceEndpointNuGet(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("feed_url", *serviceEndpoint.Url)

	switch *serviceEndpoint.Authorization.Scheme {
	case "None":
		d.Set("api_key", d.Get("api_key"))
	case "Token":
		d.Set("personal_access_token", d.Get("personal_access_token"))
	case "UsernamePassword":
		d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
		d.Set("password", d.Get("password"))
	}
}
