package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/model"
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
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_BITBUCKET_SERVICE_CONNECTION_USERNAME", nil),
			Description: "The bitbucket username which should be used.",
		},

		"password": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_BITBUCKET_SERVICE_CONNECTION_PASSWORD", nil),
			Description: "The bitbucket password which should be used.",
			Sensitive:   true,
		},
	})

	return r
}

func resourceServiceEndpointBitbucketCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointBitBucket(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

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
	if isServiceEndpointDeleted(d, err, serviceEndpoint, getArgs) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointBitBucket(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointBitbucketUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointBitBucket(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointBitbucketRead(d, m)
}

func resourceServiceEndpointBitbucketDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointBitBucket(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointBitBucket(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String(string(model.RepoTypeValues.Bitbucket))
	serviceEndpoint.Url = converter.String("https://api.bitbucket.org/")
	return serviceEndpoint, nil
}

func flattenServiceEndpointBitBucket(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
}
