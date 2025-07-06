package serviceendpoint

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	serviceEndpointTypesCache       = make(map[string]bool)
	serviceEndpointTypesList        = make([]serviceendpoint.ServiceEndpointType, 0)
	serviceEndpointTypesMutex       sync.RWMutex
	serviceEndpointTypesInitialized bool
)

// EndpointConfig represents the configuration for a service endpoint
type EndpointConfig struct {
	ServiceEndpointType string
	AuthType            string
	AuthData            map[string]string
	Data                map[string]string
	// Fields to track which fields are explicitly configured in the resource
	HasDataBlock       bool
	HasParametersBlock bool
}

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
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: false,
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
		},
	}

	return r
}

// InitServiceEndpointTypes loads all service endpoint types from Azure DevOps
// This can be called during provider initialization or when first needed
func InitServiceEndpointTypes(ctx context.Context, client *serviceendpoint.Client) error {
	serviceEndpointTypesMutex.Lock()
	defer serviceEndpointTypesMutex.Unlock()

	if serviceEndpointTypesInitialized {
		return nil // Already initialized
	}

	args := serviceendpoint.GetServiceEndpointTypesArgs{}
	serviceEndpointTypes, err := (*client).GetServiceEndpointTypes(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to retrieve service endpoint types: %v", err)
	}

	if serviceEndpointTypes == nil {
		return fmt.Errorf("no service endpoint types available")
	}

	// Clear existing cache
	serviceEndpointTypesCache = make(map[string]bool)
	serviceEndpointTypesList = make([]serviceendpoint.ServiceEndpointType, 0, len(*serviceEndpointTypes))

	// Populate cache with fetched types
	for _, availableType := range *serviceEndpointTypes {
		if availableType.Name != nil {
			typeName := *availableType.Name
			serviceEndpointTypesCache[typeName] = true
			serviceEndpointTypesList = append(serviceEndpointTypesList, availableType)
		}
	}

	serviceEndpointTypesInitialized = true
	return nil
}

func validateAuthScheme(availableType *serviceendpoint.ServiceEndpointType, config EndpointConfig) (map[string]forminput.InputDescriptor, error) {
	possibleAuthSchemes := make([]string, 0, len(*availableType.AuthenticationSchemes))
	possibleAuthData := make(map[string]forminput.InputDescriptor)

	correctAuthScheme := false

	for _, authScheme := range *availableType.AuthenticationSchemes {
		possibleAuthSchemes = append(possibleAuthSchemes, *authScheme.Scheme)
		if *authScheme.Scheme == config.AuthType {
			correctAuthScheme = true
			for _, data := range *authScheme.InputDescriptors {
				possibleAuthData[*data.Id] = data
			}
		}
	}
	if !correctAuthScheme {
		return nil, fmt.Errorf("service endpoint type '%s' does not support authentication scheme '%s'. Supported schemes: %v",
			*availableType.Name, config.AuthType, possibleAuthSchemes)
	}
	return possibleAuthData, nil
}

func validateFields(ctx context.Context, configFields map[string]string, possibleFields map[string]forminput.InputDescriptor, fieldType, endpointType string, isConfigured bool) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(configFields) == 0 {
		// Only log a warning if this field is actually configured in the resource
		if isConfigured {
			message := fmt.Sprintf("Received empty configFields for service endpoint type '%s' for '%s' field. This indicates sensitive data or known-after-apply fields in configuration. Skipping validation.", endpointType, fieldType)
			tflog.Warn(ctx, message)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Skipping field validation for sensitive or known-after-apply data",
				Detail:   message,
			})
		}
		return diags
	}
	// Check for unsupported fields
	for key := range configFields {
		if _, exists := possibleFields[key]; !exists {
			return diag.FromErr(fmt.Errorf("service endpoint type '%s' does not support %s field '%s'. Supported fields: {%s}",
				endpointType, fieldType, key, func() string {
					fields := make([]string, 0, len(possibleFields))
					for k := range possibleFields {
						fields = append(fields, k)
					}
					return fmt.Sprintf("%s", fields)
				}()))
		}
	}
	// Check for missing required fields
	missingFields := make(map[string]string)
	for key, value := range possibleFields {
		if value.Validation != nil && value.Validation.IsRequired != nil && *value.Validation.IsRequired {
			if _, exists := configFields[key]; !exists {
				missingFields[key] = *value.Name
			}
		}
	}
	if len(missingFields) > 0 {
		return diag.FromErr(fmt.Errorf("service endpoint type '%s' is missing required %s fields: {%s}", endpointType, fieldType, func() string {
			fields := make([]string, 0, len(missingFields))
			for k, v := range missingFields {
				fields = append(fields, fmt.Sprintf("Key: %s, Display Name: %s\n", k, v))
			}
			return fmt.Sprintf("%s", fields)
		}()))
	}
	return diags
}

