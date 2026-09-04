package serviceendpoint

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceServiceEndpointProjectPermissions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateOrUpdateServiceEndpointProjectPermissions,
		ReadContext:   resourceReadServiceEndpointProjectPermissions,
		UpdateContext: resourceCreateOrUpdateServiceEndpointProjectPermissions,
		DeleteContext: resourceDeleteServiceEndpointProjectPermissions,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_endpoint_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the project where the service endpoint is created (Source Project).",
			},
			"project_reference": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"service_endpoint_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceCreateOrUpdateServiceEndpointProjectPermissions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	endpointID, err := uuid.Parse(d.Get("service_endpoint_id").(string))
	if err != nil {
		return diag.Errorf("invalid service_endpoint_id: %v", err)
	}
	sourceProjectID, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("invalid project_id: %v", err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
		ctx,
		serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: &endpointID,
			Project:    converter.String(sourceProjectID.String()),
		},
	)
	if err != nil {
		return diag.Errorf("Error finding service endpoint: %+v", err)
	}

	plannedProjectReferences := d.Get("project_reference").([]interface{})
	plannedProjects := make(map[string]map[string]interface{})
	for _, raw := range plannedProjectReferences {
		obj := raw.(map[string]interface{})
		pid := strings.ToLower(obj["project_id"].(string))
		plannedProjects[pid] = obj
	}

	var projectsToRemove []string
	if serviceEndpoint.ServiceEndpointProjectReferences != nil {
		for _, ref := range *serviceEndpoint.ServiceEndpointProjectReferences {
			pid := strings.ToLower(ref.ProjectReference.Id.String())
			// Don't remove the source project
			if strings.EqualFold(pid, sourceProjectID.String()) {
				continue
			}
			// If it's not in the planned projects, remove it
			if _, ok := plannedProjects[pid]; !ok {
				projectsToRemove = append(projectsToRemove, pid)
			}
		}
	}

	// 1. Delete removed project references
	if len(projectsToRemove) > 0 {
		err = clients.ServiceEndpointClient.DeleteServiceEndpoint(ctx, serviceendpoint.DeleteServiceEndpointArgs{
			EndpointId: &endpointID,
			ProjectIds: &projectsToRemove,
		})
		if err != nil {
			return diag.Errorf("Error removing service endpoint project permissions: %+v", err)
		}
	}

	// 2. Upsert planned project references
	if len(plannedProjects) > 0 {
		var newReferences []serviceendpoint.ServiceEndpointProjectReference
		for pid, tfConfig := range plannedProjects {
			targetProjectID := uuid.MustParse(pid)
			name := tfConfig["service_endpoint_name"].(string)
			desc := tfConfig["description"].(string)

			newReferences = append(newReferences, serviceendpoint.ServiceEndpointProjectReference{
				ProjectReference: &serviceendpoint.ProjectReference{
					Id: &targetProjectID,
				},
				Name:        converter.String(name),
				Description: converter.String(desc),
			})
		}

		err = clients.ServiceEndpointClient.ShareServiceEndpoint(ctx, serviceendpoint.ShareServiceEndpointArgs{
			EndpointId:                &endpointID,
			EndpointProjectReferences: &newReferences,
		})
		if err != nil {
			return diag.Errorf("Error sharing service endpoint to projects: %+v", err)
		}
	}

	d.SetId(endpointID.String())
	return resourceReadServiceEndpointProjectPermissions(ctx, d, m)
}

func resourceReadServiceEndpointProjectPermissions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	endpointIDStr := d.Get("service_endpoint_id").(string)
	sourceProjectIDStr := d.Get("project_id").(string)
	endpointID, err := uuid.Parse(endpointIDStr)
	if err != nil {
		return diag.Errorf("invalid service_endpoint_id: %s", err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
		ctx,
		serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: &endpointID,
			Project:    converter.String(sourceProjectIDStr),
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading service endpoint: %+v", err)
	}

	var flattenedRefs []interface{}

	if serviceEndpoint.ServiceEndpointProjectReferences != nil {
		for _, ref := range *serviceEndpoint.ServiceEndpointProjectReferences {
			pid := strings.ToLower(ref.ProjectReference.Id.String())

			// Skip the source project
			if strings.EqualFold(pid, sourceProjectIDStr) {
				continue
			}

			item := map[string]interface{}{
				"project_id": ref.ProjectReference.Id.String(),
			}
			if ref.Name != nil {
				item["service_endpoint_name"] = *ref.Name
			}
			if ref.Description != nil {
				item["description"] = *ref.Description
			}
			flattenedRefs = append(flattenedRefs, item)
		}
	}

	d.Set("project_reference", flattenedRefs)
	return nil
}

func resourceDeleteServiceEndpointProjectPermissions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	endpointID := uuid.MustParse(d.Get("service_endpoint_id").(string))
	sourceProjectIDStr := d.Get("project_id").(string)

	_, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
		ctx,
		serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: &endpointID,
			Project:    converter.String(sourceProjectIDStr),
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return diag.Errorf("Error reading service endpoint: %+v", err)
	}

	var projectsToDelete []string
	if list := d.Get("project_reference").([]interface{}); list != nil {
		for _, raw := range list {
			obj := raw.(map[string]interface{})
			pid := obj["project_id"].(string)
			if !strings.EqualFold(pid, sourceProjectIDStr) {
				projectsToDelete = append(projectsToDelete, pid)
			}
		}
	}

	if len(projectsToDelete) > 0 {
		err = clients.ServiceEndpointClient.DeleteServiceEndpoint(ctx, serviceendpoint.DeleteServiceEndpointArgs{
			EndpointId: &endpointID,
			ProjectIds: &projectsToDelete,
		})
		if err != nil {
			return diag.Errorf("Error deleting service endpoint project permissions: %+v", err)
		}
	}

	return nil
}
