package crudserviceendpoint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

type flatFunc func(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string)
type expandFunc func(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string)
type importFunc func(clients *config.AggregatedClient, id string) (string, string, error)

//GenBaseServiceEndpointResource creates a Resource with the common parts
// that all Service Endpoints require.
func GenBaseServiceEndpointResource(f flatFunc, e expandFunc, i importFunc) *schema.Resource {
	return &schema.Resource{
		Create: genServiceEndpointCreateFunc(f, e),
		Read:   genServiceEndpointReadFunc(f),
		Update: genServiceEndpointUpdateFunc(f, e),
		Delete: genServiceEndpointDeleteFunc(e),
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, serviceEndpointID, err := i(meta.(*config.AggregatedClient), d.Id())
				if err != nil {
					return nil, fmt.Errorf("Error parsing the variable service endpoint ID from the Terraform resource data:  %v", err)
				}
				d.Set("project_id", projectID)
				d.SetId(serviceEndpointID)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: genBaseSchema(),
	}
}

func genBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"service_endpoint_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.NoEmptyStrings,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "Managed by Terraform",
		},
		"authorization": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// DoBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func DoBaseExpansion(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var serviceEndpointID *uuid.UUID
	parsedID, err := uuid.Parse(d.Id())
	if err == nil {
		serviceEndpointID = &parsedID
	}
	projectID := converter.String(d.Get("project_id").(string))
	serviceEndpoint := &serviceendpoint.ServiceEndpoint{
		Id:          serviceEndpointID,
		Name:        converter.String(d.Get("service_endpoint_name").(string)),
		Owner:       converter.String("library"),
		Description: converter.String(d.Get("description").(string)),
	}

	return serviceEndpoint, projectID
}

// DoBaseFlattening performs the flattening for the 'base' attributes that are defined in the schema, above
func DoBaseFlattening(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	d.SetId(serviceEndpoint.Id.String())
	d.Set("service_endpoint_name", *serviceEndpoint.Name)
	d.Set("project_id", projectID)
	d.Set("description", *serviceEndpoint.Description)
	d.Set("authorization", &map[string]interface{}{
		"scheme": *serviceEndpoint.Authorization.Scheme,
	})
}

// GetScheme allows you to get the nested scheme value
func GetScheme(d *schema.ResourceData) (string, error) {
	authorization := d.Get("authorization").(*schema.Set)
	if authorization == nil {
		return "", errors.New("authorization not set")
	}
	authorizationList := authorization.List()
	if len(authorizationList) != 1 {
		return "", errors.New("authorization is invalid")
	}
	scheme := authorizationList[0].(map[string]interface{})["scheme"].(string)
	return scheme, nil
}

// MakeProtectedSchema create protected schema
func MakeProtectedSchema(r *schema.Resource, keyName, envVarName, description string) {
	r.Schema[keyName] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		DefaultFunc:      schema.EnvDefaultFunc(envVarName, nil),
		Description:      description,
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
	}

	secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema(keyName)
	r.Schema[secretHashKey] = secretHashSchema
}

// MakeUnprotectedSchema create unprotected schema
func MakeUnprotectedSchema(r *schema.Resource, keyName, envVarName, description string) {
	r.Schema[keyName] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc(envVarName, nil),
		Description: description,
	}
}

// Make the Azure DevOps API call to create the endpoint
func createServiceEndpoint(clients *config.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
	if strings.EqualFold(*endpoint.Type, "github") && strings.EqualFold(*endpoint.Authorization.Scheme, "InstallationToken") {
		return nil, fmt.Errorf("Github Apps must be created on Github and then can be imported")
	}
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
	if strings.EqualFold(*endpoint.Type, "github") && strings.EqualFold(*endpoint.Authorization.Scheme, "InstallationToken") {
		return nil, fmt.Errorf("Github Apps can not be updated must match imported values exactly")
	}
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
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error looking up service endpoint given ID (%v) and project ID (%v): %v", serviceEndpointID, projectID, err)
		}

		if serviceEndpoint.Id == nil {
			// e.g. service endpoint has been deleted separately without TF
			d.SetId("")
		} else {
			flatFunc(d, serviceEndpoint, projectID)
		}
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
