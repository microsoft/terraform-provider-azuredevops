package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	crud "github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
)

func resourceServiceEndpointAzureRM() *schema.Resource {
	r := crud.GenBaseServiceEndpointResource(flattenServiceEndpointAzureRM, expandServiceEndpointAzureRM, parseImportedProjectIDAndServiceEndpointID)
	crud.MakeUnprotectedSchema(r, "azurerm_spn_clientid", "ARM_CLIENT_ID", "The service principal id which should be used.")
	crud.MakeProtectedSchema(r, "azurerm_spn_clientsecret", "ARM_CLIENT_SECRET", "The service principal secret which should be used.")
	crud.MakeUnprotectedSchema(r, "azurerm_spn_tenantid", "ARM_TENANT_ID", "The service principal tenant id which should be used.")
	crud.MakeUnprotectedSchema(r, "azurerm_subscription_id", "ARM_SUBSCRIPTION_ID", "The Azure subscription Id which should be used.")
	crud.MakeUnprotectedSchema(r, "azurerm_subscription_name", "ARM_SUBSCRIPTION_NAME", "The Azure subscription name which should be used.")
	crud.MakeUnprotectedSchema(r, "azurerm_scope", "ARM_SCOPE", "The Azure scope which should be used by the spn.")
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAzureRM(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	serviceEndpoint, projectID := crud.DoBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"authenticationType":  "spnKey",
			"scope":               d.Get("azurerm_scope").(string),
			"serviceprincipalid":  d.Get("azurerm_spn_clientid").(string),
			"serviceprincipalkey": d.Get("azurerm_spn_clientsecret").(string),
			"tenantid":            d.Get("azurerm_spn_tenantid").(string),
		},
		Scheme: converter.String("ServicePrincipal"),
	}
	serviceEndpoint.Data = &map[string]string{
		"creationMode":     "Manual",
		"environment":      "AzureCloud",
		"scopeLevel":       "Subscription",
		"SubscriptionId":   d.Get("azurerm_subscription_id").(string),
		"SubscriptionName": d.Get("azurerm_subscription_name").(string),
	}
	serviceEndpoint.Type = converter.String("azurerm")
	serviceEndpoint.Url = converter.String("https://management.azure.com/")
	return serviceEndpoint, projectID
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAzureRM(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	crud.DoBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("azurerm_scope", (*serviceEndpoint.Authorization.Parameters)["scope"])
	d.Set("azurerm_spn_clientid", (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"])
	tfhelper.HelpFlattenSecret(d, "azurerm_spn_clientsecret")
	d.Set("azurerm_spn_tenantid", (*serviceEndpoint.Authorization.Parameters)["tenantid"])
	d.Set("azurerm_spn_clientsecret", (*serviceEndpoint.Authorization.Parameters)["serviceprincipalkey"])
	d.Set("azurerm_subscription_id", (*serviceEndpoint.Data)["SubscriptionId"])
	d.Set("azurerm_subscription_name", (*serviceEndpoint.Data)["SubscriptionName"])
}
