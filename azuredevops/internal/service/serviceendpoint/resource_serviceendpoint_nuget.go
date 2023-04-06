package serviceendpoint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointNuget schema and implementation for Nuget service endpoint resource
func ResourceServiceEndpointNuget() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointNuget, expandServiceEndpointNuget)

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
		Description: "Url for the Nuget Feed",
	}

	at := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"token": {
				Description: "The Nuget Feed access token.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
		},
	}

	ak := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Description: "The Nuget Feed API key.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
		},
	}

	aup := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Description: "The Nuget feed user name.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"password": {
				Description: "The Nuget feed password.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
		},
	}

	r.Schema["authentication_token"] = &schema.Schema{
		Type:         schema.TypeList,
		Optional:     true,
		MinItems:     1,
		MaxItems:     1,
		Elem:         at,
		ExactlyOneOf: []string{"authentication_basic", "authentication_token", "authentication_none"},
	}

	r.Schema["authentication_none"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem:     ak,
	}

	r.Schema["authentication_basic"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem:     aup,
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointNuget(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("externalnugetfeed")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	authScheme := "Token"

	authParams := make(map[string]string)

	if x, ok := d.GetOk("authentication_token"); ok {
		authScheme = "Token"
		msi := x.([]interface{})[0].(map[string]interface{})
		authParams["apitoken"], ok = msi["token"].(string)
		if !ok {
			return nil, nil, errors.New("Unable to read 'token'")
		}
	} else if x, ok := d.GetOk("authentication_none"); ok {
		authScheme = "None"
		msi := x.([]interface{})[0].(map[string]interface{})
		authParams["nugetkey"], ok = msi["key"].(string)
		if !ok {
			return nil, nil, errors.New("Unable to read 'key'")
		}
	} else if x, ok := d.GetOk("authentication_basic"); ok {
		authScheme = "UsernamePassword"
		msi := x.([]interface{})[0].(map[string]interface{})
		authParams["username"], ok = msi["username"].(string)
		if !ok {
			return nil, nil, errors.New("Unable to read 'username'")
		}
		authParams["password"], ok = msi["password"].(string)
		if !ok {
			return nil, nil, errors.New("Unable to read 'password'")
		}
	}
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &authParams,
		Scheme:     &authScheme,
	}

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointNuget(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "UsernamePassword") {
		if _, ok := d.GetOk("authentication_basic"); !ok {
			auth := make(map[string]interface{})
			auth["username"] = (*serviceEndpoint.Authorization.Parameters)["username"]
			auth["password"] = ""
			d.Set("authentication_basic", []interface{}{auth})
		}
	} else if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
		if _, ok := d.GetOk("authentication_token"); !ok {
			auth := make(map[string]interface{})
			auth["token"] = ""
			d.Set("authentication_token", []interface{}{auth})
		}
	} else if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "None") {
		if _, ok := d.GetOk("authentication_none"); !ok {
			auth := make(map[string]interface{})
			auth["key"] = ""
			d.Set("authentication_none", []interface{}{auth})
		}
	} else {
		panic(fmt.Errorf("inconsistent authorization scheme. Expected: (Token, None, UsernamePassword)  , but got %s", *serviceEndpoint.Authorization.Scheme))
	}

	d.Set("url", *serviceEndpoint.Url)
}
