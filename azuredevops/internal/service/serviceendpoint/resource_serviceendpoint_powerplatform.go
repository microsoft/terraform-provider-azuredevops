package serviceendpoint

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointPowerPlatform schema and implementation for PowerPlatform service endpoint resource
func ResourceServiceEndpointPowerPlatform() *schema.Resource {
	r := &schema.Resource{
		CreateContext: resourceServiceEndpointPowerPlatformCreate,
		ReadContext:   resourceServiceEndpointPowerPlatformRead,
		UpdateContext: resourceServiceEndpointPowerPlatformUpdate,
		DeleteContext: resourceServiceEndpointPowerPlatformDelete,
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
			Description:  "The Server URL for the Power Platform connection (e.g. https://org.crm.dynamics.com or generic)",
			ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
		},
		"credentials": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"serviceprincipalid": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The Application (Client) ID of the Service Principal.",
						ValidateFunc: validation.IsUUID,
					},
					"serviceprincipalkey": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						Description:  "The Client Secret of the Service Principal.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tenant_id": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The Tenant ID.",
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},
		"features": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"validate": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Whether or not to validate connection with Azure after create or update operations.",
					},
				},
			},
		},
	})

	return r
}

func resourceServiceEndpointPowerPlatformCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointPowerPlatform(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	resp, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return diag.Errorf("creating service endpoint in Azure DevOps: %+v", err)
	}

	d.SetId(resp.Id.String())

	if v, ok := d.GetOk("features"); ok {
		features := v.([]interface{})[0].(map[string]interface{})
		if features["validate"].(bool) {
			projectID := d.Get("project_id").(string)
			err = validateServiceEndpoint(clients, resp, projectID, 60*time.Second)
			if err != nil {
				if delErr := deleteServiceEndpoint(clients, resp, d.Timeout(schema.TimeoutDelete)); delErr != nil {
					return diag.Errorf("Error validating service endpoint and failed to delete it: %+v", err)
				}
				return diag.Errorf("Error validating service endpoint: %+v", err)
			}
		}
	}

	return resourceServiceEndpointPowerPlatformRead(clients.Ctx, d, m)
}

func resourceServiceEndpointPowerPlatformRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return diag.Errorf("reading service endpoint in Azure DevOps: %+v", err)
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if isServiceEndpointDeleted(d, err, serviceEndpoint, getArgs) {
		return nil
	}
	if err != nil {
		return diag.Errorf("looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	flattenServiceEndpointPowerPlatform(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointPowerPlatformUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointPowerPlatform(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	resp, err := updateServiceEndpoint(clients, serviceEndpoint)
	if err != nil {
		return diag.Errorf("updating service endpoint in Azure DevOps: %+v", err)
	}

	if v, ok := d.GetOk("features"); ok {
		features := v.([]interface{})[0].(map[string]interface{})
		if features["validate"].(bool) {
			projectID := d.Get("project_id").(string)
			err = validateServiceEndpoint(clients, resp, projectID, 60*time.Second)
			if err != nil {
				return diag.Errorf("Error validating service endpoint: %+v", err)
			}
		}
	}

	return resourceServiceEndpointPowerPlatformRead(clients.Ctx, d, m)
}

func resourceServiceEndpointPowerPlatformDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointPowerPlatform(d)
	if err != nil {
		return diag.Errorf(errMsgTfConfigRead, err)
	}

	err = deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf(" Deleting service endpoint in Azure DevOps: %+v", err)
	}
	return nil
}

func expandServiceEndpointPowerPlatform(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)

	serviceEndpoint.Type = converter.String("powerplatform-spn")

	if v, ok := d.GetOk("url"); ok {
		serviceEndpoint.Url = converter.String(v.(string))
	} else {
		return nil, fmt.Errorf("url is required for PowerPlatform service endpoint")
	}

	var credentials map[string]any
	if v, ok := d.GetOk("credentials"); ok && len(v.([]any)) > 0 {
		credentials = v.([]any)[0].(map[string]any)
	} else {
		return nil, fmt.Errorf("credentials block is required for PowerPlatform service endpoint")
	}

	parameters := map[string]string{
		"tenantId":      credentials["tenant_id"].(string),
		"applicationId": credentials["serviceprincipalid"].(string),
		"clientSecret":  credentials["serviceprincipalkey"].(string),
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme:     converter.String("None"),
		Parameters: &parameters,
	}

	serviceEndpoint.Data = &map[string]string{}

	return serviceEndpoint, nil
}

func flattenServiceEndpointPowerPlatform(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)

	credentials := make(map[string]any)

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		params := *serviceEndpoint.Authorization.Parameters

		if v, ok := params["tenantId"]; ok {
			credentials["tenant_id"] = v
		}

		if v, ok := params["applicationId"]; ok {
			credentials["serviceprincipalid"] = v
		}

		if oldCredentials, ok := d.Get("credentials").([]any); ok && len(oldCredentials) > 0 {
			oldMap := oldCredentials[0].(map[string]any)
			credentials["serviceprincipalkey"] = oldMap["serviceprincipalkey"]
		}
	}

	if serviceEndpoint.Url != nil {
		d.Set("url", *serviceEndpoint.Url)
	}

	d.Set("credentials", []any{credentials})
}
