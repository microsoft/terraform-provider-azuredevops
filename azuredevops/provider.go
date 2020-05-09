package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_resource_authorization":     resourceResourceAuthorization(),
			"azuredevops_build_definition":           resourceBuildDefinition(),
			"azuredevops_project":                    resourceProject(),
			"azuredevops_variable_group":             resourceVariableGroup(),
			"azuredevops_serviceendpoint_azurerm":    resourceServiceEndpointAzureRM(),
			"azuredevops_serviceendpoint_bitbucket":  resourceServiceEndpointBitBucket(),
			"azuredevops_serviceendpoint_dockerhub":  resourceServiceEndpointDockerHub(),
			"azuredevops_serviceendpoint_github":     resourceServiceEndpointGitHub(),
			"azuredevops_serviceendpoint_kubernetes": resourceServiceEndpointKubernetes(),
			"azuredevops_git_repository":             resourceGitRepository(),
			"azuredevops_user_entitlement":           resourceUserEntitlement(),
			"azuredevops_group_membership":           resourceGroupMembership(),
			"azuredevops_agent_pool":                 resourceAzureAgentPool(),
			"azuredevops_group":                      resourceGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azuredevops_group":            dataGroup(),
			"azuredevops_project":          dataProject(),
			"azuredevops_projects":         dataProjects(),
			"azuredevops_git_repositories": dataGitRepositories(),
			"azuredevops_users":            dataUsers(),
		},
		Schema: map[string]*schema.Schema{
			"org_service_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_ORG_SERVICE_URL", nil),
				Description: "The url of the Azure DevOps instance which should be used.",
			},
			"personal_access_token": {
				Type:        schema.TypeString,
				Required:    true,
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
		client, err := config.GetAzdoClient(d.Get("personal_access_token").(string), d.Get("org_service_url").(string))
		return client, err
	}
}
