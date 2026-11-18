//go:build (all || resource_serviceendpoint_workload_identity) && !exclude_serviceendpoints

package serviceendpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func workloadIdentityTestConfig() map[string]interface{} {
	return map[string]interface{}{
		"project_id":            "11111111-1111-1111-1111-111111111111",
		"service_endpoint_name": "workload-identity",
		"org_url":               "https://dev.azure.com/example",
		"tenant_id":             "22222222-2222-2222-2222-222222222222",
		"service_principal_id":  "33333333-3333-3333-3333-333333333333",
		"object_id":             "44444444-4444-4444-4444-444444444444",
		"subject_descriptor":    "aadsp.example",
	}
}

// This fixture models the REST contract verified through the service endpoint
// types API independently of the resource expansion. The issuer and subject are opaque server-generated values.
func workloadIdentityTestEndpoint() *serviceendpoint.ServiceEndpoint {
	return &serviceendpoint.ServiceEndpoint{
		Id:          converter.ToPtr(uuid.MustParse("55555555-5555-5555-5555-555555555555")),
		Name:        converter.String("workload-identity"),
		Description: converter.String("Managed by Terraform"),
		Owner:       converter.String("library"),
		Type:        converter.String("workloadidentityuser"),
		Url:         converter.String("https://dev.azure.com/example"),
		IsReady:     converter.Bool(true),
		Authorization: &serviceendpoint.EndpointAuthorization{
			Scheme: converter.String("WorkloadIdentityFederation"),
			Parameters: &map[string]string{
				"tenantid":                          "22222222-2222-2222-2222-222222222222",
				"serviceprincipalid":                "33333333-3333-3333-3333-333333333333",
				"workloadIdentityFederationIssuer":  "https://login.microsoftonline.com/example/v2.0",
				"workloadIdentityFederationSubject": "server-generated-subject",
			},
		},
		Data: &map[string]string{
			"objectId":          "44444444-4444-4444-4444-444444444444",
			"subjectDescriptor": "aadsp.example",
			"isDraft":           "false",
			"workloadIdentityFederationServicePrincipalType": "Application",
		},
		ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: converter.ToPtr(uuid.MustParse("11111111-1111-1111-1111-111111111111")),
			},
			Name: converter.String("workload-identity"), Description: converter.String("Managed by Terraform"),
		}},
	}
}

func workloadIdentityTestClient(t *testing.T) (*azdosdkmocks.MockServiceendpointClient, *client.AggregatedClient) {
	mock := azdosdkmocks.NewMockServiceendpointClient(gomock.NewController(t))
	return mock, &client.AggregatedClient{ServiceEndpointClient: mock, Ctx: context.Background()}
}

func TestServiceEndpointWorkloadIdentity_Schema(t *testing.T) {
	r := ResourceServiceEndpointWorkloadIdentity()
	require.NoError(t, (&schema.Provider{ResourcesMap: map[string]*schema.Resource{"test_endpoint": r}}).InternalValidate())
	for _, tc := range []struct {
		name          string
		principalType string
		resourceID    string
		wantError     bool
	}{
		{name: "application", principalType: "Application"},
		{name: "managed identity", principalType: "ManagedIdentity", resourceID: "/subscriptions/example/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example"},
		{name: "managed identity without resource ID", principalType: "ManagedIdentity"},
		{name: "application with resource ID", principalType: "Application", resourceID: "/unexpected", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := workloadIdentityTestConfig()
			config["service_principal_type"] = tc.principalType
			if tc.resourceID != "" {
				config["managed_identity_resource_id"] = tc.resourceID
			}
			_, err := r.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil)
			if tc.wantError {
				require.ErrorContains(t, err, "managed_identity_resource_id")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServiceEndpointWorkloadIdentity_Create(t *testing.T) {
	for _, principalType := range []string{"Application", "ManagedIdentity"} {
		t.Run(principalType, func(t *testing.T) {
			r := ResourceServiceEndpointWorkloadIdentity()
			config := workloadIdentityTestConfig()
			config["service_principal_type"] = principalType
			ep := workloadIdentityTestEndpoint()
			(*ep.Data)["workloadIdentityFederationServicePrincipalType"] = principalType
			if principalType == "ManagedIdentity" {
				config["managed_identity_resource_id"] = "/subscriptions/example/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example"
				(*ep.Data)["managedIdentityResourceId"] = config["managed_identity_resource_id"].(string)
			}
			d := schema.TestResourceDataRaw(t, r.Schema, config)
			mock, clients := workloadIdentityTestClient(t)
			mock.EXPECT().CreateServiceEndpoint(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, args serviceendpoint.CreateServiceEndpointArgs) (*serviceendpoint.ServiceEndpoint, error) {
				request := args.Endpoint
				require.Nil(t, request.Id)
				require.Equal(t, "workloadidentityuser", *request.Type)
				require.Equal(t, config["org_url"], *request.Url)
				require.True(t, *request.IsReady)
				require.Equal(t, "WorkloadIdentityFederation", *request.Authorization.Scheme)
				require.Equal(t, map[string]string{"tenantid": config["tenant_id"].(string), "serviceprincipalid": config["service_principal_id"].(string)}, *request.Authorization.Parameters)
				require.Equal(t, *ep.Data, *request.Data)
				require.Equal(t, *ep.ServiceEndpointProjectReferences, *request.ServiceEndpointProjectReferences)
				return ep, nil
			})
			mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, args serviceendpoint.GetServiceEndpointDetailsArgs) (*serviceendpoint.ServiceEndpoint, error) {
				require.Equal(t, serviceendpoint.ServiceEndpointActionFilterValues.View, *args.ActionFilter)
				return ep, nil
			})
			require.Empty(t, r.CreateContext(context.Background(), d, clients))
			require.Equal(t, ep.Id.String(), d.Id())
			require.Equal(t, "server-generated-subject", d.Get("workload_identity_federation_subject"))
		})
	}
}

