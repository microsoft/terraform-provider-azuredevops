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
	serviceEndpoint, projectID, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		d.Set("service_endpoint_id", serviceEndpoint.Id.String())
		doBaseFlattening(d, serviceEndpoint, projectID.String())
		return nil
	}
	return fmt.Errorf(" Looking up service endpoint!")
}
