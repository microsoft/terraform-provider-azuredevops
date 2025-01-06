package azuredevops_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk"
	mock_azuredevops "github.com/microsoft/terraform-provider-azuredevops/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_HasChildResources(t *testing.T) {
	expectedResources := []string{
		"azuredevops_agent_pool",
		"azuredevops_agent_queue",
		"azuredevops_area_permissions",
		"azuredevops_branch_policy_auto_reviewers",
		"azuredevops_branch_policy_build_validation",
		"azuredevops_branch_policy_comment_resolution",
		"azuredevops_branch_policy_merge_types",
		"azuredevops_branch_policy_min_reviewers",
		"azuredevops_branch_policy_status_check",
		"azuredevops_branch_policy_work_item_linking",
		"azuredevops_build_definition",
		"azuredevops_build_definition_permissions",
		"azuredevops_build_folder",
		"azuredevops_build_folder_permissions",
		"azuredevops_check_approval",
		"azuredevops_check_branch_control",
		"azuredevops_check_business_hours",
		"azuredevops_check_exclusive_lock",
		"azuredevops_check_required_template",
		"azuredevops_elastic_pool",
		"azuredevops_environment",
		"azuredevops_environment_resource_kubernetes",
		"azuredevops_feed",
		"azuredevops_feed_permission",
		"azuredevops_feed_retention_policy",
		"azuredevops_git_permissions",
		"azuredevops_git_repository",
		"azuredevops_git_repository_branch",
		"azuredevops_git_repository_file",
		"azuredevops_group",
		"azuredevops_group_entitlement",
		"azuredevops_group_membership",
		"azuredevops_iteration_permissions",
		"azuredevops_library_permissions",
		"azuredevops_pipeline_authorization",
		"azuredevops_project",
		"azuredevops_project_features",
		"azuredevops_project_permissions",
		"azuredevops_project_pipeline_settings",
		"azuredevops_project_tags",
		"azuredevops_repository_policy_author_email_pattern",
		"azuredevops_repository_policy_case_enforcement",
		"azuredevops_repository_policy_check_credentials",
		"azuredevops_repository_policy_file_path_pattern",
		"azuredevops_repository_policy_max_file_size",
		"azuredevops_repository_policy_max_path_length",
		"azuredevops_repository_policy_reserved_names",
		"azuredevops_resource_authorization",
		"azuredevops_securityrole_assignment",
		"azuredevops_serviceendpoint_argocd",
		"azuredevops_serviceendpoint_artifactory",
		"azuredevops_serviceendpoint_aws",
		"azuredevops_serviceendpoint_azure_service_bus",
		"azuredevops_serviceendpoint_azurecr",
		"azuredevops_serviceendpoint_azuredevops",
		"azuredevops_serviceendpoint_azurerm",
		"azuredevops_serviceendpoint_bitbucket",
		"azuredevops_serviceendpoint_checkmarx_sca",
		"azuredevops_serviceendpoint_checkmarx_sast",
		"azuredevops_serviceendpoint_dockerregistry",
		"azuredevops_serviceendpoint_dynamics_lifecycle_services",
		"azuredevops_serviceendpoint_externaltfs",
		"azuredevops_serviceendpoint_gcp_terraform",
		"azuredevops_serviceendpoint_generic",
		"azuredevops_serviceendpoint_generic_git",
		"azuredevops_serviceendpoint_github",
		"azuredevops_serviceendpoint_github_enterprise",
		"azuredevops_serviceendpoint_gitlab",
		"azuredevops_serviceendpoint_incomingwebhook",
		"azuredevops_serviceendpoint_jenkins",
		"azuredevops_serviceendpoint_jfrog_artifactory_v2",
		"azuredevops_serviceendpoint_jfrog_distribution_v2",
		"azuredevops_serviceendpoint_jfrog_platform_v2",
		"azuredevops_serviceendpoint_jfrog_xray_v2",
		"azuredevops_serviceendpoint_kubernetes",
		"azuredevops_serviceendpoint_maven",
		"azuredevops_serviceendpoint_nexus",
		"azuredevops_serviceendpoint_npm",
		"azuredevops_serviceendpoint_nuget",
		"azuredevops_serviceendpoint_octopusdeploy",
		"azuredevops_serviceendpoint_permissions",
		"azuredevops_serviceendpoint_runpipeline",
		"azuredevops_serviceendpoint_servicefabric",
		"azuredevops_serviceendpoint_snyk",
		"azuredevops_serviceendpoint_sonarcloud",
		"azuredevops_serviceendpoint_sonarqube",
		"azuredevops_serviceendpoint_ssh",
		"azuredevops_serviceendpoint_visualstudiomarketplace",
		"azuredevops_servicehook_permissions",
		"azuredevops_servicehook_storage_queue_pipelines",
		"azuredevops_service_principal_entitlement",
		"azuredevops_tagging_permissions",
		"azuredevops_team",
		"azuredevops_team_administrators",
		"azuredevops_team_members",
		"azuredevops_user_entitlement",
		"azuredevops_variable_group",
		"azuredevops_variable_group_permissions",
		"azuredevops_wiki",
		"azuredevops_wiki_page",
		"azuredevops_workitem",
		"azuredevops_workitemquery_permissions",
	}

	resources := azuredevops.Provider().ResourcesMap

	for _, resource := range expectedResources {
		require.Contains(t, resources, resource, fmt.Sprintf("An expected resource (%s) was not registered", resource))
		require.NotNil(t, resources[resource], "A resource cannot have a nil schema")
	}
	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")
}

func TestProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []string{
		"azuredevops_agent_pool",
		"azuredevops_agent_pools",
		"azuredevops_agent_queue",
		"azuredevops_area",
		"azuredevops_build_definition",
		"azuredevops_client_config",
		"azuredevops_environment",
		"azuredevops_feed",
		"azuredevops_git_repositories",
		"azuredevops_git_repository",
		"azuredevops_group",
		"azuredevops_groups",
		"azuredevops_identity_group",
		"azuredevops_identity_groups",
		"azuredevops_identity_user",
		"azuredevops_iteration",
		"azuredevops_project",
		"azuredevops_projects",
		"azuredevops_securityrole_definitions",
		"azuredevops_serviceendpoint_azurecr",
		"azuredevops_serviceendpoint_azurerm",
		"azuredevops_serviceendpoint_bitbucket",
		"azuredevops_serviceendpoint_github",
		"azuredevops_serviceendpoint_npm",
		"azuredevops_serviceendpoint_sonarcloud",
		"azuredevops_service_principal",
		"azuredevops_team",
		"azuredevops_teams",
		"azuredevops_users",
		"azuredevops_variable_group",
	}

	dataSources := azuredevops.Provider().DataSourcesMap

	for _, resource := range expectedDataSources {
		require.Contains(t, dataSources, resource, "An expected data source was not registered")
		require.NotNil(t, dataSources[resource], "A data source cannot have a nil schema")
	}
	require.Equal(t, len(expectedDataSources), len(dataSources), "There are an unexpected number of registered data sources")
}

func TestProvider_SchemaIsValid(t *testing.T) {
	type testParams struct {
		name          string
		required      bool
		defaultEnvVar string
		sensitive     bool
	}

	tests := []testParams{
		{"org_service_url", false, "AZDO_ORG_SERVICE_URL", false},
		{"personal_access_token", false, "AZDO_PERSONAL_ACCESS_TOKEN", true},

		{"client_id", false, "ARM_CLIENT_ID", false},
		{"tenant_id", false, "ARM_TENANT_ID", false},
		{"client_id_plan", false, "ARM_CLIENT_ID_PLAN", false},
		{"tenant_id_plan", false, "ARM_TENANT_ID_PLAN", false},
		{"client_id_apply", false, "ARM_CLIENT_ID_APPLY", false},
		{"tenant_id_apply", false, "ARM_TENANT_ID_APPLY", false},
		{"oidc_request_token", false, "ARM_OIDC_REQUEST_TOKEN", false},
		{"oidc_request_url", false, "ARM_OIDC_REQUEST_URL", false},
		{"oidc_token", false, "ARM_OIDC_TOKEN", true},
		{"oidc_token_file_path", false, "ARM_oidc_token_file_path", false},
		{"use_oidc", false, "ARM_USE_OIDC", false},
		{"oidc_audience", false, "ARM_OIDC_AUDIENCE", false},
		{"oidc_tfc_tag", false, "ARM_OIDC_TFC_TAG", false},
		{"client_certificate_path", false, "ARM_CLIENT_CERTIFICATE_PATH", false},
		{"client_certificate", false, "ARM_CLIENT_CERTIFICATE", true},
		{"client_certificate_password", false, "ARM_CLIENT_CERTIFICATE_PASSWORD", true},
		{"client_secret", false, "ARM_CLIENT_SECRET", true},
		{"client_secret_path", false, "ARM_CLIENT_SECRET_PATH", false},
		{"use_msi", false, "ARM_USE_MSI", false},
	}

	schema := azuredevops.Provider().Schema
	require.Equal(t, len(tests), len(schema), "There are an unexpected number of properties in the schema")

	for _, test := range tests {
		require.Contains(t, schema, test.name, "An expected property was not found in the schema")
		require.NotNil(t, schema[test.name], "A property in the schema cannot have a nil value")
		require.Equal(t, test.sensitive, schema[test.name].Sensitive, "A property in the schema has an incorrect sensitivity value")
		require.Equal(t, test.required, schema[test.name].Required, "A property in the schema has an incorrect required value")

		if test.defaultEnvVar != "" {
			expectedValue := os.Getenv(test.defaultEnvVar)

			actualValue, err := schema[test.name].DefaultFunc()
			if actualValue == nil {
				actualValue = ""
			}

			require.Nil(t, err, "An error occurred when getting the default value from the environment")
			require.Equal(t, expectedValue, actualValue, "The default value pulled from the environment has the wrong value")
		}
	}
}

