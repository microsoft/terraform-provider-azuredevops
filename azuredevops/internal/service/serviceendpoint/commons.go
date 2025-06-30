package serviceendpoint

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

const (
	errMsgTfConfigRead  = " Reading terraform configuration: %+v"
	errMsgServiceCreate = " Looking up service endpoint given ID (%s) and project ID (%s): %v "
	errMsgServiceDelete = " Delete service endpoint. ServiceEndpointID: %s, projectID: %s. %v "
)

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

func baseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
		"service_endpoint_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"description": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "Managed by Terraform",
			ValidateFunc: validation.StringLenBetween(0, 1024),
		},
		"authorization": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func createServiceEndpoint(d *schema.ResourceData, clients *client.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint) (*serviceendpoint.ServiceEndpoint, error) {
	if endpoint.ServiceEndpointProjectReferences == nil || len(*endpoint.ServiceEndpointProjectReferences) == 0 {
		return nil, fmt.Errorf("A ServiceEndpoint requires at least one ServiceEndpointProjectReference")
	}

	createdServiceEndpoint, err := clients.ServiceEndpointClient.CreateServiceEndpoint(
		clients.Ctx,
		serviceendpoint.CreateServiceEndpointArgs{
			Endpoint: endpoint,
		})
	if err != nil {
		return nil, fmt.Errorf("Error creating service endpoint in Azure DevOps: %+v", err)
	}

	projectID := (*endpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id

	stateConf := &retry.StateChangeConf{
		ContinuousTargetOccurence: 1,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{opState.InProgress},
		Target:                    []string{opState.Ready, opState.Failed},
		Refresh:                   getServiceEndpoint(clients, createdServiceEndpoint.Id, projectID),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
	}

	if _, err := stateConf.WaitForState(); err != nil { //nolint:staticcheck
		if delErr := deleteServiceEndpoint(clients, createdServiceEndpoint, d.Timeout(schema.TimeoutDelete)); delErr != nil {
			log.Printf("[DEBUG] Failed to delete the failed service endpoint: %v ", delErr)
		}
		return nil, fmt.Errorf("waiting for service endpoint ready. %v ", err)
	}

	return createdServiceEndpoint, err
}

func updateServiceEndpoint(clients *client.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint) (*serviceendpoint.ServiceEndpoint, error) {
	updatedServiceEndpoint, err := clients.ServiceEndpointClient.UpdateServiceEndpoint(
		clients.Ctx,
		serviceendpoint.UpdateServiceEndpointArgs{
			Endpoint:   endpoint,
			EndpointId: endpoint.Id,
		})

	return updatedServiceEndpoint, err
}

func deleteServiceEndpoint(clients *client.AggregatedClient, serviceEndpoint *serviceendpoint.ServiceEndpoint, timeout time.Duration) error {
	projectID := (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id
	if err := clients.ServiceEndpointClient.DeleteServiceEndpoint(
		clients.Ctx,
		serviceendpoint.DeleteServiceEndpointArgs{
			ProjectIds: &[]string{
				projectID.String(),
			},
			EndpointId: serviceEndpoint.Id,
		}); err != nil {
		return fmt.Errorf("Delete service endpoint error %v", err)
	}

	stateConf := &retry.StateChangeConf{
		ContinuousTargetOccurence: 1,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{opState.InProgress},
		Target:                    []string{opState.Ready, opState.Failed},
		Refresh:                   checkServiceEndpointStatus(clients, projectID, serviceEndpoint.Id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf("Wait for service endpoint to be deleted error. %v ", err)
	}
	return nil
}

func validateServiceEndpoint(clients *client.AggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, projectId string, retryTimeout time.Duration) error {
	reqArgs := serviceendpoint.ExecuteServiceEndpointRequestArgs{
		ServiceEndpointRequest: &serviceendpoint.ServiceEndpointRequest{
			DataSourceDetails: &serviceendpoint.DataSourceDetails{
				DataSourceName: converter.String("TestConnection"),
			},
			ResultTransformationDetails: &serviceendpoint.ResultTransformationDetails{},
			ServiceEndpointDetails: &serviceendpoint.ServiceEndpointDetails{
				Data:          endpoint.Data,
				Authorization: endpoint.Authorization,
				Url:           endpoint.Url,
				Type:          endpoint.Type,
			},
		},
		Project:    &projectId,
		EndpointId: converter.String(endpoint.Id.String()),
	}

	log.Printf(":: %s :: Initiating validation", *endpoint.Name)
	err := retry.RetryContext(clients.Ctx, retryTimeout, func() *retry.RetryError {
		reqResult, err := clients.ServiceEndpointClient.ExecuteServiceEndpointRequest(clients.Ctx, reqArgs)
		if err != nil {
			log.Printf(":: %s :: error during endpoint validation request", *endpoint.Name)
			return retry.NonRetryableError(err)
		}
		if !strings.EqualFold(*reqResult.StatusCode, "ok") {
			log.Printf(":: %s :: validation failed with StatusCode '%s', retrying...", *endpoint.Name, *reqResult.StatusCode)
			return retry.RetryableError(fmt.Errorf("Error validating connection: (type: %s, name: %s, code: %s, message: %s)", *endpoint.Type, *endpoint.Name, *reqResult.StatusCode, *reqResult.ErrorMessage))
		}
		log.Printf(":: %s :: successfully validated connection", *endpoint.Name)
		return nil
	})
	return err
}

func serviceEndpointGetArgs(d *schema.ResourceData) (*serviceendpoint.GetServiceEndpointDetailsArgs, error) {
	var serviceEndpointID *uuid.UUID
	parsedServiceEndpointID, err := uuid.Parse(d.Id())
	if err != nil {
		return nil, fmt.Errorf("parsing the service endpoint ID from the Terraform resource data: %v", err)
	}
	serviceEndpointID = &parsedServiceEndpointID
	projectID, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return nil, err
	}
	return &serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: serviceEndpointID,
		Project:    converter.String(projectID.String()),
	}, nil
}

// Service endpoint delete is an async operation, make sure service endpoint is deleted.
func checkServiceEndpointStatus(clients *client.AggregatedClient, projectID *uuid.UUID, endPointID *uuid.UUID) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
			clients.Ctx,
			serviceendpoint.GetServiceEndpointDetailsArgs{
				Project:    converter.String(projectID.String()),
				EndpointId: endPointID,
			})
		if err != nil {
			return nil, opState.Failed, fmt.Errorf(errMsgServiceDelete, endPointID, *projectID, err)
		}
		if serviceEndpoint != nil && serviceEndpoint.OperationStatus != nil {
			opStatus := (serviceEndpoint.OperationStatus).(map[string]interface{})["state"]
			if opStatus == opState.Failed {
				return nil, opState.Failed, fmt.Errorf(errMsgServiceDelete, endPointID, *projectID, serviceEndpoint.OperationStatus)
			}
			return serviceendpoint.ServiceEndpoint{}, opStatus.(string), nil
		}
		return serviceendpoint.ServiceEndpoint{}, opState.Ready, nil
	}
}

