package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointNpm() *schema.Resource {
	resource := &schema.Resource{
		Read: dataSourceServiceEndpointNpmRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}

	maps.Copy(resource.Schema, map[string]*schema.Schema{
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
	})

	return resource
}

func dataSourceServiceEndpointNpmRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, projectID, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		doBaseFlattening(d, serviceEndpoint, projectID.String())
		d.Set("url", serviceEndpoint.Url)

		return nil
	}
	return fmt.Errorf(" Looking up service endpoint!")
}