func TestServiceEndpointWorkloadIdentity_CreateDraft(t *testing.T) {
	for _, draft := range []bool{false, true} {
		t.Run(fmtBool(draft), func(t *testing.T) {
			r := ResourceServiceEndpointWorkloadIdentity()
			config := workloadIdentityTestConfig()
			config["draft"] = draft
			d := schema.TestResourceDataRaw(t, r.Schema, config)
			ep := workloadIdentityTestEndpoint()
			ep.IsReady = converter.Bool(false)
			(*ep.Data)["isAutoFICCreationFailed"] = "true"
			(*ep.Data)["isDraft"] = "true"
			ep.OperationStatus = map[string]interface{}{"state": "Failed", "statusMessage": "A valid refresh token was not found"}
			mock, clients := workloadIdentityTestClient(t)
			mock.EXPECT().CreateServiceEndpoint(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, args serviceendpoint.CreateServiceEndpointArgs) (*serviceendpoint.ServiceEndpoint, error) {
				require.Equal(t, !draft, *args.Endpoint.IsReady)
				require.Equal(t, map[bool]string{true: "true", false: "false"}[draft], (*args.Endpoint.Data)["isDraft"])
				return ep, nil
			})
			mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).Return(ep, nil)
			diagnostics := r.CreateContext(context.Background(), d, clients)
			require.Len(t, diagnostics, 1)
			require.Equal(t, !draft, diagnostics.HasError())
			if draft {
				require.Equal(t, diag.Warning, diagnostics[0].Severity)
			}
			require.Equal(t, ep.Id.String(), d.Id())
			require.Equal(t, false, d.Get("is_ready"))
			require.Equal(t, "server-generated-subject", d.Get("workload_identity_federation_subject"))
		})
	}
}

func TestServiceEndpointWorkloadIdentity_UpdatePreservesFederation(t *testing.T) {
	r := ResourceServiceEndpointWorkloadIdentity()
	config := workloadIdentityTestConfig()
	config["service_endpoint_name"] = "renamed"
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	ep := workloadIdentityTestEndpoint()
	ep.IsReady = converter.Bool(false)
	(*ep.Data)["serverMetadata"] = "preserve-me"
	d.SetId(ep.Id.String())
	mock, clients := workloadIdentityTestClient(t)
	gomock.InOrder(
		mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).Return(ep, nil),
		mock.EXPECT().UpdateServiceEndpoint(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, args serviceendpoint.UpdateServiceEndpointArgs) (*serviceendpoint.ServiceEndpoint, error) {
			require.Equal(t, ep.Id, args.EndpointId)
			require.Equal(t, "renamed", *args.Endpoint.Name)
			require.True(t, *args.Endpoint.IsReady)
			require.Equal(t, "preserve-me", (*args.Endpoint.Data)["serverMetadata"])
			require.Equal(t, "server-generated-subject", (*args.Endpoint.Authorization.Parameters)["workloadIdentityFederationSubject"])
			require.Equal(t, (*ep.Authorization.Parameters)["workloadIdentityFederationIssuer"], (*args.Endpoint.Authorization.Parameters)["workloadIdentityFederationIssuer"])
			require.Len(t, *args.Endpoint.Authorization.Parameters, 4)
			ep = args.Endpoint
			return ep, nil
		}),
		mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).DoAndReturn(func(context.Context, serviceendpoint.GetServiceEndpointDetailsArgs) (*serviceendpoint.ServiceEndpoint, error) {
			return ep, nil
		}),
	)
	require.Empty(t, r.UpdateContext(context.Background(), d, clients))
	require.Equal(t, true, d.Get("is_ready"))
	require.Equal(t, "renamed", d.Get("service_endpoint_name"))
}

