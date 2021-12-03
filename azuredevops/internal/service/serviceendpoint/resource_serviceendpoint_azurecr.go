package serviceendpoint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointAzureCR schema and implementation for ACR service endpoint resource
func ResourceServiceEndpointAzureCR() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAzureCR, expandServiceEndpointAzureCR)
	makeUnprotectedSchema(r, "azurecr_spn_tenantid", "ACR_TENANT_ID", "The service principal tenant id which should be used.")
	makeUnprotectedSchema(r, "azurecr_subscription_id", "ACR_SUBSCRIPTION_ID", "The Azure subscription Id which should be used.")
	makeUnprotectedSchema(r, "azurecr_subscription_name", "ACR_SUBSCRIPTION_NAME", "The Azure subscription name which should be used.")

	r.Schema["resource_group"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Scope Resource Group",
	}

	r.Schema["azurecr_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_AZURECR_SERVICE_CONNECTION_REGISTRY", nil),
		Description: "The AzureContainerRegistry registry which should be used.",
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

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAzureCR(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	subscriptionID := d.Get("azurecr_subscription_id")
	scope := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s",
		subscriptionID, d.Get("resource_group"), d.Get("azurecr_name"),
	)
	loginServer := fmt.Sprintf("%s.azurecr.io", strings.ToLower(d.Get("azurecr_name").(string)))
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"authenticationType": "spnKey",
			"tenantId":           d.Get("azurecr_spn_tenantid").(string),
			"loginServer":        loginServer,
			"scope":              scope,
			"serviceprincipalid": d.Get("service_principal_id").(string),
		},
		Scheme: converter.String("ServicePrincipal"),
	}
	serviceEndpoint.Data = &map[string]string{
		"registryId":               scope,
		"subscriptionId":           subscriptionID.(string),
		"subscriptionName":         d.Get("azurecr_subscription_name").(string),
		"registrytype":             "ACR",
		"appObjectId":              d.Get("app_object_id").(string),
		"spnObjectId":              d.Get("spn_object_id").(string),
		"azureSpnPermissions":      d.Get("az_spn_role_permissions").(string),
		"azureSpnRoleAssignmentId": d.Get("az_spn_role_assignment_id").(string),
	}
	serviceEndpoint.Type = converter.String("dockerregistry")
	azureContainerRegistryURL := fmt.Sprintf("https://%s", loginServer)
	serviceEndpoint.Url = converter.String(azureContainerRegistryURL)

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAzureCR(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
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
}
