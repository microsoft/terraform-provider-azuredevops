package serviceendpoint

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointAzureRM schema and implementation for AzureRM service endpoint resource
func ResourceServiceEndpointAzureRM() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAzureRM, expandServiceEndpointAzureRM)
	makeUnprotectedSchema(r, "azurerm_spn_tenantid", "ARM_TENANT_ID", "The service principal tenant id which should be used.")
	makeUnprotectedSchema(r, "azurerm_subscription_id", "ARM_SUBSCRIPTION_ID", "The Azure subscription Id which should be used.")
	makeUnprotectedSchema(r, "azurerm_subscription_name", "ARM_SUBSCRIPTION_NAME", "The Azure subscription name which should be used.")

	r.Schema["resource_group"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		Description:   "Scope Resource Group",
		ConflictsWith: []string{"credentials"},
	}

	secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema("serviceprincipalkey")
	r.Schema["credentials"] = &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"resource_group"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serviceprincipalid": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The service principal id which should be used.",
				},
				"serviceprincipalkey": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "The service principal secret which should be used.",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
				},
				secretHashKey: secretHashSchema,
			},
		},
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAzureRM(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)

	scope := fmt.Sprintf("/subscriptions/%s", d.Get("azurerm_subscription_id"))
	scopeLevel := "Subscription"
	if _, ok := d.GetOk("resource_group"); ok {
		scope += fmt.Sprintf("/resourcegroups/%s", d.Get("resource_group"))
		scopeLevel = "ResourceGroup"
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"authenticationType":  "spnKey",
			"serviceprincipalid":  "",
			"serviceprincipalkey": "",
			"tenantid":            d.Get("azurerm_spn_tenantid").(string),
		},
		Scheme: converter.String("ServicePrincipal"),
	}
	serviceEndpoint.Data = &map[string]string{
		"creationMode":     "Automatic",
		"environment":      "AzureCloud",
		"scopeLevel":       "Subscription",
		"subscriptionId":   d.Get("azurerm_subscription_id").(string),
		"subscriptionName": d.Get("azurerm_subscription_name").(string),
	}

	if scopeLevel == "ResourceGroup" {
		(*serviceEndpoint.Authorization.Parameters)["scope"] = scope
	}

	if _, ok := d.GetOk("credentials"); ok {
		credentials := d.Get("credentials").([]interface{})[0].(map[string]interface{})
		(*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"] = credentials["serviceprincipalid"].(string)
		(*serviceEndpoint.Authorization.Parameters)["serviceprincipalkey"] = expandSpnKey(credentials)
		(*serviceEndpoint.Data)["creationMode"] = "Manual"
	}

	serviceEndpoint.Type = converter.String("azurerm")
	serviceEndpoint.Url = converter.String("https://management.azure.com/")
	return serviceEndpoint, projectID, nil
}

func expandSpnKey(credentials map[string]interface{}) string {
	// Note: if this is an update for a field other than `serviceprincipalkey`, the `serviceprincipalkey` will be
	// set to `""`. Without catching this case and setting the value to `"null"`, the `serviceprincipalkey` will
	// actually be set to `""` by the Azure DevOps service.
	//
	// This step is critical in order to ensure that the service connection can update without loosing its password!
	//
	// This behavior is unfortunately not documented in the API documentation.
	spnKey, ok := credentials["serviceprincipalkey"]
	if !ok || spnKey.(string) == "" {
		return "null"
	}

	return spnKey.(string)
}

func flattenCredentials(serviceEndpoint *serviceendpoint.ServiceEndpoint, hashKey string, hashValue string) interface{} {
	return []map[string]interface{}{{
		"serviceprincipalid":  (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"],
		"serviceprincipalkey": (*serviceEndpoint.Authorization.Parameters)["serviceprincipalkey"],
		hashKey:               hashValue,
	}}
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAzureRM(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	scope := (*serviceEndpoint.Authorization.Parameters)["scope"]

	if (*serviceEndpoint.Data)["creationMode"] == "Manual" {
		newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "credentials", d.Get("credentials.0").(map[string]interface{}), "serviceprincipalkey")
		credentials := flattenCredentials(serviceEndpoint, hashKey, newHash)
		d.Set("credentials", credentials)
	}

	s := strings.SplitN(scope, "/", -1)
	if len(s) == 5 {
		d.Set("resource_group", s[4])
	}

	d.Set("azurerm_spn_tenantid", (*serviceEndpoint.Authorization.Parameters)["tenantid"])
	d.Set("azurerm_subscription_id", (*serviceEndpoint.Data)["subscriptionId"])
	d.Set("azurerm_subscription_name", (*serviceEndpoint.Data)["subscriptionName"])
}
