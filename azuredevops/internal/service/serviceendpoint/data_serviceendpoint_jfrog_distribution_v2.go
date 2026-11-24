package serviceendpoint

import (
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceServiceEndpointJFrogDistributionV2 schema and implementation for JFrog Distribution service endpoint resource
func DataSourceServiceEndpointJFrogDistributionV2() *schema.Resource {
	r := &schema.Resource{
		Read: DataSourceServiceEndpointJFrogDistributionV2Read,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
	maps.Copy(r.Schema, map[string]*schema.Schema{
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
	})

	return r
}

func DataSourceServiceEndpointJFrogDistributionV2Read(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	doBaseFlattening(d, serviceEndpoint)
	d.Set("url", serviceEndpoint.Url)
	return nil
}