func TestServiceEndpointWorkloadIdentity_ReadImport(t *testing.T) {
	r := ResourceServiceEndpointWorkloadIdentity()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"project_id": workloadIdentityTestConfig()["project_id"]})
	ep := workloadIdentityTestEndpoint()
	d.SetId(ep.Id.String())
	mock, clients := workloadIdentityTestClient(t)
	mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).Return(ep, nil)
	require.Empty(t, r.ReadContext(context.Background(), d, clients))
	for key, value := range workloadIdentityTestConfig() {
		require.Equal(t, value, d.Get(key), key)
	}
	require.Equal(t, "Application", d.Get("service_principal_type"))
	require.Equal(t, "", d.Get("managed_identity_resource_id"))
}

func TestServiceEndpointWorkloadIdentity_ReadRejectsPartialOrWrongEndpoint(t *testing.T) {
	for _, mutate := range []func(*serviceendpoint.ServiceEndpoint){
		func(ep *serviceendpoint.ServiceEndpoint) { ep.Authorization = nil },
		func(ep *serviceendpoint.ServiceEndpoint) { ep.Authorization.Parameters = nil },
		func(ep *serviceendpoint.ServiceEndpoint) { delete(*ep.Authorization.Parameters, "tenantid") },
		func(ep *serviceendpoint.ServiceEndpoint) { ep.Data = nil },
		func(ep *serviceendpoint.ServiceEndpoint) { delete(*ep.Data, "objectId") },
		func(ep *serviceendpoint.ServiceEndpoint) { ep.IsReady = nil },
		func(ep *serviceendpoint.ServiceEndpoint) { ep.Type = converter.String("azurerm") },
		func(ep *serviceendpoint.ServiceEndpoint) { ep.Authorization.Scheme = converter.String("Token") },
	} {
		r := ResourceServiceEndpointWorkloadIdentity()
		d := schema.TestResourceDataRaw(t, r.Schema, workloadIdentityTestConfig())
		ep := workloadIdentityTestEndpoint()
		mutate(ep)
		d.SetId(ep.Id.String())
		mock, clients := workloadIdentityTestClient(t)
		mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).Return(ep, nil)
		require.True(t, r.ReadContext(context.Background(), d, clients).HasError())
		require.Equal(t, workloadIdentityTestConfig()["service_principal_id"], d.Get("service_principal_id"))
	}
}

func TestServiceEndpointWorkloadIdentity_ReadNotFoundAndErrors(t *testing.T) {
	for _, code := range []int{404, 403, 500} {
		r := ResourceServiceEndpointWorkloadIdentity()
		d := schema.TestResourceDataRaw(t, r.Schema, workloadIdentityTestConfig())
		d.SetId(workloadIdentityTestEndpoint().Id.String())
		mock, clients := workloadIdentityTestClient(t)
		mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).Return(nil, azuredevops.WrappedError{StatusCode: &code})
		diagnostics := r.ReadContext(context.Background(), d, clients)
		require.Equal(t, code != 404, diagnostics.HasError())
		require.Equal(t, code == 404, d.Id() == "")
	}
}

func TestServiceEndpointWorkloadIdentity_CreateFailure(t *testing.T) {
	for _, apiError := range []error{errors.New("access denied"), nil} {
		r := ResourceServiceEndpointWorkloadIdentity()
		d := schema.TestResourceDataRaw(t, r.Schema, workloadIdentityTestConfig())
		mock, clients := workloadIdentityTestClient(t)
		mock.EXPECT().CreateServiceEndpoint(gomock.Any(), gomock.Any()).Return(nil, apiError)
		require.True(t, r.CreateContext(context.Background(), d, clients).HasError())
		require.Empty(t, d.Id())
	}
}

func TestServiceEndpointWorkloadIdentity_DeleteKeepsIdentity(t *testing.T) {
	r := ResourceServiceEndpointWorkloadIdentity()
	d := schema.TestResourceDataRaw(t, r.Schema, workloadIdentityTestConfig())
	ep := workloadIdentityTestEndpoint()
	d.SetId(ep.Id.String())
	mock, clients := workloadIdentityTestClient(t)
	mock.EXPECT().DeleteServiceEndpoint(gomock.Any(), serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: ep.Id,
		ProjectIds: &[]string{workloadIdentityTestConfig()["project_id"].(string)},
		Deep:       converter.Bool(false),
	}).Return(nil)
	mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).Return(nil, azuredevops.WrappedError{StatusCode: converter.Int(404)})
	require.Empty(t, r.DeleteContext(context.Background(), d, clients))
	require.Empty(t, d.Id())
}

