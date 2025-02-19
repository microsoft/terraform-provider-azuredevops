package service

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataClientConfig schema and implementation for AzDO client configuration
func DataClientConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: clientConfigRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"organization_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func clientConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	parts := strings.Split(m.(*client.AggregatedClient).OrganizationURL, "/")

	orgMeta, err := clients.OrganizationClient.GetOrganization(clients.Ctx, parts[3])
	if err != nil {
		return diag.Errorf(" Getting organization metadata: %s", err)
	}

	d.SetId(*orgMeta.Id)
	d.Set("organization_url", m.(*client.AggregatedClient).OrganizationURL)
	d.Set("status", orgMeta.Status)
	d.Set("name", orgMeta.Name)
	d.Set("tenant_id", orgMeta.TenantId)
	d.Set("owner_id", orgMeta.Owner)
	return nil
}
