package serviceendpoint

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointAzureCR schema and implementation for ACR service endpoint resource
func ResourceServiceEndpointAzureCR() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointAzureCRCreate,
		Read:   resourceServiceEndpointAzureCRRead,
		Update: resourceServiceEndpointAzureCRUpdate,
		Delete: resourceServiceEndpointAzureCRDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	r.Schema["azurecr_spn_tenantid"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("ACR_TENANT_ID", nil),
		Description: "The service principal tenant id which should be used.",
	}

	r.Schema["azurecr_subscription_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("ACR_SUBSCRIPTION_ID", nil),
		Description: "The Azure subscription Id which should be used.",
	}

	r.Schema["azurecr_subscription_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("ACR_SUBSCRIPTION_NAME", nil),
		Description: "The Azure subscription name which should be used.",
	}

	r.Schema["resource_group"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
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

	r.Schema["credentials"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serviceprincipalid": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The service principal id which should be used.",
				},
			},
		},
	}

	r.Schema["service_endpoint_authentication_scheme"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "The AzureCR Service Endpoint Authentication Scheme, this can be 'WorkloadIdentityFederation', 'ManagedServiceIdentity' or 'ServicePrincipal'.",
		Default:      "ServicePrincipal",
		ValidateFunc: validation.StringInSlice([]string{"WorkloadIdentityFederation", "ManagedServiceIdentity", "ServicePrincipal"}, false),
	}

	r.Schema["workload_identity_federation_issuer"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The issuer of the workload identity federation service principal.",
	}

	r.Schema["workload_identity_federation_subject"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The subject of the workload identity federation service principal.",
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

func resourceServiceEndpointAzureCRCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointAzureCR(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointAzureCRRead(d, m)
}

func resourceServiceEndpointAzureCRRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	flattenServiceEndpointAzureCR(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	return nil
}

func resourceServiceEndpointAzureCRUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointAzureCR(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointAzureCR(d, updatedServiceEndpoint, projectID.String())
	return resourceServiceEndpointAzureCRRead(d, m)
}

func resourceServiceEndpointAzureCRDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointAzureCR(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAzureCR(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndPointAuthenticationScheme := EndpointAuthenticationScheme(d.Get("service_endpoint_authentication_scheme").(string))

	subscriptionID := d.Get("azurecr_subscription_id")
	scope := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s",
		subscriptionID, d.Get("resource_group"), d.Get("azurecr_name"),
	)

	loginServer := fmt.Sprintf("%s.azurecr.io", strings.ToLower(d.Get("azurecr_name").(string)))

	var credentials map[string]interface{}

	if _, ok := d.GetOk("credentials"); ok {
		credentials = d.Get("credentials").([]interface{})[0].(map[string]interface{})
	}

	hasCredentials := len(credentials) > 0

	serviceEndPointAuthenticationSchemeHasCreationMode := serviceEndPointAuthenticationScheme == ServicePrincipal || serviceEndPointAuthenticationScheme == WorkloadIdentityFederation

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
					"authenticationType": "spnKey",
					"tenantid":           d.Get("azurecr_spn_tenantid").(string),
					"loginServer":        loginServer,
					"scope":              scope,
					"serviceprincipalid": d.Get("service_principal_id").(string),
				},
				Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
			}
		}
		if serviceEndpointCreationMode == Manual {
			return nil, nil, fmt.Errorf("ServicePrincipal Manual EndpointCreationMode is not supported yet")
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
	case ManagedServiceIdentity:
		return nil, nil, fmt.Errorf("ManagedServiceIdentity AuthenticationScheme is not supported yet")
	case WorkloadIdentityFederation:
		if serviceEndpointCreationMode == Automatic {
			serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
				Parameters: &map[string]string{
					"serviceprincipalid": "",
					"tenantid":           d.Get("azurecr_spn_tenantid").(string),
					"loginServer":        loginServer,
					"scope":              scope,
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
					"tenantid":           d.Get("azurecr_spn_tenantid").(string),
					"loginServer":        loginServer,
					"scope":              scope,
				},
				Scheme: converter.String(string(serviceEndPointAuthenticationScheme)),
			}
		}

		serviceEndpoint.Data = &map[string]string{
			"registryId":       scope,
			"registrytype":     "ACR",
			"subscriptionId":   subscriptionID.(string),
			"subscriptionName": d.Get("azurecr_subscription_name").(string),
			"creationMode":     string(serviceEndpointCreationMode),
		}
	}

	serviceEndpoint.Type = converter.String("dockerregistry")
	azureContainerRegistryURL := fmt.Sprintf("https://%s", loginServer)
	serviceEndpoint.Url = converter.String(azureContainerRegistryURL)

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAzureCR(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	serviceEndPointType := EndpointAuthenticationScheme(*serviceEndpoint.Authorization.Scheme)

	if serviceEndPointType == WorkloadIdentityFederation {
		d.Set("workload_identity_federation_issuer", (*serviceEndpoint.Authorization.Parameters)["workloadIdentityFederationIssuer"])
		d.Set("workload_identity_federation_subject", (*serviceEndpoint.Authorization.Parameters)["workloadIdentityFederationSubject"])
	}

	d.Set("service_endpoint_authentication_scheme", string(serviceEndPointType))

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		d.Set("azurecr_spn_tenantid", (*serviceEndpoint.Authorization.Parameters)["tenantid"])
		d.Set("service_principal_id", (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"])

		if scope, ok := (*serviceEndpoint.Authorization.Parameters)["scope"]; ok {
			s := strings.SplitN(scope, "/", -1)
			d.Set("resource_group", s[4])
			d.Set("azurecr_name", s[8])
		}
	}

	if (*serviceEndpoint.Data)["creationMode"] == "Manual" {
		if _, ok := d.GetOk("credentials"); !ok {
			credentials := make(map[string]interface{})
			credentials["serviceprincipalid"] = (*serviceEndpoint.Authorization.Parameters)["serviceprincipalid"]
			credentials["serviceprincipalkey"] = d.Get("credentials.0.serviceprincipalkey").(string)
			d.Set("credentials", []interface{}{credentials})
		}
	}

	if serviceEndPointType == WorkloadIdentityFederation {
		if serviceEndpoint.Data != nil {
			d.Set("azurecr_subscription_id", (*serviceEndpoint.Data)["subscriptionId"])
			d.Set("azurecr_subscription_name", (*serviceEndpoint.Data)["subscriptionName"])
		}
	}

	if serviceEndPointType == ServicePrincipal {
		if serviceEndpoint.Data != nil {
			d.Set("azurecr_subscription_id", (*serviceEndpoint.Data)["subscriptionId"])
			d.Set("azurecr_subscription_name", (*serviceEndpoint.Data)["subscriptionName"])
			d.Set("app_object_id", (*serviceEndpoint.Data)["appObjectId"])
			d.Set("spn_object_id", (*serviceEndpoint.Data)["spnObjectId"])
			d.Set("az_spn_role_permissions", (*serviceEndpoint.Data)["azureSpnPermissions"])
			d.Set("az_spn_role_assignment_id", (*serviceEndpoint.Data)["azureSpnRoleAssignmentId"])
		}
	}
}
