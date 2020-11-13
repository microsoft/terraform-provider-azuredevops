package serviceendpoint

// rapid7 endpoint
import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointRapid7 schema and implementation for rapid7 service endpoint resource
func ResourceServiceEndpointRapid7() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointRapid7, expandServiceEndpointRapid7)
	r.Schema["region"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_RAPID7_SERVICE_CONNECTION_REGION", nil),
		Description:  "The Rapid7 Server region which should be used.",
		Sensitive:    false,
		ValidateFunc: validation.StringIsNotWhiteSpace,
	}
	r.Schema["auth_token"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_RAPID7_SERVICE_CONNECTION_TOKEN", nil),
		Description:  "The Rapid7 Server auth Token which should be used.",
		Sensitive:    true,
		ValidateFunc: validation.StringIsNotWhiteSpace,
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema("auth_token")
	r.Schema[patHashKey] = patHashSchema
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointRapid7(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("auth_token").(string),
		},
		Scheme: converter.String("Token"),
	}
	serviceEndpoint.Data = &map[string]string{
		"region": d.Get("region").(string),
	}
	serviceEndpoint.Type = converter.String("ias")

	serviceEndpoint.Url = converter.String("https://rapid7.com")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointRapid7(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	tfhelper.HelpFlattenSecret(d, "auth_token")
	d.Set("auth_token", (*serviceEndpoint.Authorization.Parameters)["apitoken"])
	d.Set("region", (*serviceEndpoint.Data)["regionValues"])
}
