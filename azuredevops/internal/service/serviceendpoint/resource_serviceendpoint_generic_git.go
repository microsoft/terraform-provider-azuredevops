package serviceendpoint

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
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
	r.Schema["repository_url"] = &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Required:     true,
		Description:  "The server URL of the GenericGit git service connection.",
	}
	r.Schema["username"] = &schema.Schema{
		Type:        schema.TypeString,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GenericGit_GIT_SERVICE_CONNECTION_USERNAME", nil),
		Description: "The username to use for the GenericGit service git connection.",
		Optional:    true,
	}
	r.Schema["password"] = &schema.Schema{
		Type:        schema.TypeString,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GenericGit_GIT_SERVICE_CONNECTION_PASSWORD", nil),
		Description: "The password or token key to use for the GenericGit git service connection.",
		Sensitive:   true,
		Optional:    true,
	}
	r.Schema["enable_pipelines_access"] = &schema.Schema{
		Type:        schema.TypeBool,
		Default:     true,
		Description: "A value indicating whether or not to attempt accessing this git server from Azure Pipelines.",
		Optional:    true,
	}
	return r
}

func resourceServiceEndpointGenericGitCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointGenericGit(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

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
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	flattenServiceEndpointGeneric(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	return nil
}

func resourceServiceEndpointGenericGitUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointGenericGit(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointGeneric(d, updatedServiceEndpoint, projectID)
	return resourceServiceEndpointGenericGitRead(d, m)
}

func resourceServiceEndpointGenericGitDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointGenericGit(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointGenericGit(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
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
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointGenericGitGit(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("repository_url", *serviceEndpoint.Url)
	if v, err := strconv.ParseBool((*serviceEndpoint.Data)["accessExternalGitServer"]); err != nil {
		d.Set("enable_pipelines_access", v)
	}
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
}
