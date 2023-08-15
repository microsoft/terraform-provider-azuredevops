package serviceendpoint

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/serviceendpoint/migration"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointAzureRM schema and implementation for AzureRM service endpoint resource
func ResourceServiceEndpointAzureRM() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAzureRM, expandServiceEndpointAzureRM)
	makeUnprotectedSchema(r, "azurerm_spn_tenantid", "ARM_TENANT_ID", "The service principal tenant id which should be used.")

	r.Schema["resource_group"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		Description:   "Scope Resource Group",
		ConflictsWith: []string{"credentials", "azurerm_management_group_id"},
	}

	// Subscription scopeLevel
	makeUnprotectedOptionalSchema(r, "azurerm_subscription_id", "ARM_SUBSCRIPTION_ID", "The Azure subscription Id which should be used.", []string{"azurerm_management_group_id"})
	makeUnprotectedOptionalSchema(r, "azurerm_subscription_name", "ARM_SUBSCRIPTION_NAME", "The Azure subscription name which should be used.", []string{"azurerm_management_group_id"})

	// ManagementGroup scopeLevel
	makeUnprotectedOptionalSchema(r, "azurerm_management_group_id", "ARM_MGMT_GROUP_ID", "The Azure managementGroup Id which should be used.", []string{"azurerm_subscription_id", "resource_group"})
	makeUnprotectedOptionalSchema(r, "azurerm_management_group_name", "ARM_MGMT_GROUP_NAME", "The Azure managementGroup name which should be used.", []string{"azurerm_subscription_id", "resource_group"})

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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The service principal secret which should be used.",
					Sensitive:   true,
				},
			},
		},
	}
	r.Schema["environment"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Environment (Azure Cloud type)",
		Default:      "AzureCloud",
		ValidateFunc: validation.StringInSlice([]string{"AzureCloud", "AzureChinaCloud"}, false),
	}

	r.Schema["service_endpoint_authentication_scheme"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "The AzureRM Service Endpoint Authentication Scheme, this can be 'WorkloadIdentityFederation', 'ManagedServiceIdentity' or 'ServicePrincipal'.",
		Default:      "ServicePrincipal",
		ValidateFunc: validation.StringInSlice([]string{"WorkloadIdentityFederation", "ManagedServiceIdentity", "ServicePrincipal"}, false),
	}

	r.SchemaVersion = 2
	r.StateUpgraders = []schema.StateUpgrader{
		{
			Type:    migration.ServiceEndpointAzureRmSchemaV0ToV1().CoreConfigSchema().ImpliedType(),
			Upgrade: migration.ServiceEndpointAzureRmStateUpgradeV0ToV1(),
			Version: 0,
		},
		{
			Type:    migration.ServiceEndpointAzureRmSchemaV1ToV2().CoreConfigSchema().ImpliedType(),
			Upgrade: migration.ServiceEndpointAzureRmStateUpgradeV1ToV2(),
			Version: 1,
		},
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAzureRM(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)

	serviceEndPointAuthenticationScheme := AzureRmEndpointAuthenticationScheme(d.Get("service_endpoint_authentication_scheme").(string))

	// NOTE: This is a temporary workaround for a bug in the Azure DevOps API. This will be removed once the API is fixed.
	if serviceEndPointAuthenticationScheme == WorkloadIdentityFederation {
		(*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Name = converter.String("doesntmatter")
	}

	// Validate one of either subscriptionId or managementGroupId is set
	subId := d.Get("azurerm_subscription_id").(string)
	subName := d.Get("azurerm_subscription_name").(string)

	mgmtGrpId := d.Get("azurerm_management_group_id").(string)
	mgmtGrpName := d.Get("azurerm_management_group_name").(string)
	environment := d.Get("environment").(string)

	scopeLevelMap := map[string][]string{
		"subscription":    {subId, subName},
		"managementGroup": {mgmtGrpId, mgmtGrpName},
	}

	if err := validateScopeLevel(scopeLevelMap); err != nil {
		return nil, nil, err
	}

	var scope string
	var scopeLevel string

	serviceEndPointAuthenticationSchemeHasCreationMode := serviceEndPointAuthenticationScheme == ServicePrincipal || serviceEndPointAuthenticationScheme == WorkloadIdentityFederation

	if _, ok := d.GetOk("azurerm_subscription_id"); ok {
		scope = fmt.Sprintf("/subscriptions/%s", d.Get("azurerm_subscription_id"))
		scopeLevel = "Subscription"
		if serviceEndPointAuthenticationSchemeHasCreationMode {
			if _, ok := d.GetOk("resource_group"); ok {
				scope += fmt.Sprintf("/resourcegroups/%s", d.Get("resource_group"))
				scopeLevel = "ResourceGroup"
			}
		}
	}

	var credentials map[string]interface{}

	if _, ok := d.GetOk("credentials"); ok {
		credentials = d.Get("credentials").([]interface{})[0].(map[string]interface{})
	}

	hasCredentials := credentials != nil && len(credentials) > 0

	var serviceEndpointCreationMode AzureRmEndpointCreationMode

	if serviceEndPointAuthenticationSchemeHasCreationMode {
		if hasCredentials {
			serviceEndpointCreationMode = Manual
		} else {
			serviceEndpointCreationMode = Automatic
		}
	}

	switch serviceEndPointAuthenticationScheme {
	case ServicePrincipal:
		if serviceEndpointCreationMode == Automatic {
			serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
				Parameters: &map[string]string{
					"authenticationType":  "spnKey",
					"serviceprincipalid":  "",
					"serviceprincipalkey": "",
					"tenantid":            d.Get("azurerm_spn_tenantid").(string),
				},
				Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
			}
		}
		if serviceEndpointCreationMode == Manual {
			serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
				Parameters: &map[string]string{
					"authenticationType":  "spnKey",
					"serviceprincipalid":  credentials["serviceprincipalid"].(string),
					"serviceprincipalkey": credentials["serviceprincipalkey"].(string),
					"tenantid":            d.Get("azurerm_spn_tenantid").(string),
				},
				Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
			}
		}

		serviceEndpoint.Data = &map[string]string{
			"creationMode": string(serviceEndpointCreationMode),
			"environment":  environment,
		}

	case ManagedServiceIdentity:
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"tenantid": d.Get("azurerm_spn_tenantid").(string),
			},
			Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
		}

		serviceEndpoint.Data = &map[string]string{
			"environment": environment,
		}

	case WorkloadIdentityFederation:
		if serviceEndpointCreationMode == Automatic {
			serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
				Parameters: &map[string]string{
					"serviceprincipalid": "",
					"tenantid":           d.Get("azurerm_spn_tenantid").(string),
				},
				Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
			}
		}
		if serviceEndpointCreationMode == Manual {
			servicePrincipalId := credentials["serviceprincipalid"].(string)
			if servicePrincipalId == "" {
				return nil, nil, fmt.Errorf("serviceprincipalid is required for WorkloadIdentityFederation")
			}
			serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
				Parameters: &map[string]string{
					"serviceprincipalid": servicePrincipalId,
					"tenantid":           d.Get("azurerm_spn_tenantid").(string),
				},
				Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
			}
		}

		serviceEndpoint.Data = &map[string]string{
			"creationMode": string(serviceEndpointCreationMode),
			"environment":  environment,
		}
	}

	var endpointUrl string
	if environment == "AzureCloud" {
		endpointUrl = "https://management.azure.com/"
	} else if environment == "AzureChinaCloud" {
		endpointUrl = "https://management.chinacloudapi.cn/"
	}

	if scopeLevel == "Subscription" || scopeLevel == "ResourceGroup" {
		(*serviceEndpoint.Data)["scopeLevel"] = "Subscription"
		(*serviceEndpoint.Data)["subscriptionId"] = d.Get("azurerm_subscription_id").(string)
		(*serviceEndpoint.Data)["subscriptionName"] = d.Get("azurerm_subscription_name").(string)
	}

	if scopeLevel == "ResourceGroup" {
		(*serviceEndpoint.Authorization.Parameters)["scope"] = scope
	}

	if _, ok := d.GetOk("azurerm_management_group_id"); ok {
		(*serviceEndpoint.Data)["scopeLevel"] = "ManagementGroup"
		(*serviceEndpoint.Data)["managementGroupId"] = d.Get("azurerm_management_group_id").(string)
		(*serviceEndpoint.Data)["managementGroupName"] = d.Get("azurerm_management_group_name").(string)
	}

	serviceEndpoint.Type = converter.String("azurerm")
	serviceEndpoint.Url = converter.String(endpointUrl)
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAzureRM(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	scope := (*serviceEndpoint.Authorization.Parameters)["scope"]

	serviceEndPointType := AzureRmEndpointAuthenticationScheme(*serviceEndpoint.Authorization.Scheme)
	d.Set("service_endpoint_authentication_scheme", string(serviceEndPointType))
	if v, ok := (*serviceEndpoint.Data)["environment"]; ok {
		d.Set("environment", v)
	}

	if (*serviceEndpoint.Data)["creationMode"] == "Manual" {
		if _, ok := d.GetOk("credentials"); !ok {
			credentials := make(map[string]interface{})
			credentials["serviceprincipalid"] = (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"]
			credentials["serviceprincipalkey"] = d.Get("credentials.0.serviceprincipalkey").(string)
			d.Set("credentials", []interface{}{credentials})
		}
	}

	s := strings.SplitN(scope, "/", -1)
	if len(s) == 5 {
		d.Set("resource_group", s[4])
	}

	d.Set("azurerm_spn_tenantid", (*serviceEndpoint.Authorization.Parameters)["tenantid"])

	if _, ok := (*serviceEndpoint.Data)["managementGroupId"]; ok {
		d.Set("azurerm_management_group_id", (*serviceEndpoint.Data)["managementGroupId"])
		d.Set("azurerm_management_group_name", (*serviceEndpoint.Data)["managementGroupName"])
	}

	if _, ok := (*serviceEndpoint.Data)["subscriptionId"]; ok {
		d.Set("azurerm_subscription_id", (*serviceEndpoint.Data)["subscriptionId"])
		d.Set("azurerm_subscription_name", (*serviceEndpoint.Data)["subscriptionName"])
	}
}

// Validation function to ensure either Subscription or ManagementGroup scopeLevels are set correctly
func validateScopeLevel(scopeMap map[string][]string) error {
	// Check for empty
	if strings.TrimSpace(strings.Join(scopeMap["subscription"], "")) == "" && strings.TrimSpace(strings.Join(scopeMap["managementGroup"], "")) == "" {
		return fmt.Errorf("One of either subscription scoped (azurerm_subscription_id, azurerm_subscription_name) or managementGroup scoped (azurerm_management_ggroup_id, azurerm_management_group_name) details must be provided")
	}

	// check for valid subscription details
	var subElementCount int
	for _, ele := range scopeMap["subscription"] {
		if ele == "" {
			subElementCount = subElementCount + 1
		}
	}

	if subElementCount == 1 {
		return fmt.Errorf("azurerm_subscription_id and azurerm_subscription_name must be provided")
	}

	// check for valid managementGroup details
	var mgmtElementCount int
	for _, ele := range scopeMap["managementGroup"] {
		if ele == "" {
			mgmtElementCount = mgmtElementCount + 1
		}
	}

	if mgmtElementCount == 1 {
		return fmt.Errorf("azurerm_management_group_id and azurerm_management_group_name must be provided")
	}

	return nil
}

type AzureRmEndpointAuthenticationScheme string

const (
	ServicePrincipal           AzureRmEndpointAuthenticationScheme = "ServicePrincipal"
	ManagedServiceIdentity     AzureRmEndpointAuthenticationScheme = "ManagedServiceIdentity"
	WorkloadIdentityFederation AzureRmEndpointAuthenticationScheme = "WorkloadIdentityFederation"
)

type AzureRmEndpointCreationMode string

const (
	Automatic AzureRmEndpointCreationMode = "Automatic"
	Manual    AzureRmEndpointCreationMode = "Manual"
)
