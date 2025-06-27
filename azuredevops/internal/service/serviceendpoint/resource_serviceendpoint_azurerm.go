package serviceendpoint

import (
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/serviceendpoint/migration"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const endpointValidationTimeoutSeconds = 60 * time.Second

// ResourceServiceEndpointAzureRM schema and implementation for AzureRM service endpoint resource
func ResourceServiceEndpointAzureRM() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointAzureRMCreate,
		Read:   resourceServiceEndpointAzureRMRead,
		Update: resourceServiceEndpointAzureRMUpdate,
		Delete: resourceServiceEndpointAzureRMDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	maps.Copy(r.Schema, map[string]*schema.Schema{
		"azurerm_spn_tenantid": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
			Description: "The service principal tenant id which should be used.",
		},

		"resource_group": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Description:   "Scope Resource Group",
			ConflictsWith: []string{"credentials", "azurerm_management_group_id"},
		},

		// Subscription scopeLevel
		"azurerm_subscription_id": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", nil),
			Description: "The Azure subscription Id which should be used.",
		},

		"azurerm_subscription_name": {
			Type:          schema.TypeString,
			Optional:      true,
			DefaultFunc:   schema.EnvDefaultFunc("ARM_SUBSCRIPTION_NAME", nil),
			Description:   "The Azure subscription name which should be used.",
			ConflictsWith: []string{"azurerm_management_group_id"},
		},

		// ManagementGroup scopeLevel
		"azurerm_management_group_id": {
			Type:          schema.TypeString,
			Optional:      true,
			DefaultFunc:   schema.EnvDefaultFunc("ARM_MGMT_GROUP_ID", nil),
			Description:   "The Azure managementGroup Id which should be used.",
			ConflictsWith: []string{"azurerm_subscription_id", "resource_group"},
		},

		"azurerm_management_group_name": {
			Type:          schema.TypeString,
			Optional:      true,
			DefaultFunc:   schema.EnvDefaultFunc("ARM_MGMT_GROUP_NAME", nil),
			Description:   "The Azure managementGroup name which should be used.",
			ConflictsWith: []string{"azurerm_subscription_id", "resource_group"},
		},

		"credentials": {
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
						Type:          schema.TypeString,
						Optional:      true,
						ConflictsWith: []string{"credentials.0.serviceprincipalcertificate"},
						Sensitive:     true,
						ValidateFunc:  validation.StringIsNotEmpty,
					},
					"serviceprincipalcertificate": {
						Type:          schema.TypeString,
						Optional:      true,
						ConflictsWith: []string{"credentials.0.serviceprincipalkey"},
						Sensitive:     true,
						ValidateFunc:  validation.StringIsNotEmpty,
					},
				},
			},
		},

		"environment": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "Environment (Azure Cloud type)",
			Default:      "AzureCloud",
			ValidateFunc: validation.StringInSlice([]string{"AzureCloud", "AzureChinaCloud", "AzureUSGovernment", "AzureGermanCloud", "AzureStack"}, false),
		},

		"service_endpoint_authentication_scheme": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "The AzureRM Service Endpoint Authentication Scheme, this can be 'WorkloadIdentityFederation', 'ManagedServiceIdentity' or 'ServicePrincipal'.",
			Default:      "ServicePrincipal",
			ValidateFunc: validation.StringInSlice([]string{"WorkloadIdentityFederation", "ManagedServiceIdentity", "ServicePrincipal"}, false),
		},

		"server_url": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"workload_identity_federation_issuer": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The issuer of the workload identity federation service principal.",
		},

		"workload_identity_federation_subject": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The subject of the workload identity federation service principal.",
		},

		"service_principal_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"features": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"validate": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Whether or not to validate connection with azure after create or update operations",
					},
				},
			},
		},
	})

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

func resourceServiceEndpointAzureRMCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointAzureRM(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	resp, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	serviceEndpoint.Id = resp.Id
	if shouldValidate(endpointFeatures(d)) {
		if err := validateServiceEndpoint(clients, serviceEndpoint, d.Get("project_id").(string), endpointValidationTimeoutSeconds); err != nil {
			if delErr := clients.ServiceEndpointClient.DeleteServiceEndpoint(
				clients.Ctx,
				serviceendpoint.DeleteServiceEndpointArgs{
					ProjectIds: &[]string{
						(*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String(),
					},
					EndpointId: resp.Id,
				}); delErr != nil {
				return fmt.Errorf("Delete service endpoint error %v", delErr)
			}
			return err
		}
	}

	d.SetId(resp.Id.String())
	return resourceServiceEndpointAzureRMRead(d, m)
}

func resourceServiceEndpointAzureRMRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if isServiceEndpointDeleted(d, err, serviceEndpoint, getArgs) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	if serviceEndpoint == nil || serviceEndpoint.Id == nil {
		d.SetId("")
		return nil
	}
	d.Set("features", d.Get("features"))

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointAzureRM(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointAzureRMUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointAzureRM(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if shouldValidate(endpointFeatures(d)) {
		if err := validateServiceEndpoint(clients, serviceEndpoint, d.Get("project_id").(string), endpointValidationTimeoutSeconds); err != nil {
			return err
		}
	}
	_, err = updateServiceEndpoint(clients, serviceEndpoint)
	if err != nil {
		return fmt.Errorf("updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointAzureRMRead(d, m)
}

func resourceServiceEndpointAzureRMDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointAzureRM(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAzureRM(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)

	serviceEndPointAuthenticationScheme := EndpointAuthenticationScheme(d.Get("service_endpoint_authentication_scheme").(string))

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
		return nil, err
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

	hasCredentials := len(credentials) > 0

	var serviceEndpointCreationMode EndpointCreationMode

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
					"serviceprincipalid": credentials["serviceprincipalid"].(string),
					"tenantid":           d.Get("azurerm_spn_tenantid").(string),
				},
				Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
			}

			if spnKey := credentials["serviceprincipalkey"].(string); spnKey != "" {
				(*serviceEndpoint.Authorization.Parameters)["authenticationType"] = "spnKey"
				(*serviceEndpoint.Authorization.Parameters)["serviceprincipalkey"] = spnKey
			}
			if spnCert := credentials["serviceprincipalcertificate"].(string); spnCert != "" {
				(*serviceEndpoint.Authorization.Parameters)["authenticationType"] = "spnCertificate"
				(*serviceEndpoint.Authorization.Parameters)["servicePrincipalCertificate"] = spnCert
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
				return nil, fmt.Errorf("serviceprincipalid is required for WorkloadIdentityFederation")
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
	switch environment {
	case "AzureCloud":
		endpointUrl = "https://management.azure.com/"
	case "AzureChinaCloud":
		endpointUrl = "https://management.chinacloudapi.cn/"
	case "AzureUSGovernment":
		endpointUrl = "https://management.usgovcloudapi.net/"
	case "AzureGermanCloud":
		endpointUrl = "https://management.microsoftazure.de"
	case "AzureStack":
		if serverUrl, ok := d.GetOk("server_url"); ok {
			endpointUrl = serverUrl.(string)
		} else {
			return nil, fmt.Errorf("`server_url` is required when `environment` is `AzureStack`")
		}
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
	return serviceEndpoint, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAzureRM(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	scope := (*serviceEndpoint.Authorization.Parameters)["scope"]

	serviceEndPointType := EndpointAuthenticationScheme(*serviceEndpoint.Authorization.Scheme)
	d.Set("service_endpoint_authentication_scheme", string(serviceEndPointType))
	if v, ok := (*serviceEndpoint.Data)["environment"]; ok {
		d.Set("environment", v)
	}

	if serviceEndPointType == WorkloadIdentityFederation {
		d.Set("workload_identity_federation_issuer", (*serviceEndpoint.Authorization.Parameters)["workloadIdentityFederationIssuer"])
		d.Set("workload_identity_federation_subject", (*serviceEndpoint.Authorization.Parameters)["workloadIdentityFederationSubject"])
	}

	if (*serviceEndpoint.Data)["creationMode"] == "Manual" {
		if _, ok := d.GetOk("credentials"); !ok {
			credentials := make(map[string]interface{})
			credentials["serviceprincipalid"] = (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"]
			credentials["serviceprincipalkey"] = ""
			credentials["serviceprincipalcertificate"] = ""
			d.Set("credentials", []interface{}{credentials})
		}
	}

	if serviceEndpoint.Url != nil {
		d.Set("server_url", serviceEndpoint.Url)
	}

	s := strings.Split(scope, "/")
	if len(s) == 5 {
		d.Set("resource_group", s[4])
	}

	d.Set("azurerm_spn_tenantid", (*serviceEndpoint.Authorization.Parameters)["tenantid"])

	if v, ok := (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"]; ok {
		d.Set("service_principal_id", v)
	}

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

	// check for valid managementGroup details
	var mgmtElementCount int
	for _, ele := range scopeMap["managementGroup"] {
		if ele == "" {
			mgmtElementCount++
		}
	}

	if mgmtElementCount == 1 {
		return fmt.Errorf("azurerm_management_group_id and azurerm_management_group_name must be provided")
	}

	return nil
}

func endpointFeatures(d *schema.ResourceData) map[string]interface{} {
	features := d.Get("features").([]interface{})
	if len(features) != 0 {
		return features[0].(map[string]interface{})
	}
	return nil
}

func shouldValidate(features map[string]interface{}) bool {
	validate, ok := features["validate"].(bool)
	if !ok {
		return false
	}
	return validate
}
