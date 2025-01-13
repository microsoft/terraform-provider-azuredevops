package approvalsandchecks

import (
	"fmt"
	"maps"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

type CompleteEvent string

type CompleteEventValuesType struct {
	Callback    CompleteEvent
	ApiResponse CompleteEvent
}

var CompleteEventValues = CompleteEventValuesType{
	Callback:    "Callback",
	ApiResponse: "ApiResponse",
}

func ResourceCheckRestAPI() *schema.Resource {
	r := genBaseCheckResource(flattenRestAPI, expandRestAPI)

	maps.Copy(r.Schema, map[string]*schema.Schema{
		"display_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"connected_service_name_selector": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"connectedServiceName",
				"connectedServiceNameARM",
			}, false),
		},

		"connected_service_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"method": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "PATCH",
			}, false),
		},

		"body": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"headers": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"retry_interval": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"success_criteria": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"url_suffix": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"variable_group_name": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"completion_event": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(CompleteEventValues.Callback), string(CompleteEventValues.ApiResponse),
			}, false),
			Default: string(CompleteEventValues.Callback),
		},

		"timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1440,
			ValidateFunc: validation.IntBetween(1, 43200),
		},
	})

	return r
}

func flattenRestAPI(d *schema.ResourceData, check *pipelineschecksextras.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, check, projectID)
	if err != nil {
		return err
	}

	if check.Timeout != nil {
		d.Set("timeout", *check.Timeout)
	}

	if check.Settings != nil {
		settings := check.Settings.(map[string]interface{})
		if v, ok := settings["displayName"]; ok {
			d.Set("display_name", v.(string))
		}

		if v, exist := settings["retryInterval"]; exist {
			d.Set("retry_interval", v.(float64))
		}

		if v, exist := settings["linkedVariableGroup"]; exist {
			d.Set("variable_group_name", v.(string))
		}

		if v, ok := settings["inputs"]; ok {
			inputs := v.(map[string]interface{})
			if v, exist := inputs["connectedServiceNameSelector"]; exist {
				serviceNameSelector := v.(string)
				d.Set("connected_service_name_selector", serviceNameSelector)
				if serviceNameSelector == "connectedServiceName" {
					if v, exist := inputs["connectedServiceName"]; exist {
						d.Set("connected_service_name", v.(string))
					}
				} else if serviceNameSelector == "connectedServiceNameARM" {
					if v, exist := inputs["connectedServiceNameARM"]; exist {
						d.Set("connected_service_name", v.(string))
					}
				}
			}

			if v, exist := inputs["method"]; exist {
				d.Set("method", v.(string))
			}

			if v, exist := inputs["headers"]; exist {
				d.Set("headers", v.(string))
			}

			if v, exist := inputs["body"]; exist {
				d.Set("body", v.(string))
			}

			if v, exist := inputs["waitForCompletion"]; exist {
				waitForCompletion, err := strconv.ParseBool(v.(string))
				if err != nil {
					return fmt.Errorf(" parsing `waitForCompletion`: %v", err)
				}
				d.Set("completion_event", string(CompleteEventValues.Callback))
				if !waitForCompletion {
					d.Set("completion_event", string(CompleteEventValues.ApiResponse))
				}
			}

			if v, exist := inputs["successCriteria"]; exist {
				d.Set("success_criteria", v.(string))
			}
		}
	}
	return nil
}

func expandRestAPI(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	settings := map[string]interface{}{
		"definitionRef": map[string]string{
			"id":      "9c3e8943-130d-4c78-ac63-8af81df62dfb",
			"name":    "InvokeRESTAPI",
			"version": "1.220.0",
		},
		"displayName": d.Get("display_name").(string),
	}

	// inputs
	serviceNameSelector := d.Get("connected_service_name_selector").(string)
	input := map[string]interface{}{
		"connectedServiceNameSelector": serviceNameSelector,
		"method":                       d.Get("method").(string),
	}
	if serviceNameSelector == "connectedServiceName" {
		input["connectedServiceName"] = d.Get("connected_service_name").(string)
	} else if serviceNameSelector == "connectedServiceNameARM" {
		input["connectedServiceNameARM"] = d.Get("connected_service_name").(string)
	}

	if v, ok := d.GetOk("headers"); ok {
		input["headers"] = v.(string)
	}

	if v, ok := d.GetOk("body"); ok {
		input["body"] = v.(string)
	}

	if v, ok := d.GetOk("url_suffix"); ok {
		input["urlSuffix"] = v.(string)
	}

	completionEvent := d.Get("completion_event").(string)
	input["waitForCompletion"] = "true"
	if strings.EqualFold(completionEvent, string(CompleteEventValues.ApiResponse)) {
		input["waitForCompletion"] = "false"
		if v, ok := d.GetOk("success_criteria"); ok {
			input["successCriteria"] = v.(string)
		}
	}
	settings["inputs"] = input
	// inputs end

	// There is no need to retry a Callback check. https://devblogs.microsoft.com/devops/updates-to-approvals-and-checks/
	if v, ok := d.GetOk("retry_interval"); ok {
		retryInterval := v.(int)
		if completionEvent == string(CompleteEventValues.Callback) {
			return nil, "", fmt.Errorf(" Does not need to set `retry_interval` when `completion_event=Callback`.")
		}

		timeout := d.Get("timeout").(int)
		minRetryInterval := timeout / 10
		if minRetryInterval > retryInterval {
			return nil, "", fmt.Errorf(" We require you enter a value of 0 or at least %d,"+
				" to keep the number of retries below 10. Starting Autumn 2023, non-compliant "+
				"checks will fail automatically. Timeout: %d, retryInterval: %d", minRetryInterval, timeout, retryInterval)
		}
		settings["retryInterval"] = retryInterval
	}

	if v, ok := d.GetOk("variable_group_name"); ok {
		settings["linkedVariableGroup"] = v.(string)
	}
	return doBaseExpansion(d, approvalAndCheckType.TaskCheck, settings, converter.ToPtr(d.Get("timeout").(int)))
}
