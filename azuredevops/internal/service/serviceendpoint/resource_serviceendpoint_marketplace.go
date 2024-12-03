package serviceendpoint

import (
	"context"
	"errors"
	"maps"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointMarketplace() *schema.Resource {
	r := &schema.Resource{
		CreateContext: resourceServiceEndpointMarketplaceCreate,
		ReadContext:   resourceServiceEndpointMarketplaceRead,
		UpdateContext: resourceServiceEndpointMarketplaceUpdate,
		DeleteContext: resourceServiceEndpointMarketplaceDelete,
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
		"url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"authentication_token": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"token": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
			ExactlyOneOf: []string{"authentication_basic", "authentication_token"},
		},

		"authentication_basic": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"username": {
						Type:     schema.TypeString,
						Required: true,
					},
					"password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},
	})
	return r
}

func resourceServiceEndpointMarketplaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointMarketplace(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointMarketplaceRead(ctx, d, m)
}

func resourceServiceEndpointMarketplaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return diag.FromErr(err)
	}
	flattenServiceEndpointMarketplace(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointMarketplaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointMarketplace(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return diag.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointMarketplaceRead(ctx, d, m)
}

func resourceServiceEndpointMarketplaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointMarketplace(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	if err = deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func expandServiceEndpointMarketplace(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("TFSMarketplacePublishing")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	authScheme := "Token"

	authParams := make(map[string]string)
	if tokenAuth, ok := d.GetOk("authentication_token"); ok {
		authScheme = "Token"
		token := tokenAuth.([]interface{})[0].(map[string]interface{})
		authParams["apitoken"], ok = token["token"].(string)
		if !ok {
			return nil, errors.New(" unable to get `apitoken`.")
		}
	} else if basicAuth, ok := d.GetOk("authentication_basic"); ok {
		authScheme = "UsernamePassword"
		unamePwd := basicAuth.([]interface{})[0].(map[string]interface{})
		if v, exist := unamePwd["username"].(string); exist {
			authParams["username"] = v
		}

		if v, exist := unamePwd["password"].(string); exist {
			authParams["password"] = v
		}
	}
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &authParams,
		Scheme:     &authScheme,
	}
	return serviceEndpoint, nil
}

func flattenServiceEndpointMarketplace(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Scheme != nil {
		if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "UsernamePassword") {
			if _, ok := d.GetOk("authentication_basic"); !ok {
				auth := make(map[string]interface{})
				if serviceEndpoint.Authorization.Parameters != nil {
					if v, exist := (*serviceEndpoint.Authorization.Parameters)["username"]; exist {
						auth["username"] = v
					}
				}
				d.Set("authentication_basic", []interface{}{auth})
			}
		}
		// ignore scheme=`Token` since the service does not return sensitive data
	}
	d.Set("url", *serviceEndpoint.Url)
}
