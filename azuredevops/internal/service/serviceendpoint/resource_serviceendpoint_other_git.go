package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointOtherGit schema and implementation for bitbucket service endpoint resource
func ResourceServiceEndpointOtherGit() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointOtherGit, expandServiceEndpointOtherGit)
	makeUnprotectedSchema(r, "username", "AZDO_OTHER_GIT_SERVICE_CONNECTION_USERNAME", "The bitbucket username which should be used.")
	makeUnprotectedSchema(r, "url", "AZDO_OTHER_GIT_SERVICE_CONNECTION_URL", "The HTTPS URL for the repo.")
	makeProtectedSchema(r, "password", "AZDO_OTHER_GIT_SERVICE_CONNECTION_PASSWORD", "The bitbucket password whi|ch should be used.")
	return r
}

func expandServiceEndpointOtherGit(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String(string("git"))
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointOtherGit(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	tfhelper.HelpFlattenSecret(d, "password")
	d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
}
