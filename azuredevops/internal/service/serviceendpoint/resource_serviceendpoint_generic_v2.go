package serviceendpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointGenericV2 schema and implementation for generic service endpoint resource
func ResourceServiceEndpointGenericV2() *schema.Resource {
	r := &schema.Resource{
		CreateContext: resourceServiceEndpointGenericV2Create,
		ReadContext:   resourceServiceEndpointGenericV2Read,
		UpdateContext: resourceServiceEndpointGenericV2Update,
		DeleteContext: resourceServiceEndpointGenericV2Delete,
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
				ValidateFunc: validation.IsUUID,
			},
			"service_endpoint_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"service_endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"server_url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.IsURLWithHTTPorHTTPS,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"authorization": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"parameters": {
							Type:      schema.TypeMap,
							Optional:  true,
							Sensitive: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"data": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"authorization_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return r
}

func resourceServiceEndpointGenericV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	name := d.Get("service_endpoint_name").(string)
	projectID := d.Get("project_id").(string)
	description := d.Get("description").(string)
	serviceEndpointType := d.Get("service_endpoint_type").(string)
	serverURL := d.Get("server_url").(string)

	// Get authorization details
	authScheme, authParams, err := getAuthorizationDetails(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get additional data
	data := make(map[string]string)
	if dataRaw := d.Get("data"); dataRaw != nil {
		dataMap, ok := dataRaw.(map[string]interface{})
		if !ok {
			return diag.FromErr(fmt.Errorf("invalid data format"))
		}

		for k, v := range dataMap {
			strVal, ok := v.(string)
			if !ok {
				return diag.FromErr(fmt.Errorf("data value for key %q is not a string", k))
			}
			data[k] = strVal
		}
	}

	serviceEndpoint, err := createGenericV2ServiceEndpoint(ctx, clients, name, projectID, description, serviceEndpointType, serverURL, authScheme, authParams, data)
	if err != nil {
		return diag.FromErr(err)
	}

	if serviceEndpoint == nil || serviceEndpoint.Id == nil {
		return diag.FromErr(fmt.Errorf("service endpoint creation failed: endpoint or ID is nil"))
	}

	d.SetId(serviceEndpoint.Id.String())
	err = d.Set("authorization_type", authScheme)
	if err != nil {
		return nil
	}

	return resourceServiceEndpointGenericV2Read(ctx, d, m)
}

func resourceServiceEndpointGenericV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpointID := d.Id()
	projectID := d.Get("project_id").(string)

	if serviceEndpointID == "" {
		return diag.FromErr(fmt.Errorf("service endpoint ID is required"))
	}

	if projectID == "" {
		return diag.FromErr(fmt.Errorf("project ID is required"))
	}

	serviceEndpoint, err := getServiceEndpointGenericV2(ctx, clients, serviceEndpointID, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if serviceEndpoint == nil {
		d.SetId("")
		return nil
	}

	// Update state with the latest information from the service endpoint
	if err := d.Set("service_endpoint_name", serviceEndpoint.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting service_endpoint_name: %v", err))
	}

	if err := d.Set("description", serviceEndpoint.Description); err != nil {
		return diag.FromErr(fmt.Errorf("error setting description: %v", err))
	}

	if err := d.Set("service_endpoint_type", serviceEndpoint.Type); err != nil {
		return diag.FromErr(fmt.Errorf("error setting service_endpoint_type: %v", err))
	}

	if err := d.Set("server_url", serviceEndpoint.Url); err != nil {
		return diag.FromErr(fmt.Errorf("error setting server_url: %v", err))
	}

	// Handle authorization
	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Scheme != nil {
		if err := d.Set("authorization_type", *serviceEndpoint.Authorization.Scheme); err != nil {
			return diag.FromErr(fmt.Errorf("error setting authorization_type: %v", err))
		}
		// We don't update the authorization parameters as they may contain sensitive information
		// that we don't get back from the API
	}

	// Handle data - copy non-sensitive values
	if serviceEndpoint.Data != nil {
		data := make(map[string]string)
		for k, v := range *serviceEndpoint.Data {
			data[k] = v
		}
		if len(data) > 0 {
			if err := d.Set("data", data); err != nil {
				return diag.FromErr(fmt.Errorf("error setting data: %v", err))
			}
		}
	}

	return nil
}

