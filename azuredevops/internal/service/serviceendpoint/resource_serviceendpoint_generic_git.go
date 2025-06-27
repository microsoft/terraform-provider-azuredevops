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

func ResourceServiceEndpointGenericGit() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointGenericGitCreate,
		Read:   resourceServiceEndpointGenericGitRead,
		Update: resourceServiceEndpointGenericGitUpdate,
		Delete: resourceServiceEndpointGenericGitDelete,
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
		"repository_url": {
			Type:         schema.TypeString,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Required:     true,
			Description:  "The server URL of the GenericGit git service connection.",
		},
		"username": {
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_GIT_SERVICE_CONNECTION_USERNAME", nil),
			Description: "The username to use for the GenericGit service git connection.",
			Optional:    true,
		},
		"password": {
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_GIT_SERVICE_CONNECTION_PASSWORD", nil),
			Description: "The password or token key to use for the GenericGit git service connection.",
			Sensitive:   true,
			Optional:    true,
		},
		"enable_pipelines_access": {
			Type:        schema.TypeBool,
			Default:     true,
			Description: "A value indicating whether or not to attempt accessing this git server from Azure Pipelines.",
			Optional:    true,
		},
	})

	return r
}

func resourceServiceEndpointGenericGitCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointGenericGit(d)
	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointGenericGitRead(d, m)
}

func resourceServiceEndpointGenericGitRead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointGenericGit(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointGenericGitUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointGenericGit(d)
	if _, err := updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Upating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointGenericGitRead(d, m)
}

func resourceServiceEndpointGenericGitDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointGenericGit(d)
	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointGenericGit(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("git")
	serviceEndpoint.Url = converter.String(d.Get("repository_url").(string))
	serviceEndpoint.Data = &map[string]string{
		"accessExternalGitServer": strconv.FormatBool(d.Get("enable_pipelines_access").(bool)),
	}
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	return serviceEndpoint
}

func flattenServiceEndpointGenericGit(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("repository_url", *serviceEndpoint.Url)
	if v, err := strconv.ParseBool((*serviceEndpoint.Data)["accessExternalGitServer"]); err != nil {
		d.Set("enable_pipelines_access", v)
	}
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
}
