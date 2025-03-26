package serviceendpoint

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointOpenshift() *schema.Resource {
	r := &schema.Resource{
		CreateContext: resourceServiceEndpointOpenshiftCreate,
		ReadContext:   resourceServiceEndpointOpenshiftRead,
		UpdateContext: resourceServiceEndpointOpenshiftUpdate,
		DeleteContext: resourceServiceEndpointOpenshiftDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}
	maps.Copy(r.Schema, map[string]*schema.Schema{
		"server_url": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"accept_untrusted_certs": {
			Type:     schema.TypeBool,
			Optional: true,
		},

		"certificate_authority_file": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"auth_basic": {
			Type:          schema.TypeList,
			Optional:      true,
			MinItems:      1,
			MaxItems:      1,
			ConflictsWith: []string{"auth_token", "auth_none"},
			AtLeastOneOf:  []string{"auth_basic", "auth_token", "auth_none"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"username": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
					"password": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
				},
			},
		},

		"auth_token": {
			Type:          schema.TypeList,
			Optional:      true,
			MinItems:      1,
			MaxItems:      1,
			ConflictsWith: []string{"auth_basic", "auth_none"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"token": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
				},
			},
		},

		"auth_none": {
			Type:          schema.TypeList,
			Optional:      true,
			MinItems:      1,
			MaxItems:      1,
			ConflictsWith: []string{"auth_basic", "auth_token"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kube_config": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	})
	return r
}

func resourceServiceEndpointOpenshiftCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	args, err := expandServiceEndpointOpenshift(d)
	if err != nil {
		return diag.Errorf(" Expanding service connection: %+v", err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, args)
	if err != nil {
		return diag.Errorf("Creating service connection: %+v", err)
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointOpenshiftRead(ctx, d, m)
}

func resourceServiceEndpointOpenshiftRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.Errorf("Getting service endpoint args: %+v", err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if isServiceEndpointDeleted(d, err, serviceEndpoint, getArgs) {
		return nil
	}
	if err != nil {
		return diag.Errorf(" looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return diag.Errorf(" Checking service connection permissions: %v", err)
	}
	if err := flattenServiceEndpointOpenshift(d, serviceEndpoint); err != nil {
		return diag.Errorf(" Flattening service endpoint configuration: %+v", err)
	}
	return nil
}

func resourceServiceEndpointOpenshiftUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointOpenshift(d)
	if err != nil {
		return diag.Errorf(" Expanding service connection: %+v", err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return diag.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointOpenshiftRead(ctx, d, m)
}

func resourceServiceEndpointOpenshiftDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	args, err := expandServiceEndpointOpenshift(d)
	if err != nil {
		return diag.Errorf(" Expanding service connection: %+v", err)
	}

	err = deleteServiceEndpoint(clients, args, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf(" Deleting service endpoint in Azure DevOps: %+v", err)
	}
	return nil
}

func expandServiceEndpointOpenshift(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	authType := ""
	params := make(map[string]string)

	if config, ok := d.GetOk("auth_basic"); ok {
		authType = "UsernamePassword"
		val := config.([]interface{})[0].(map[string]interface{})
		params = map[string]string{
			"username":             val["username"].(string),
			"password":             val["password"].(string),
			"acceptUntrustedCerts": "false",
		}

		if !d.GetRawConfig().IsNull() {
			raw := d.GetRawConfig().AsValueMap()
			acceptUntrustedCerts := raw["accept_untrusted_certs"]
			if !acceptUntrustedCerts.IsNull() && !acceptUntrustedCerts.False() {
				params["acceptUntrustedCerts"] = "true"
			}
		}

		if v, ok := d.GetOk("certificate_authority_file"); ok || v.(string) != "" {
			params["certificateAuthorityFile"] = v.(string)
		}
	}

	if config, ok := d.GetOk("auth_token"); ok {
		authType = "Token"
		val := config.([]interface{})[0].(map[string]interface{})

		params = map[string]string{
			"apitoken":             val["token"].(string),
			"acceptUntrustedCerts": "false",
		}

		if !d.GetRawConfig().IsNull() {
			raw := d.GetRawConfig().AsValueMap()
			acceptUntrustedCerts := raw["accept_untrusted_certs"]
			if !acceptUntrustedCerts.IsNull() && !acceptUntrustedCerts.False() {
				params["acceptUntrustedCerts"] = "true"
			}
		}

		if v, ok := d.GetOk("certificate_authority_file"); ok || v.(string) != "" {
			params["certificateAuthorityFile"] = v.(string)
		}
	}
	if config, ok := d.GetOk("auth_none"); ok {
		authType = "None"
		val := config.([]interface{})[0].(map[string]interface{})

		params = map[string]string{}
		if v, ok := val["kube_config"]; ok {
			params["kubeConfig"] = v.(string)
		}
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &params,
		Scheme:     &authType,
	}

	serviceEndpoint.Type = converter.String("openshift")
	serviceEndpoint.Url = converter.String(d.Get("server_url").(string))
	return serviceEndpoint, nil
}

func flattenServiceEndpointOpenshift(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) error {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("server_url", serviceEndpoint.Url)
	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		params := serviceEndpoint.Authorization.Parameters
		if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "UsernamePassword") {
			acceptUntrustedCerts, err := strconv.ParseBool((*params)["acceptUntrustedCerts"])
			if err != nil {
				return fmt.Errorf(" Parse `acceptUntrustedCerts`: %v", err)
			}
			d.Set("accept_untrusted_certs", acceptUntrustedCerts)
			d.Set("certificate_authority_file", (*params)["certificateAuthorityFile"])
		}
		if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
			acceptUntrustedCerts, err := strconv.ParseBool((*params)["acceptUntrustedCerts"])
			if err != nil {
				return fmt.Errorf(" Parse `acceptUntrustedCerts`: %v", err)
			}
			d.Set("accept_untrusted_certs", acceptUntrustedCerts)
			d.Set("certificate_authority_file", (*params)["certificateAuthorityFile"])
		}
		if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "None") {
			// params not return, ignore
		}
	}
	return nil
}