func resourceServiceEndpointGenericV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpointID := d.Id()
	projectID := d.Get("project_id").(string)

	// Get current service endpoint to preserve any fields we're not updating
	currentEndpoint, err := getServiceEndpointGenericV2(ctx, clients, serviceEndpointID, projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update fields that have changed
	if d.HasChange("service_endpoint_name") {
		currentEndpoint.Name = converter.String(d.Get("service_endpoint_name").(string))
	}

	if d.HasChange("description") {
		currentEndpoint.Description = converter.String(d.Get("description").(string))
	}

	if d.HasChange("server_url") {
		currentEndpoint.Url = converter.String(d.Get("server_url").(string))
	}

	// Handle authorization updates if changed
	if d.HasChange("authorization") {
		authScheme, authParams, err := getAuthorizationDetails(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error processing authorization details: %v", err))
		}

		authorization := &serviceendpoint.EndpointAuthorization{
			Scheme:     converter.String(authScheme),
			Parameters: &authParams,
		}

		currentEndpoint.Authorization = authorization
	}

	// Handle data updates if changed
	if d.HasChange("data") {
		data := make(map[string]string)

		// If there are existing data values, preserve them
		if currentEndpoint.Data != nil {
			for k, v := range *currentEndpoint.Data {
				data[k] = v
			}
		}

		// Update with new values
		if dataRaw := d.Get("data"); dataRaw != nil {
			dataMap, ok := dataRaw.(map[string]interface{})
			if !ok {
				return diag.FromErr(fmt.Errorf("invalid data format"))
			}

			for k, v := range dataMap {
				strVal, ok := v.(string)
				if !ok {
					return diag.FromErr(fmt.Errorf("data value for key %q is not a string", k))
				}
				data[k] = strVal
			}
		}

		currentEndpoint.Data = &data
	}

	// Update service endpoint
	_, err = updateServiceEndpointGenericV2(ctx, clients, currentEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceEndpointGenericV2Read(ctx, d, m)
}

func resourceServiceEndpointGenericV2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpointID := d.Id()
	projectID := d.Get("project_id").(string)

	err := deleteServiceEndpointGenericV2(ctx, clients, serviceEndpointID, projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getAuthorizationDetails(d *schema.ResourceData) (string, map[string]string, error) {
	authSet := d.Get("authorization").(*schema.Set)
	if authSet.Len() == 0 {
		return "", nil, fmt.Errorf("no authorization configuration found")
	}

	authData, ok := authSet.List()[0].(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("invalid authorization configuration format")
	}

	scheme, ok := authData["scheme"].(string)
	if !ok || scheme == "" {
		return "", nil, fmt.Errorf("missing or invalid authorization scheme")
	}

	params := make(map[string]string)
	if paramsRaw, ok := authData["parameters"].(map[string]interface{}); ok {
		for k, v := range paramsRaw {
			strValue, ok := v.(string)
			if !ok {
				return "", nil, fmt.Errorf("parameter %q has invalid type, expected string", k)
			}
			params[k] = strValue
		}
	}

	return scheme, params, nil
}

func createGenericV2ServiceEndpoint(ctx context.Context, clients *client.AggregatedClient, name, projectID, description, serviceEndpointType, serverURL, authScheme string, authParams map[string]string, data map[string]string) (*serviceendpoint.ServiceEndpoint, error) {
	// Validate service endpoint type exists
	if err := validateServiceEndpointType(ctx, clients, serviceEndpointType); err != nil {
		return nil, err
	}

	// Generate the project reference
	projectReference := serviceendpoint.ServiceEndpointProjectReference{
		ProjectReference: &serviceendpoint.ProjectReference{
			Id: &uuid.UUID{}, // Will be filled later
		},
		Name:        converter.String(name),
		Description: converter.String(description),
	}

	// Fill the project ID
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %v", err)
	}
	projectReference.ProjectReference.Id = &projectUUID

	// Create service endpoint object
	serviceEndpoint := &serviceendpoint.ServiceEndpoint{
		Name:                             converter.String(name),
		Description:                      converter.String(description),
		Type:                             converter.String(serviceEndpointType),
		Url:                              converter.String(serverURL),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{projectReference},
	}

	// Handle authentication
	authorization := &serviceendpoint.EndpointAuthorization{
		Scheme:     converter.String(authScheme),
		Parameters: &authParams,
	}
	serviceEndpoint.Authorization = authorization

	// Handle additional data
	if len(data) > 0 {
		serviceEndpoint.Data = &data
	}

	// Create service endpoint in Azure DevOps
	args := serviceendpoint.CreateServiceEndpointArgs{
		Endpoint: serviceEndpoint,
	}

	createdEndpoint, err := clients.ServiceEndpointClient.CreateServiceEndpoint(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("error creating service endpoint: %v", err)
	}

	return createdEndpoint, nil
}

func getServiceEndpointGenericV2(ctx context.Context, clients *client.AggregatedClient, endpointID, projectID string) (*serviceendpoint.ServiceEndpoint, error) {
	ProjectUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return nil, fmt.Errorf("invalid service endpoint ID: %v", err)
	}

	args := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: &ProjectUUID,
		Project:    converter.String(projectID),
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("error looking up service endpoint: %v", err)
	}

	return serviceEndpoint, nil
}

func updateServiceEndpointGenericV2(ctx context.Context, clients *client.AggregatedClient, serviceEndpoint *serviceendpoint.ServiceEndpoint) (*serviceendpoint.ServiceEndpoint, error) {
	args := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   serviceEndpoint,
		EndpointId: serviceEndpoint.Id,
	}

	updatedEndpoint, err := clients.ServiceEndpointClient.UpdateServiceEndpoint(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("error updating service endpoint: %v", err)
	}

	return updatedEndpoint, nil
}

func deleteServiceEndpointGenericV2(ctx context.Context, clients *client.AggregatedClient, endpointID, projectID string) error {
	ProjectUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return fmt.Errorf("invalid service endpoint ID: %v", err)
	}

	args := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: &ProjectUUID,
		ProjectIds: &[]string{projectID},
	}

	err = clients.ServiceEndpointClient.DeleteServiceEndpoint(ctx, args)
	if err != nil {
		return fmt.Errorf("error deleting service endpoint: %v", err)
	}

	return nil
}

func validateServiceEndpointType(ctx context.Context, clients *client.AggregatedClient, serviceEndpointType string) error {
	// Get available service endpoint types from Azure DevOps
	args := serviceendpoint.GetServiceEndpointTypesArgs{}

	serviceEndpointTypes, err := clients.ServiceEndpointClient.GetServiceEndpointTypes(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to retrieve service endpoint types: %v", err)
	}

	if serviceEndpointTypes == nil {
		return fmt.Errorf("no service endpoint types available")
	}

	// Check if the requested type exists
	for _, availableType := range *serviceEndpointTypes {
		if availableType.Name != nil && *availableType.Name == serviceEndpointType {
			return nil
		}
	}

	// If we reach here, the type wasn't found
	return fmt.Errorf("service endpoint type '%s' is not available", serviceEndpointType)
}
