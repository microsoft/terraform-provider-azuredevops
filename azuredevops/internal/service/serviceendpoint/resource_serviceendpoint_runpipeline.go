package serviceendpoint

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointRunPipeline schema and implementation for Azure DevOps service endpoint resource
func ResourceServiceEndpointRunPipeline() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointRunPipeline, expandServiceEndpointRunPipeline)
	r.Schema["organization_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Azure DevOps organization name",
	}

	r.Schema["auth_personal"] = &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem:     rpPersonalAccessTokenField(),
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure:
func expandServiceEndpointRunPipeline(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("azdoapi")

	scheme := "Token"
	parameters := rpExpandAuthPersonalSet(d.Get("auth_personal").(*schema.Set))

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}

	org := d.Get("organization_name").(string)
	serviceUrl := fmt.Sprint("https://dev.azure.com/", org)
	serviceEndpoint.Url = &serviceUrl

	data := map[string]string{}
	releaseUrl := fmt.Sprint("https://vsrm.dev.azure.com/", org)
	data["releaseUrl"] = releaseUrl
	serviceEndpoint.Data = &data
	return serviceEndpoint, projectID
}

func rpExpandAuthPersonalSet(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	if len(d.List()) == 1 {
		val := d.List()[0].(map[string]interface{}) //auth_personal block may have only one element inside
		authPerson["apitoken"] = val["personal_access_token"].(string)
	}
	return authPerson
}

func rpPersonalAccessTokenField() *schema.Resource {
	fieldName := "personal_access_token"
	personalAccessToken := &schema.Resource{
		Schema: map[string]*schema.Schema{
			fieldName: {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_PERSONAL_ACCESS_TOKEN", nil),
				Description:  "The Azure DevOps personal access token which should be used.",
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(fieldName)
	personalAccessToken.Schema[patHashKey] = patHashSchema

	return personalAccessToken
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointRunPipeline(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
	authPersonal := rpFlattenAuthPersonal(d, authPersonalSet)
	if authPersonal != nil {
		d.Set("auth_personal", authPersonal)
	}
}

func rpFlattenAuthPersonal(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "auth_personal", authPersonal, "personal_access_token")
			authPersonal[hashKey] = newHash
			return []interface{}{authPersonal}
		}
	}
	return nil
}
