package serviceendpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataServiceEndpointGeneric() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceEndpointGenericRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
}

func dataSourceServiceEndpointGenericRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil && serviceEndpoint.Id != nil {
		if err = checkServiceConnection(serviceEndpoint); err != nil {
			return err
		}
		doBaseFlattening(d, serviceEndpoint)
		d.Set("service_endpoint_id", serviceEndpoint.Id.String())
		return nil
	}
	return fmt.Errorf("Looking up service endpoint!")
}
