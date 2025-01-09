package service

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataClientConfig schema and implementation for AzDO client configuration
func DataClientConfig() *schema.Resource {
	return &schema.Resource{
		Read: clientConfigRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"organization_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func clientConfigRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(time.Now().UTC().String())
	d.Set("organization_url", m.(*client.AggregatedClient).OrganizationURL)
	return nil
}
