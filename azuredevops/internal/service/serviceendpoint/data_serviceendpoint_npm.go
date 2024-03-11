package serviceendpoint

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointNpm() *schema.Resource {
	r := dataSourceGenBaseServiceEndpointResource(dataSourceServiceEndpointNpmRead)

	r.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return r
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
	return fmt.Errorf("Error looking up service endpoint!")
}
