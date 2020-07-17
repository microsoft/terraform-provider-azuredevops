package serviceendpoint

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	personalAccessToken = "personal_access_token"
)

// ResourceServiceEndpointGitHub schema and implementation for github service endpoint resource
func ResourceServiceEndpointGitHub() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointGitHub, expandServiceEndpointGitHub)
	authPersonal := &schema.Resource{
		Schema: map[string]*schema.Schema{
			personalAccessToken: {
				Type:             schema.TypeString,
				Required:         true,
				DefaultFunc:      schema.EnvDefaultFunc("AZDO_GITHUB_SERVICE_CONNECTION_PAT", nil),
				Description:      "The GitHub personal access token which should be used.",
				Sensitive:        true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
			},
		},
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(personalAccessToken)
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

func expandAuthOauth(d map[string]interface{}) map[string]string {
	return map[string]string{
		"ConfigurationId": d["oauth_configuration_id"].(string),
	}
}
func expandAuthOauthList(d []interface{}) []map[string]string {
	vs := make([]map[string]string, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandAuthOauth(val))
		}
	}
	return vs
}
func expandAuthOauthSet(d *schema.Set) map[string]string {
	d2 := expandAuthOauthList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandAuthPersonal(d map[string]interface{}) map[string]string {
	return map[string]string{
		"accessToken": d[personalAccessToken].(string),
	}
}
func expandAuthPersonalList(d []interface{}) []map[string]string {
	vs := make([]map[string]string, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandAuthPersonal(val))
		}
	}
	return vs
}
func expandAuthPersonalSet(d *schema.Set) map[string]string {
	d2 := expandAuthPersonalList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGitHub(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	scheme := "InstallationToken"

	parameters := &map[string]string{}
	authPersonal := expandAuthPersonalSet(d.Get("auth_personal").(*schema.Set))
	authGrant := expandAuthOauthSet(d.Get("auth_oauth").(*schema.Set))

	if authPersonal != nil {
		scheme = "PersonalAccessToken"
		parameters = &authPersonal
	}

	if authGrant != nil {
		scheme = "OAuth"
		parameters = &authGrant
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: parameters,
		Scheme:     &scheme,
	}
	serviceEndpoint.Type = converter.String("github")
	serviceEndpoint.Url = converter.String("http://github.com")

	return serviceEndpoint, projectID, nil
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
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "PersonalAccessToken") {
		authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
		authPersonal := flattenAuthPerson(d, authPersonalSet)
		if authPersonal != nil {
			d.Set("auth_personal", authPersonal)
		}
	}
}

func flattenAuthPerson(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "auth_personal", authPersonal, personalAccessToken)
			authPersonal[hashKey] = newHash
			return []interface{}{authPersonal}
		}
	}
	return nil
}
