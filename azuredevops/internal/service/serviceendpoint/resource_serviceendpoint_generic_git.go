package serviceendpoint

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointGenericGit schema and implementation for generic git service endpoint resource
func ResourceServiceEndpointGenericGit() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointGenericGit, expandServiceEndpointGenericGit)
	r.Schema["repository_url"] = &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Required:     true,
		Description:  "The server URL of the generic git service connection.",
	}
	r.Schema["username"] = &schema.Schema{
		Type:        schema.TypeString,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_GIT_SERVICE_CONNECTION_USERNAME", nil),
		Description: "The username to use for the generic service git connection.",
		Optional:    true,
	}
	r.Schema["password"] = &schema.Schema{
		Type:             schema.TypeString,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_GENERIC_GIT_SERVICE_CONNECTION_PASSWORD", nil),
		Description:      "The password or token key to use for the generic git service connection.",
		Sensitive:        true,
		Optional:         true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
	}
	r.Schema["enable_pipelines_access"] = &schema.Schema{
		Type:        schema.TypeBool,
		Default:     true,
		Description: "A value indicating whether or not to attempt accessing this git server from Azure Pipelines.",
		Optional:    true,
	}
	secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema("password")
	r.Schema[secretHashKey] = secretHashSchema
	return r
}

func expandServiceEndpointGenericGit(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("git")
	serviceEndpoint.Url = converter.String(d.Get("repository_url").(string))
	serviceEndpoint.Data = &map[string]string{
		"accessExternalGitServer": strconv.FormatBool(d.Get("enable_pipelines_access").(bool)),
	}
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointGenericGit(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("repository_url", *serviceEndpoint.Url)
	if v, err := strconv.ParseBool((*serviceEndpoint.Data)["accessExternalGitServer"]); err != nil {
		d.Set("enable_pipelines_access", v)
	}
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	tfhelper.HelpFlattenSecret(d, "password")
	d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
}
