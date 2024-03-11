package serviceendpoint

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointSonarCloud() *schema.Resource {
	return dataSourceGenBaseServiceEndpointResource(dataSourceServiceEndpointSonarCloudRead)
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
	return fmt.Errorf("Error looking up Sonar Cloud service endpoint !")
}
