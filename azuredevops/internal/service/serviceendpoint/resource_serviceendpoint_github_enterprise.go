package serviceendpoint

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	personalAccessTokenGithubEnterprise = "personal_access_token"
)

// ResourceServiceEndpointGitHubEnterprise schema and implementation for github-enterprise service endpoint resource
func ResourceServiceEndpointGitHubEnterprise() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointGitHubEnterprise, expandServiceEndpointGitHubEnterprise)
	authPersonal := &schema.Resource{
		Schema: map[string]*schema.Schema{
			personalAccessTokenGithubEnterprise: {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_GITHUB_ENTERPRISE_SERVICE_CONNECTION_PAT", nil),
				Description:  "The GitHub personal access token which should be used.",
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(personalAccessTokenGithubEnterprise)
	authPersonal.Schema[patHashKey] = patHashSchema
	r.Schema["auth_personal"] = &schema.Schema{
		Type:     schema.TypeSet,
		MinItems: 1,
		MaxItems: 1,
		Elem:     authPersonal,
		Required: true,
	}

	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Computed:     false,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Required:     true,
	}

	return r
}

func flattenServiceEndpointGitHubEnterprise(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
		authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
		authPersonal := flattenAuthPersonGithubEnterprise(d, authPersonalSet)
		if authPersonal != nil {
			d.Set("auth_personal", authPersonal)
		}
	}

	d.Set("type", "githubenterprise")
	d.Set("url", *serviceEndpoint.Url)
}

func flattenAuthPersonGithubEnterprise(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "auth_personal", authPersonal, personalAccessTokenGithub)
			authPersonal[hashKey] = newHash
			return []interface{}{authPersonal}
		}
	}
	return nil
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGitHubEnterprise(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)

	serviceEndpoint.Type = converter.String("githubenterprise")

	seUrl := d.Get("url").(string)
	serviceEndpoint.Url = converter.String(seUrl)

	scheme := "InstallationToken"
	parameters := map[string]string{}

	if config, ok := d.GetOk("auth_personal"); ok {
		scheme = "Token"
		parameters = expandAuthPersonalSetGithubEnterprise(config.(*schema.Set))
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}

	return serviceEndpoint, projectID, nil
}

func expandAuthPersonalSetGithubEnterprise(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure

	authPerson["apitoken"] = val[personalAccessTokenGithub].(string)
	return authPerson
}
