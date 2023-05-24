package serviceendpoint

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointAwsForTerraform schema and implementation for "Terraform for AWS" service endpoint resource
func ResourceServiceEndpointAwsForTerraform() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAwsForTerraform, expandServiceEndpointAwsForTerraform)
	r.Schema["access_key_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_TF_AWS_SERVICE_CONNECTION_ACCESS_KEY_ID", nil),
		Description:  "The AWS access key ID for signing programmatic requests.",
	}
	r.Schema["secret_access_key"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateFunc:     validation.StringIsNotEmpty,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_TF_AWS_SERVICE_CONNECTION_SECRET_ACCESS_KEY", nil),
		Description:      "The AWS secret access key for signing programmatic requests.",
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
	}
	saSecretHashKey, saSecretHashSchema := tfhelper.GenerateSecreteMemoSchema("secret_access_key")
	r.Schema[saSecretHashKey] = saSecretHashSchema
	r.Schema["region"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_TF_AWS_SERVICE_CONNECTION_REGION", nil),
		Description:  "The AWS region to use for programmatic requests.",
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAwsForTerraform(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("access_key_id").(string),
			"password": d.Get("secret_access_key").(string),
			"region":   d.Get("region").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("AWSServiceEndpoint")

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAwsForTerraform(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	tfhelper.HelpFlattenSecret(d, "secret_access_key")

	d.Set("access_key_id", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("secret_access_key", (*serviceEndpoint.Authorization.Parameters)["password"])
	d.Set("region", (*serviceEndpoint.Authorization.Parameters)["region"])
}
