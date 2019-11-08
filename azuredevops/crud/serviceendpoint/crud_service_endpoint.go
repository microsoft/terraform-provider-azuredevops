package crudserviceendpoint

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

type flatFunc func(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string)
type expandFunc func(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string)

//GenBaseServiceEndpointResource creates a Resource with the common parts
// that all Service Endpoints require.
func GenBaseServiceEndpointResource(f flatFunc, e expandFunc) *schema.Resource {
	return &schema.Resource{
		Create: genServiceEndpointCreateFunc(f, e),
		Read:   genServiceEndpointReadFunc(f),
		Update: genServiceEndpointUpdateFunc(f, e),
		Delete: genServiceEndpointDeleteFunc(e),
		Schema: genBaseScema(),
	}
}

func genBaseScema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"service_endpoint_name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

//DoBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func DoBaseExpansion(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var serviceEndpointID *uuid.UUID
	parsedID, err := uuid.Parse(d.Id())
	if err == nil {
		serviceEndpointID = &parsedID
	}
	projectID := converter.String(d.Get("project_id").(string))
	serviceEndpoint := &serviceendpoint.ServiceEndpoint{
		Id:    serviceEndpointID,
		Name:  converter.String(d.Get("service_endpoint_name").(string)),
		Owner: converter.String("library"),
	}

	return serviceEndpoint, projectID
}

//DoBaseFlattening performs the flattening for the 'base' attributes that are defined in the schema, above
func DoBaseFlattening(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	d.SetId(serviceEndpoint.Id.String())
	d.Set("service_endpoint_name", *serviceEndpoint.Name)
	d.Set("project_id", projectID)
}

// Make the Azure DevOps API call to create the endpoint
func createServiceEndpoint(clients *config.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
	createdServiceEndpoint, err := clients.ServiceEndpointClient.CreateServiceEndpoint(
		clients.Ctx,
		serviceendpoint.CreateServiceEndpointArgs{
			Endpoint: endpoint,
			Project:  project,
		})

	return createdServiceEndpoint, err
}

func deleteServiceEndpoint(clients *config.AggregatedClient, project *string, endPointID *uuid.UUID) error {
	err := clients.ServiceEndpointClient.DeleteServiceEndpoint(
		clients.Ctx,
		serviceendpoint.DeleteServiceEndpointArgs{
			Project:    project,
			EndpointId: endPointID,
		})

	return err
}

func updateServiceEndpoint(clients *config.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
	updatedServiceEndpoint, err := clients.ServiceEndpointClient.UpdateServiceEndpoint(
		clients.Ctx,
		serviceendpoint.UpdateServiceEndpointArgs{
			Endpoint:   endpoint,
			Project:    project,
			EndpointId: endpoint.Id,
		})

	return updatedServiceEndpoint, err
}

func genServiceEndpointCreateFunc(flatFunc flatFunc, expandFunc expandFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*config.AggregatedClient)
		serviceEndpoint, projectID := expandFunc(d)

		createdServiceEndpoint, err := createServiceEndpoint(clients, serviceEndpoint, projectID)
		if err != nil {
			return fmt.Errorf("Error creating service endpoint in Azure DevOps: %+v", err)
		}

		flatFunc(d, createdServiceEndpoint, projectID)
		return nil
	}
}

func genServiceEndpointReadFunc(flatFunc flatFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*config.AggregatedClient)

		var serviceEndpointID *uuid.UUID
		parsedServiceEndpointID, err := uuid.Parse(d.Id())
		if err != nil {
			return fmt.Errorf("Error parsing the service endpoint ID from the Terraform resource data: %v", err)
		}
		serviceEndpointID = &parsedServiceEndpointID
		projectID := converter.String(d.Get("project_id").(string))

		serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
			clients.Ctx,
			serviceendpoint.GetServiceEndpointDetailsArgs{
				EndpointId: serviceEndpointID,
				Project:    projectID,
			},
		)
		if err != nil {
			return fmt.Errorf("Error looking up service endpoint given ID (%v) and project ID (%v): %v", serviceEndpointID, projectID, err)
		}

		flatFunc(d, serviceEndpoint, projectID)
		return nil
	}
}

func genServiceEndpointUpdateFunc(flatFunc flatFunc, expandFunc expandFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*config.AggregatedClient)
		serviceEndpoint, projectID := expandFunc(d)

		updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint, projectID)
		if err != nil {
			return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
		}

		flatFunc(d, updatedServiceEndpoint, projectID)
		return nil
	}
}

func genServiceEndpointDeleteFunc(expandFunc expandFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*config.AggregatedClient)
		serviceEndpoint, projectID := expandFunc(d)

		return deleteServiceEndpoint(clients, projectID, serviceEndpoint.Id)
	}
}
