package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_build_definition":     resourceBuildDefinition(),
			"azuredevops_project":              resourceProject(),
			"azuredevops_serviceendpoint":      resourceServiceEndpoint(),
			"azuredevops_azure_git_repository": resourceAzureGitRepository(),
			"azuredevops_user_entitlement":     resourceUserEntitlement(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azuredevops_group": dataGroup(),
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
		client, err := getAzdoClient(d.Get("personal_access_token").(string), d.Get("org_service_url").(string))
		return client, err
	}
}
