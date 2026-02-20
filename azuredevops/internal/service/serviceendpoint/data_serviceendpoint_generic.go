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
		Schema: dataSourceServiceEndpointGenericSchema(),
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
		if serviceEndpoint.Authorization != nil {
			auth := make(map[string]interface{})
			if serviceEndpoint.Authorization.Scheme != nil {
				auth["scheme"] = *serviceEndpoint.Authorization.Scheme
			}
			if serviceEndpoint.Authorization.Parameters != nil {
				for k, v := range *serviceEndpoint.Authorization.Parameters {
					auth[k] = v
				}
			}
			d.Set("authorization", auth)
		}

		if serviceEndpoint.Data != nil {
			d.Set("data", *serviceEndpoint.Data)
		}

		d.Set("service_endpoint_id", serviceEndpoint.Id.String())
		return nil
	}
	return fmt.Errorf("Looking up service endpoint!")
}

func dataSourceServiceEndpointGenericSchema() map[string]*schema.Schema {
	d := dataSourceGenBaseSchema()
	d["data"] = &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return d
}
