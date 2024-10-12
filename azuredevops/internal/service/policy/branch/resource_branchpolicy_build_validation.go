package branch

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
)

type buildValidationPolicySettings struct {
	BuildDefinitionID       int      `json:"buildDefinitionId"`
	PolicyDisplayName       string   `json:"displayName"`
	ManualQueueOnly         bool     `json:"manualQueueOnly"`
	QueueOnSourceUpdateOnly bool     `json:"queueOnSourceUpdateOnly"`
	ValidDuration           int      `json:"validDuration"`
	FilenamePatterns        []string `json:"filenamePatterns"`
}

func ResourceBranchPolicyBuildValidation() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: buildValidationFlattenFunc,
		ExpandFunc:  buildValidationExpandFunc,
		PolicyType:  BuildValidation,
	})

	settingsSchema := resource.Schema["settings"].Elem.(*schema.Resource).Schema
	maps.Copy(settingsSchema, map[string]*schema.Schema{
		"build_definition_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"display_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"manual_queue_only": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
		},
		"queue_on_source_update_only": {
			Type:     schema.TypeBool,
			Default:  true,
			Optional: true,
		},
		"valid_duration": {
			Type:         schema.TypeInt,
			Default:      720,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"filename_patterns": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	})
	return resource
}

func buildValidationFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf(" Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := buildValidationPolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf(" Unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings["build_definition_id"] = policySettings.BuildDefinitionID
	settings["display_name"] = policySettings.PolicyDisplayName
	settings["manual_queue_only"] = policySettings.ManualQueueOnly
	settings["queue_on_source_update_only"] = policySettings.QueueOnSourceUpdateOnly
	settings["valid_duration"] = policySettings.ValidDuration
	settings["filename_patterns"] = policySettings.FilenamePatterns

	d.Set("settings", settingsList)
	return nil
}

func buildValidationExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})
	policySettings := policyConfig.Settings.(map[string]interface{})

	policySettings["buildDefinitionId"] = settings["build_definition_id"].(int)
	policySettings["displayName"] = settings["display_name"].(string)
	policySettings["manualQueueOnly"] = settings["manual_queue_only"].(bool)
	policySettings["queueOnSourceUpdateOnly"] = settings["queue_on_source_update_only"].(bool)
	policySettings["validDuration"] = settings["valid_duration"].(int)
	policySettings["filenamePatterns"] = expandFilenamePatterns(settings["filename_patterns"].([]interface{}))

	return policyConfig, projectID, nil
}

func expandFilenamePatterns(patterns []interface{}) *[]string {
	patternsArray := make([]string, len(patterns))

	for i, variableGroup := range patterns {
		patternsArray[i] = variableGroup.(string)
	}

	return &patternsArray
}
