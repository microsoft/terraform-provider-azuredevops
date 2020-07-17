package policy

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

const (
	buildDefinitionID       = "build_definition_id"
	policyDisplayName       = "display_name"
	manualQueueOnly         = "manual_queue_only"
	queueOnSourceUpdateOnly = "queue_on_source_update_only"
	validDuration           = "valid_duration"
	filenamePatterns        = "filename_patterns"
)

type buildValidationPolicySettings struct {
	BuildDefinitionID       int      `json:"buildDefinitionId"`
	PolicyDisplayName       string   `json:"displayName"`
	ManualQueueOnly         bool     `json:"manualQueueOnly"`
	QueueOnSourceUpdateOnly bool     `json:"queueOnSourceUpdateOnly"`
	ValidDuration           int      `json:"validDuration"`
	FilenamePatterns        []string `json:"filenamePatterns"`
}

// ResourceBranchPolicyBuildValidation schema and implementation for build validation policy resource
func ResourceBranchPolicyBuildValidation() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: buildValidationFlattenFunc,
		ExpandFunc:  buildValidationExpandFunc,
		PolicyType:  BuildValidation,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema
	settingsSchema[buildDefinitionID] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	settingsSchema[policyDisplayName] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
	settingsSchema[manualQueueOnly] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  false,
		Optional: true,
	}
	settingsSchema[queueOnSourceUpdateOnly] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  true,
		Optional: true,
	}
	settingsSchema[validDuration] = &schema.Schema{
		Type:         schema.TypeInt,
		Default:      720,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
	}
	settingsSchema[filenamePatterns] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	return resource
}

func buildValidationFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf("Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := buildValidationPolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings[buildDefinitionID] = policySettings.BuildDefinitionID
	settings[policyDisplayName] = policySettings.PolicyDisplayName
	settings[manualQueueOnly] = policySettings.ManualQueueOnly
	settings[queueOnSourceUpdateOnly] = policySettings.QueueOnSourceUpdateOnly
	settings[validDuration] = policySettings.ValidDuration
	settings[filenamePatterns] = policySettings.FilenamePatterns

	d.Set(SchemaSettings, settingsList)
	return nil
}

func buildValidationExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})
	policySettings := policyConfig.Settings.(map[string]interface{})

	policySettings["buildDefinitionId"] = settings[buildDefinitionID].(int)
	policySettings["displayName"] = settings[policyDisplayName].(string)
	policySettings["manualQueueOnly"] = settings[manualQueueOnly].(bool)
	policySettings["queueOnSourceUpdateOnly"] = settings[queueOnSourceUpdateOnly].(bool)
	policySettings["validDuration"] = settings[validDuration].(int)
	policySettings["filenamePatterns"] = expandFilenamePatterns(settings[filenamePatterns].(*schema.Set))

	return policyConfig, projectID, nil
}

func expandFilenamePatterns(patterns *schema.Set) *[]string {
	patternsList := patterns.List()
	patternsArray := make([]string, len(patternsList))

	for i, variableGroup := range patternsList {
		patternsArray[i] = variableGroup.(string)
	}

	return &patternsArray
}
