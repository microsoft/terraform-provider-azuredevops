package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointAws schema and implementation for aws service endpoint resource
func ResourceServiceEndpointAws() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAws, expandServiceEndpointAws)
	r.Schema["access_key_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_ACCESS_KEY_ID", nil),
		Description: "The AWS access key ID for signing programmatic requests.",
	}
	r.Schema["secret_access_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_SECRET_ACCESS_KEY", nil),
		Description: "The AWS secret access key for signing programmatic requests.",
		Sensitive:   true,
	}
	r.Schema["session_token"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_SESSION_TOKEN", nil),
		Description: "The AWS session token for signing programmatic requests.",
		Sensitive:   true,
	}
	r.Schema["role_to_assume"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_RTA", nil),
		Description: "The Amazon Resource Name (ARN) of the role to assume.",
	}
	r.Schema["role_session_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_RSN", nil),
		Description: "Optional identifier for the assumed role session.",
	}
	r.Schema["external_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_EXTERNAL_ID", nil),
		Description: "A unique identifier that is used by third parties when assuming roles in their customers' accounts, aka cross-account role access.",
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAws(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username":        d.Get("access_key_id").(string),
			"password":        d.Get("secret_access_key").(string),
			"sessionToken":    d.Get("session_token").(string),
			"assumeRoleArn":   d.Get("role_to_assume").(string),
			"roleSessionName": d.Get("role_session_name").(string),
			"externalId":      d.Get("external_id").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("aws")
	serviceEndpoint.Url = converter.String("https://aws.amazon.com/")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAws(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	d.Set("access_key_id", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("role_to_assume", (*serviceEndpoint.Authorization.Parameters)["assumeRoleArn"])
	d.Set("role_session_name", (*serviceEndpoint.Authorization.Parameters)["roleSessionName"])
	d.Set("external_id", (*serviceEndpoint.Authorization.Parameters)["externalId"])
}
