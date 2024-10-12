package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointDockerRegistry schema and implementation for docker registry service endpoint resource
func ResourceServiceEndpointDockerRegistry() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointDockerRegistryCreate,
		Read:   resourceServiceEndpointDockerRegistryRead,
		Update: resourceServiceEndpointDockerRegistryUpdate,
		Delete: resourceServiceEndpointDockerRegistryDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	maps.Copy(r.Schema, map[string]*schema.Schema{
		"docker_registry": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_REGISTRY", "https://index.docker.io/v1/"),
			Description: "The DockerRegistry registry which should be used.",
		},
		"docker_username": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_USERNAME", nil),
			Description: "The DockerRegistry username which should be used.",
		},
		"docker_password": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_PASSWORD", nil),
			Description: "The DockerRegistry password which should be used.",
			Sensitive:   true,
		},
		"docker_email": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_EMAIL", nil),
			Description: "The DockerRegistry email address which should be used.",
		},
		"registry_type": {
			Type:         schema.TypeString,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_REGISTRY_TYPE", "DockerHub"),
			ValidateFunc: validation.StringInSlice([]string{"DockerHub", "Others"}, false),
			ForceNew:     true,
		},
	})
	return r
}

func resourceServiceEndpointDockerRegistryCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointDockerRegistry(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointDockerRegistryRead(d, m)
}

func resourceServiceEndpointDockerRegistryRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	flattenServiceEndpointDockerRegistry(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	return nil
}

func resourceServiceEndpointDockerRegistryUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointDockerRegistry(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointDockerRegistry(d, updatedServiceEndpoint, projectID.String())
	return resourceServiceEndpointDockerRegistryRead(d, m)
}

func resourceServiceEndpointDockerRegistryDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointDockerRegistry(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
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
func flattenServiceEndpointDockerRegistry(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	if serviceEndpoint.Authorization != nil {
		if serviceEndpoint.Authorization.Parameters != nil {
			if v, ok := (*serviceEndpoint.Authorization.Parameters)["registry"]; ok {
				d.Set("docker_registry", v)
			}
			if v, ok := (*serviceEndpoint.Authorization.Parameters)["email"]; ok {
				d.Set("docker_email", v)
			}
			if v, ok := (*serviceEndpoint.Authorization.Parameters)["username"]; ok {
				d.Set("docker_username", v)
			}
		}
	}
	if serviceEndpoint.Data != nil {
		if v, ok := (*serviceEndpoint.Data)["registrytype"]; ok {
			d.Set("registry_type", v)
		}
	}
}