// Check if the service endpoint has been deleted
// 1) Service response 404
// 2) Service response 200 but service endpoint is null
// 3) Service response 200 but service endpoint is not null but ID is null
// 4) Service response 200 but service endpoint is not null and ID is null but other data is null
// 5) Service response 200 with state property? deleted = true/false or state = deleted/success/updating
func isServiceEndpointDeleted(d *schema.ResourceData, err error, serviceEndpoint *serviceendpoint.ServiceEndpoint, args *serviceendpoint.GetServiceEndpointDetailsArgs) bool {
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			log.Printf(" [INFO] Service endpoint not found. ID: (%v)", args.EndpointId)
			d.SetId("")
			return true
		}
		return false
	}

	if serviceEndpoint == nil || serviceEndpoint.Id == nil ||
		(serviceEndpoint.Id != nil && serviceEndpoint.Authorization == nil) {
		log.Printf(" [INFO] Service endpoint not found. ID: (%v)", args.EndpointId)
		d.SetId("")
		return true
	}
	return false
}

func getServiceEndpoint(client *client.AggregatedClient, serviceEndpointID *uuid.UUID, projectID *uuid.UUID) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		serviceEndpoint, err := client.ServiceEndpointClient.GetServiceEndpointDetails(
			client.Ctx,
			serviceendpoint.GetServiceEndpointDetailsArgs{
				EndpointId: serviceEndpointID,
				Project:    converter.String(projectID.String()),
			},
		)
		if err != nil {
			return nil, opState.Failed, fmt.Errorf(errMsgServiceCreate, serviceEndpointID, *projectID, err)
		}

		if *serviceEndpoint.IsReady {
			return serviceEndpoint, opState.Ready, nil
		} else if serviceEndpoint.OperationStatus != nil {
			opStatus := (serviceEndpoint.OperationStatus).(map[string]interface{})["state"]
			if opStatus == opState.Failed {
				return nil, opState.Failed, fmt.Errorf(errMsgServiceCreate, serviceEndpointID, *projectID, serviceEndpoint.OperationStatus)
			}
			return nil, opStatus.(string), nil
		}
		return nil, opState.Failed, fmt.Errorf(errMsgServiceCreate, serviceEndpointID, *projectID, serviceEndpoint.OperationStatus)
	}
}

// doBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func doBaseExpansion(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var serviceEndpointID *uuid.UUID
	parsedID, err := uuid.Parse(d.Id())
	if err == nil {
		serviceEndpointID = &parsedID
	}
	projectID := uuid.MustParse(d.Get("project_id").(string))
	name := converter.String(d.Get("service_endpoint_name").(string))
	serviceEndpoint := &serviceendpoint.ServiceEndpoint{
		Id:          serviceEndpointID,
		Name:        name,
		Owner:       converter.String("library"),
		Description: converter.String(d.Get("description").(string)),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
			{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: &projectID,
				},
				Name:        name,
				Description: converter.String(d.Get("description").(string)),
			},
		},
	}

	return serviceEndpoint
}

