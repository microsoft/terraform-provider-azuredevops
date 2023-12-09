package azuredevops

import (
	"context"
	"encoding/base64"

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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/dynamiccredentialproviders"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	servicePrincipalAuthFields := []string{"oidc_token", "oidc_token_file_path", "oidc_github_actions", "oidc_hcp", "client_certificate_path", "client_certificate", "client_secret", "client_secret_path"}
	allAuthFields := append([]string{"personal_access_token"}, servicePrincipalAuthFields...)

	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_resource_authorization":                 build.ResourceResourceAuthorization(),
			"azuredevops_pipeline_authorization":                 build.ResourcePipelineAuthorization(),
			"azuredevops_branch_policy_build_validation":         branch.ResourceBranchPolicyBuildValidation(),
			"azuredevops_branch_policy_min_reviewers":            branch.ResourceBranchPolicyMinReviewers(),
			"azuredevops_branch_policy_auto_reviewers":           branch.ResourceBranchPolicyAutoReviewers(),
			"azuredevops_branch_policy_work_item_linking":        branch.ResourceBranchPolicyWorkItemLinking(),
			"azuredevops_branch_policy_comment_resolution":       branch.ResourceBranchPolicyCommentResolution(),
			"azuredevops_branch_policy_merge_types":              branch.ResourceBranchPolicyMergeTypes(),
			"azuredevops_branch_policy_status_check":             branch.ResourceBranchPolicyStatusCheck(),
			"azuredevops_build_definition":                       build.ResourceBuildDefinition(),
			"azuredevops_build_folder":                           build.ResourceBuildFolder(),
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
			"azuredevops_check_approval":                         approvalsandchecks.ResourceCheckApproval(),
			"azuredevops_check_exclusive_lock":                   approvalsandchecks.ResourceCheckExclusiveLock(),
			"azuredevops_check_branch_control":                   approvalsandchecks.ResourceCheckBranchControl(),
			"azuredevops_check_business_hours":                   approvalsandchecks.ResourceCheckBusinessHours(),
			"azuredevops_check_required_template":                approvalsandchecks.ResourceCheckRequiredTemplate(),
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
			"azuredevops_serviceendpoint_gcp_terraform":          serviceendpoint.ResourceServiceEndpointGcp(),
			"azuredevops_serviceendpoint_incomingwebhook":        serviceendpoint.ResourceServiceEndpointIncomingWebhook(),
			"azuredevops_serviceendpoint_github_enterprise":      serviceendpoint.ResourceServiceEndpointGitHubEnterprise(),
			"azuredevops_serviceendpoint_kubernetes":             serviceendpoint.ResourceServiceEndpointKubernetes(),
			"azuredevops_serviceendpoint_maven":                  serviceendpoint.ResourceServiceEndpointMaven(),
			"azuredevops_serviceendpoint_nuget":                  serviceendpoint.ResourceServiceEndpointNuGet(),
			"azuredevops_serviceendpoint_nexus":                  serviceendpoint.ResourceServiceEndpointNexus(),
			"azuredevops_serviceendpoint_jenkins":                serviceendpoint.ResourceServiceEndpointJenkins(),
			"azuredevops_serviceendpoint_octopusdeploy":          serviceendpoint.ResourceServiceEndpointOctopusDeploy(),
			"azuredevops_serviceendpoint_runpipeline":            serviceendpoint.ResourceServiceEndpointRunPipeline(),
			"azuredevops_serviceendpoint_servicefabric":          serviceendpoint.ResourceServiceEndpointServiceFabric(),
			"azuredevops_serviceendpoint_sonarqube":              serviceendpoint.ResourceServiceEndpointSonarQube(),
			"azuredevops_serviceendpoint_sonarcloud":             serviceendpoint.ResourceServiceEndpointSonarCloud(),
			"azuredevops_serviceendpoint_ssh":                    serviceendpoint.ResourceServiceEndpointSSH(),
			"azuredevops_serviceendpoint_npm":                    serviceendpoint.ResourceServiceEndpointNpm(),
			"azuredevops_serviceendpoint_generic":                serviceendpoint.ResourceServiceEndpointGeneric(),
			"azuredevops_serviceendpoint_generic_git":            serviceendpoint.ResourceServiceEndpointGenericGit(),
			"azuredevops_serviceendpoint_externaltfs":            serviceendpoint.ResourceServiceEndpointExternalTFS(),
			"azuredevops_git_repository":                         git.ResourceGitRepository(),
			"azuredevops_git_repository_branch":                  git.ResourceGitRepositoryBranch(),
			"azuredevops_git_repository_file":                    git.ResourceGitRepositoryFile(),
			"azuredevops_user_entitlement":                       memberentitlementmanagement.ResourceUserEntitlement(),
			"azuredevops_group_entitlement":                      memberentitlementmanagement.ResourceGroupEntitlement(),
			"azuredevops_group_membership":                       graph.ResourceGroupMembership(),
			"azuredevops_agent_pool":                             taskagent.ResourceAgentPool(),
			"azuredevops_elastic_pool":                           taskagent.ResourceAgentPoolVMSS(),
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
			"azuredevops_environment":                            taskagent.ResourceEnvironment(),
			"azuredevops_workitem":                               workitemtracking.ResourceWorkItem(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azuredevops_build_definition":           build.DataBuildDefinition(),
			"azuredevops_agent_pool":                 taskagent.DataAgentPool(),
			"azuredevops_agent_pools":                taskagent.DataAgentPools(),
			"azuredevops_agent_queue":                taskagent.DataAgentQueue(),
			"azuredevops_client_config":              service.DataClientConfig(),
			"azuredevops_environment":                taskagent.DataEnvironment(),
			"azuredevops_group":                      graph.DataGroup(),
			"azuredevops_project":                    core.DataProject(),
			"azuredevops_projects":                   core.DataProjects(),
			"azuredevops_git_repositories":           git.DataGitRepositories(),
			"azuredevops_git_repository":             git.DataGitRepository(),
			"azuredevops_users":                      graph.DataUsers(),
			"azuredevops_area":                       workitemtracking.DataArea(),
			"azuredevops_iteration":                  workitemtracking.DataIteration(),
			"azuredevops_team":                       core.DataTeam(),
			"azuredevops_teams":                      core.DataTeams(),
			"azuredevops_groups":                     graph.DataGroups(),
			"azuredevops_variable_group":             taskagent.DataVariableGroup(),
			"azuredevops_serviceendpoint_azurerm":    serviceendpoint.DataServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_github":     serviceendpoint.DataServiceEndpointGithub(),
			"azuredevops_serviceendpoint_npm":        serviceendpoint.DataResourceServiceEndpointNpm(),
			"azuredevops_serviceendpoint_azurecr":    serviceendpoint.DataResourceServiceEndpointAzureCR(),
			"azuredevops_serviceendpoint_sonarcloud": serviceendpoint.DataResourceServiceEndpointSonarCloud(),
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
			"client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID", nil),
				Description:  "The service principal client id which should be used.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"client_id", "tenant_id"},
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
				Description:  "The service principal tenant id which should be used.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"client_id", "tenant_id"},
			},
			"client_id_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID_PLAN", nil),
				Description:  "The service principal client id which should be used during a plan operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"client_id_plan", "tenant_id_plan", "client_id_apply", "tenant_id_apply"},
			},
			"tenant_id_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID_PLAN", nil),
				Description:  "The service principal tenant id which should be used during a plan operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"client_id_plan", "tenant_id_plan", "client_id_apply", "tenant_id_apply"},
			},
			"client_id_apply": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID_APPLY", nil),
				Description:  "The service principal client id which should be used during an apply operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"client_id_plan", "tenant_id_plan", "client_id_apply", "tenant_id_apply"},
			},
			"tenant_id_apply": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID_APPLY", nil),
				Description:  "The service principal tenant id which should be used during an apply operation in Terraform Cloud..",
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"client_id_plan", "tenant_id_plan", "client_id_apply", "tenant_id_apply"},
			},
			"oidc_request_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}, ""),
				Description: "The bearer token for the request to the OIDC provider. For use When authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_request_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"}, ""),
				Description: "The URL for the OIDC provider from which to request an ID token. For use When authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_OIDC_TOKEN", nil),
				Description:  "OIDC token to authenticate as a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"oidc_token", "client_id", "tenant_id"},
			},
			"oidc_token_file_path": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_oidc_token_file_path", nil),
				Description:  "OIDC token from file to authenticate as a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"oidc_token_file_path", "client_id", "tenant_id"},
			},
			"oidc_github_actions": {
				Type:         schema.TypeBool,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_OIDC_GITHUB_ACTIONS", nil),
				Description:  "Use the GitHub Actions OIDC token to authenticate to a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"oidc_github_actions", "client_id", "tenant_id"},
			},
			"oidc_github_actions_audience": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_OIDC_GITHUB_ACTIONS_AUDIENCE", "api://AzureADTokenExchange"),
				Description:  "Set the audience for the github actions ODIC token.",
				RequiredWith: []string{"oidc_github_actions_audience", "oidc_github_actions"},
			},
			"oidc_hcp": {
				Type:         schema.TypeBool,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_OIDC_HCP", nil),
				Description:  "Use dynamic provider credentials in HCP to authenticate as a service principal.",
				ExactlyOneOf: allAuthFields,
			},
			"client_certificate_path": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", nil),
				Description:  "Path to a certificate to use to authenticate to the service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"client_certificate_path", "client_id", "tenant_id"},
			},
			"client_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE", nil),
				Description:  "Base64 encoded certificate to use to authenticate to the service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"client_certificate", "client_id", "tenant_id"},
			},
			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", nil),
				Description: "Password for a client certificate password.",
			},
			"client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_SECRET", nil),
				Description:  "Client secret for authenticating to  a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"client_secret", "client_id", "tenant_id"},
			},
			"client_secret_path": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_SECRET_PATH", nil),
				Description:  "Path to a file containing a client secret for authenticating to  a service principal.",
				ExactlyOneOf: allAuthFields,
				RequiredWith: []string{"client_secret_path", "client_id", "tenant_id"},
			},
		},
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		var azdoClient *client.AggregatedClient
		var err error
		var tokenFunction func() (string, error)
		// Personal Access Token
		if personal_access_token, ok := d.GetOk("personal_access_token"); ok {
			tokenFunction = func() (string, error) {
				auth := "_:" + personal_access_token.(string)
				return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth)), nil
			}
		} else {
			// Service Principal
			tokenFunction, err = dynamiccredentialproviders.GetAuthToken(ctx, d, dynamiccredentialproviders.AzIdentityFuncsImpl{})
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		azdoClient, err = client.GetAzdoClient(tokenFunction, d.Get("org_service_url").(string), terraformVersion)
		return azdoClient, diag.FromErr(err)
	}
}
