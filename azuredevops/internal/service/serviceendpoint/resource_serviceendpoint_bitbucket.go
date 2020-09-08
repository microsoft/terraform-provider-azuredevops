package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/model"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointBitBucket schema and implementation for bitbucket service endpoint resource
func ResourceServiceEndpointBitBucket() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointBitBucket, expandServiceEndpointBitBucket)
	makeUnprotectedSchema(r, "username", "AZDO_BITBUCKET_SERVICE_CONNECTION_USERNAME", "The bitbucket username which should be used.")
	makeProtectedSchema(r, "password", "AZDO_BITBUCKET_SERVICE_CONNECTION_PASSWORD", "The bitbucket password whi|ch should be used.")
	return r
}

func expandServiceEndpointBitBucket(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String(string(model.RepoTypeValues.Bitbucket))
	serviceEndpoint.Url = converter.String("https://api.bitbucket.org/")
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointBitBucket(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	tfhelper.HelpFlattenSecret(d, "password")
	d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
}
