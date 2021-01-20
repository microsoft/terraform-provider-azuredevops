package serviceendpoint

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const errMsgTfConfigRead = "Error reading terraform configuration: %+v"

type flatFunc func(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string)
type expandFunc func(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error)

type operationState struct {
	Ready      string
	Failed     string
	InProgress string
}

var opState = operationState{
	Ready:      "Ready",
	Failed:     "Failed",
	InProgress: "InProgress",
}

// genBaseServiceEndpointResource creates a Resource with the common parts
// that all Service Endpoints require.
func genBaseServiceEndpointResource(f flatFunc, e expandFunc) *schema.Resource {
	return &schema.Resource{
		Create: genServiceEndpointCreateFunc(f, e),
		Read:   genServiceEndpointReadFunc(f),
		Update: genServiceEndpointUpdateFunc(f, e),
		Delete: genServiceEndpointDeleteFunc(e),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"service_endpoint_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"authorization": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// doBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func doBaseExpansion(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
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

// doBaseFlattening performs the flattening for the 'base' attributes that are defined in the schema, above
func doBaseFlattening(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	d.SetId(serviceEndpoint.Id.String())
	d.Set("service_endpoint_name", serviceEndpoint.Name)
	d.Set("project_id", projectID)
	d.Set("description", serviceEndpoint.Description)
	d.Set("authorization", &map[string]interface{}{
		"scheme": *serviceEndpoint.Authorization.Scheme,
	})
}

// makeProtectedSchema create protected schema
func makeProtectedSchema(r *schema.Resource, keyName, envVarName, description string) {
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

// makeUnprotectedSchema create unprotected schema
func makeUnprotectedSchema(r *schema.Resource, keyName, envVarName, description string) {
	r.Schema[keyName] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc(envVarName, nil),
		Description: description,
	}
}

// Make the Azure DevOps API call to create the endpoint
func createServiceEndpoint(clients *client.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
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

func deleteServiceEndpoint(clients *client.AggregatedClient, project *string, endPointID *uuid.UUID) error {
	err := clients.ServiceEndpointClient.DeleteServiceEndpoint(
		clients.Ctx,
		serviceendpoint.DeleteServiceEndpointArgs{
			Project:    project,
			EndpointId: endPointID,
		})

	return err
}

func updateServiceEndpoint(clients *client.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
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
		clients := m.(*client.AggregatedClient)
		serviceEndpoint, projectID, err := expandFunc(d)
		if err != nil {
			return fmt.Errorf(errMsgTfConfigRead, err)
		}

		createdServiceEndpoint, err := createServiceEndpoint(clients, serviceEndpoint, projectID)
		if err != nil {
			return fmt.Errorf("Error creating service endpoint in Azure DevOps: %+v", err)
		}

		log.Printf("[DEBUG] Waiting service endpoint ready")
		stateConf := &resource.StateChangeConf{
			ContinuousTargetOccurence: 1,
			Delay:                     10 * time.Second,
			MinTimeout:                10 * time.Second,
			Pending:                   []string{opState.Failed},
			Target:                    []string{opState.Ready},
			Refresh:                   getServiceEndpoint(clients, createdServiceEndpoint.Id, projectID),
			Timeout:                   d.Timeout(schema.TimeoutCreate),
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf(" waiting for service endpoint ready. %v ", err)
		}

		flatFunc(d, createdServiceEndpoint, projectID)
		return genServiceEndpointReadFunc(flatFunc)(d, m)
	}
}

func genServiceEndpointReadFunc(flatFunc flatFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)

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
		clients := m.(*client.AggregatedClient)
		serviceEndpoint, projectID, err := expandFunc(d)
		if err != nil {
			return fmt.Errorf(errMsgTfConfigRead, err)
		}

		updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint, projectID)
		if err != nil {
			return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
		}

		flatFunc(d, updatedServiceEndpoint, projectID)
		return genServiceEndpointReadFunc(flatFunc)(d, m)
	}
}

func genServiceEndpointDeleteFunc(expandFunc expandFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		serviceEndpoint, projectID, err := expandFunc(d)
		if err != nil {
			return fmt.Errorf(errMsgTfConfigRead, err)
		}

		return deleteServiceEndpoint(clients, projectID, serviceEndpoint.Id)
	}
}

func getServiceEndpoint(client *client.AggregatedClient, serviceEndpointID *uuid.UUID, projectID *string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		serviceEndpoint, err := client.ServiceEndpointClient.GetServiceEndpointDetails(
			client.Ctx,
			serviceendpoint.GetServiceEndpointDetailsArgs{
				EndpointId: serviceEndpointID,
				Project:    projectID,
			},
		)

		if err != nil {
			return nil, opState.Failed, fmt.Errorf("Error looking up service endpoint given ID (%v) and project ID (%v): %v ", serviceEndpointID, projectID, err)
		}

		if *serviceEndpoint.IsReady {
			return serviceEndpoint, opState.Ready, nil
		} else if serviceEndpoint.OperationStatus != nil {
			opStatus := (serviceEndpoint.OperationStatus).(map[string]interface{})
			if opStatus["state"] == opState.Failed {
				return nil, opState.Failed, fmt.Errorf("Error looking up service endpoint given ID (%v) and project ID (%v): %v ", serviceEndpointID, projectID, serviceEndpoint.OperationStatus)
			}
		}
		return nil, opState.Failed, nil
	}
}
