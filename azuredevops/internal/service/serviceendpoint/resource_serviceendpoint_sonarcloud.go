package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointSonarCloud schema and implementation for SonarCloud service endpoint resource
func ResourceServiceEndpointSonarCloud() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointSonarCloud, expandServiceEndpointSonarCloud)

	r.Schema["token"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Sensitive:    true,
		ValidateFunc: validation.StringIsNotWhiteSpace,
		Description:  "Authentication Token generated through SonarCloud (go to My Account > Security > Generate Tokens)",
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointSonarCloud(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("Token"),
		Parameters: &map[string]string{
			"apitoken": d.Get("token").(string),
		},
	}
	serviceEndpoint.Type = converter.String("sonarcloud")
	serviceEndpoint.Url = converter.String("https://sonarcloud.io")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointSonarCloud(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	d.Get("")
}
