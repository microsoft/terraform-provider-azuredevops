package azuredevops

import (
	"fmt"
	"log"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
)

func resourceServiceEndpoint() *schema.Resource {

	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema("github_service_endpoint_pat")

	return &schema.Resource{
		Create: resourceServiceEndpointCreate,
		Read:   resourceServiceEndpointRead,
		Update: resourceServiceEndpointUpdate,
		Delete: resourceServiceEndpointDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_endpoint_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_endpoint_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_endpoint_owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"github_service_endpoint_pat": {
				Type:             schema.TypeString,
				Required:         true,
				DefaultFunc:      schema.EnvDefaultFunc("AZDO_GITHUB_SERVICE_CONNECTION_PAT", nil),
				Description:      "The GitHub personal access token which should be used.",
				Sensitive:        true,
				DiffSuppressFunc: tfhelper.DiffFuncSupressSecretChanged,
			},
			patHashKey: patHashSchema,
		},
	}
}

func resourceServiceEndpointCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	serviceEndpoint, projectID := expandServiceEndpoint(d)

	createdServiceEndpoint, err := createServiceEndpoint(clients, serviceEndpoint, projectID)
	if err != nil {
		return fmt.Errorf("Error creating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpoint(d, createdServiceEndpoint, projectID)
	return nil
}

func resourceServiceEndpointRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	var serviceEndpointID *uuid.UUID
	parsedServiceEndpointID, err := uuid.Parse(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing the service endpoint ID from the Terraform resource data: %v", err)
	}
	serviceEndpointID = &parsedServiceEndpointID
	projectID := converter.String(d.Get("project_id").(string))

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
		clients.ctx,
		serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: serviceEndpointID,
			Project:    projectID,
		},
	)
	if err != nil {
		return fmt.Errorf("Error looking up service endpoint given ID (%v) and project ID (%v): %v", serviceEndpointID, projectID, err)
	}

	flattenServiceEndpoint(d, serviceEndpoint, projectID)
	return nil
}

func resourceServiceEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	serviceEndpoint, projectID := expandServiceEndpoint(d)

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint, projectID)
	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpoint(d, updatedServiceEndpoint, projectID)
	return nil
}

func resourceServiceEndpointDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	serviceEndpoint, projectID := expandServiceEndpoint(d)

	return deleteServiceEndpoint(clients, projectID, serviceEndpoint.Id)
}

// Make the Azure DevOps API call to create the endpoint
func createServiceEndpoint(clients *aggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
	createdServiceEndpoint, err := clients.ServiceEndpointClient.CreateServiceEndpoint(
		clients.ctx,
		serviceendpoint.CreateServiceEndpointArgs{
			Endpoint: endpoint,
			Project:  project,
		})

	return createdServiceEndpoint, err
}

func deleteServiceEndpoint(clients *aggregatedClient, project *string, endPointID *uuid.UUID) error {
	err := clients.ServiceEndpointClient.DeleteServiceEndpoint(
		clients.ctx,
		serviceendpoint.DeleteServiceEndpointArgs{
			Project:    project,
			EndpointId: endPointID,
		})

	return err
}

func updateServiceEndpoint(clients *aggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
	updatedServiceEndpoint, err := clients.ServiceEndpointClient.UpdateServiceEndpoint(
		clients.ctx,
		serviceendpoint.UpdateServiceEndpointArgs{
			Endpoint:   endpoint,
			Project:    project,
			EndpointId: endpoint.Id,
		})

	return updatedServiceEndpoint, err
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpoint(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var serviceEndpointID *uuid.UUID
	parsedID, err := uuid.Parse(d.Id())
	if err == nil {
		serviceEndpointID = &parsedID
	}
	log.Printf("Updating github_service_endpoint_pat to %s", d.Get("github_service_endpoint_pat").(string))
	projectID := converter.String(d.Get("project_id").(string))
	serviceEndpoint := &serviceendpoint.ServiceEndpoint{
		Id:    serviceEndpointID,
		Name:  converter.String(d.Get("service_endpoint_name").(string)),
		Type:  converter.String(d.Get("service_endpoint_type").(string)),
		Url:   converter.String(d.Get("service_endpoint_url").(string)),
		Owner: converter.String(d.Get("service_endpoint_owner").(string)),
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"accessToken": d.Get("github_service_endpoint_pat").(string),
			},
			Scheme: converter.String("PersonalAccessToken"),
		},
	}

	return serviceEndpoint, projectID
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpoint(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	d.SetId(serviceEndpoint.Id.String())
	d.Set("service_endpoint_name", *serviceEndpoint.Name)
	d.Set("service_endpoint_type", *serviceEndpoint.Type)
	d.Set("service_endpoint_url", *serviceEndpoint.Url)
	d.Set("service_endpoint_owner", *serviceEndpoint.Owner)
	tfhelper.HelpFlattenSecret(d, "github_service_endpoint_pat")
	d.Set("github_service_endpoint_pat", (*serviceEndpoint.Authorization.Parameters)["accessToken"])
	d.Set("project_id", projectID)
}
