package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointAquaPlatform() *schema.Resource {
	resource := &schema.Resource{
		Read: dataSourceServiceEndpointAquaPlatformRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}

	maps.Copy(resource.Schema, map[string]*schema.Schema{
		"aqua_platform_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aqua_auth_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
	})

	return resource
}

func dataSourceServiceEndpointAquaPlatformRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}

	if serviceEndpoint != nil && serviceEndpoint.Id != nil {
		if err = checkServiceConnection(serviceEndpoint); err != nil {
			return err
		}
		doBaseFlattening(d, serviceEndpoint)
		flattenServiceEndpointAquaPlatform(d, serviceEndpoint)
		return nil
	}

	return fmt.Errorf("Looking up service endpoint!")
}
