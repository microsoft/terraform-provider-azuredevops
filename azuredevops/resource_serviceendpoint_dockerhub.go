package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	crud "github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
)

func resourceServiceEndpointDockerHub() *schema.Resource {
	r := crud.GenBaseServiceEndpointResource(flattenServiceEndpointDockerHub, expandServiceEndpointDockerHub, parseImportedProjectIDAndServiceEndpointID)
	crud.MakeUnprotectedSchema(r, "docker_username", "AZDO_DOCKERHUB_SERVICE_CONNECTION_USERNAME", "The DockerHub username which should be used.")
	crud.MakeUnprotectedSchema(r, "docker_email", "AZDO_DOCKERHUB_SERVICE_CONNECTION_EMAIL", "The DockerHub email address which should be used.")
	crud.MakeProtectedSchema(r, "docker_password", "AZDO_DOCKERHUB_SERVICE_CONNECTION_PASSWORD", "The DockerHub password which should be used.")
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointDockerHub(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := crud.DoBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("docker_username").(string),
			"password": d.Get("docker_password").(string),
			"email":    d.Get("docker_email").(string),
			"registry": "https://index.docker.io/v1/",
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Data = &map[string]string{
		"registrytype": "DockerHub",
	}
	serviceEndpoint.Type = converter.String("dockerregistry")
	serviceEndpoint.Url = converter.String("https://hub.docker.com/")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointDockerHub(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	crud.DoBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("docker_email", (*serviceEndpoint.Authorization.Parameters)["email"])
	d.Set("docker_username", (*serviceEndpoint.Authorization.Parameters)["username"])
	tfhelper.HelpFlattenSecret(d, "docker_password")
	d.Set("docker_password", (*serviceEndpoint.Authorization.Parameters)["password"])
}

// parseImportedProjectIDAndServiceEndpointID : Parse the Id (projectId/serviceEndpointId) or (projectName/serviceEndpointId)
func parseImportedProjectIDAndServiceEndpointID(clients *config.AggregatedClient, id string) (string, string, error) {
	project, resourceID, err := tfhelper.ParseImportedUUID(id)
	if err != nil {
		return "", "", err
	}

	// Get the project ID
	currentProject, err := ProjectRead(clients, project, project)
	if err != nil {
		return "", "", err
	}

	return currentProject.Id.String(), resourceID, nil
}
