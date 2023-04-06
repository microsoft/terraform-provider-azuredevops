package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointGeneric schema and implementation for generic service endpoint resource
func ResourceServiceEndpointGeneric() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointGeneric, expandServiceEndpointGeneric)
	r.Schema["server_url"] = &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Required:     true,
		Description:  "The server URL of the generic service connection.",
	}
	r.Schema["username"] = &schema.Schema{
		Type:        schema.TypeString,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_SERVICE_CONNECTION_USERNAME", nil),
		Description: "The username to use for the generic service connection.",
		Optional:    true,
	}
	r.Schema["password"] = &schema.Schema{
		Type:        schema.TypeString,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GENERIC_SERVICE_CONNECTION_PASSWORD", nil),
		Description: "The password or token key to use for the generic service connection.",
		Sensitive:   true,
		Optional:    true,
	}
	return r
}

func expandServiceEndpointGeneric(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("generic")
	serviceEndpoint.Url = converter.String(d.Get("server_url").(string))
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointGeneric(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("server_url", *serviceEndpoint.Url)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
}