func TestAuthPAT(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	testToken := "thepassword"
	resourceData.Set("personal_access_token", testToken)

	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Basic "+base64.StdEncoding.EncodeToString([]byte("_:"+testToken)), token)
}

type simpleTokenGetter struct {
	token string
}

func (s simpleTokenGetter) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     s.token,
		ExpiresOn: time.Now(),
	}, nil
}

func TestAuthOIDCToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("oidc_token", "buffalo123")
	resourceData.Set("use_oidc", true)

	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string,
			getAssertion func(context.Context) (string, error),
			options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func TestAuthOIDCTokenFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	oidcToken := "buffalo123"
	tempFile := t.TempDir() + "/clientSecret.txt"
	err := os.WriteFile(tempFile, []byte(oidcToken), 0644)
	assert.Nil(t, err)

	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("oidc_token_file_path", tempFile)
	resourceData.Set("use_oidc", true)

	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, token func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func TestAuthClientSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	clientSecret := "buffalo123"
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("client_secret", clientSecret)

	mockIdentityClient.EXPECT().NewClientSecretCredential(tenantId, clientId, clientSecret, nil).DoAndReturn(
		func(tenantID, clientID, secret string, options *azidentity.ClientSecretCredentialOptions) (*simpleTokenGetter, error) {
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func TestAuthClientSecretFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	clientSecret := "buffalo123"
	tempFile := t.TempDir() + "/clientSecret.txt"
	err := os.WriteFile(tempFile, []byte(clientSecret), 0644)
	assert.Nil(t, err)

	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("client_secret_path", tempFile)

	mockIdentityClient.EXPECT().NewClientSecretCredential(tenantId, clientId, clientSecret, nil).DoAndReturn(
		func(tenantID, clientID, secret string, options *azidentity.ClientSecretCredentialOptions) (*simpleTokenGetter, error) {
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func TestAuthTrfm(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	fakeTokenValue := "tokenvalue"
	os.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN", fakeTokenValue)
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("use_oidc", true)

	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func TestAuthTrfmPlanApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId_apply := "00000000-0000-0000-0000-000000000003"
	tenantId_apply := "00000000-0000-0000-0000-000000000004"
	clientId_plan := "00000000-0000-0000-0000-000000000005"
	tenantId_plan := "00000000-0000-0000-0000-000000000006"
	trfm_fake_token_plan := fmt.Sprintf("header.%s.signature", base64.StdEncoding.EncodeToString([]byte("{\"terraform_run_phase\":\"plan\"}")))
	trfm_fake_token_apply := fmt.Sprintf("header.%s.signature", base64.StdEncoding.EncodeToString([]byte("{\"terraform_run_phase\":\"apply\"}")))
	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	accessToken := "thepassword"
	resourceData.Set("client_id_apply", clientId_apply)
	resourceData.Set("tenant_id_apply", tenantId_apply)
	resourceData.Set("client_id_plan", clientId_plan)
	resourceData.Set("tenant_id_plan", tenantId_plan)
	resourceData.Set("use_oidc", true)

	// Apply phase test
	os.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN", trfm_fake_token_apply)
	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId_apply, clientId_apply, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)

	// Plan phase test
	os.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN", trfm_fake_token_plan)
	mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId_plan, clientId_plan, gomock.Any(), nil).DoAndReturn(
		func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err = sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err = resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func generateCert() []byte {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate certificate private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetUint64(20),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Minute * 5),
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to create private key: %v", err)
	}

	publicBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privateBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	return append(publicBytes[:], privateBytes[:]...)
}

func TestAuthClientCert(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	cert := generateCert()
	accessToken := "thepassword"

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("client_certificate", base64.StdEncoding.EncodeToString(cert))

	theseCerts, theseKey, err := azidentity.ParseCertificates(cert, nil)
	assert.Nil(t, err)

	mockIdentityClient.EXPECT().NewClientCertificateCredential(tenantId, clientId, gomock.Any(), gomock.Any(), nil).DoAndReturn(
		func(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (*simpleTokenGetter, error) {
			assert.Equal(t, theseCerts, certs)
			assert.Equal(t, theseKey, key)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func TestAuthClientCertFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000001"
	tenantId := "00000000-0000-0000-0000-000000000002"
	cert := generateCert()
	accessToken := "thepassword"
	tempFile := t.TempDir() + "/clientCerts.pem"
	err := os.WriteFile(tempFile, cert, 0644)
	assert.Nil(t, err)

	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("client_certificate_path", tempFile)

	theseCerts, theseKey, err := azidentity.ParseCertificates(cert, nil)
	assert.Nil(t, err)

	mockIdentityClient.EXPECT().NewClientCertificateCredential(tenantId, clientId, gomock.Any(), gomock.Any(), nil).DoAndReturn(
		func(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (*simpleTokenGetter, error) {
			assert.Equal(t, theseCerts, certs)
			assert.Equal(t, theseKey, key)
			getter := simpleTokenGetter{token: accessToken}
			return &getter, nil
		}).Times(1)
	resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
	assert.Nil(t, err)
	token, err := resp()
	assert.Nil(t, err)
	assert.Equal(t, "Bearer "+accessToken, token)
}

func TestGHActionsNoAudience(t *testing.T) {
	testCases := []struct {
		testAudience     string
		expectedAudience string
	}{
		{
			testAudience:     "",
			expectedAudience: "api://AzureADTokenExchange",
		},
		{
			testAudience:     "my-test-audience",
			expectedAudience: "my-test-audience",
		},
	}

	ctrl := gomock.NewController(t)
	mockIdentityClient := mock_azuredevops.NewMockIdentityFuncsI(ctrl)
	clientId := "00000000-0000-0000-0000-000000000003"
	tenantId := "00000000-0000-0000-0000-000000000004"
	resourceData := schema.TestResourceDataRaw(t, azuredevops.Provider().Schema, nil)
	ghFakeOIDCaccessToken := "the_gh_oidc_identity_token"
	accessToken := "thepassword"
	ghToken := "gh_oidc_token"
	resourceData.Set("client_id", clientId)
	resourceData.Set("tenant_id", tenantId)
	resourceData.Set("use_oidc", true)

	for _, testCase := range testCases {
		resourceData.Set("oidc_audience", testCase.testAudience)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, testCase.expectedAudience, r.URL.Query().Get("audience"))
			assert.Equal(t, int64(0), r.ContentLength)
			assert.Equal(t, "Bearer "+ghToken, r.Header.Get("Authorization"))
			w.Header().Add("content-type", "application/json")
			fmt.Fprintln(w, "{\"value\":\""+ghFakeOIDCaccessToken+"\"}")
		}))
		defer ts.Close()

		resourceData.Set("oidc_request_url", ts.URL)
		resourceData.Set("oidc_request_token", ghToken)

		mockIdentityClient.EXPECT().NewClientAssertionCredential(tenantId, clientId, gomock.Any(), nil).DoAndReturn(
			func(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (*simpleTokenGetter, error) {
				getter := simpleTokenGetter{token: accessToken}
				return &getter, nil
			}).Times(1)
		resp, err := sdk.GetAuthTokenProvider(context.Background(), resourceData, mockIdentityClient)
		assert.Nil(t, err)
		token, err := resp()
		assert.Nil(t, err)
		assert.Equal(t, "Bearer "+accessToken, token)
	}
}
