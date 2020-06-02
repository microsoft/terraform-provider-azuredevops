package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_resource_authorization":         resourceResourceAuthorization(),
			"azuredevops_branch_policy_build_validation": resourceBranchPolicyBuildValidation(),
			"azuredevops_branch_policy_min_reviewers":    resourceBranchPolicyMinReviewers(),
			"azuredevops_build_definition":               resourceBuildDefinition(),
			"azuredevops_project":                        resourceProject(),
			"azuredevops_variable_group":                 resourceVariableGroup(),
			"azuredevops_serviceendpoint_azurerm":        resourceServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_bitbucket":      resourceServiceEndpointBitBucket(),
			"azuredevops_serviceendpoint_dockerregistry": resourceServiceEndpointDockerRegistry(),
			"azuredevops_serviceendpoint_github":         resourceServiceEndpointGitHub(),
			"azuredevops_serviceendpoint_kubernetes":     resourceServiceEndpointKubernetes(),
			"azuredevops_git_repository":                 resourceGitRepository(),
			"azuredevops_user_entitlement":               resourceUserEntitlement(),
			"azuredevops_group_membership":               resourceGroupMembership(),
			"azuredevops_agent_pool":                     resourceAzureAgentPool(),
			"azuredevops_agent_queue":                    resourceAgentQueue(),
			"azuredevops_group":                          resourceGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azuredevops_agent_pool":       dataAzureAgentPool(),
			"azuredevops_agent_pools":      dataAzureAgentPools(),
			"azuredevops_client_config":    dataClientConfig(),
			"azuredevops_group":            dataGroup(),
			"azuredevops_project":          dataProject(),
			"azuredevops_projects":         dataProjects(),
			"azuredevops_git_repositories": dataGitRepositories(),
			"azuredevops_users":            dataUsers(),
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

		client, err := config.GetAzdoClient(d.Get("personal_access_token").(string), d.Get("org_service_url").(string), terraformVersion)

		return client, err
	}
}
