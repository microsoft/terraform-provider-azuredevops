package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// DataSourceServiceEndpointJFrogDistributionV2 schema and implementation for JFrog Distribution service endpoint resource
func DataSourceServiceEndpointJFrogDistributionV2() *schema.Resource {
	r := &schema.Resource{
		Read: DataSourceServiceEndpointJFrogDistributionV2Read,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
	maps.Copy(r.Schema, map[string]*schema.Schema{
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
	})

	return r
}

func DataSourceServiceEndpointJFrogDistributionV2Read(d *schema.ResourceData, m interface{}) error {
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
