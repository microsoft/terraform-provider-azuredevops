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
				Type:     schema.TypeSet,
				Optional: true,
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
	targetProjectID, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("invalid project_id: %v", err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
		ctx,
		serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: &endpointID,
			Project:    converter.String(targetProjectID.String()),
		},
	)
	if err != nil {
		return diag.Errorf("Error finding service endpoint: %+v", err)
	}

	oldSet, newSet := d.GetChange("project_reference")

	projectsToRemove := make(map[string]bool)
	if oldSet != nil {
		for _, raw := range oldSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			pid := obj["project_id"].(string)
			projectsToRemove[strings.ToLower(pid)] = true
		}
	}

	projectsToUpsert := make(map[string]map[string]interface{})
	if newSet != nil {
		for _, raw := range newSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			pid := obj["project_id"].(string)

			delete(projectsToRemove, strings.ToLower(pid))

			projectsToUpsert[strings.ToLower(pid)] = obj
		}
	}

	var newReferences []serviceendpoint.ServiceEndpointProjectReference

	if serviceEndpoint.ServiceEndpointProjectReferences != nil {
		for _, existingRef := range *serviceEndpoint.ServiceEndpointProjectReferences {
			existingPid := strings.ToLower(existingRef.ProjectReference.Id.String())

			if _, shouldRemove := projectsToRemove[existingPid]; shouldRemove {
				continue
			}

			if tfConfig, found := projectsToUpsert[existingPid]; found {
				name := tfConfig["service_endpoint_name"].(string)
				desc := tfConfig["description"].(string)

				existingRef.Name = converter.String(name)
				existingRef.Description = converter.String(desc)

				newReferences = append(newReferences, existingRef)

				delete(projectsToUpsert, existingPid)
			} else {
				newReferences = append(newReferences, existingRef)
			}
		}
	}

	for pid, tfConfig := range projectsToUpsert {
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

	serviceEndpoint.ServiceEndpointProjectReferences = &newReferences
	_, err = clients.ServiceEndpointClient.UpdateServiceEndpoint(
		clients.Ctx,
		serviceendpoint.UpdateServiceEndpointArgs{
			EndpointId: &endpointID,
			Endpoint:   serviceEndpoint,
		},
	)
	if err != nil {
		return diag.Errorf("Error updating service endpoint references: %+v", err)
	}

	d.SetId(endpointID.String())
	return resourceReadServiceEndpointProjectPermissions(clients.Ctx, d, m)
}

func resourceReadServiceEndpointProjectPermissions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	endpointIDStr := d.Get("service_endpoint_id").(string)
	targetProjectIDStr := d.Get("project_id").(string)
	endpointID, err := uuid.Parse(endpointIDStr)
	if err != nil {
		return diag.Errorf("f%s", err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
		ctx,
		serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: &endpointID,
			Project:    converter.String(targetProjectIDStr),
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading service endpoint: %+v", err)
	}

	expectedProjects := make(map[string]bool)
	if set := d.Get("project_reference").(*schema.Set); set != nil {
		for _, raw := range set.List() {
			obj := raw.(map[string]interface{})
			expectedProjects[strings.ToLower(obj["project_id"].(string))] = true
		}
	}

	var flattenedRefs []interface{}

	if serviceEndpoint.ServiceEndpointProjectReferences != nil {
		for _, ref := range *serviceEndpoint.ServiceEndpointProjectReferences {
			pid := strings.ToLower(ref.ProjectReference.Id.String())

			if _, expected := expectedProjects[pid]; expected {
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
	}

	d.Set("project_reference", flattenedRefs)
	return nil
}

func resourceDeleteServiceEndpointProjectPermissions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	endpointID := uuid.MustParse(d.Get("service_endpoint_id").(string))
	targetProjectIDStr := d.Get("project_id").(string)

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
		clients.Ctx,
		serviceendpoint.GetServiceEndpointDetailsArgs{
			EndpointId: &endpointID,
			Project:    converter.String(targetProjectIDStr),
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return diag.Errorf("f%s", err)
	}

	projectsToDelete := make(map[string]bool)
	if set := d.Get("project_reference").(*schema.Set); set != nil {
		for _, raw := range set.List() {
			obj := raw.(map[string]interface{})
			projectsToDelete[strings.ToLower(obj["project_id"].(string))] = true
		}
	}

	var newRefs []serviceendpoint.ServiceEndpointProjectReference
	if serviceEndpoint.ServiceEndpointProjectReferences != nil {
		for _, ref := range *serviceEndpoint.ServiceEndpointProjectReferences {
			pid := strings.ToLower(ref.ProjectReference.Id.String())

			if _, found := projectsToDelete[pid]; !found {
				newRefs = append(newRefs, ref)
			}
		}
	}

	serviceEndpoint.ServiceEndpointProjectReferences = &newRefs

	_, err = clients.ServiceEndpointClient.UpdateServiceEndpoint(
		ctx,
		serviceendpoint.UpdateServiceEndpointArgs{
			EndpointId: &endpointID,
			Endpoint:   serviceEndpoint,
		},
	)

	return diag.Errorf("f%s", err)
}
