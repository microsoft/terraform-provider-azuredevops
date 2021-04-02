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
	personalAccessTokenGithub = "personal_access_token"
)

// ResourceServiceEndpointGitHub schema and implementation for github service endpoint resource
func ResourceServiceEndpointGitHub() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointGitHub, expandServiceEndpointGitHub)
	authPersonal := &schema.Resource{
		Schema: map[string]*schema.Schema{
			personalAccessTokenGithub: {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_GITHUB_SERVICE_CONNECTION_PAT", nil),
				Description:  "The GitHub personal access token which should be used.",
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(personalAccessTokenGithub)
	authPersonal.Schema[patHashKey] = patHashSchema
	r.Schema["auth_personal"] = &schema.Schema{
		Type:          schema.TypeSet,
		Optional:      true,
		MinItems:      1,
		MaxItems:      1,
		Elem:          authPersonal,
		ConflictsWith: []string{"auth_oauth"},
	}

	r.Schema["auth_oauth"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oauth_configuration_id": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
		ConflictsWith: []string{"auth_personal"},
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGitHub(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)

	serviceEndpoint.Type = converter.String("github")
	serviceEndpoint.Url = converter.String("https://github.com")

	scheme := "InstallationToken"
	parameters := map[string]string{}

	if config, ok := d.GetOk("auth_personal"); ok {
		scheme = "Token"
		parameters = expandAuthPersonalSetGithub(config.(*schema.Set))
	}

	if config, ok := d.GetOk("auth_oauth"); ok {
		scheme = "OAuth"
		parameters = expandAuthOauthSet(config.(*schema.Set))
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}

	return serviceEndpoint, projectID, nil
}

func expandAuthPersonalSetGithub(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure
	authPerson["AccessToken"] = val[personalAccessTokenGithub].(string)
	return authPerson
}

func expandAuthOauthSet(d *schema.Set) map[string]string {
	authConfig := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure
	authConfig["ConfigurationId"] = val["oauth_configuration_id"].(string)
	authConfig["AccessToken"] = ""
	return authConfig
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointGitHub(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "OAuth") {
		d.Set("auth_oauth", &[]map[string]interface{}{
			{
				"oauth_configuration_id": (*serviceEndpoint.Authorization.Parameters)["ConfigurationId"],
			},
		})
	}
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
		authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
		authPersonal := flattenAuthPerson(d, authPersonalSet)
		if authPersonal != nil {
			d.Set("auth_personal", authPersonal)
		}
	}

	d.Set("type", *serviceEndpoint.Type)
	d.Set("url", *serviceEndpoint.Url)
}

func flattenAuthPerson(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "auth_personal", authPersonal, personalAccessTokenGithub)
			authPersonal[hashKey] = newHash
			return []interface{}{authPersonal}
		}
	}
	return nil
}
