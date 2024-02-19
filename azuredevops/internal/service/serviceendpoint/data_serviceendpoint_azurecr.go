package serviceendpoint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointAzureCR() *schema.Resource {
	r := dataSourceGenBaseServiceEndpointResource(dataSourceServiceEndpointAzureCRRead)

	r.Schema["azurecr_spn_tenantid"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["azurecr_subscription_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["azurecr_subscription_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["resource_group"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["azurecr_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["app_object_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["spn_object_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["az_spn_role_assignment_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["az_spn_role_permissions"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	r.Schema["service_principal_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return r
}

func dataSourceServiceEndpointAzureCRRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, projectID, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}
	if serviceEndpoint != nil {
		doBaseFlattening(d, serviceEndpoint, projectID.String())
		d.Set("azurecr_spn_tenantid", (*serviceEndpoint.Authorization.Parameters)["tenantId"])
		d.Set("azurecr_subscription_id", (*serviceEndpoint.Data)["subscriptionId"])
		d.Set("azurecr_subscription_name", (*serviceEndpoint.Data)["subscriptionName"])

		d.Set("app_object_id", (*serviceEndpoint.Data)["appObjectId"])
		d.Set("spn_object_id", (*serviceEndpoint.Data)["spnObjectId"])
		d.Set("az_spn_role_permissions", (*serviceEndpoint.Data)["azureSpnPermissions"])
		d.Set("az_spn_role_assignment_id", (*serviceEndpoint.Data)["azureSpnRoleAssignmentId"])
		d.Set("service_principal_id", (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"])

		scope := (*serviceEndpoint.Authorization.Parameters)["scope"]
		s := strings.SplitN(scope, "/", -1)
		d.Set("resource_group", s[4])
		d.Set("azurecr_name", s[8])

		return nil
	}
	return fmt.Errorf(" looking up service endpoint!")
}
