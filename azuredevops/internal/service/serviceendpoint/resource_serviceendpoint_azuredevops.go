package serviceendpoint

import (
	"fmt"
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

func ResourceServiceEndpointAzureDevOps() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointAzureDevOpsCreate,
		Read:   resourceServiceEndpointAzureDevOpsRead,
		Update: resourceServiceEndpointAzureDevOpsUpdate,
		Delete: resourceServiceEndpointAzureDevOpsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}
	r.DeprecationMessage = "This resource is duplicate with azuredevops_serviceendpoint_runpipeline,  will be removed in the future, use azuredevops_serviceendpoint_runpipeline instead."

	r.Schema["org_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_DEVOPS_ORG_URL", "https://dev.azure.com/[organization]"),
		Description:  "The Organization Url.",
	}

	r.Schema["release_api_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_DEVOPS_RELEASE_API_URL", "https://vsrm.dev.azure.com/[organization]"),
	}

	r.Schema["personal_access_token"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_DEVOPS_PAT", nil),
		Description: "The Azure DevOps personal access token.",
	}

	return r
}

func resourceServiceEndpointAzureDevOpsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointAzureDevOps(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint111(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointAzureDevOpsRead(d, m)
}

func resourceServiceEndpointAzureDevOpsRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServiceEndpointAzureDevOps(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	return nil
}

func resourceServiceEndpointAzureDevOpsUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointAzureDevOps(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointAzureDevOps(d, updatedServiceEndpoint, projectID)
	return resourceServiceEndpointAzureDevOpsRead(d, m)
}

func resourceServiceEndpointAzureDevOpsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointAzureDevOps(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointAzureDevOps(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("personal_access_token").(string),
		},
		Scheme: converter.String("Token"),
	}
	serviceEndpoint.Type = converter.String("AZDOAPI")
	serviceEndpoint.Url = converter.String(d.Get("org_url").(string))
	serviceEndpoint.Data = &map[string]string{
		"releaseUrl": d.Get("release_api_url").(string),
	}
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointAzureDevOps(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("org_url", serviceEndpoint.Url)
	d.Set("release_api_url", (*serviceEndpoint.Data)["releaseUrl"])
}
