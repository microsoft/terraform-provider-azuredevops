package serviceendpoint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointJenkins schema and implementation for Jenkins service endpoint resource
func ResourceServiceEndpointJenkins() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointJenkins, expandServiceEndpointJenkins)

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
		Description: "Url for the Jenkins Repository",
	}

	r.Schema["accept_untrusted_certs"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Allows the Jenkins clients to accept self-signed SSL server certificates without installing them into the TFS service role and/or Build Agent computers.",
	}

	patHashKeyU, patHashSchemaU := tfhelper.GenerateSecreteMemoSchema("username")
	patHashKeyP, patHashSchemaP := tfhelper.GenerateSecreteMemoSchema("password")
	aup := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Description:      "The Jenkins user name.",
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
			},
			patHashKeyU: patHashSchemaU,
			"password": {
				Description:      "The Jenkins password.",
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
func expandServiceEndpointJenkins(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("jenkins")
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

	data := map[string]string{}
	data["AcceptUntrustedCerts"] = strconv.FormatBool(d.Get("accept_untrusted_certs").(bool))

	serviceEndpoint.Data = &data

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointJenkins(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
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
	unsecured, err := strconv.ParseBool((*serviceEndpoint.Data)["AcceptUntrustedCerts"])
	if err != nil {
		return
	}
	d.Set("accept_untrusted_certs", unsecured)
}
