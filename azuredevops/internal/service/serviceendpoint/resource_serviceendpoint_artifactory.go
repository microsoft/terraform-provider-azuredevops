package serviceendpoint

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointArtifactory schema and implementation for Artifactory service endpoint resource
func ResourceServiceEndpointArtifactory() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointArtifactory, expandServiceEndpointArtifactory)

	r.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: func(i interface{}, key string) (_ []string, errors []error) {
			url, ok := i.(string)
			if !ok {
				errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
				return
			}
			if strings.HasSuffix(url, "/") {
				errors = append(errors, fmt.Errorf("%q should not end with slash, got %q.", key, url))
				return
			}
			return validation.IsURLWithHTTPorHTTPS(url, key)
		},
		Description: "Url for the Artifactory Server",
	}

	r.Schema["token"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		ConflictsWith:    []string{"username", "password"},
		ExactlyOneOf:     []string{"token", "password"},
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "Authentication Token generated through Artifactory",
	}
	r.Schema["username"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ConflictsWith:    []string{"token"},
		RequiredWith:     []string{"password"},
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "Artifactory Username",
	}
	r.Schema["password"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		ConflictsWith:    []string{"token"},
		RequiredWith:     []string{"username"},
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "Artifactory Password",
	}
	// Add spots in the schema to store the token/password hashes
	for _, key := range []string{"token", "username", "password"} {
		secretHashKey, secretHashSchema := tfhelper.GenerateSecreteMemoSchema(key)
		secretHashSchema.Optional = true
		r.Schema[secretHashKey] = secretHashSchema
	}
	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointArtifactory(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("artifactoryService")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	if u, ok := d.GetOk("username"); ok {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Scheme: converter.String("UsernamePassword"),
			Parameters: &map[string]string{
				"username": u.(string),
				"password": d.Get("password").(string),
			},
		}
	} else {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Scheme: converter.String("Token"),
			Parameters: &map[string]string{
				"apitoken": d.Get("token").(string),
			},
		}
	}
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointArtifactory(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	d.Set("url", *serviceEndpoint.Url)
	if *serviceEndpoint.Authorization.Scheme == "UsernamePassword" {
		tfhelper.HelpFlattenSecret(d, "password")
		tfhelper.HelpFlattenSecret(d, "username")
		d.Set("password", (*serviceEndpoint.Authorization.Parameters)["password"])
		d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])

	} else if *serviceEndpoint.Authorization.Scheme == "Token" {
		tfhelper.HelpFlattenSecret(d, "token")
		d.Set("token", (*serviceEndpoint.Authorization.Parameters)["apitoken"])
	} else {
		log.Fatalf("Scheme %q unknown.", *serviceEndpoint.Authorization.Scheme)
	}
}
