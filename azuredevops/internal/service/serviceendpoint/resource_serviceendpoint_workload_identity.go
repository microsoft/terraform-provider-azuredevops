package serviceendpoint

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointWorkloadIdentity manages the native Azure DevOps
// connection, not the AZDOAPI endpoint supplied by the RunPipelines extension.
func ResourceServiceEndpointWorkloadIdentity() *schema.Resource {
	r := &schema.Resource{
		CreateContext: resourceServiceEndpointWorkloadIdentityCreate,
		ReadContext:   resourceServiceEndpointWorkloadIdentityRead,
		UpdateContext: resourceServiceEndpointWorkloadIdentityUpdate,
		DeleteContext: resourceServiceEndpointWorkloadIdentityDelete,
		Importer:      tfhelper.ImportProjectQualifiedResourceUUID(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: baseSchema(),
	}
	maps.Copy(r.Schema, map[string]*schema.Schema{
		"org_url": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  "URL of the Azure DevOps organization to connect to.",
		},
		"tenant_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
		"service_principal_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "Application (client) ID of the app registration or managed identity.",
		},
		"object_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "Object ID of the service principal in the Microsoft Entra tenant.",
		},
		"subject_descriptor": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
			Description:  "Azure DevOps graph descriptor of the identity added to the organization.",
		},
		"service_principal_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "Application",
			ValidateFunc: validation.StringInSlice([]string{"Application", "ManagedIdentity"}, false),
		},
		"managed_identity_resource_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"draft": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Allow the connection to remain a draft while configuring its federated credential. Set to false to finalize it.",
		},
		"is_ready": {
			Type:     schema.TypeBool,
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
	r.CustomizeDiff = func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
		// Validate an explicitly configured resource ID, not a computed value
		// retained from a managed identity being replaced with an application.
		raw := d.GetRawConfig()
		if !raw.IsNull() && raw.IsKnown() && raw.Type().HasAttribute("managed_identity_resource_id") {
			id := raw.GetAttr("managed_identity_resource_id")
			if id.IsNull() || !id.IsKnown() {
				return nil
			}
		}
		if !d.NewValueKnown("service_principal_type") || !d.NewValueKnown("managed_identity_resource_id") {
			return nil
		}
		managed := d.Get("service_principal_type").(string) == "ManagedIdentity"
		_, hasResourceID := d.GetOk("managed_identity_resource_id")
		if !managed && hasResourceID {
			return fmt.Errorf("managed_identity_resource_id can only be set when service_principal_type is ManagedIdentity")
		}
		return nil
	}
	return r
}

func expandServiceEndpointWorkloadIdentity(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	ep := doBaseExpansion(d)
	ep.Type = converter.String("workloadidentityuser")
	ep.Url = converter.String(d.Get("org_url").(string))
	ep.IsReady = converter.Bool(!d.Get("draft").(bool))
	ep.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("WorkloadIdentityFederation"),
		Parameters: &map[string]string{
			"tenantid":           d.Get("tenant_id").(string),
			"serviceprincipalid": d.Get("service_principal_id").(string),
		},
	}
	ep.Data = &map[string]string{
		"objectId": d.Get("object_id").(string),
		"isDraft":  strconv.FormatBool(d.Get("draft").(bool)),
		"workloadIdentityFederationServicePrincipalType": d.Get("service_principal_type").(string),
	}
	if descriptor, ok := d.GetOk("subject_descriptor"); ok {
		(*ep.Data)["subjectDescriptor"] = descriptor.(string)
	}
	if id, ok := d.GetOk("managed_identity_resource_id"); ok {
		(*ep.Data)["managedIdentityResourceId"] = id.(string)
	}
	return ep
}

func resourceServiceEndpointWorkloadIdentityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	ep, err := clients.ServiceEndpointClient.CreateServiceEndpoint(ctx, serviceendpoint.CreateServiceEndpointArgs{
		Endpoint: expandServiceEndpointWorkloadIdentity(d),
	})
	if err != nil {
		return diag.Errorf("creating Azure DevOps workload identity service endpoint: %v", err)
	}
	if ep == nil || ep.Id == nil {
		return diag.Errorf("creating Azure DevOps workload identity service endpoint: response has no ID")
	}
	// Keep the ID even when federation setup fails. The service can successfully
	// create a draft before discovering that it cannot create the Entra credential.
	d.SetId(ep.Id.String())
	return finishServiceEndpointWorkloadIdentity(ctx, d, clients, schema.TimeoutCreate)
}

func resourceServiceEndpointWorkloadIdentityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	args, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ep, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(ctx, *args)
	if utils.ResponseWasNotFound(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("reading Azure DevOps workload identity service endpoint: %v", err)
	}
	return diag.FromErr(flattenServiceEndpointWorkloadIdentity(d, ep))
}

func resourceServiceEndpointWorkloadIdentityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	args, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.FromErr(err)
	}
	current, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(ctx, *args)
	if err != nil {
		return diag.Errorf("reading Azure DevOps workload identity service endpoint before update: %v", err)
	}
	if err := validateServiceEndpointWorkloadIdentity(current); err != nil {
		return diag.FromErr(err)
	}
	ep := expandServiceEndpointWorkloadIdentity(d)
	// Preserve the server-generated federation details and identity metadata when
	// renaming or finalizing a draft. Do not ask the service to generate new trust.
	parameters := maps.Clone(*current.Authorization.Parameters)
	for key, value := range *ep.Authorization.Parameters {
		for existing := range parameters {
			if strings.EqualFold(existing, key) {
				delete(parameters, existing)
			}
		}
		parameters[key] = value
	}
	data := maps.Clone(*current.Data)
	maps.Copy(data, *ep.Data)
	ep.Authorization.Parameters = &parameters
	ep.Data = &data
	_, err = clients.ServiceEndpointClient.UpdateServiceEndpoint(ctx, serviceendpoint.UpdateServiceEndpointArgs{
		EndpointId: ep.Id,
		Endpoint:   ep,
	})
	if err != nil {
		return diag.Errorf("updating Azure DevOps workload identity service endpoint: %v", err)
	}
	return finishServiceEndpointWorkloadIdentity(ctx, d, clients, schema.TimeoutUpdate)
}

func finishServiceEndpointWorkloadIdentity(ctx context.Context, d *schema.ResourceData, clients *client.AggregatedClient, timeout string) diag.Diagnostics {
	args, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.FromErr(err)
	}
	wait := &retry.StateChangeConf{
		Pending: []string{opState.InProgress}, Target: []string{opState.Ready},
		Timeout: d.Timeout(timeout), MinTimeout: time.Second,
		Refresh: func() (interface{}, string, error) {
			ep, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(ctx, *args)
			if err != nil {
				return nil, "", err
			}
			if err := flattenServiceEndpointWorkloadIdentity(d, ep); err != nil {
				return nil, "", err
			}
			if status, ok := ep.OperationStatus.(map[string]interface{}); ok {
				if status["state"] == opState.InProgress {
					return ep, opState.InProgress, nil
				}
				// Only a recognized federation draft can defer a failed operation.
				// Unrelated backend failures must still fail the apply.
				if status["state"] == opState.Failed && !isWorkloadIdentityFederationDraft(ep) {
					return nil, "", fmt.Errorf("Azure DevOps failed to configure the workload identity service endpoint: %v", status["statusMessage"])
				}
			}
			return ep, opState.Ready, nil
		},
	}
	if _, err := wait.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("waiting for Azure DevOps workload identity service endpoint: %v", err)
	}
	if !d.Get("is_ready").(bool) {
		message := "Azure DevOps saved the service connection as a draft. Configure its federated identity credential and finalize the connection by applying with draft = false. The calling identity may lack permission to configure federation automatically."
		if !d.Get("draft").(bool) {
			return diag.Errorf("%s", message)
		}
		return diag.Diagnostics{{Severity: diag.Warning, Summary: "Service connection is not ready", Detail: message}}
	}
	return nil
}

