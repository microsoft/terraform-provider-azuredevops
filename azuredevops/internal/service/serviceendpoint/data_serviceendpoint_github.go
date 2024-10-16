package serviceendpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataServiceEndpointGithub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceEndpointGithubRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
}

func dataSourceServiceEndpointGithubRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		if err = checkServiceConnection(serviceEndpoint); err != nil {
			return err
		}
		doBaseFlattening(d, serviceEndpoint)
		d.Set("service_endpoint_id", serviceEndpoint.Id.String())
		return nil
	}
	return fmt.Errorf(" Looking up service endpoint!")
}
