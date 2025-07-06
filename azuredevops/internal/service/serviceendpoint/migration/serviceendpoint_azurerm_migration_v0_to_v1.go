package migration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/state-migration
func ServiceEndpointAzureRmSchemaV0ToV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorization": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"resource_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"azurerm_subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"azurerm_subscription_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"azurerm_management_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"azurerm_management_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"credentials": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"resource_group"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"serviceprincipalid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"serviceprincipalkey": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"serviceprincipalkey_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"environment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func ServiceEndpointAzureRmStateUpgradeV0ToV1() schema.StateUpgradeFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		if _, ok := rawState["environment"]; !ok {
			rawState["environment"] = "AzureCloud"
		}

		return rawState, nil
	}
}
