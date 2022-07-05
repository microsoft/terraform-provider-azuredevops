package serviceendpoint

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointNexus schema and implementation for Nexus service endpoint resource
func ResourceServiceEndpointNexus() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointNexus, expandServiceEndpointNexus)

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
		Description: "Url for the Nexus Repository",
	}

	patHashKeyU, patHashSchemaU := tfhelper.GenerateSecreteMemoSchema("username")
	patHashKeyP, patHashSchemaP := tfhelper.GenerateSecreteMemoSchema("password")
	aup := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Description:      "The Nexus user name.",
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
			},
			patHashKeyU: patHashSchemaU,
			"password": {
				Description:      "The Nexus password.",
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
			},
			patHashKeyP: patHashSchemaP,
		},
	}

	r.Schema["authentication_basic"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem:     aup,
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointNexus(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("NexusIqServiceConnection")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))

	authScheme := "UsernamePassword"
	authParams := make(map[string]string)

	if x, ok := d.GetOk("authentication_basic"); ok {
		authScheme = "UsernamePassword"
		msi := x.([]interface{})[0].(map[string]interface{})
		authParams["username"] = expandSecret(msi, "username")
		authParams["password"] = expandSecret(msi, "password")
	}
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &authParams,
		Scheme:     &authScheme,
	}

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointNexus(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "UsernamePassword") {
		auth := make(map[string]interface{})
		if old, ok := d.GetOk("authentication_basic"); ok {
			oldAuthList := old.([]interface{})[0].(map[string]interface{})
			if len(oldAuthList) > 0 {
				newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "authentication_basic", oldAuthList, "password")
				auth[hashKey] = newHash
				newHash, hashKey = tfhelper.HelpFlattenSecretNested(d, "authentication_basic", oldAuthList, "username")
				auth[hashKey] = newHash
			}
		}
		if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
			auth["password"] = (*serviceEndpoint.Authorization.Parameters)["password"]
			auth["username"] = (*serviceEndpoint.Authorization.Parameters)["username"]
		}
		d.Set("authentication_basic", []interface{}{auth})
	} else {
		panic(fmt.Errorf("inconsistent authorization scheme %s", *serviceEndpoint.Authorization.Scheme))
	}

	d.Set("url", *serviceEndpoint.Url)
}
