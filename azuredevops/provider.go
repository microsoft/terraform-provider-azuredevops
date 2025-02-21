package azuredevops

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/approvalsandchecks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/dashboard"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/memberentitlementmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/policy/branch"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/policy/repository"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/securityroles"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/servicehook"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_agent_pool":                                  taskagent.ResourceAgentPool(),
			"azuredevops_agent_queue":                                 taskagent.ResourceAgentQueue(),
			"azuredevops_area_permissions":                            permissions.ResourceAreaPermissions(),
			"azuredevops_branch_policy_auto_reviewers":                branch.ResourceBranchPolicyAutoReviewers(),
			"azuredevops_branch_policy_build_validation":              branch.ResourceBranchPolicyBuildValidation(),
			"azuredevops_branch_policy_comment_resolution":            branch.ResourceBranchPolicyCommentResolution(),
			"azuredevops_branch_policy_merge_types":                   branch.ResourceBranchPolicyMergeTypes(),
			"azuredevops_branch_policy_min_reviewers":                 branch.ResourceBranchPolicyMinReviewers(),
			"azuredevops_branch_policy_status_check":                  branch.ResourceBranchPolicyStatusCheck(),
			"azuredevops_branch_policy_work_item_linking":             branch.ResourceBranchPolicyWorkItemLinking(),
			"azuredevops_build_definition":                            build.ResourceBuildDefinition(),
			"azuredevops_build_definition_permissions":                permissions.ResourceBuildDefinitionPermissions(),
			"azuredevops_build_folder":                                build.ResourceBuildFolder(),
			"azuredevops_build_folder_permissions":                    permissions.ResourceBuildFolderPermissions(),
			"azuredevops_check_approval":                              approvalsandchecks.ResourceCheckApproval(),
			"azuredevops_check_branch_control":                        approvalsandchecks.ResourceCheckBranchControl(),
			"azuredevops_check_business_hours":                        approvalsandchecks.ResourceCheckBusinessHours(),
			"azuredevops_check_exclusive_lock":                        approvalsandchecks.ResourceCheckExclusiveLock(),
			"azuredevops_check_required_template":                     approvalsandchecks.ResourceCheckRequiredTemplate(),
			"azuredevops_check_rest_api":                              approvalsandchecks.ResourceCheckRestAPI(),
			"azuredevops_dashboard":                                   dashboard.ResourceDashboard(),
			"azuredevops_elastic_pool":                                taskagent.ResourceAgentPoolVMSS(),
			"azuredevops_environment":                                 taskagent.ResourceEnvironment(),
			"azuredevops_environment_resource_kubernetes":             taskagent.ResourceEnvironmentKubernetes(),
			"azuredevops_feed":                                        feed.ResourceFeed(),
			"azuredevops_feed_permission":                             feed.ResourceFeedPermission(),
			"azuredevops_feed_retention_policy":                       feed.ResourceFeedRetentionPolicy(),
			"azuredevops_git_permissions":                             permissions.ResourceGitPermissions(),
			"azuredevops_git_repository":                              git.ResourceGitRepository(),
			"azuredevops_git_repository_branch":                       git.ResourceGitRepositoryBranch(),
			"azuredevops_git_repository_file":                         git.ResourceGitRepositoryFile(),
			"azuredevops_group":                                       graph.ResourceGroup(),
			"azuredevops_group_entitlement":                           memberentitlementmanagement.ResourceGroupEntitlement(),
			"azuredevops_group_membership":                            graph.ResourceGroupMembership(),
			"azuredevops_iteration_permissions":                       permissions.ResourceIterationPermissions(),
			"azuredevops_library_permissions":                         permissions.ResourceLibraryPermissions(),
			"azuredevops_pipeline_authorization":                      build.ResourcePipelineAuthorization(),
			"azuredevops_project":                                     core.ResourceProject(),
			"azuredevops_project_features":                            core.ResourceProjectFeatures(),
			"azuredevops_project_permissions":                         permissions.ResourceProjectPermissions(),
			"azuredevops_project_pipeline_settings":                   core.ResourceProjectPipelineSettings(),
			"azuredevops_project_tags":                                core.ResourceProjectTag(),
			"azuredevops_repository_policy_author_email_pattern":      repository.ResourceRepositoryPolicyAuthorEmailPatterns(),
			"azuredevops_repository_policy_case_enforcement":          repository.ResourceRepositoryEnforceConsistentCase(),
			"azuredevops_repository_policy_check_credentials":         repository.ResourceRepositoryPolicyCheckCredentials(),
			"azuredevops_repository_policy_file_path_pattern":         repository.ResourceRepositoryFilePathPatterns(),
			"azuredevops_repository_policy_max_file_size":             repository.ResourceRepositoryMaxFileSize(),
			"azuredevops_repository_policy_max_path_length":           repository.ResourceRepositoryMaxPathLength(),
			"azuredevops_repository_policy_reserved_names":            repository.ResourceRepositoryReservedNames(),
			"azuredevops_resource_authorization":                      build.ResourceResourceAuthorization(),
			"azuredevops_securityrole_assignment":                     securityroles.ResourceSecurityRoleAssignment(),
			"azuredevops_serviceendpoint_argocd":                      serviceendpoint.ResourceServiceEndpointArgoCD(),
			"azuredevops_serviceendpoint_artifactory":                 serviceendpoint.ResourceServiceEndpointArtifactory(),
			"azuredevops_serviceendpoint_aws":                         serviceendpoint.ResourceServiceEndpointAws(),
			"azuredevops_serviceendpoint_azure_service_bus":           serviceendpoint.ResourceServiceEndpointAzureServiceBus(),
			"azuredevops_serviceendpoint_azurecr":                     serviceendpoint.ResourceServiceEndpointAzureCR(),
			"azuredevops_serviceendpoint_azuredevops":                 serviceendpoint.ResourceServiceEndpointAzureDevOps(),
			"azuredevops_serviceendpoint_azurerm":                     serviceendpoint.ResourceServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_bitbucket":                   serviceendpoint.ResourceServiceEndpointBitBucket(),
			"azuredevops_serviceendpoint_checkmarx_one":               serviceendpoint.ResourceServiceEndpointCheckMarxOneService(),
			"azuredevops_serviceendpoint_checkmarx_sca":               serviceendpoint.ResourceServiceEndpointCheckMarxSCA(),
			"azuredevops_serviceendpoint_checkmarx_sast":              serviceendpoint.ResourceServiceEndpointCheckMarxSAST(),
			"azuredevops_serviceendpoint_dockerregistry":              serviceendpoint.ResourceServiceEndpointDockerRegistry(),
			"azuredevops_serviceendpoint_dynamics_lifecycle_services": serviceendpoint.ResourceServiceEndpointDynamicsLifecycleServices(),
			"azuredevops_serviceendpoint_externaltfs":                 serviceendpoint.ResourceServiceEndpointExternalTFS(),
			"azuredevops_serviceendpoint_gcp_terraform":               serviceendpoint.ResourceServiceEndpointGcp(),
			"azuredevops_serviceendpoint_generic":                     serviceendpoint.ResourceServiceEndpointGeneric(),
			"azuredevops_serviceendpoint_generic_git":                 serviceendpoint.ResourceServiceEndpointGenericGit(),
			"azuredevops_serviceendpoint_github":                      serviceendpoint.ResourceServiceEndpointGitHub(),
			"azuredevops_serviceendpoint_github_enterprise":           serviceendpoint.ResourceServiceEndpointGitHubEnterprise(),
			"azuredevops_serviceendpoint_gitlab":                      serviceendpoint.ResourceServiceEndpointGitLab(),
			"azuredevops_serviceendpoint_incomingwebhook":             serviceendpoint.ResourceServiceEndpointIncomingWebhook(),
			"azuredevops_serviceendpoint_jenkins":                     serviceendpoint.ResourceServiceEndpointJenkins(),
			"azuredevops_serviceendpoint_jfrog_artifactory_v2":        serviceendpoint.ResourceServiceEndpointJFrogArtifactoryV2(),
			"azuredevops_serviceendpoint_jfrog_distribution_v2":       serviceendpoint.ResourceServiceEndpointJFrogDistributionV2(),
			"azuredevops_serviceendpoint_jfrog_platform_v2":           serviceendpoint.ResourceServiceEndpointJFrogPlatformV2(),
			"azuredevops_serviceendpoint_jfrog_xray_v2":               serviceendpoint.ResourceServiceEndpointJFrogXRayV2(),
			"azuredevops_serviceendpoint_kubernetes":                  serviceendpoint.ResourceServiceEndpointKubernetes(),
			"azuredevops_serviceendpoint_maven":                       serviceendpoint.ResourceServiceEndpointMaven(),
			"azuredevops_serviceendpoint_nexus":                       serviceendpoint.ResourceServiceEndpointNexus(),
			"azuredevops_serviceendpoint_npm":                         serviceendpoint.ResourceServiceEndpointNpm(),
			"azuredevops_serviceendpoint_nuget":                       serviceendpoint.ResourceServiceEndpointNuGet(),
			"azuredevops_serviceendpoint_octopusdeploy":               serviceendpoint.ResourceServiceEndpointOctopusDeploy(),
			"azuredevops_serviceendpoint_permissions":                 permissions.ResourceServiceEndpointPermissions(),
			"azuredevops_serviceendpoint_runpipeline":                 serviceendpoint.ResourceServiceEndpointRunPipeline(),
			"azuredevops_serviceendpoint_servicefabric":               serviceendpoint.ResourceServiceEndpointServiceFabric(),
			"azuredevops_serviceendpoint_snyk":                        serviceendpoint.ResourceServiceEndpointSnyk(),
			"azuredevops_serviceendpoint_sonarcloud":                  serviceendpoint.ResourceServiceEndpointSonarCloud(),
			"azuredevops_serviceendpoint_sonarqube":                   serviceendpoint.ResourceServiceEndpointSonarQube(),
			"azuredevops_serviceendpoint_ssh":                         serviceendpoint.ResourceServiceEndpointSSH(),
			"azuredevops_serviceendpoint_visualstudiomarketplace":     serviceendpoint.ResourceServiceEndpointMarketplace(),
			"azuredevops_servicehook_permissions":                     permissions.ResourceServiceHookPermissions(),
			"azuredevops_servicehook_storage_queue_pipelines":         servicehook.ResourceServicehookStorageQueuePipelines(),
			"azuredevops_service_principal_entitlement":               memberentitlementmanagement.ResourceServicePrincipalEntitlement(),
			"azuredevops_tagging_permissions":                         permissions.ResourceTaggingPermissions(),
			"azuredevops_team":                                        core.ResourceTeam(),
			"azuredevops_team_administrators":                         core.ResourceTeamAdministrators(),
			"azuredevops_team_members":                                core.ResourceTeamMembers(),
			"azuredevops_user_entitlement":                            memberentitlementmanagement.ResourceUserEntitlement(),
			"azuredevops_variable_group":                              taskagent.ResourceVariableGroup(),
			"azuredevops_variable_group_permissions":                  permissions.ResourceVariableGroupPermissions(),
			"azuredevops_wiki":                                        wiki.ResourceWiki(),
			"azuredevops_wiki_page":                                   wiki.ResourceWikiPage(),
			"azuredevops_workitem":                                    workitemtracking.ResourceWorkItem(),
			"azuredevops_workitemquery_permissions":                   permissions.ResourceWorkItemQueryPermissions(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azuredevops_agent_pool":                 taskagent.DataAgentPool(),
			"azuredevops_agent_pools":                taskagent.DataAgentPools(),
			"azuredevops_agent_queue":                taskagent.DataAgentQueue(),
			"azuredevops_area":                       workitemtracking.DataArea(),
			"azuredevops_build_definition":           build.DataBuildDefinition(),
			"azuredevops_client_config":              service.DataClientConfig(),
			"azuredevops_descriptor":                 graph.DataDescriptor(),
			"azuredevops_environment":                taskagent.DataEnvironment(),
			"azuredevops_feed":                       feed.DataFeed(),
			"azuredevops_git_repositories":           git.DataGitRepositories(),
			"azuredevops_git_repository":             git.DataGitRepository(),
			"azuredevops_group":                      graph.DataGroup(),
			"azuredevops_group_membership":           graph.DataGroupMembership(),
			"azuredevops_groups":                     graph.DataGroups(),
			"azuredevops_identity_group":             identity.DataIdentityGroup(),
			"azuredevops_identity_groups":            identity.DataIdentityGroups(),
			"azuredevops_identity_user":              identity.DataIdentityUser(),
			"azuredevops_iteration":                  workitemtracking.DataIteration(),
			"azuredevops_project":                    core.DataProject(),
			"azuredevops_projects":                   core.DataProjects(),
			"azuredevops_securityrole_definitions":   securityroles.DataSecurityRoleDefinitions(),
			"azuredevops_serviceendpoint_azurecr":    serviceendpoint.DataResourceServiceEndpointAzureCR(),
			"azuredevops_serviceendpoint_azurerm":    serviceendpoint.DataServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_bitbucket":  serviceendpoint.DataResourceServiceEndpointBitbucket(),
			"azuredevops_serviceendpoint_github":     serviceendpoint.DataServiceEndpointGithub(),
			"azuredevops_serviceendpoint_npm":        serviceendpoint.DataResourceServiceEndpointNpm(),
			"azuredevops_serviceendpoint_sonarcloud": serviceendpoint.DataResourceServiceEndpointSonarCloud(),
			"azuredevops_service_principal":          graph.DataServicePrincipal(),
			"azuredevops_storage_key":                graph.DataStorageKey(),
			"azuredevops_team":                       core.DataTeam(),
			"azuredevops_teams":                      core.DataTeams(),
			"azuredevops_user":                       graph.DataUser(),
			"azuredevops_users":                      graph.DataUsers(),
			"azuredevops_variable_group":             taskagent.DataVariableGroup(),
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
				Description:  "The service principal client or managed service principal id which should be used.",
				ValidateFunc: validation.IsUUID,
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
				Description:  "The service principal tenant id which should be used.",
				ValidateFunc: validation.IsUUID,
			},
			"client_id_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID_PLAN", nil),
				Description:  "The service principal client id which should be used during a plan operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
			},
			"tenant_id_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID_PLAN", nil),
				Description:  "The service principal tenant id which should be used during a plan operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
			},
			"client_id_apply": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID_APPLY", nil),
				Description:  "The service principal client id which should be used during an apply operation in Terraform Cloud.",
				ValidateFunc: validation.IsUUID,
			},
			"tenant_id_apply": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID_APPLY", nil),
				Description:  "The service principal tenant id which should be used during an apply operation in Terraform Cloud..",
				ValidateFunc: validation.IsUUID,
			},
			"oidc_request_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}, nil),
				Description: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_request_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"}, nil),
				Description: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN", nil),
				Description: "OIDC token to authenticate as a service principal.",
			},
			"oidc_token_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN_FILE_PATH", nil),
				Description: "OIDC token from file to authenticate as a service principal.",
			},
			"use_oidc": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_OIDC", nil),
				Description: "Use an OIDC token to authenticate to a service principal.",
			},
			"oidc_audience": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_AUDIENCE", nil),
				Description: "Set the audience when requesting OIDC tokens.",
			},
			"oidc_tfc_tag": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TFC_TAG", nil),
				Description: "Terraform Cloud dynamic credential provider tag.",
			},
			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", nil),
				Description: "Path to a certificate to use to authenticate to the service principal.",
			},
			"client_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE", nil),
				Description: "Base64 encoded certificate to use to authenticate to the service principal.",
			},
			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", nil),
				Description: "Password for a client certificate password.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", nil),
				Description: "Client secret for authenticating to  a service principal.",
			},
			"client_secret_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET_PATH", nil),
				Description: "Path to a file containing a client secret for authenticating to  a service principal.",
			},
			"use_msi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", nil),
				Description: "Use an Azure Managed Service Identity.",
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

		tokenFunction, err := sdk.GetAuthTokenProvider(ctx, d, sdk.AzIdentityFuncsImpl{})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		organizationUrl := d.Get("org_service_url").(string)
		azdoClient, err := client.GetAzdoClient(tokenFunction, organizationUrl, terraformVersion)

		if err != nil {
			return nil, diag.FromErr(clientErrorHandle(err, organizationUrl))
		}

		return azdoClient, diag.FromErr(err)
	}
}

func clientErrorHandle(err error, orgUrl string) error {
	switch err.(type) {
	case azuredevops.WrappedError:
		var wrapperError azuredevops.WrappedError
		if errors.As(err, &wrapperError) {
			if clientError := buildError(*wrapperError.StatusCode, orgUrl); clientError != nil {
				return clientError
			}
		}

	case *azuredevops.WrappedError:
		var wrapperError *azuredevops.WrappedError
		if errors.As(err, &wrapperError) {
			if clientError := buildError(*wrapperError.StatusCode, orgUrl); clientError != nil {
				return clientError
			}
		}
	}
	return err
}

func buildError(statusCode int, orgUrl string) error {
	if statusCode == http.StatusNotFound {
		return fmt.Errorf(" Azure DevOps Organization: %s doesn't exist or can't be found. Make sure the URL is correct.", orgUrl)
	} else if statusCode == http.StatusUnauthorized {
		return fmt.Errorf(" You are not authorized to access Azure DevOps Organization %s", orgUrl)
	}
	return nil
}
