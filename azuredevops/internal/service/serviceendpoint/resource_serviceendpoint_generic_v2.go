package serviceendpoint

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/forminput"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// Cache to store validated service endpoint types
var (
	serviceEndpointTypesList        = make(map[string]serviceendpoint.ServiceEndpointType)
	serviceEndpointTypesMutex       sync.RWMutex
	serviceEndpointTypesInitialized bool
)

// EndpointConfig represents the configuration for a service endpoint
type EndpointConfig struct {
	ServiceEndpointType string
	AuthType            string
	AuthData            map[string]string
	Data                map[string]string
}

// ResourceServiceEndpointGenericV2 schema and implementation for generic service endpoint resource
func ResourceServiceEndpointGenericV2() *schema.Resource {
	return &schema.Resource{
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
		Importer:      tfhelper.ImportProjectQualifiedResourceUUID(),
		CustomizeDiff: customizeServiceEndpointGenericV2Diff,
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
			"authorization_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"authorization_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// InitServiceEndpointTypes loads all service endpoint types from Azure DevOps
func InitServiceEndpointTypes(ctx context.Context, client *serviceendpoint.Client) error {
	serviceEndpointTypesMutex.Lock()
	defer serviceEndpointTypesMutex.Unlock()

	if serviceEndpointTypesInitialized {
		return nil // Already initialized
	}

	args := serviceendpoint.GetServiceEndpointTypesArgs{}
	serviceEndpointTypes, err := (*client).GetServiceEndpointTypes(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to retrieve service endpoint types: %w", err)
	}

	if serviceEndpointTypes == nil {
		return fmt.Errorf("no service endpoint types available")
	}

	// Populate cache with fetched types
	for _, availableType := range *serviceEndpointTypes {
		if availableType.Name != nil {
			typeName := *availableType.Name
			serviceEndpointTypesList[typeName] = availableType
		}
	}

	serviceEndpointTypesInitialized = true
	return nil
}

// validateAuthScheme checks if the provided auth scheme is valid for the endpoint type
// and returns the list of possible input descriptors for that auth scheme
func validateAuthScheme(availableType *serviceendpoint.ServiceEndpointType, config EndpointConfig) (map[string]forminput.InputDescriptor, error) {
	if availableType == nil || availableType.AuthenticationSchemes == nil {
		return nil, fmt.Errorf("invalid service endpoint type definition")
	}

	possibleAuthSchemes := make([]string, 0, len(*availableType.AuthenticationSchemes))
	possibleAuthData := make(map[string]forminput.InputDescriptor)

	for _, authScheme := range *availableType.AuthenticationSchemes {
		if authScheme.Scheme == nil {
			continue
		}

		possibleAuthSchemes = append(possibleAuthSchemes, *authScheme.Scheme)

		if *authScheme.Scheme == config.AuthType && authScheme.InputDescriptors != nil {
			for _, data := range *authScheme.InputDescriptors {
				if data.Id != nil {
					possibleAuthData[*data.Id] = data
				}
			}
			return possibleAuthData, nil
		}
	}

	return nil, fmt.Errorf("service endpoint type '%s' does not support authentication scheme '%s'. Supported schemes: %v",
		*availableType.Name, config.AuthType, possibleAuthSchemes)
}

// validateFields ensures that the provided configuration fields match the expected fields
func validateFields(configFields map[string]string, possibleFields map[string]forminput.InputDescriptor, fieldType, endpointType string, planTime bool) error {

	// Skip validation at plan time if fields are empty (known-after-apply)
	if planTime && len(configFields) == 0 {
		return nil
	}

	// Check for unsupported fields
	invalidFields := make([]string, 0)
	for key := range configFields {
		if _, exists := possibleFields[key]; !exists {
			invalidFields = append(invalidFields, key)
		}
	}

	if len(invalidFields) > 0 {
		validFields := make([]string, 0, len(possibleFields))
		for k := range possibleFields {
			validFields = append(validFields, k)
		}
		return fmt.Errorf("service endpoint type '%s' does not support %s field(s): %v. Supported fields: %v",
			endpointType, fieldType, invalidFields, validFields)
	}

	// Check for missing required fields
	missingFields := make(map[string]string)
	for key, value := range possibleFields {
		if value.Validation != nil && value.Validation.IsRequired != nil && *value.Validation.IsRequired {
			if _, exists := configFields[key]; !exists {
				if value.Name != nil {
					missingFields[key] = *value.Name
				} else {
					missingFields[key] = key
				}
			}
		}
	}

	if len(missingFields) > 0 {
		missingFieldList := make([]string, 0, len(missingFields))
		for k, v := range missingFields {
			missingFieldList = append(missingFieldList, fmt.Sprintf("%s: %s", k, v))
		}
		return fmt.Errorf("service endpoint type '%s' is missing required %s fields: %v",
			endpointType, fieldType, missingFieldList)
	}

	return nil
}

// validateServiceEndpointType validates that the configuration matches the endpoint type requirements
func validateServiceEndpointType(availableType *serviceendpoint.ServiceEndpointType, config EndpointConfig, planTime bool) error {
	if availableType == nil || availableType.Name == nil {
		return fmt.Errorf("invalid service endpoint type definition")
	}

	// Validate Data fields
	possibleData := make(map[string]forminput.InputDescriptor)
	if availableType.InputDescriptors != nil {
		for _, data := range *availableType.InputDescriptors {
			if data.Id != nil {
				possibleData[*data.Id] = data
			}
		}
	}

	if err := validateFields(config.Data, possibleData, "data", *availableType.Name, planTime); err != nil {
		return err
	}

	// Validate AuthData fields
	possibleAuthData, err := validateAuthScheme(availableType, config)
	if err != nil {
		return err
	}

	if err := validateFields(config.AuthData, possibleAuthData, "auth", *availableType.Name, planTime); err != nil {
		return err
	}

	return nil
}

// validateServiceEndpointSchema validates that the service endpoint type exists and configuration is valid
func validateServiceEndpointSchema(clients *client.AggregatedClient, serviceEndpoint EndpointConfig, planTime bool) error {
	serviceEndpointType := serviceEndpoint.ServiceEndpointType

	// Check if types have been initialized yet
	serviceEndpointTypesMutex.RLock()
	initialized := serviceEndpointTypesInitialized
	serviceEndpointTypesMutex.RUnlock()

	// If not initialized, initialize the cache
	if !initialized {
		if err := InitServiceEndpointTypes(clients.Ctx, &clients.ServiceEndpointClient); err != nil {
			return fmt.Errorf("error initializing service endpoint types: %w", err)
		}
	}

	// Check if the requested type exists in the cache
	serviceEndpointTypesMutex.RLock()
	foundType, ok := serviceEndpointTypesList[serviceEndpointType]
	serviceEndpointTypesMutex.RUnlock()

	// If the type was found in the cache, validate it
	if ok {
		return validateServiceEndpointType(&foundType, serviceEndpoint, planTime)
	}

	// If the type wasn't found, prepare an error message with all valid types
	serviceEndpointTypesMutex.RLock()
	availableTypes := make([]string, 0, len(serviceEndpointTypesList))
	for _, endpoint := range serviceEndpointTypesList {
		if endpoint.DisplayName != nil && endpoint.Name != nil {
			availableTypes = append(availableTypes, fmt.Sprintf("%s: %s", *endpoint.DisplayName, *endpoint.Name))
		}
	}
	serviceEndpointTypesMutex.RUnlock()

	if len(availableTypes) == 0 {
		return fmt.Errorf("service endpoint type '%s' is not available. No service endpoint types available",
			serviceEndpointType)
	}

	return fmt.Errorf(
		"service endpoint type '%s' is not available.\nValid types are:\n%s",
		serviceEndpointType,
		strings.Join(availableTypes, "\n"),
	)
}

// resourceServiceEndpointGenericV2Create creates a new service endpoint
func resourceServiceEndpointGenericV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	// Get configuration values from the resource data
	config, err := getEndpointConfigFromResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate the service endpoint configuration
	if err := validateServiceEndpointSchema(clients, *config, false); err != nil {
		return diag.FromErr(fmt.Errorf("service endpoint validation failed: %w", err))
	}

	// Create the service endpoint
	serviceEndpoint, err := createGenericV2ServiceEndpoint(ctx, d, clients, config)
	if err != nil {
		return diag.FromErr(err)
	}

	if serviceEndpoint == nil || serviceEndpoint.Id == nil {
		return diag.FromErr(fmt.Errorf("service endpoint creation failed: endpoint or ID is nil"))
	}

	d.SetId(serviceEndpoint.Id.String())
	return resourceServiceEndpointGenericV2Read(ctx, d, m)
}

// getEndpointConfigFromResource extracts endpoint configuration from resource data
func getEndpointConfigFromResource(d *schema.ResourceData) (*EndpointConfig, error) {
	// Get authorization details
	authScheme, authParams, err := getAuthorizationDetails(d)
	if err != nil {
		return nil, fmt.Errorf("error processing authorization details: %w", err)
	}

	// Get additional parameters
	data, err := toStringMap(d.Get("parameters"), "parameters")
	if err != nil {
		return nil, err
	}

	return &EndpointConfig{
		ServiceEndpointType: d.Get("service_endpoint_type").(string),
		AuthType:            authScheme,
		AuthData:            authParams,
		Data:                data,
	}, nil
}

// createGenericV2ServiceEndpoint creates a service endpoint in Azure DevOps
func createGenericV2ServiceEndpoint(ctx context.Context, d *schema.ResourceData, clients *client.AggregatedClient, config *EndpointConfig) (*serviceendpoint.ServiceEndpoint, error) {
	name := d.Get("service_endpoint_name").(string)
	projectID := d.Get("project_id").(string)
	description := d.Get("description").(string)
	serverURL := d.Get("server_url").(string)

	// Generate the project reference
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	projectReference := serviceendpoint.ServiceEndpointProjectReference{
		ProjectReference: &serviceendpoint.ProjectReference{
			Id: &projectUUID,
		},
		Name:        converter.String(name),
		Description: converter.String(description),
	}

	// Create service endpoint object
	serviceEndpoint := &serviceendpoint.ServiceEndpoint{
		Name:                             converter.String(name),
		Description:                      converter.String(description),
		Type:                             converter.String(config.ServiceEndpointType),
		Url:                              converter.String(serverURL),
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{projectReference},
		Authorization: &serviceendpoint.EndpointAuthorization{
			Scheme:     converter.String(config.AuthType),
			Parameters: &config.AuthData,
		},
	}

	// Handle additional data
	if len(config.Data) > 0 {
		serviceEndpoint.Data = &config.Data
	}

	// Create service endpoint in Azure DevOps
	args := serviceendpoint.CreateServiceEndpointArgs{
		Endpoint: serviceEndpoint,
	}

	createdEndpoint, err := clients.ServiceEndpointClient.CreateServiceEndpoint(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("error creating service endpoint: %w", err)
	}

	return createdEndpoint, nil
}

// resourceServiceEndpointGenericV2Read reads a service endpoint
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

	return updateResourceDataFromServiceEndpoint(d, serviceEndpoint)
}

