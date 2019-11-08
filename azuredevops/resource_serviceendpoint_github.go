package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"

	crud "github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
)

const (
	githubSchemaKey = "github_service_endpoint_pat"
)

func resourceServiceEndpointGitHub() *schema.Resource {
	r := crud.GenBaseServiceEndpointResource(flattenServiceEndpointGitHub, expandServiceEndpointGitHub)
	r.Schema[githubSchemaKey] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_GITHUB_SERVICE_CONNECTION_PAT", nil),
		Description:      "The GitHub personal access token which should be used.",
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSupressSecretChanged,
	}

	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(githubSchemaKey)
	r.Schema[patHashKey] = patHashSchema

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGitHub(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	serviceEndpoint, projectID := crud.DoBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"accessToken": d.Get(githubSchemaKey).(string),
		},
		Scheme: converter.String("PersonalAccessToken"),
	}
	serviceEndpoint.Type = converter.String("github")
	serviceEndpoint.Url = converter.String("http://github.com")

	return serviceEndpoint, projectID
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointGitHub(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	crud.DoBaseFlattening(d, serviceEndpoint, projectID)
	tfhelper.HelpFlattenSecret(d, githubSchemaKey)
	d.Set(githubSchemaKey, (*serviceEndpoint.Authorization.Parameters)["accessToken"])
}
