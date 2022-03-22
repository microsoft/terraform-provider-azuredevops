package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointAzureDevOps() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAzureDevOps, expandServiceEndpointAzureDevOps)
	r.DeprecationMessage = "This resource is duplicate with azuredevops_serviceendpoint_runpipeline,  will be removed in the future, use azuredevops_serviceendpoint_runpipeline instead."
	r.Schema["org_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_DEVOPS_ORG_URL", "https://dev.azure.com/[organization]"),
		Description:  "The Organization Url.",
	}
	r.Schema["release_api_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_DEVOPS_RELEASE_API_URL", "https://vsrm.dev.azure.com/[organization]"),
	}

	makeProtectedSchema(r, "personal_access_token", "AZDO_DEVOPS_PAT", "The Azure DevOps personal access token.")
	return r
}

func expandServiceEndpointAzureDevOps(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("personal_access_token").(string),
		},
		Scheme: converter.String("Token"),
	}
	serviceEndpoint.Type = converter.String("AZDOAPI")
	serviceEndpoint.Url = converter.String(d.Get("org_url").(string))
	serviceEndpoint.Data = &map[string]string{
		"releaseUrl": d.Get("release_api_url").(string),
	}
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointAzureDevOps(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("org_url", serviceEndpoint.Url)
	tfhelper.HelpFlattenSecret(d, "password")
	d.Set("release_api_url", (*serviceEndpoint.Data)["releaseUrl"])
}
