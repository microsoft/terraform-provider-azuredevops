package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointGcp schema and implementation for gcp service endpoint resource
func ResourceServiceEndpointGcp() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointGcp, expandServiceEndpointGcp)
	r.Schema["client_email"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    false,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_CLIENT_EMAIL", nil),
		Description: "The client email field in the JSON key file for creating the JSON Web Token.",
	}
	r.Schema["private_key"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_PRIVATE_KEY", nil),
		Description:      "Private Key for connecting to the endpoint.",
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
	}
	saSecretHashKey, saSecretHashSchema := tfhelper.GenerateSecreteMemoSchema("private_key")
	r.Schema[saSecretHashKey] = saSecretHashSchema
	r.Schema["token_uri"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    false,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_TOKEN_URI", nil),
		Description: "The token uri field in the JSON key file for creating the JSON Web Token.",
		Sensitive:   false,
	}
	r.Schema["scope"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_SCOPE", nil),
		Description: "Scope to be provided",
	}
	r.Schema["gcp_project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    false,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_GCP_PROJECT_ID", nil),
		Description: "Scope to be provided",
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGcp(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"Issuer":     d.Get("client_email").(string),
			"Audience":   d.Get("token_uri").(string),
			"Scope":      d.Get("scope").(string),
			"PrivateKey": d.Get("private_key").(string),
		},
		Scheme: converter.String("JWT"),
	}
	serviceEndpoint.Data = &map[string]string{
		"project": d.Get("gcp_project_id").(string),
	}
	serviceEndpoint.Type = converter.String("GoogleCloudServiceEndpoint")
	serviceEndpoint.Url = converter.String("https://www.googleapis.com/")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointGcp(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	tfhelper.HelpFlattenSecret(d, "private_key")

	d.Set("client_email", (*serviceEndpoint.Authorization.Parameters)["Issuer"])
	d.Set("token_uri", (*serviceEndpoint.Authorization.Parameters)["Audience"])
	d.Set("scope", (*serviceEndpoint.Authorization.Parameters)["Scope"])
	d.Set("gcp_project_id", (*serviceEndpoint.Data)["project"])
}
