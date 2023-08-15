package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointDockerRegistry schema and implementation for docker registry service endpoint resource
func ResourceServiceEndpointDockerRegistry() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointDockerRegistry, expandServiceEndpointDockerRegistry)
	r.Schema["docker_registry"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_REGISTRY", "https://index.docker.io/v1/"),
		Description: "The DockerRegistry registry which should be used.",
	}
	r.Schema["docker_username"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_USERNAME", nil),
		Description: "The DockerRegistry username which should be used.",
	}
	r.Schema["docker_password"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_PASSWORD", nil),
		Description: "The DockerRegistry password which should be used.",
		Sensitive:   true,
	}
	r.Schema["docker_email"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_EMAIL", nil),
		Description: "The DockerRegistry email address which should be used.",
	}
	r.Schema["registry_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_REGISTRY_TYPE", "DockerHub"),
		ValidateFunc: validation.StringInSlice([]string{
			string("DockerHub"),
			string("Others"),
		}, false),
		ForceNew: true,
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointDockerRegistry(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"registry": d.Get("docker_registry").(string),
			"username": d.Get("docker_username").(string),
			"password": d.Get("docker_password").(string),
			"email":    d.Get("docker_email").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Data = &map[string]string{
		"registrytype": d.Get("registry_type").(string),
	}
	serviceEndpoint.Type = converter.String("dockerregistry")
	serviceEndpoint.Url = converter.String("https://hub.docker.com/") // DevOps UI sets hub.docker.com for both DockerHub and Others types
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointDockerRegistry(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("docker_registry", (*serviceEndpoint.Authorization.Parameters)["registry"])
	d.Set("docker_email", (*serviceEndpoint.Authorization.Parameters)["email"])
	d.Set("docker_username", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("registry_type", (*serviceEndpoint.Data)["registrytype"])
}
