package serviceendpoint

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataServiceEndpointType schema and implementation for service endpoint type data source
func DataServiceEndpointType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataServiceEndpointTypeRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the service endpoint type",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"authorization_scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The authorization scheme to retrieve parameters for",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the service endpoint type",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the service endpoint type",
			},
			"ui_contribution_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UI contribution ID for this service endpoint type",
			},
			"authentication_schemes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Available authentication schemes for this service endpoint type",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of parameter descriptors for the service endpoint",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter name (descriptor ID)",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The default value for this parameter, if provided by the API",
						},
						"possible_values": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "List of possible values for this parameter, if provided by the API",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"authorization_parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of authorization parameter descriptors (only set if authorization_scheme is provided)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The authorization parameter name (descriptor ID)",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The default value for this parameter, if provided by the API",
						},
						"possible_values": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "List of possible values for this parameter, if provided by the API",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataServiceEndpointTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	serviceEndpointTypesOnceCache.Do(func() {
		args := serviceendpoint.GetServiceEndpointTypesArgs{}
		serviceEndpointTypesCache, serviceEndpointTypesCacheErr = clients.ServiceEndpointClient.GetServiceEndpointTypes(clients.Ctx, args)
	})

	if serviceEndpointTypesCacheErr != nil {
		return diag.FromErr(fmt.Errorf("querying service endpoint types: %v", serviceEndpointTypesCacheErr))
	}

	if serviceEndpointTypesCache == nil || len(*serviceEndpointTypesCache) == 0 {
		return diag.FromErr(fmt.Errorf("no service endpoint types found"))
	}

	name := d.Get("name").(string)
	authScheme := d.Get("authorization_scheme").(string)

	var foundType *serviceendpoint.ServiceEndpointType
	for _, t := range *serviceEndpointTypesCache {
		if t.Name != nil && strings.EqualFold(*t.Name, name) {
			foundType = &t
			break
		}
	}

	if foundType == nil {
		return diag.FromErr(fmt.Errorf("service endpoint type not found with name %s", name))
	}

	// Set basic attributes
	if foundType.Name != nil {
		d.SetId(*foundType.Name)
		if err := d.Set("name", *foundType.Name); err != nil {
			return diag.FromErr(err)
		}
	} else {
		d.SetId(name)
	}

	if foundType.DisplayName != nil {
		if err := d.Set("display_name", *foundType.DisplayName); err != nil {
			return diag.FromErr(err)
		}
	}
	if foundType.Description != nil {
		if err := d.Set("description", *foundType.Description); err != nil {
			return diag.FromErr(err)
		}
	}
	if foundType.UiContributionId != nil {
		if err := d.Set("ui_contribution_id", *foundType.UiContributionId); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set authentication schemes list
	authSchemes := make([]interface{}, 0)
	if foundType.AuthenticationSchemes != nil {
		for _, scheme := range *foundType.AuthenticationSchemes {
			if scheme.Scheme != nil {
				authSchemes = append(authSchemes, *scheme.Scheme)
			}
		}
	}
	if err := d.Set("authentication_schemes", authSchemes); err != nil {
		return diag.FromErr(err)
	}

	// Extract parameters from InputDescriptors
	parameters := make([]map[string]interface{}, 0)
	if foundType.InputDescriptors != nil {
		for _, descriptor := range *foundType.InputDescriptors {
			param := make(map[string]interface{})
			if descriptor.Id != nil {
				param["name"] = *descriptor.Id
			}
			if descriptor.Values != nil {
				if descriptor.Values.DefaultValue != nil {
					param["default_value"] = *descriptor.Values.DefaultValue
				}
				if descriptor.Values.PossibleValues != nil {
					possValues := make([]string, 0, len(*descriptor.Values.PossibleValues))
					for _, v := range *descriptor.Values.PossibleValues {
						if v.Value != nil {
							possValues = append(possValues, *v.Value)
						}
					}
					param["possible_values"] = possValues
				}
			}
			parameters = append(parameters, param)
		}
	}
	if err := d.Set("parameters", parameters); err != nil {
		return diag.FromErr(err)
	}

	// Extract authorization parameters if authorization_scheme is provided
	if authScheme != "" {
		authParameters, err := extractAuthorizationParameters(foundType, authScheme)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("authorization_parameters", authParameters); err != nil {
			return diag.FromErr(err)
		}
	} else {
		// Set to null/empty if no authorization_scheme provided
		if err := d.Set("authorization_parameters", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func extractAuthorizationParameters(endpointType *serviceendpoint.ServiceEndpointType, authScheme string) ([]map[string]interface{}, error) {
	if endpointType.AuthenticationSchemes == nil {
		return nil, fmt.Errorf("no authentication schemes available for service endpoint type '%s'", *endpointType.Name)
	}

	for _, scheme := range *endpointType.AuthenticationSchemes {
		if scheme.Scheme != nil && strings.EqualFold(*scheme.Scheme, authScheme) {
			authParams := make([]map[string]interface{}, 0)
			if scheme.InputDescriptors != nil {
				for _, descriptor := range *scheme.InputDescriptors {
					param := make(map[string]interface{})
					if descriptor.Id != nil {
						param["name"] = *descriptor.Id
					}
					if descriptor.Values != nil {
						if descriptor.Values.DefaultValue != nil {
							param["default_value"] = *descriptor.Values.DefaultValue
						}
						if descriptor.Values.PossibleValues != nil {
							possValues := make([]string, 0, len(*descriptor.Values.PossibleValues))
							for _, v := range *descriptor.Values.PossibleValues {
								if v.Value != nil {
									possValues = append(possValues, *v.Value)
								}
							}
							param["possible_values"] = possValues
						}
					}
					authParams = append(authParams, param)
				}
			}
			return authParams, nil
		}
	}

	// Build list of available schemes for error message
	availableSchemes := make([]string, 0)
	for _, scheme := range *endpointType.AuthenticationSchemes {
		if scheme.Scheme != nil {
			availableSchemes = append(availableSchemes, *scheme.Scheme)
		}
	}

	return nil, fmt.Errorf("authorization scheme '%s' not found for service endpoint type '%s'. Available schemes: %v",
		authScheme, *endpointType.Name, availableSchemes)
}
