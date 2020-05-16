package azuredevops

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

const (
	organizationURL = "organization_url"
)

func dataClientConfig() *schema.Resource {
	return &schema.Resource{
		Read: clientConfigRead,
		Schema: map[string]*schema.Schema{
			organizationURL: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func clientConfigRead(d *schema.ResourceData, m interface{}) error {
	// The ID is meaningless for this data source, so ID can act as a
	// point in time snapshot
	d.SetId(time.Now().UTC().String())
	d.Set(organizationURL, m.(*config.AggregatedClient).OrganizationURL)
	return nil
}