func isWorkloadIdentityFederationDraft(ep *serviceendpoint.ServiceEndpoint) bool {
	if *ep.IsReady {
		return false
	}
	data := *ep.Data
	return strings.EqualFold(data["isAutoFICCreationFailed"], "true") ||
		data["federatedCredentialDraftState"] == "actionRequired" ||
		data["federatedCredentialDraftState"] == "pendingValidation"
}

func resourceServiceEndpointWorkloadIdentityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	args, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = clients.ServiceEndpointClient.DeleteServiceEndpoint(ctx, serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: args.EndpointId,
		ProjectIds: &[]string{*args.Project},
		Deep:       converter.Bool(false),
	})
	if err != nil && !utils.ResponseWasNotFound(err) {
		return diag.Errorf("deleting Azure DevOps workload identity service endpoint: %v", err)
	}
	if err == nil {
		wait := &retry.StateChangeConf{
			Pending: []string{"Present"}, Target: []string{"Deleted"},
			Timeout: d.Timeout(schema.TimeoutDelete), MinTimeout: time.Second,
			Refresh: func() (interface{}, string, error) {
				_, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(ctx, *args)
				if utils.ResponseWasNotFound(err) {
					return true, "Deleted", nil
				}
				return false, "Present", err
			},
		}
		if _, err := wait.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("waiting for Azure DevOps workload identity service endpoint deletion: %v", err)
		}
	}
	d.SetId("")
	return nil
}

func validateServiceEndpointWorkloadIdentity(ep *serviceendpoint.ServiceEndpoint) error {
	if ep == nil || ep.Id == nil || ep.Type == nil || ep.Url == nil || ep.Data == nil || ep.Authorization == nil || ep.Authorization.Scheme == nil || ep.Authorization.Parameters == nil || ep.IsReady == nil {
		return fmt.Errorf("Azure DevOps workload identity service endpoint response is incomplete; check the caller's service connection permissions")
	}
	if !strings.EqualFold(*ep.Type, "workloadidentityuser") || !strings.EqualFold(*ep.Authorization.Scheme, "WorkloadIdentityFederation") {
		return fmt.Errorf("expected a workloadidentityuser endpoint with WorkloadIdentityFederation authentication")
	}
	for _, key := range []string{"tenantid", "serviceprincipalid"} {
		if workloadIdentityAuthorizationParameter(ep, key) == "" {
			return fmt.Errorf("Azure DevOps workload identity service endpoint response is missing %s; check the caller's service connection permissions", key)
		}
	}
	for _, key := range []string{"objectId", "workloadIdentityFederationServicePrincipalType"} {
		if (*ep.Data)[key] == "" {
			return fmt.Errorf("Azure DevOps workload identity service endpoint response is missing %s; check the caller's service connection permissions", key)
		}
	}
	return nil
}

func workloadIdentityAuthorizationParameter(ep *serviceendpoint.ServiceEndpoint, key string) string {
	for parameter, value := range *ep.Authorization.Parameters {
		if strings.EqualFold(parameter, key) {
			return value
		}
	}
	return ""
}

func flattenServiceEndpointWorkloadIdentity(d *schema.ResourceData, ep *serviceendpoint.ServiceEndpoint) error {
	if err := validateServiceEndpointWorkloadIdentity(ep); err != nil {
		return err
	}
	doBaseFlattening(d, ep)
	d.Set("org_url", *ep.Url)
	d.Set("is_ready", *ep.IsReady)
	for attribute, key := range map[string]string{
		"tenant_id":                            "tenantid",
		"service_principal_id":                 "serviceprincipalid",
		"workload_identity_federation_issuer":  "workloadIdentityFederationIssuer",
		"workload_identity_federation_subject": "workloadIdentityFederationSubject",
	} {
		d.Set(attribute, workloadIdentityAuthorizationParameter(ep, key))
	}
	for attribute, key := range map[string]string{
		"object_id": "objectId", "subject_descriptor": "subjectDescriptor",
		"service_principal_type":       "workloadIdentityFederationServicePrincipalType",
		"managed_identity_resource_id": "managedIdentityResourceId",
	} {
		d.Set(attribute, (*ep.Data)[key])
	}
	return nil
}