// updateResourceDataFromServiceEndpoint updates the resource data with values from the service endpoint
func updateResourceDataFromServiceEndpoint(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) diag.Diagnostics {
	if serviceEndpoint.Name != nil {
		if err := d.Set("service_endpoint_name", *serviceEndpoint.Name); err != nil {
			return diag.FromErr(fmt.Errorf("error setting service_endpoint_name: %w", err))
		}
	}

	if serviceEndpoint.Description != nil {
		if err := d.Set("description", *serviceEndpoint.Description); err != nil {
			return diag.FromErr(fmt.Errorf("error setting description: %w", err))
		}
	}

	if serviceEndpoint.Type != nil {
		if err := d.Set("service_endpoint_type", *serviceEndpoint.Type); err != nil {
			return diag.FromErr(fmt.Errorf("error setting service_endpoint_type: %w", err))
		}
	}

	if serviceEndpoint.Url != nil {
		if err := d.Set("server_url", *serviceEndpoint.Url); err != nil {
			return diag.FromErr(fmt.Errorf("error setting server_url: %w", err))
		}
	}

	// Handle authorization
	if serviceEndpoint.Authorization != nil {
		if serviceEndpoint.Authorization.Scheme != nil {
			if err := d.Set("authorization_type", *serviceEndpoint.Authorization.Scheme); err != nil {
				return diag.FromErr(fmt.Errorf("error setting authorization_type: %w", err))
			}
		}

		if serviceEndpoint.Authorization.Parameters != nil {
			authParams := make(map[string]string)
			for k, v := range *serviceEndpoint.Authorization.Parameters {
				authParams[k] = v
			}
			if err := d.Set("authorization_parameters", authParams); err != nil {
				return diag.FromErr(fmt.Errorf("error setting authorization_parameters: %w", err))
			}
		} else {
			if err := d.Set("authorization_parameters", nil); err != nil {
				return diag.FromErr(fmt.Errorf("error setting authorization_parameters to nil: %w", err))
			}
		}
	} else {
		if err := d.Set("authorization_type", ""); err != nil {
			return diag.FromErr(fmt.Errorf("error setting authorization_type to empty: %w", err))
		}
		if err := d.Set("authorization_parameters", nil); err != nil {
			return diag.FromErr(fmt.Errorf("error setting authorization_parameters to nil: %w", err))
		}
	}

	// Handle data parameters
	if serviceEndpoint.Data != nil {
		data := make(map[string]string)
		for k, v := range *serviceEndpoint.Data {
			data[k] = v
		}
		if len(data) > 0 {
			if err := d.Set("parameters", data); err != nil {
				return diag.FromErr(fmt.Errorf("error setting parameters: %w", err))
			}
		}
	}

	return nil
}

