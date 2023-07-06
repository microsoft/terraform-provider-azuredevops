package serviceendpoint

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataServiceEndpointAzureRM() *schema.Resource {
	r := dataSourceGenBaseServiceEndpointResource(dataSourceServiceEndpointAzureRMRead)
	schemaKeys := []string{"azurerm_management_group_id", "azurerm_management_group_name", "azurerm_subscription_id", "azurerm_subscription_name", "resource_group", "azurerm_spn_tenantid", "azurerm_service_endpoint_type"}
	for _, k := range schemaKeys {
		dataSourceMakeUnprotectedComputedSchema(r, k)
	}
	return r
}

func dataSourceServiceEndpointAzureRMRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, projectID, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		(*serviceEndpoint.Data)["creationMode"] = ""
		d.Set("service_endpoint_id", serviceEndpoint.Id.String())
		flattenServiceEndpointAzureRM(d, serviceEndpoint, projectID)
		return nil
	}
	return fmt.Errorf("Error looking up service endpoint!")
}
