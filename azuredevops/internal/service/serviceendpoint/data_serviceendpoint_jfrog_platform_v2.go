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

// DataSourceServiceEndpointJFrogPlatformV2 schema and implementation for JFrog Platform service endpoint resource
func DataSourceServiceEndpointJFrogPlatformV2() *schema.Resource {
	r := &schema.Resource{
		Read: DataSourceServiceEndpointJFrogPlatformV2Read,
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

func DataSourceServiceEndpointJFrogPlatformV2Read(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
        doBaseFlattening(d, serviceEndpoint)
	d.Set("url", serviceEndpoint.Url)
	return nil
}
