package serviceendpoint

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

var (
	serviceEndpointTypesCache     *[]serviceendpoint.ServiceEndpointType
	serviceEndpointTypesOnceCache sync.Once
	serviceEndpointTypesCacheErr  error
)

// DataServiceEndpointTypes schema and implementation for service endpoint types data source
func DataServiceEndpointTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataServiceEndpointTypesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the service endpoint type",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the service endpoint type",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the service endpoint type",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the service endpoint type",
						},
						"ui_contribution_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UI contribution ID for this service endpoint type",
						},
						"authentication_schemes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available authentication schemes for this service endpoint type",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataServiceEndpointTypesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	serviceEndpointTypesOnceCache.Do(func() {
		args := serviceendpoint.GetServiceEndpointTypesArgs{}
		serviceEndpointTypesCache, serviceEndpointTypesCacheErr = clients.ServiceEndpointClient.GetServiceEndpointTypes(clients.Ctx, args)
	})

	if serviceEndpointTypesCacheErr != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("querying service endpoint types: %v", serviceEndpointTypesCacheErr))
	}

	if serviceEndpointTypesCache == nil || len(*serviceEndpointTypesCache) == 0 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("no service endpoint types found"))
	}

	flattenedTypes := flattenServiceEndpointTypes(serviceEndpointTypesCache)
	d.SetId(uuid.New().String())
	err := d.Set("types", flattenedTypes)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenServiceEndpointTypes(types *[]serviceendpoint.ServiceEndpointType) []interface{} {
	if types == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, t := range *types {
		serviceType := map[string]interface{}{}

		if t.Name != nil {
			serviceType["name"] = *t.Name
			// Use name as id since there's no dedicated Id field
			serviceType["id"] = *t.Name
		}
		if t.DisplayName != nil {
			serviceType["display_name"] = *t.DisplayName
		}
		if t.Description != nil {
			serviceType["description"] = *t.Description
		}
		if t.UiContributionId != nil {
			serviceType["ui_contribution_id"] = *t.UiContributionId
		}

		// Flatten authentication schemes
		authSchemes := make([]interface{}, 0)
		if t.AuthenticationSchemes != nil {
			for _, authScheme := range *t.AuthenticationSchemes {
				if authScheme.Scheme != nil {
					authSchemes = append(authSchemes, *authScheme.Scheme)
				}
			}
		}
		serviceType["authentication_schemes"] = authSchemes

		results = append(results, serviceType)
	}
	return results
}
