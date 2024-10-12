package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataServiceEndpointAzureRM() *schema.Resource {
	resource := &schema.Resource{
		Read: dataSourceServiceEndpointAzureRMRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}

	maps.Copy(resource.Schema, map[string]*schema.Schema{
		"azurerm_management_group_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"azurerm_management_group_name": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"azurerm_subscription_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"azurerm_subscription_name": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"resource_group": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"azurerm_spn_tenantid": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"service_endpoint_authentication_scheme": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"environment": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"service_principal_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"workload_identity_federation_issuer": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"workload_identity_federation_subject": {
			Type:     schema.TypeString,
			Computed: true,
		},
	})
	return resource
}

func dataSourceServiceEndpointAzureRMRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, projectID, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		(*serviceEndpoint.Data)["creationMode"] = ""
		d.Set("service_endpoint_id", serviceEndpoint.Id.String())
		flattenServiceEndpointAzureRM(d, serviceEndpoint, projectID.String())
		return nil
	}
	return fmt.Errorf(" Looking up service endpoint!")
}
