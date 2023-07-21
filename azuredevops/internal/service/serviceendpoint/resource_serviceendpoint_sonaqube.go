package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointSonarQube schema and implementation for SonarQube service endpoint resource
func ResourceServiceEndpointSonarQube() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointSonarQube, expandServiceEndpointSonarQube)

	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Description:  "Url for the SonarQube Server",
	}

	r.Schema["token"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Sensitive:    true,
		ValidateFunc: validation.StringIsNotWhiteSpace,
		Description:  "Authentication Token generated through SonarQube (go to My Account > Security > Generate Tokens)",
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointSonarQube(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
		Parameters: &map[string]string{
			"username": d.Get("token").(string),
		},
	}
	serviceEndpoint.Type = converter.String("sonarqube")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointSonarQube(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	d.Set("url", *serviceEndpoint.Url)
}
