package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointArgoCD schema and implementation for ArgoCD service endpoint resource
func ResourceServiceEndpointArgoCD() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointArgoCD, expandServiceEndpointArgoCD)

	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Description:  "Url for the ArgoCD Server",
	}

	r.Schema["token"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "Authentication Token generated through ArgoCD (go to My Account > Security > Generate Tokens)",
	}
	// Add a spot in the schema to store the token secretly
	stSecretHashKey, stSecretHashSchema := tfhelper.GenerateSecreteMemoSchema("token")
	r.Schema[stSecretHashKey] = stSecretHashSchema

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointArgoCD(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
		Parameters: &map[string]string{
			"username": d.Get("token").(string),
		},
	}
	serviceEndpoint.Type = converter.String("argocd")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointArgoCD(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	tfhelper.HelpFlattenSecret(d, "token")

	d.Set("url", *serviceEndpoint.Url)
	d.Set("token", (*serviceEndpoint.Authorization.Parameters)["username"])
}
