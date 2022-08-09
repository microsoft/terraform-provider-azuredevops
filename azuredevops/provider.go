package azuredevops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service"
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

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
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
			"azuredevops_serviceendpoint_argocd":                 serviceendpoint.ResourceServiceEndpointArgoCD(),
			"azuredevops_serviceendpoint_artifactory":            serviceendpoint.ResourceServiceEndpointArtifactory(),
			"azuredevops_serviceendpoint_aws":                    serviceendpoint.ResourceServiceEndpointAws(),
			"azuredevops_serviceendpoint_azurerm":                serviceendpoint.ResourceServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_bitbucket":              serviceendpoint.ResourceServiceEndpointBitBucket(),
			"azuredevops_serviceendpoint_azuredevops":            serviceendpoint.ResourceServiceEndpointAzureDevOps(),
			"azuredevops_serviceendpoint_custom":                serviceendpoint.ResourceServiceEndpointCustom(),
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
			"azuredevops_serviceendpoint_generic":                serviceendpoint.ResourceServiceEndpointGeneric(),
			"azuredevops_serviceendpoint_generic_git":            serviceendpoint.ResourceServiceEndpointGenericGit(),
			"azuredevops_git_repository":                         git.ResourceGitRepository(),
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
			"azuredevops_environment":                            taskagent.ResourceEnvironment(),
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

		client, err := client.GetAzdoClient(d.Get("personal_access_token").(string), d.Get("org_service_url").(string), terraformVersion)

		return client, diag.FromErr(err)
	}
}
