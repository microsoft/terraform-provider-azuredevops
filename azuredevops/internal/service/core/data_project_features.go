package core

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func DataProjectFeatures() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectFeaturesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"features": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataProjectFeaturesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	id := d.Get("project_id").(string)

	identifier := id
	if identifier == "" {
		identifier = name
	}

	projectID := d.Get("project_id").(string)
	featureStates := d.Get("features").(map[string]interface{})
	currentFeatureStates, err := getConfiguredProjectFeatureStates(ctx, clients.FeatureManagementClient, &featureStates, projectID)
	if err != nil {
		return diag.FromErr(err)
	}
	if currentFeatureStates == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to retrieve current feature states for project: %s", projectID))
	}

	// Convert the upstream feature states to the expected map format
	upstreamFeatures := make(map[string]interface{})
	for featureType, enabledValue := range *currentFeatureStates {
		upstreamFeatures[string(featureType)] = string(enabledValue)
	}

	d.SetId(projectID)
	d.Set("features", upstreamFeatures)
	return nil
}