// resourceServiceEndpointGenericV2Update updates an existing service endpoint
func resourceServiceEndpointGenericV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpointID := d.Id()
	projectID := d.Get("project_id").(string)

	// Get current service endpoint to preserve any fields we're not updating
	currentEndpoint, err := getServiceEndpointGenericV2(ctx, clients, serviceEndpointID, projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the service endpoint with values from resource data
	updatedEndpoint, err := updateServiceEndpointFromResourceData(d, currentEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update service endpoint in Azure DevOps
	_, err = updateServiceEndpointGenericV2(ctx, clients, updatedEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceEndpointGenericV2Read(ctx, d, m)
}

// updateServiceEndpointFromResourceData updates a service endpoint with values from resource data
func updateServiceEndpointFromResourceData(d *schema.ResourceData, endpoint *serviceendpoint.ServiceEndpoint) (*serviceendpoint.ServiceEndpoint, error) {
	// Update fields that have changed
	if d.HasChange("service_endpoint_name") {
		endpoint.Name = converter.String(d.Get("service_endpoint_name").(string))
	}

	if d.HasChange("description") {
		endpoint.Description = converter.String(d.Get("description").(string))
	}

	if d.HasChange("server_url") {
		endpoint.Url = converter.String(d.Get("server_url").(string))
	}

	// Handle authorization updates if any auth fields changed
	if d.HasChange("authorization_type") || d.HasChange("authorization_parameters") {
		authScheme, authParams, err := getAuthorizationDetails(d)
		if err != nil {
			return nil, fmt.Errorf("error processing authorization details: %w", err)
		}
		endpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Scheme:     converter.String(authScheme),
			Parameters: &authParams,
		}
	}

	// Handle data updates if changed
	if d.HasChange("parameters") {
		data, err := toStringMap(d.Get("parameters"), "parameters")
		if err != nil {
			return nil, err
		}
		endpoint.Data = &data
	}

	return endpoint, nil
}

// resourceServiceEndpointGenericV2Delete deletes a service endpoint
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

// toStringMap converts a raw interface{} to map[string]string
func toStringMap(raw interface{}, fieldName string) (map[string]string, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid format for %s", fieldName)
	}

	result := make(map[string]string, len(m))
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%s value for key %q is not a string", fieldName, k)
		}
		result[k] = s
	}
	return result, nil
}

