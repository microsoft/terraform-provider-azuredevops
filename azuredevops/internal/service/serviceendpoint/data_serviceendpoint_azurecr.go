package serviceendpoint

import (
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointAzureCR() *schema.Resource {
	resource := &schema.Resource{
		Read: dataSourceServiceEndpointAzureCRRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}
	maps.Copy(resource.Schema, map[string]*schema.Schema{
		"azurecr_spn_tenantid": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"azurecr_subscription_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"azurecr_subscription_name": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"resource_group": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"azurecr_name": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"app_object_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"spn_object_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"az_spn_role_assignment_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"az_spn_role_permissions": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"service_principal_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	})

	return resource
}

func dataSourceServiceEndpointAzureCRRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil && serviceEndpoint.Id != nil {
		if err = checkServiceConnection(serviceEndpoint); err != nil {
			return err
		}
		doBaseFlattening(d, serviceEndpoint)
		d.Set("azurecr_spn_tenantid", (*serviceEndpoint.Authorization.Parameters)["tenantId"])
		d.Set("azurecr_subscription_id", (*serviceEndpoint.Data)["subscriptionId"])
		d.Set("azurecr_subscription_name", (*serviceEndpoint.Data)["subscriptionName"])

		d.Set("app_object_id", (*serviceEndpoint.Data)["appObjectId"])
		d.Set("spn_object_id", (*serviceEndpoint.Data)["spnObjectId"])
		d.Set("az_spn_role_permissions", (*serviceEndpoint.Data)["azureSpnPermissions"])
		d.Set("az_spn_role_assignment_id", (*serviceEndpoint.Data)["azureSpnRoleAssignmentId"])
		d.Set("service_principal_id", (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"])

		scope := (*serviceEndpoint.Authorization.Parameters)["scope"]
		s := strings.Split(scope, "/")
		d.Set("resource_group", s[4])
		d.Set("azurecr_name", s[8])

		return nil
	}
	return fmt.Errorf("Looking up service endpoint!")
}