func TestServiceEndpointWorkloadIdentity_ServerResolvedMetadata(t *testing.T) {
	r := ResourceServiceEndpointWorkloadIdentity()
	config := workloadIdentityTestConfig()
	delete(config, "subject_descriptor")
	config["service_principal_type"] = "ManagedIdentity"
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	request := expandServiceEndpointWorkloadIdentity(d)
	require.NotContains(t, *request.Data, "subjectDescriptor")
	require.NotContains(t, *request.Data, "managedIdentityResourceId")
	ep := workloadIdentityTestEndpoint()
	(*ep.Data)["workloadIdentityFederationServicePrincipalType"] = "ManagedIdentity"
	(*ep.Data)["managedIdentityResourceId"] = "/subscriptions/example/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example"
	// Both spellings occur in endpoint payloads.
	(*ep.Authorization.Parameters)["ServicePrincipalId"] = (*ep.Authorization.Parameters)["serviceprincipalid"]
	delete(*ep.Authorization.Parameters, "serviceprincipalid")
	require.NoError(t, flattenServiceEndpointWorkloadIdentity(d, ep))
	require.Equal(t, "aadsp.example", d.Get("subject_descriptor"))
	require.Equal(t, (*ep.Data)["managedIdentityResourceId"], d.Get("managed_identity_resource_id"))
	require.Equal(t, config["service_principal_id"], d.Get("service_principal_id"))
}

func TestServiceEndpointWorkloadIdentity_ReplaceManagedIdentityWithApplication(t *testing.T) {
	r := ResourceServiceEndpointWorkloadIdentity()
	config := workloadIdentityTestConfig()
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	ep := workloadIdentityTestEndpoint()
	(*ep.Data)["workloadIdentityFederationServicePrincipalType"] = "ManagedIdentity"
	(*ep.Data)["managedIdentityResourceId"] = "/subscriptions/example/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example"
	require.NoError(t, flattenServiceEndpointWorkloadIdentity(d, ep))
	state := d.State()
	state.RawConfig = cty.ObjectVal(map[string]cty.Value{
		"service_principal_type":       cty.StringVal("Application"),
		"managed_identity_resource_id": cty.NullVal(cty.String),
	})
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
	require.NoError(t, err)
	require.True(t, diff.RequiresNew())
}

func TestServiceEndpointWorkloadIdentity_DraftDoesNotHideUnrelatedFailure(t *testing.T) {
	r := ResourceServiceEndpointWorkloadIdentity()
	config := workloadIdentityTestConfig()
	config["draft"] = true
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	ep := workloadIdentityTestEndpoint()
	ep.IsReady = converter.Bool(false)
	ep.OperationStatus = map[string]interface{}{"state": "Failed", "statusMessage": "unexpected backend error"}
	mock, clients := workloadIdentityTestClient(t)
	mock.EXPECT().CreateServiceEndpoint(gomock.Any(), gomock.Any()).Return(ep, nil)
	mock.EXPECT().GetServiceEndpointDetails(gomock.Any(), gomock.Any()).Return(ep, nil)
	diagnostics := r.CreateContext(context.Background(), d, clients)
	require.True(t, diagnostics.HasError())
	require.Contains(t, diagnostics[0].Summary, "unexpected backend error")
	require.Equal(t, ep.Id.String(), d.Id())
}

func TestServiceEndpointWorkloadIdentity_DeleteNotFoundAndErrors(t *testing.T) {
	for _, code := range []int{404, 403, 500} {
		r := ResourceServiceEndpointWorkloadIdentity()
		d := schema.TestResourceDataRaw(t, r.Schema, workloadIdentityTestConfig())
		d.SetId(workloadIdentityTestEndpoint().Id.String())
		mock, clients := workloadIdentityTestClient(t)
		mock.EXPECT().DeleteServiceEndpoint(gomock.Any(), gomock.Any()).Return(azuredevops.WrappedError{StatusCode: &code})
		diagnostics := r.DeleteContext(context.Background(), d, clients)
		require.Equal(t, code != 404, diagnostics.HasError())
		require.Equal(t, code == 404, d.Id() == "")
	}
}

func fmtBool(value bool) string {
	if value {
		return "draft"
	}
	return "ready"
}
