package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointCustom schema and implementation for Custom service endpoint resource
func ResourceServiceEndpointCustom() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointCustom, expandServiceEndpointCustom)
	r.Schema["service_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Service Type of the Custom service connection.",
	}
	r.Schema["server_url"] = &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Required:     true,
		Description:  "The server URL of the Custom service connection.",
	}
	r.Schema["username"] = &schema.Schema{
		Type:        schema.TypeString,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_CUSTOM_SERVICE_CONNECTION_USERNAME", nil),
		Description: "The username to use for the Custom service connection.",
		Optional:    true,
	}
	r.Schema["password"] = &schema.Schema{
		Type:             schema.TypeString,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_CUSTOM_SERVICE_CONNECTION_PASSWORD", nil),
		Description:      "The password or token key to use for the Custom service connection.",
		Sensitive:        true,
		Optional:         true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
	}
	r.Schema["data"] = &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Optional payload required for the creation of the endpoint",
	}
	secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema("password")
	r.Schema[secretHashKey] = secretHashSchema
	return r
}

func expandServiceEndpointCustom(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String(d.Get("service_type").(string))
	serviceEndpoint.Url = converter.String(d.Get("server_url").(string))
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Data = d.Get("data").(*map[string]string)
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointCustom(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("service_type", *serviceEndpoint.Type)
	d.Set("server_url", *serviceEndpoint.Url)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("data", *serviceEndpoint.Data)
	tfhelper.HelpFlattenSecret(d, "password")
}
