package serviceendpoint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
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

	r.Schema["username"] = &schema.Schema{
		Description: "The Jenkins user name.",
		Type:        schema.TypeString,
		Required:    true,
	}
	r.Schema["password"] = &schema.Schema{
		Description: "The Jenkins password.",
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointJenkins(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("jenkins")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}

	data := map[string]string{}
	data["AcceptUntrustedCerts"] = strconv.FormatBool(d.Get("accept_untrusted_certs").(bool))

	serviceEndpoint.Data = &data

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointJenkins(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("url", *serviceEndpoint.Url)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	unsecured, err := strconv.ParseBool((*serviceEndpoint.Data)["AcceptUntrustedCerts"])
	if err != nil {
		fmt.Println(err)
		return
	}
	d.Set("accept_untrusted_certs", unsecured)
}