// getAuthorizationDetailsRaw extracts auth scheme and params from either ResourceData or ResourceDiff
func getAuthorizationDetailsRaw(schemeVal interface{}, paramsVal interface{}) (string, map[string]string, error) {
	scheme, _ := schemeVal.(string)
	if scheme == "" {
		return "", nil, fmt.Errorf("missing or invalid authorization scheme")
	}

	params, err := toStringMap(paramsVal, "authorization_parameters")
	if err != nil {
		return "", nil, err
	}
	return scheme, params, nil
}

// getAuthorizationDetails extracts auth scheme and params from ResourceData
func getAuthorizationDetails(d *schema.ResourceData) (string, map[string]string, error) {
	return getAuthorizationDetailsRaw(d.Get("authorization_type"), d.Get("authorization_parameters"))
}

// getAuthorizationDetailsFromDiff extracts auth scheme and params from ResourceDiff
func getAuthorizationDetailsFromDiff(d *schema.ResourceDiff) (string, map[string]string, error) {
	return getAuthorizationDetailsRaw(d.Get("authorization_type"), d.Get("authorization_parameters"))
}

// getServiceEndpointGenericV2 retrieves a service endpoint from Azure DevOps
func getServiceEndpointGenericV2(ctx context.Context, clients *client.AggregatedClient, endpointID, projectID string) (*serviceendpoint.ServiceEndpoint, error) {
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return nil, fmt.Errorf("invalid service endpoint ID: %w", err)
	}

	args := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: &endpointUUID,
		Project:    converter.String(projectID),
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("error looking up service endpoint: %w", err)
	}

	return serviceEndpoint, nil
}

