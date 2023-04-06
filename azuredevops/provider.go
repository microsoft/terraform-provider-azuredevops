package azuredevops

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/approvalsandchecks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/memberentitlementmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/policy/branch"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/policy/repository"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/workitemtracking"
)

type GHIdTokenResponse struct {
	Value string `json:"value"`
}

type HCPWorkloadToken struct {
	RunPhase string `json:"terraform_run_phase"`
}

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	servicePrincipalAuthFields := []string{"sp_oidc_token", "sp_oidc_token_path", "sp_oidc_github_actions", "sp_oidc_hcp", "sp_client_certificate_path", "sp_client_certificate", "sp_client_secret", "sp_client_secret_path"}
	allAuthFields := append([]string{"personal_access_token"}, servicePrincipalAuthFields...)

	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_resource_authorization":                 build.ResourceResourceAuthorization(),
			"azuredevops_branch_policy_build_validation":         branch.ResourceBranchPolicyBuildValidation(),
			"azuredevops_branch_policy_min_reviewers":            branch.ResourceBranchPolicyMinReviewers(),
			"azuredevops_branch_policy_auto_reviewers":           branch.ResourceBranchPolicyAutoReviewers(),
			"azuredevops_branch_policy_work_item_linking":        branch.ResourceBranchPolicyWorkItemLinking(),
			"azuredevops_branch_policy_comment_resolution":       branch.ResourceBranchPolicyCommentResolution(),
			"azuredevops_branch_policy_merge_types":              branch.ResourceBranchPolicyMergeTypes(),
			"azuredevops_branch_policy_status_check":             branch.ResourceBranchPolicyStatusCheck(),
			"azuredevops_build_definition":                       build.ResourceBuildDefinition(),
			"azuredevops_build_folder":                           build.ResourceBuildFolder(),
			"azuredevops_library_permissions":                    permissions.ResourceLibraryPermissions(),
			"azuredevops_project":                                core.ResourceProject(),
			"azuredevops_project_features":                       core.ResourceProjectFeatures(),
			"azuredevops_project_pipeline_settings":              core.ResourceProjectPipelineSettings(),
			"azuredevops_variable_group":                         taskagent.ResourceVariableGroup(),
			"azuredevops_repository_policy_author_email_pattern": repository.ResourceRepositoryPolicyAuthorEmailPatterns(),
			"azuredevops_repository_policy_file_path_pattern":    repository.ResourceRepositoryFilePathPatterns(),
			"azuredevops_repository_policy_case_enforcement":     repository.ResourceRepositoryEnforceConsistentCase(),
			"azuredevops_repository_policy_reserved_names":       repository.ResourceRepositoryReservedNames(),
			"azuredevops_repository_policy_max_path_length":      repository.ResourceRepositoryMaxPathLength(),
			"azuredevops_repository_policy_max_file_size":        repository.ResourceRepositoryMaxFileSize(),
			"azuredevops_repository_policy_check_credentials":    repository.ResourceRepositoryPolicyCheckCredentials(),
			"azuredevops_check_branch_control":                   approvalsandchecks.ResourceCheckBranchControl(),
			"azuredevops_check_business_hours":                   approvalsandchecks.ResourceCheckBusinessHours(),
			"azuredevops_serviceendpoint_argocd":                 serviceendpoint.ResourceServiceEndpointArgoCD(),
			"azuredevops_serviceendpoint_artifactory":            serviceendpoint.ResourceServiceEndpointArtifactory(),
			"azuredevops_serviceendpoint_jfrog_artifactory_v2":   serviceendpoint.ResourceServiceEndpointJFrogArtifactoryV2(),
			"azuredevops_serviceendpoint_jfrog_distribution_v2":  serviceendpoint.ResourceServiceEndpointJFrogDistributionV2(),
			"azuredevops_serviceendpoint_jfrog_platform_v2":      serviceendpoint.ResourceServiceEndpointJFrogPlatformV2(),
			"azuredevops_serviceendpoint_jfrog_xray_v2":          serviceendpoint.ResourceServiceEndpointJFrogXRayV2(),
			"azuredevops_serviceendpoint_aws":                    serviceendpoint.ResourceServiceEndpointAws(),
			"azuredevops_serviceendpoint_azurerm":                serviceendpoint.ResourceServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_bitbucket":              serviceendpoint.ResourceServiceEndpointBitBucket(),
			"azuredevops_serviceendpoint_azuredevops":            serviceendpoint.ResourceServiceEndpointAzureDevOps(),
			"azuredevops_serviceendpoint_dockerregistry":         serviceendpoint.ResourceServiceEndpointDockerRegistry(),
			"azuredevops_serviceendpoint_azurecr":                serviceendpoint.ResourceServiceEndpointAzureCR(),
			"azuredevops_serviceendpoint_github":                 serviceendpoint.ResourceServiceEndpointGitHub(),
			"azuredevops_serviceendpoint_incomingwebhook":        serviceendpoint.ResourceServiceEndpointIncomingWebhook(),
			"azuredevops_serviceendpoint_github_enterprise":      serviceendpoint.ResourceServiceEndpointGitHubEnterprise(),
			"azuredevops_serviceendpoint_kubernetes":             serviceendpoint.ResourceServiceEndpointKubernetes(),
			"azuredevops_serviceendpoint_octopusdeploy":          serviceendpoint.ResourceServiceEndpointOctopusDeploy(),
			"azuredevops_serviceendpoint_runpipeline":            serviceendpoint.ResourceServiceEndpointRunPipeline(),
			"azuredevops_serviceendpoint_servicefabric":          serviceendpoint.ResourceServiceEndpointServiceFabric(),
			"azuredevops_serviceendpoint_sonarqube":              serviceendpoint.ResourceServiceEndpointSonarQube(),
			"azuredevops_serviceendpoint_sonarcloud":             serviceendpoint.ResourceServiceEndpointSonarCloud(),
			"azuredevops_serviceendpoint_ssh":                    serviceendpoint.ResourceServiceEndpointSSH(),
			"azuredevops_serviceendpoint_npm":                    serviceendpoint.ResourceServiceEndpointNpm(),
			"azuredevops_serviceendpoint_nuget":                  serviceendpoint.ResourceServiceEndpointNuget(),
			"azuredevops_serviceendpoint_generic":                serviceendpoint.ResourceServiceEndpointGeneric(),
			"azuredevops_serviceendpoint_generic_git":            serviceendpoint.ResourceServiceEndpointGenericGit(),
			"azuredevops_serviceendpoint_externaltfs":            serviceendpoint.ResourceServiceEndpointExternalTFS(),
			"azuredevops_git_repository":                         git.ResourceGitRepository(),
			"azuredevops_git_repository_branch":                  git.ResourceGitRepositoryBranch(),
			"azuredevops_git_repository_file":                    git.ResourceGitRepositoryFile(),
			"azuredevops_user_entitlement":                       memberentitlementmanagement.ResourceUserEntitlement(),
			"azuredevops_group_membership":                       graph.ResourceGroupMembership(),
			"azuredevops_agent_pool":                             taskagent.ResourceAgentPool(),
			"azuredevops_agent_queue":                            taskagent.ResourceAgentQueue(),
			"azuredevops_group":                                  graph.ResourceGroup(),
			"azuredevops_project_permissions":                    permissions.ResourceProjectPermissions(),
			"azuredevops_git_permissions":                        permissions.ResourceGitPermissions(),
			"azuredevops_workitemquery_permissions":              permissions.ResourceWorkItemQueryPermissions(),
			"azuredevops_area_permissions":                       permissions.ResourceAreaPermissions(),
			"azuredevops_iteration_permissions":                  permissions.ResourceIterationPermissions(),
			"azuredevops_build_definition_permissions":           permissions.ResourceBuildDefinitionPermissions(),
			"azuredevops_build_folder_permissions":               permissions.ResourceBuildFolderPermissions(),
			"azuredevops_team":                                   core.ResourceTeam(),
			"azuredevops_team_members":                           core.ResourceTeamMembers(),
			"azuredevops_team_administrators":                    core.ResourceTeamAdministrators(),
			"azuredevops_serviceendpoint_permissions":            permissions.ResourceServiceEndpointPermissions(),
			"azuredevops_servicehook_permissions":                permissions.ResourceServiceHookPermissions(),
			"azuredevops_tagging_permissions":                    permissions.ResourceTaggingPermissions(),
			"azuredevops_variable_group_permissions":             permissions.ResourceVariableGroupPermissions(),
			"azuredevops_environment":                            taskagent.ResourceEnvironment(),
			"azuredevops_workitem":                               workitemtracking.ResourceWorkItem(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azuredevops_build_definition":        build.DataBuildDefinition(),
			"azuredevops_agent_pool":              taskagent.DataAgentPool(),
			"azuredevops_agent_pools":             taskagent.DataAgentPools(),
			"azuredevops_agent_queue":             taskagent.DataAgentQueue(),
			"azuredevops_client_config":           service.DataClientConfig(),
			"azuredevops_group":                   graph.DataGroup(),
			"azuredevops_project":                 core.DataProject(),
			"azuredevops_projects":                core.DataProjects(),
			"azuredevops_git_repositories":        git.DataGitRepositories(),
			"azuredevops_git_repository":          git.DataGitRepository(),
			"azuredevops_users":                   graph.DataUsers(),
			"azuredevops_area":                    workitemtracking.DataArea(),
			"azuredevops_iteration":               workitemtracking.DataIteration(),
			"azuredevops_team":                    core.DataTeam(),
			"azuredevops_teams":                   core.DataTeams(),
			"azuredevops_groups":                  graph.DataGroups(),
			"azuredevops_variable_group":          taskagent.DataVariableGroup(),
			"azuredevops_serviceendpoint_azurerm": serviceendpoint.DataServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_github":  serviceendpoint.DataServiceEndpointGithub(),
		},
		Schema: map[string]*schema.Schema{
			"org_service_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_ORG_SERVICE_URL", nil),
				Description: "The url of the Azure DevOps instance which should be used.",
			},
			"personal_access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_PERSONAL_ACCESS_TOKEN", nil),
				Description: "The personal access token which should be used.",
				Sensitive:   true,
			},
			"sp_client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_CLIENT_ID", nil),
				Description:  "The service principal client id which should be used.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"sp_client_id", "sp_tenant_id"},
			},
			"sp_tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_TENANT_ID", nil),
				Description:  "The service principal tenant id which should be used.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"sp_client_id", "sp_tenant_id"},
			},
			"sp_client_id_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_CLIENT_ID_PLAN", nil),
				Description:  "The service principal client id which should be used during a plan operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"sp_client_id_plan", "sp_tenant_id_plan", "sp_client_id_apply", "sp_tenant_id_apply"},
			},
			"sp_tenant_id_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_TENANT_ID_PLAN", nil),
				Description:  "The service principal tenant id which should be used during a plan operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"sp_client_id_plan", "sp_tenant_id_plan", "sp_client_id_apply", "sp_tenant_id_apply"},
			},
			"sp_client_id_apply": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_CLIENT_ID_APPLY", nil),
				Description:  "The service principal client id which should be used during an apply operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"sp_client_id_plan", "sp_tenant_id_plan", "sp_client_id_apply", "sp_tenant_id_apply"},
			},
			"sp_tenant_id_apply": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_TENANT_ID_APPLY", nil),
				Description:  "The service principal tenant id which should be used during an apply operation in Terraform Cloud..",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"sp_client_id_plan", "sp_tenant_id_plan", "sp_client_id_apply", "sp_tenant_id_apply"},
			},
			"sp_oidc_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_OIDC_TOKEN", nil),
				Description:  "OIDC token to authenticate as a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"sp_oidc_token", "sp_client_id", "sp_tenant_id"},
			},
			"sp_oidc_token_path": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_OIDC_TOKEN_PATH", nil),
				Description:  "OIDC token from file to authenticate as a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"sp_oidc_token_path", "sp_client_id", "sp_tenant_id"},
			},
			"sp_oidc_github_actions": {
				Type:         schema.TypeBool,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_OIDC_GITHUB_ACTIONS", nil),
				Description:  "Use the GitHub Actions OIDC token to authenticate to a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"sp_oidc_github_actions", "sp_client_id", "sp_tenant_id"},
			},
			"sp_oidc_github_actions_audience": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_OIDC_GITHUB_ACTIONS_AUDIENCE", nil),
				Description:  "Set the audience for the github actions ODIC token.",
				RequiredWith: []string{"sp_oidc_github_actions_audience", "sp_oidc_github_actions"},
			},
			"sp_oidc_hcp": {
				Type:         schema.TypeBool,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_OIDC_HCP", nil),
				Description:  "Use dynamic provider credentials in HCP to authenticate as a service principal.",
				ExactlyOneOf: allAuthFields,
			},
			"sp_client_certificate_path": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_CLIENT_CERTIFICATE_PATH", nil),
				Description:  "Path to a certificate to use to authenticate to the service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"sp_client_certificate_path", "sp_client_id", "sp_tenant_id"},
			},
			"sp_client_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_CLIENT_CERTIFICATE", nil),
				Description:  "Base64 encoded certificate to use to authenticate to the service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"sp_client_certificate", "sp_client_id", "sp_tenant_id"},
			},
			"sp_client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_SP_CLIENT_CERTIFICATE_PASSWORD", nil),
				Description: "Password for a client certificate password.",
			},
			"sp_client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_CLIENT_SECRET", nil),
				Description:  "Client secret for authenticating to  a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"sp_client_secret", "sp_client_id", "sp_tenant_id"},
			},
			"sp_client_secret_path": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SP_CLIENT_SECRET_PATH", nil),
				Description:  "Path to a file containing a client secret for authenticating to  a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"sp_client_secret_path", "sp_client_id", "sp_tenant_id"},
			},
		},
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func getGitHubOIDCToken(d *schema.ResourceData) (string, error) {
	requestUrl := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	client := &http.Client{}
	audience := "api://AzureADTokenExchange"

	if userAudience, ok := d.GetOk("sp_oidc_github_actions_audience"); ok {
		audience = userAudience.(string)
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}
	query := parsedUrl.Query()
	query.Add("audience", audience)
	parsedUrl.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", parsedUrl.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+requestToken)
	req.Header.Add("Accept", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	response_interface := GHIdTokenResponse{}
	err = json.NewDecoder(response.Body).Decode(&response_interface)
	if err != nil {
		return "", err
	}

	return response_interface.Value, nil
}

type TokenGetter interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type AzIdentityFuncs interface {
	NewClientAssertionCredential(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (TokenGetter, error)
	NewClientCertificateCredential(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (TokenGetter, error)
	NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenGetter, error)
}

type AzIdentityFuncsReal struct{}

func (a AzIdentityFuncsReal) NewClientAssertionCredential(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientAssertionCredential(tenantID, clientID, getAssertion, options)
}

func (a AzIdentityFuncsReal) NewClientCertificateCredential(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientCertificateCredential(tenantID, clientID, certs, key, options)
}

func (a AzIdentityFuncsReal) NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}

func GetAuthToken(ctx context.Context, d *schema.ResourceData, azIdentityFuncs AzIdentityFuncs) (string, error) {
	// Personal Access Token
	if personal_access_token, ok := d.GetOk("personal_access_token"); ok {
		return personal_access_token.(string), nil
	}

	tenantId := d.Get("sp_tenant_id").(string)
	clientId := d.Get("sp_client_id").(string)
	AzureDevOpsAppDefaultScope := "499b84ac-1321-427f-aa17-267ca6975798/.default"
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{AzureDevOpsAppDefaultScope},
	}

	var cred TokenGetter
	var err error

	// OIDC Token
	if sp_oidc_token, ok := d.GetOk("sp_oidc_token"); ok {
		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return sp_oidc_token.(string), nil }, nil)
		if err != nil {
			return "", err
		}
	}

	// OIDC Token From File
	if sp_oidc_token_path, ok := d.GetOk("sp_oidc_token_path"); ok {
		fileBytes, err := ioutil.ReadFile(sp_oidc_token_path.(string))
		if err != nil {
			return "", err
		}
		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return strings.TrimSpace(string(fileBytes)), nil }, nil)
		if err != nil {
			return "", err
		}
	}

	// OIDC Token in a GitHub Action Workflow
	if sp_oidc_github_actions, ok := d.GetOk("sp_oidc_github_actions"); ok && sp_oidc_github_actions.(bool) {
		gitHubToken, err := getGitHubOIDCToken(d)
		if err != nil {
			return "", err
		}
		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return gitHubToken, nil }, nil)
		if err != nil {
			return "", err
		}
	}

	// OIDC Token in a HashiCorp Vault run
	if sp_oidc_hcp, ok := d.GetOk("sp_oidc_hcp"); ok && sp_oidc_hcp.(bool) {
		workloadIdentityToken := os.Getenv("TFC_WORKLOAD_IDENTITY_TOKEN")

		// Check if plan & apply phases use different service principals
		if clientIdPlan, ok := d.GetOk("sp_client_id_plan"); ok {
			clientIdApply := d.Get("sp_client_id_apply").(string)
			tenantIdPlan := d.Get("sp_tenant_id_plan").(string)
			tenantIdApply := d.Get("sp_tenant_id_apply").(string)

			workloadIdentityTokenUnmarshalled := HCPWorkloadToken{}
			jwtParts := strings.Split(workloadIdentityToken, ".")
			if len(jwtParts) != 3 {
				return "", errors.New("Unable to split TFC_WORKLOAD_IDENTITY_TOKEN jwt")
			}
			tokenClaims, err := base64.StdEncoding.DecodeString(jwtParts[1])
			if err != nil {
				return "", err
			}
			err = json.Unmarshal(tokenClaims, &workloadIdentityTokenUnmarshalled)
			if err != nil {
				return "", err
			}

			if strings.EqualFold(workloadIdentityTokenUnmarshalled.RunPhase, "apply") {
				clientId = clientIdApply
				tenantId = tenantIdApply
			} else if strings.EqualFold(workloadIdentityTokenUnmarshalled.RunPhase, "plan") {
				clientId = clientIdPlan.(string)
				tenantId = tenantIdPlan
			} else {
				return "", errors.New(fmt.Sprintf("Unrecognized workspace run phase: %s", workloadIdentityTokenUnmarshalled.RunPhase))
			}
		} else if clientId == "" {
			return "", errors.New(fmt.Sprintf("Either sp_client_id or sp_client_id_plan must be set when using Terraform Cloud Workload Identity Token authentication."))
		}

		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return workloadIdentityToken, nil }, nil)
		if err != nil {
			return "", err
		}
	}

	// Certificate from a file on disk
	if sp_client_certificate_path, ok := d.GetOk("sp_client_certificate_path"); ok {
		fileBytes, err := ioutil.ReadFile(sp_client_certificate_path.(string))
		if err != nil {
			return "", err
		}

		certPassword := ([]byte)(nil)
		if password, ok := d.GetOk("sp_client_certificate_password"); ok {
			certPassword = []byte(password.(string))
		}

		certs, key, err := azidentity.ParseCertificates(fileBytes, certPassword)
		if err != nil {
			return "", err
		}

		cred, err = azIdentityFuncs.NewClientCertificateCredential(tenantId, clientId, certs, key, nil)
		if err != nil {
			return "", err
		}
	}

	// Certificate from a base64 encoded string
	if sp_client_certificate, ok := d.GetOk("sp_client_certificate"); ok {
		cert_bytes, err := base64.StdEncoding.DecodeString(sp_client_certificate.(string))
		if err != nil {
			return "", err
		}
		certPassword := ([]byte)(nil)
		if password, ok := d.GetOk("sp_client_certificate_password"); ok {
			certPassword = []byte(password.(string))
		}
		certs, key, err := azidentity.ParseCertificates(cert_bytes, certPassword)
		if err != nil {
			return "", err
		}
		cred, err = azIdentityFuncs.NewClientCertificateCredential(tenantId, clientId, certs, key, nil)
		if err != nil {
			return "", err
		}
	}

	// Client Secret
	if sp_client_secret, ok := d.GetOk("sp_client_secret"); ok {
		cred, err = azIdentityFuncs.NewClientSecretCredential(tenantId, clientId, sp_client_secret.(string), nil)
		if err != nil {
			return "", err
		}
	}

	// Client Secret from a file on disk
	if sp_client_secret_path, ok := d.GetOk("sp_client_secret_path"); ok {

		fileBytes, err := ioutil.ReadFile(sp_client_secret_path.(string))
		if err != nil {
			return "", err
		}
		cred, err = azIdentityFuncs.NewClientSecretCredential(tenantId, clientId, strings.TrimSpace(string(fileBytes)), nil)
		if err != nil {
			return "", err
		}
	}

	token, err := cred.GetToken(context.Background(), tokenOptions)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		token, err := GetAuthToken(ctx, d, AzIdentityFuncsReal{})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		azdo_client, err := client.GetAzdoClient(token, d.Get("org_service_url").(string), terraformVersion)

		return azdo_client, diag.FromErr(err)
	}
}