// doBaseFlattening performs the flattening for the 'base' attributes that are defined in the schema, above
func doBaseFlattening(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	if serviceEndpoint.Id != nil {
		d.SetId(serviceEndpoint.Id.String())
	}

	if serviceEndpoint.Name != nil {
		d.Set("service_endpoint_name", serviceEndpoint.Name)
	}

	if serviceEndpoint.Description != nil {
		d.Set("description", serviceEndpoint.Description)
	}

	if serviceEndpoint.ServiceEndpointProjectReferences != nil && len(*serviceEndpoint.ServiceEndpointProjectReferences) > 0 {
		for _, project := range *serviceEndpoint.ServiceEndpointProjectReferences {
			if strings.EqualFold(project.ProjectReference.Id.String(), d.Get("project_id").(string)) {
				d.Set("project_id", project.ProjectReference.Id.String())
				break
			}
		}
	}

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Scheme != nil {
		d.Set("authorization", &map[string]interface{}{
			"scheme": *serviceEndpoint.Authorization.Scheme,
		})
	}
}

// data resources

func dataSourceGenBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},

		"service_endpoint_name": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"service_endpoint_name", "service_endpoint_id"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"service_endpoint_id": {
			Description:  "The ID of the serviceendpoint",
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"service_endpoint_name", "service_endpoint_id"},
			ValidateFunc: validation.IsUUID,
		},

		"authorization": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceGetBaseServiceEndpoint(d *schema.ResourceData, m interface{}) (*serviceendpoint.ServiceEndpoint, error) {
	clients := m.(*client.AggregatedClient)

	var projectID *uuid.UUID
	projectIDString := d.Get("project_id").(string)
	parsedProjectID, err := uuid.Parse(projectIDString)
	if err != nil {
		return nil, fmt.Errorf("Parsing projectID from the Terraform data source declaration: %v", err)
	}

	projectID = &parsedProjectID

	if serviceEndpointIDString, ok := d.GetOk("service_endpoint_id"); ok {
		var serviceEndpointID *uuid.UUID
		parsedServiceEndpointID, err := uuid.Parse(serviceEndpointIDString.(string))
		if err != nil {
			return nil, fmt.Errorf("Parsing serviceEndpointID from the Terraform data source declaration: %v", err)
		}
		serviceEndpointID = &parsedServiceEndpointID

		serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
			clients.Ctx,
			serviceendpoint.GetServiceEndpointDetailsArgs{
				EndpointId: serviceEndpointID,
				Project:    converter.String(projectID.String()),
			},
		)
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil, nil
			}
			return nil, fmt.Errorf("Looking up service endpoint with ID (%v) and projectID (%v): %v", serviceEndpointID, projectID, err)
		}

		return serviceEndpoint, nil
	}

	if serviceEndpointName, ok := d.GetOk("service_endpoint_name"); ok {
		serviceEndpoint, err := dataSourceGetServiceEndpointByNameAndProject(clients, serviceEndpointName.(string), projectID.String())
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil, nil
			}
			return nil, fmt.Errorf("Looking up service endpoint with name (%v) and projectID (%v): %v", serviceEndpointName, projectID, err)
		}

		return serviceEndpoint, nil
	}
	return nil, nil
}

func dataSourceGetServiceEndpointByNameAndProject(clients *client.AggregatedClient, serviceEndpointName string, projectID string) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpointNameList := &[]string{serviceEndpointName}

	serviceEndpoints, err := clients.ServiceEndpointClient.GetServiceEndpointsByNames(
		clients.Ctx,
		serviceendpoint.GetServiceEndpointsByNamesArgs{
			Project:       &projectID,
			EndpointNames: serviceEndpointNameList,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(*serviceEndpoints) == 0 {
		return nil, fmt.Errorf("%v not found!", serviceEndpointName)
	}

	if len(*serviceEndpoints) > 1 {
		return nil, fmt.Errorf("%v returns more than one serviceEndpoint!", serviceEndpointName)
	}

	return &(*serviceEndpoints)[0], nil
}

type EndpointAuthenticationScheme string

const (
	ServicePrincipal           EndpointAuthenticationScheme = "ServicePrincipal"
	ManagedServiceIdentity     EndpointAuthenticationScheme = "ManagedServiceIdentity"
	WorkloadIdentityFederation EndpointAuthenticationScheme = "WorkloadIdentityFederation"
)

type EndpointCreationMode string

const (
	Automatic EndpointCreationMode = "Automatic"
	Manual    EndpointCreationMode = "Manual"
)

func checkServiceConnection(endpoint *serviceendpoint.ServiceEndpoint) error {
	if endpoint.Id != nil && (endpoint.Data == nil || endpoint.Type == nil) {
		return fmt.Errorf("Service connection not fully returned, this appears to be a permission issue with PAT/SPN/identity etc.")
	}
	return nil
}