func validateServiceEndpointType(ctx context.Context, availableType *serviceendpoint.ServiceEndpointType, config EndpointConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Validate Data fields
	possibleData := make(map[string]forminput.InputDescriptor)
	for _, data := range *availableType.InputDescriptors {
		possibleData[*data.Id] = data
	}

	dataDiags := validateFields(ctx, config.Data, possibleData, "data", *availableType.Name, config.HasDataBlock)
	diags = append(diags, dataDiags...)

	// Validate AuthData fields
	possibleAuthData, err := validateAuthScheme(availableType, config)
	if err != nil {
		return diag.FromErr(err)
	}

	authDiags := validateFields(ctx, config.AuthData, possibleAuthData, "auth", *availableType.Name, config.HasParametersBlock)
	diags = append(diags, authDiags...)

	return diags
}

// validateServiceEndpointSchema validates that the service endpoint type exists using the Azure DevOps API
func validateServiceEndpointSchema(ctx context.Context, clients *client.AggregatedClient, serviceEndpoint EndpointConfig) diag.Diagnostics {
	serviceEndpointType := serviceEndpoint.ServiceEndpointType

	// Check if types have been initialized yet
	serviceEndpointTypesMutex.RLock()
	initialized := serviceEndpointTypesInitialized
	serviceEndpointTypesMutex.RUnlock()

	// If not initialized, initialize the cache
	if !initialized {
		if err := InitServiceEndpointTypes(clients.Ctx, &clients.ServiceEndpointClient); err != nil {
			return diag.FromErr(fmt.Errorf("error initializing service endpoint types: %v", err))
		}
	}

	// Check if the requested type exists in the cache
	serviceEndpointTypesMutex.RLock()
	typeFound := serviceEndpointTypesCache[serviceEndpointType]
	var foundType *serviceendpoint.ServiceEndpointType

	// Find the type details if it exists
	if typeFound {
		for _, availableType := range serviceEndpointTypesList {
			if availableType.Name != nil && *availableType.Name == serviceEndpointType {
				typeCopy := availableType // Create a copy to avoid race conditions
				foundType = &typeCopy
				break
			}
		}
	}
	serviceEndpointTypesMutex.RUnlock()

	// If the type was found in the cache, validate it
	if foundType != nil {
		return validateServiceEndpointType(ctx, foundType, serviceEndpoint)
	}

	// If the type wasn't found, prepare an error message with all valid types
	serviceEndpointTypesMutex.RLock()
	result := make([]string, 0, len(serviceEndpointTypesList))
	for _, endpoint := range serviceEndpointTypesList {
		if endpoint.DisplayName != nil && endpoint.Name != nil {
			result = append(result, fmt.Sprintf("%s: %s", *endpoint.DisplayName, *endpoint.Name))
		}
	}
	serviceEndpointTypesMutex.RUnlock()

	return diag.FromErr(fmt.Errorf(
		"service endpoint type '%s' is not available.\nValid types are:\n%s",
		serviceEndpointType,
		func() string {
			if len(result) == 0 {
				return "No service endpoint types available"
			}
			return strings.Join(result, "\n")
		}(),
	))
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

	// Has to be called again to validate known-after-apply fields
	config := EndpointConfig{
		ServiceEndpointType: serviceEndpointType,
		AuthType:            authScheme,
		AuthData:            authParams,
		Data:                data,
	}

	// Validate the service endpoint schema
	diags := validateServiceEndpointSchema(ctx, clients, config)
	if diags.HasError() {
		return diags
	}

	serviceEndpoint, err := createGenericV2ServiceEndpoint(ctx, clients, name, projectID, description, serviceEndpointType, serverURL, authScheme, authParams, data)
	if err != nil {
		return diag.FromErr(err)
	}

	if serviceEndpoint == nil || serviceEndpoint.Id == nil {
		return diag.FromErr(fmt.Errorf("service endpoint creation failed: endpoint or ID is nil"))
	}

	d.SetId(serviceEndpoint.Id.String())

	// Return any warnings that were generated during validation
	readDiags := resourceServiceEndpointGenericV2Read(ctx, d, m)
	return append(diags, readDiags...)
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
		// Get current authorization block
		authSet := d.Get("authorization").(*schema.Set)
		if authSet.Len() > 0 {
			authData := authSet.List()[0].(map[string]interface{})

			// Update the scheme if it's changed
			authScheme := *serviceEndpoint.Authorization.Scheme
			if authData["scheme"] != authScheme {
				// We need to recreate the authorization block with the updated scheme
				newAuth := map[string]interface{}{
					"scheme":     authScheme,
					"parameters": authData["parameters"], // Keep existing parameters
				}

				// Convert to a set and set in state
				newAuthSet := schema.NewSet(schema.HashResource(&schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme":     {Type: schema.TypeString},
						"parameters": {Type: schema.TypeMap},
					},
				}), []interface{}{newAuth})

				if err := d.Set("authorization", newAuthSet); err != nil {
					return diag.FromErr(fmt.Errorf("error setting authorization: %v", err))
				}
			}
		}
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

