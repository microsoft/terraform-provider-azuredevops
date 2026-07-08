package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/model"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointBitBucket() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointBitbucketCreate,
		Read:   resourceServiceEndpointBitbucketRead,
		Update: resourceServiceEndpointBitbucketUpdate,
		Delete: resourceServiceEndpointBitbucketDelete,
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
		"username": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZDO_BITBUCKET_SERVICE_CONNECTION_USERNAME", nil),
			Description:  "The bitbucket username which should be used.",
			RequiredWith: []string{"password"},
			Deprecated:   "Bitbucket Cloud has deprecated app password (username and password) authentication. Use `email` and `api_token` instead.",
		},

		"password": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZDO_BITBUCKET_SERVICE_CONNECTION_PASSWORD", nil),
			Description:  "The bitbucket password which should be used.",
			Sensitive:    true,
			RequiredWith: []string{"username"},
			Deprecated:   "Bitbucket Cloud has deprecated app password (username and password) authentication. Use `email` and `api_token` instead.",
		},

		"email": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZDO_BITBUCKET_SERVICE_CONNECTION_EMAIL", nil),
			Description:  "The bitbucket account email which should be used.",
			RequiredWith: []string{"api_token"},
			ExactlyOneOf: []string{"username", "email"},
		},

		"api_token": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZDO_BITBUCKET_SERVICE_CONNECTION_API_TOKEN", nil),
			Description:  "The bitbucket API token which should be used.",
			Sensitive:    true,
			RequiredWith: []string{"email"},
		},
	})

	return r
}

func resourceServiceEndpointBitbucketCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointBitBucket(d)
	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointBitbucketRead(d, m)
}

func resourceServiceEndpointBitbucketRead(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}
	if serviceEndpoint == nil || serviceEndpoint.Id == nil {
		return fmt.Errorf("unexpected nil service endpoint, ID: (%v), project ID: (%v)", getArgs.EndpointId, getArgs.Project)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointBitBucket(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointBitbucketUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointBitBucket(d)
	if _, err := updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointBitbucketRead(d, m)
}

func resourceServiceEndpointBitbucketDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointBitBucket(d)
	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointBitBucket(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := doBaseExpansion(d)
	if apiToken, ok := d.GetOk("api_token"); ok {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"email":    d.Get("email").(string),
				"apitoken": apiToken.(string),
			},
			Scheme: converter.String("Token"),
		}
	} else {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"username": d.Get("username").(string),
				"password": d.Get("password").(string),
			},
			Scheme: converter.String("UsernamePassword"),
		}
	}
	serviceEndpoint.Type = converter.String(string(model.RepoTypeValues.Bitbucket))
	serviceEndpoint.Url = converter.String("https://api.bitbucket.org/")
	return serviceEndpoint
}

func flattenServiceEndpointBitBucket(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	if serviceEndpoint.Authorization == nil || serviceEndpoint.Authorization.Scheme == nil || serviceEndpoint.Authorization.Parameters == nil {
		return
	}
	params := *serviceEndpoint.Authorization.Parameters
	switch *serviceEndpoint.Authorization.Scheme {
	case "Token":
		d.Set("email", params["email"])
	case "UsernamePassword":
		d.Set("username", params["username"])
	}
}
