package serviceendpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func DataResourceServiceEndpointJFrogPlatformV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceEndpointJFrogPlatformV2Read,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
}

func dataSourceServiceEndpointJFrogPlatformV2Read(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointArtifactoryV2(d, serviceEndpoint)
	return nil
}
