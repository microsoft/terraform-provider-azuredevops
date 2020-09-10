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
	azdoPersonalAccessToken = "personal_access_token"
)

// ResourceServiceEndpointAzureDevOps schema and implementation for Azure DevOps service endpoint resource
func ResourceServiceEndpointAzureDevOps() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAzureDevOps, expandServiceEndpointAzureDevOps)
	authPersonal := &schema.Resource{
		Schema: map[string]*schema.Schema{
			personalAccessToken: {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_SERVICE_CONNECTION_PAT", nil),
				Description:  "The Azure DevOps personal access token which should be used.",
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(azdoPersonalAccessToken)
	authPersonal.Schema[patHashKey] = patHashSchema
	r.Schema["auth_personal"] = &schema.Schema{
		Type:          schema.TypeSet,
		Optional:      true,
		MinItems:      1,
		MaxItems:      1,
		Elem:          authPersonal,
		ConflictsWith: []string{},
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAzureDevOps(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	scheme := "InstallationToken"

	parameters := map[string]string{}

	if config, ok := d.GetOk("auth_personal"); ok {
		scheme = "Token"
		parameters = azdoExpandAuthPersonalSet(config.(*schema.Set))
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}
	serviceEndpoint.Type = converter.String("Azure DevOps")
	serviceEndpoint.Url = converter.String("https://dev.azure.com/")

	return serviceEndpoint, projectID, nil
}

func azdoExpandAuthPersonalSet(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure
	authPerson["AccessToken"] = val[azdoPersonalAccessToken].(string)
	return authPerson
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAzureDevOps(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
		authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
		authPersonal := azdoFlattenAuthPerson(d, authPersonalSet)
		if authPersonal != nil {
			d.Set("auth_personal", authPersonal)
		}
	}
}

func azdoFlattenAuthPerson(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "auth_personal", authPersonal, azdoPersonalAccessToken)
			authPersonal[hashKey] = newHash
			return []interface{}{authPersonal}
		}
	}
	return nil
}
