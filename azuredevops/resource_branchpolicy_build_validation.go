package azuredevops

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/branchpolicy"
)

const (
	buildDefinitionID       = "build_definition_id"
	policyDisplayName       = "display_name"
	manualQueueOnly         = "manual_queue_only"
	queueOnSourceUpdateOnly = "queue_on_source_update_only"
	validDuration           = "valid_duration"
)

type buildValidationPolicySettings struct {
	BuildDefinitionID       int    `json:"buildDefinitionId"`
	PolicyDisplayName       string `json:"displayName"`
	ManualQueueOnly         bool   `json:"manualQueueOnly"`
	QueueOnSourceUpdateOnly bool   `json:"queueOnSourceUpdateOnly"`
	ValidDuration           int    `json:"validDuration"`
}

func resourceBranchPolicyBuildValidation() *schema.Resource {
	resource := branchpolicy.GenBasePolicyResource(&branchpolicy.PolicyCrudArgs{
		FlattenFunc: buildValidationFlattenFunc,
		ExpandFunc:  buildValidationExpandFunc,
		PolicyType:  branchpolicy.BuildValidation,
	})

	settingsSchema := resource.Schema[branchpolicy.SchemaSettings].Elem.(*schema.Resource).Schema
	settingsSchema[buildDefinitionID] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	settingsSchema[policyDisplayName] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.NoZeroValues,
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
	return resource
}

func buildValidationFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := branchpolicy.BaseFlattenFunc(d, policyConfig, projectID)
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

	settingsList := d.Get(branchpolicy.SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings[buildDefinitionID] = policySettings.BuildDefinitionID
	settings[policyDisplayName] = policySettings.PolicyDisplayName
	settings[manualQueueOnly] = policySettings.ManualQueueOnly
	settings[queueOnSourceUpdateOnly] = policySettings.QueueOnSourceUpdateOnly
	settings[validDuration] = policySettings.ValidDuration

	d.Set(branchpolicy.SchemaSettings, settingsList)
	return nil
}

func buildValidationExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := branchpolicy.BaseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(branchpolicy.SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})
	policySettings := policyConfig.Settings.(map[string]interface{})

	policySettings["buildDefinitionId"] = settings[buildDefinitionID].(int)
	policySettings["displayName"] = settings[policyDisplayName].(string)
	policySettings["manualQueueOnly"] = settings[manualQueueOnly].(bool)
	policySettings["queueOnSourceUpdateOnly"] = settings[queueOnSourceUpdateOnly].(bool)
	policySettings["validDuration"] = settings[validDuration].(int)

	return policyConfig, projectID, nil
}
