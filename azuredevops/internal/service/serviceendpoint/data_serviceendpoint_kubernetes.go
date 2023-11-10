package serviceendpoint

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointKubernetes() *schema.Resource {
	return dataSourceGenBaseServiceEndpointResource(dataSourceServiceEndpointKubernetesRead)
}

func dataSourceServiceEndpointKubernetesRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, projectID, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		doBaseFlattening(d, serviceEndpoint, projectID)

		return nil
	}
	return fmt.Errorf("Error looking up Kubernetes service endpoint !")
}
