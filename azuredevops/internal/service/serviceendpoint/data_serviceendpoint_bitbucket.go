package serviceendpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointBitbucket() *schema.Resource {
	resource := &schema.Resource{
		Read: dataSourceServiceEndpointBitbucketRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
	return resource
}

func dataSourceServiceEndpointBitbucketRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)

	if err != nil {
		return err
	}

	if serviceEndpoint != nil && serviceEndpoint.Id != nil {
		if err = checkServiceConnection(serviceEndpoint); err != nil {
			return err
		}
		doBaseFlattening(d, serviceEndpoint)
		d.Set("service_endpoint_id", serviceEndpoint.Url)
		return nil
	}

	return fmt.Errorf("Looking up Bitbucket Service Endpoint!")
}
