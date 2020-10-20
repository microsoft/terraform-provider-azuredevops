package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointAppcenter schema and implementation for bitbucket service endpoint resource
func ResourceServiceEndpointAppcenter() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAppcenter, expandServiceEndpointAppcenter)
	makeProtectedSchema(r, "apitoken", "AZDO_APPCENTER_SERVICE_CONNECTION_API_TOKEN", "The API token to connect to app center.")
	return r
}

func expandServiceEndpointAppcenter(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("apitoken").(string),
		},
		Scheme: converter.String("Token"),
	}
	serviceEndpoint.Type = converter.String(string("vsmobilecenter"))
	serviceEndpoint.Url = converter.String("https://api.appcenter.ms/v0.1")
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointAppcenter(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	tfhelper.HelpFlattenSecret(d, "apitoken")
	d.Set("apitoken", (*serviceEndpoint.Authorization.Parameters)["apitoken"])
}
