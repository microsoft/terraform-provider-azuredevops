package serviceendpoint

import (
	"fmt"
	"maps"
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

// ResourceServiceEndpointRunPipeline schema and implementation for Azure DevOps service endpoint resource
func ResourceServiceEndpointRunPipeline() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointRunPipelineCreate,
		Read:   resourceServiceEndpointRunPipelineRead,
		Update: resourceServiceEndpointRunPipelineUpdate,
		Delete: resourceServiceEndpointRunPipelineDelete,
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
		"organization_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Azure DevOps organization name",
		},

		"auth_personal": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"personal_access_token": {
						Type:         schema.TypeString,
						Required:     true,
						DefaultFunc:  schema.EnvDefaultFunc("AZDO_PERSONAL_ACCESS_TOKEN", nil),
						Description:  "The Azure DevOps personal access token which should be used.",
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
				},
			},
		},
	})

	return r
}

func resourceServiceEndpointRunPipelineCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointRunPipeline(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointRunPipelineRead(d, m)
}

func resourceServiceEndpointRunPipelineRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServiceEndpointRunPipeline(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	return nil
}

func resourceServiceEndpointRunPipelineUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointRunPipeline(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointRunPipeline(d, updatedServiceEndpoint, projectID.String())
	return resourceServiceEndpointRunPipelineRead(d, m)
}

func resourceServiceEndpointRunPipelineDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointRunPipeline(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure:
func expandServiceEndpointRunPipeline(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("azdoapi")

	scheme := "Token"
	parameters := rpExpandAuthPersonalSet(d.Get("auth_personal").(*schema.Set))

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}

	org := d.Get("organization_name").(string)
	serviceUrl := fmt.Sprint("https://dev.azure.com/", org)
	serviceEndpoint.Url = &serviceUrl

	data := map[string]string{}
	releaseUrl := fmt.Sprint("https://vsrm.dev.azure.com/", org)
	data["releaseUrl"] = releaseUrl
	serviceEndpoint.Data = &data
	return serviceEndpoint, projectID, nil
}

func rpExpandAuthPersonalSet(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	if len(d.List()) == 1 {
		val := d.List()[0].(map[string]interface{}) //auth_personal block may have only one element inside
		authPerson["apitoken"] = val["personal_access_token"].(string)
	}
	return authPerson
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointRunPipeline(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			d.Set("auth_personal", []interface{}{authPersonal})
		}
	}
}
