package serviceendpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointSonarCloud() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceEndpointSonarCloudRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
}

func dataSourceServiceEndpointSonarCloudRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, projectID, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		doBaseFlattening(d, serviceEndpoint, projectID.String())

		return nil
	}
	return fmt.Errorf(" Looking up Sonar Cloud service endpoint !")
}