// updateServiceEndpointGenericV2 updates a service endpoint in Azure DevOps
func updateServiceEndpointGenericV2(ctx context.Context, clients *client.AggregatedClient, serviceEndpoint *serviceendpoint.ServiceEndpoint) (*serviceendpoint.ServiceEndpoint, error) {
	args := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   serviceEndpoint,
		EndpointId: serviceEndpoint.Id,
	}

	updatedEndpoint, err := clients.ServiceEndpointClient.UpdateServiceEndpoint(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("error updating service endpoint: %w", err)
	}

	return updatedEndpoint, nil
}

// deleteServiceEndpointGenericV2 deletes a service endpoint from Azure DevOps
func deleteServiceEndpointGenericV2(ctx context.Context, clients *client.AggregatedClient, endpointID, projectID string) error {
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return fmt.Errorf("invalid service endpoint ID: %w", err)
	}

	args := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: &endpointUUID,
		ProjectIds: &[]string{projectID},
	}

	err = clients.ServiceEndpointClient.DeleteServiceEndpoint(ctx, args)
	if err != nil {
		return fmt.Errorf("error deleting service endpoint: %w", err)
	}

	return nil
}

// customizeServiceEndpointGenericV2Diff validates the service endpoint configuration during the planning phase
func customizeServiceEndpointGenericV2Diff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	// Only validate on resource creation changes
	if d.Id() != "" {
		return nil
	}

	serviceEndpointType := d.Get("service_endpoint_type").(string)
	authScheme, authParams, err := getAuthorizationDetailsFromDiff(d)
	if err != nil {
		return err
	}

	data, err := toStringMap(d.Get("parameters"), "parameters")
	if err != nil {
		return err
	}

	config := EndpointConfig{
		ServiceEndpointType: serviceEndpointType,
		AuthType:            authScheme,
		AuthData:            authParams,
		Data:                data,
	}

	return validateServiceEndpointSchema(m.(*client.AggregatedClient), config, true)
}
