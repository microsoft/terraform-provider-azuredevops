package serviceendpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	SE "github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

// DataServiceEndpointGenericV2 schema and implementation for generic service endpoint data source
func DataServiceEndpointGenericV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataServiceEndpointGenericV2Read,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"shared_project_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Use id to look up by ID
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"id", "name"},
			},
			// Use name to look up by name
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ExactlyOneOf: []string{"id", "name"},
			},
			// Read-only fields that are returned
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_scheme": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_parameters": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataServiceEndpointGenericV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)

	var serviceEndpointID string
	var err error

	// Look up service endpoint by ID if provided
	if id, ok := d.GetOk("id"); ok {
		serviceEndpointID = id.(string)
	} else {
		// Look up service endpoint by name if ID not provided
		serviceEndpointName := d.Get("name").(string)
		serviceEndpointID, err = findServiceEndpointByName(ctx, clients, serviceEndpointName, projectID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	serviceEndpoint, err := getServiceEndpointGenericV2(ctx, clients, serviceEndpointID, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return diag.FromErr(fmt.Errorf("service endpoint with ID %s does not exist in project %s", serviceEndpointID, projectID))
		}
		return diag.FromErr(err)
	}

	if serviceEndpoint == nil {
		return diag.FromErr(fmt.Errorf("service endpoint with ID %s does not exist in project %s", serviceEndpointID, projectID))
	}

	// Set the ID and computed properties
	d.SetId(serviceEndpointID)

	if serviceEndpoint.Name != nil {
		d.Set("name", *serviceEndpoint.Name)
	}

	if serviceEndpoint.Type != nil {
		d.Set("type", *serviceEndpoint.Type)
	}

	if serviceEndpoint.Description != nil {
		d.Set("description", *serviceEndpoint.Description)
	}

	if serviceEndpoint.Url != nil {
		d.Set("server_url", *serviceEndpoint.Url)
	}

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Scheme != nil {
		d.Set("authorization_scheme", *serviceEndpoint.Authorization.Scheme)
	}

	// Set data fields
	if serviceEndpoint.Data != nil {
		data := make(map[string]string)
		for k, v := range *serviceEndpoint.Data {
			data[k] = v
		}
		if len(data) > 0 {
			d.Set("data", data)
		}
	}

	// Set authorization fields
	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		authorization := make(map[string]string)
		for k, v := range *serviceEndpoint.Authorization.Parameters {
			authorization[k] = v
		}
		if len(authorization) > 0 {
			d.Set("authorization_parameters", authorization)
		}
	}

	// Populate shared_project_ids
	var sharedProjectIDs []string
	for _, ref := range *serviceEndpoint.ServiceEndpointProjectReferences {
		if ref.ProjectReference != nil && ref.ProjectReference.Id != nil &&
			ref.ProjectReference.Id.String() != d.Get("project_id").(string) {
			sharedProjectIDs = append(sharedProjectIDs, ref.ProjectReference.Id.String())
		}
	}
	err = d.Set("shared_project_ids", sharedProjectIDs)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting shared_project_ids: %w", err))
	}

	return nil
}

// findServiceEndpointByName retrieves a service endpoint ID by name
func findServiceEndpointByName(ctx context.Context, clients *client.AggregatedClient, endpointName, projectID string) (string, error) {
	serviceEndpoints, err := clients.ServiceEndpointClient.GetServiceEndpoints(ctx, SE.GetServiceEndpointsArgs{
		Project: &projectID,
	})
	if err != nil {
		return "", fmt.Errorf("error looking up service endpoints in project: %v", err)
	}

	if serviceEndpoints == nil {
		return "", fmt.Errorf("no service endpoints found in project %s", projectID)
	}

	for _, endpoint := range *serviceEndpoints {
		if endpoint.Name != nil && *endpoint.Name == endpointName && endpoint.Id != nil {
			return endpoint.Id.String(), nil
		}
	}

	return "", fmt.Errorf("service endpoint with name '%s' not found in project %s", endpointName, projectID)
}
