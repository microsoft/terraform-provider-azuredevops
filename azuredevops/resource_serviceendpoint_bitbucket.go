package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	crud "github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
)

func resourceServiceEndpointBitBucket() *schema.Resource {
	r := crud.GenBaseServiceEndpointResource(flattenServiceEndpointBitBucket, expandServiceEndpointBitBucket, parseImportedProjectIDAndServiceEndpointID)
	crud.MakeUnprotectedSchema(r, "username", "AZDO_BITBUCKET_SERVICE_CONNECTION_USERNAME", "The bitbucket username which should be used.")
	crud.MakeProtectedSchema(r, "password", "AZDO_BITBUCKET_SERVICE_CONNECTION_PASSWORD", "The bitbucket password whi|ch should be used.")
	return r
}

func expandServiceEndpointBitBucket(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	serviceEndpoint, projectID := crud.DoBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String(string(RepoTypeValues.Bitbucket))
	serviceEndpoint.Url = converter.String("https://api.bitbucket.org/")
	return serviceEndpoint, projectID
}

func flattenServiceEndpointBitBucket(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	crud.DoBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	tfhelper.HelpFlattenSecret(d, "password")
	d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
}
