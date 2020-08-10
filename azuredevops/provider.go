package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/build"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/core"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/git"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/graph"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/memberentitlementmanagement"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/policy"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/serviceendpoint"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/taskagent"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/workitemtracking"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_resource_authorization":          build.ResourceResourceAuthorization(),
			"azuredevops_branch_policy_build_validation":  policy.ResourceBranchPolicyBuildValidation(),
			"azuredevops_branch_policy_min_reviewers":     policy.ResourceBranchPolicyMinReviewers(),
			"azuredevops_branch_policy_auto_reviewers":    policy.ResourceBranchPolicyAutoReviewers(),
			"azuredevops_branch_policy_work_item_linking": policy.ResourceBranchPolicyWorkItemLinking(),
			"azuredevops_build_definition":                build.ResourceBuildDefinition(),
			"azuredevops_project":                         core.ResourceProject(),
			"azuredevops_project_features":                core.ResourceProjectFeatures(),
			"azuredevops_variable_group":                  taskagent.ResourceVariableGroup(),
			"azuredevops_serviceendpoint_aws":             serviceendpoint.ResourceServiceEndpointAws(),
			"azuredevops_serviceendpoint_azurerm":         serviceendpoint.ResourceServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_bitbucket":       serviceendpoint.ResourceServiceEndpointBitBucket(),
			"azuredevops_serviceendpoint_dockerregistry":  serviceendpoint.ResourceServiceEndpointDockerRegistry(),
			"azuredevops_serviceendpoint_azurecr":         serviceendpoint.ResourceServiceEndpointAzureCR(),
			"azuredevops_serviceendpoint_github":          serviceendpoint.ResourceServiceEndpointGitHub(),
			"azuredevops_serviceendpoint_kubernetes":      serviceendpoint.ResourceServiceEndpointKubernetes(),
			"azuredevops_git_repository":                  git.ResourceGitRepository(),
			"azuredevops_user_entitlement":                memberentitlementmanagement.ResourceUserEntitlement(),
			"azuredevops_group_membership":                graph.ResourceGroupMembership(),
			"azuredevops_agent_pool":                      taskagent.ResourceAgentPool(),
			"azuredevops_agent_queue":                     taskagent.ResourceAgentQueue(),
			"azuredevops_group":                           graph.ResourceGroup(),
			"azuredevops_project_permissions":             permissions.ResourceProjectPermissions(),
			"azuredevops_git_permissions":                 permissions.ResourceGitPermissions(),
			"azuredevops_workitemquery_permissions":       permissions.ResourceWorkItemQueryPermissions(),
			"azuredevops_area_permissions":                permissions.ResourceAreaPermissions(),
			"azuredevops_iteration_permissions":           permissions.ResourceIterationPermissions(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azuredevops_agent_pool":       taskagent.DataAgentPool(),
			"azuredevops_agent_pools":      taskagent.DataAgentPools(),
			"azuredevops_client_config":    service.DataClientConfig(),
			"azuredevops_group":            graph.DataGroup(),
			"azuredevops_project":          core.DataProject(),
			"azuredevops_projects":         core.DataProjects(),
			"azuredevops_git_repositories": git.DataGitRepositories(),
			"azuredevops_git_repository":   git.DataGitRepository(),
			"azuredevops_users":            graph.DataUsers(),
			"azuredevops_area":             workitemtracking.DataArea(),
			"azuredevops_iteration":        workitemtracking.DataIteration(),
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

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		client, err := client.GetAzdoClient(d.Get("personal_access_token").(string), d.Get("org_service_url").(string), terraformVersion)

		return client, err
	}
}