// customizeServiceEndpointGenericV2Diff validates the service endpoint configuration during the planning phase
func customizeServiceEndpointGenericV2Diff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	// Only validate on resource creation changes
	if d.Id() != "" {
		return nil
	}

	serviceEndpointType := d.Get("service_endpoint_type").(string)
	authScheme, authParams, err := getAuthorizationDetailsFromDiff(d)
	if err != nil {
		return fmt.Errorf("error retrieving authorization details: %v", err)
	}

	// Convert additional data to the required format
	data := make(map[string]string)
	if dataRaw := d.Get("data"); dataRaw != nil {
		dataMap, ok := dataRaw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid data format")
		}

		for k, v := range dataMap {
			strVal, ok := v.(string)
			if !ok {
				return fmt.Errorf("data value for key %q is not a string", k)
			}
			data[k] = strVal
		}
	}

	// Prepare endpoint config for validation
	config := EndpointConfig{
		ServiceEndpointType: serviceEndpointType,
		AuthType:            authScheme,
		AuthData:            authParams,
		Data:                data,
	}

	// Validate service endpoint schema - only return actual errors, not warnings
	diags := validateServiceEndpointSchema(ctx, clients, config)
	if diags.HasError() {
		for _, diagValue := range diags {
			if diagValue.Severity == diag.Error {
				return fmt.Errorf(diagValue.Summary + ": " + diagValue.Detail)
			}
		}
	}

	return nil
}

// getAuthorizationDetailsFromDiff extracts authorization details from ResourceDiff
func getAuthorizationDetailsFromDiff(d *schema.ResourceDiff) (string, map[string]string, bool, error) {
	authSet, ok := d.Get("authorization").(*schema.Set)
	if !ok || authSet.Len() == 0 {
		return "", nil, false, fmt.Errorf("no authorization configuration found")
	}

	authData, ok := authSet.List()[0].(map[string]interface{})
	if !ok {
		return "", nil, false, fmt.Errorf("invalid authorization configuration format")
	}

	scheme, ok := authData["scheme"].(string)
	if !ok || scheme == "" {
		return "", nil, false, fmt.Errorf("missing or invalid authorization scheme")
	}

	// Check if parameters block is explicitly configured
	hasParametersBlock := false
	rawConfig := d.GetRawConfig()
	hasParametersBlock = !rawConfig.AsValueMap()["authorization"].GetAttr("parameters").IsNull()

	params := make(map[string]string)
	if paramsRaw, ok := authData["parameters"].(map[string]interface{}); ok {
		for k, v := range paramsRaw {
			strValue, ok := v.(string)
			if !ok {
				return "", nil, hasParametersBlock, fmt.Errorf("parameter %q has invalid type, expected string", k)
			}
			params[k] = strValue
		}
	}

	return scheme, params, hasParametersBlock, nil
}
