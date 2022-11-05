package serviceendpoint

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	personalAccessTokenExternalTFS = "personal_access_token"
)

func ResourceServiceEndpointExternalTFS() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointExternalTFS, expandServiceEndpointExternalTFS)
	r.Schema["connection_url"] = &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Required:     true,
		Description:  "URL of the Azure DevOps organization or the TFS Project Collection to connect to.",
	}
	authPersonal := &schema.Resource{
		Schema: map[string]*schema.Schema{
			personalAccessTokenExternalTFS: {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_PERSONAL_ACCESS_TOKEN", nil),
				Description:  "Personal access tokens are applicable only for connections targeting Azure DevOps organization or TFS 2017 (and higher)",
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(personalAccessTokenExternalTFS)
	authPersonal.Schema[patHashKey] = patHashSchema
	r.Schema["auth_personal"] = &schema.Schema{
		Type:     schema.TypeSet,
		MinItems: 1,
		MaxItems: 1,
		Elem:     authPersonal,
		Required: true,
	}
	return r
}

func expandServiceEndpointExternalTFS(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("externaltfs")
	serviceEndpoint.Url = converter.String(d.Get("connection_url").(string))

	scheme := "Token"
	parameters := map[string]string{}

	if config, ok := d.GetOk("auth_personal"); ok {
		scheme = "Token"
		parameters = expandAuthPersonalSetExternalTFS(config.(*schema.Set))
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}
	return serviceEndpoint, projectID, nil
}

func expandAuthPersonalSetExternalTFS(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	val := d.List()[0].(map[string]interface{})

	authPerson["apitoken"] = val[personalAccessTokenExternalTFS].(string)
	return authPerson
}

func flattenServiceEndpointExternalTFS(
	d *schema.ResourceData,
	serviceEndpoint *serviceendpoint.ServiceEndpoint,
	projectID *uuid.UUID,
) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
		authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
		authPersonal := flattenAuthPersonExternalTFS(d, authPersonalSet)
		if authPersonal != nil {
			d.Set("auth_personal", authPersonal)
		}
	}

	d.Set("connection_url", *serviceEndpoint.Url)
}

func flattenAuthPersonExternalTFS(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			newHash, hashKey := tfhelper.HelpFlattenSecretNested(
				d,
				"auth_personal",
				authPersonal,
				personalAccessTokenExternalTFS,
			)
			authPersonal[hashKey] = newHash
			return []interface{}{authPersonal}
		}
	}
	return nil
}
