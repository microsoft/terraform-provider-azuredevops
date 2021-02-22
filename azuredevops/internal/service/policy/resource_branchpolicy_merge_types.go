package policy

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

type mergeTypePolicySettings struct {
	AllowSquash        bool `json:"allowSquash" tf:"allow_squash"`
	AllowRebase        bool `json:"allowRebase" tf:"allow_rebase_and_fast_forward"`
	AllowNoFastForward bool `json:"allowNoFastForward" tf:"allow_basic_no_fast_forward"`
	AllowRebaseMerge   bool `json:"allowRebaseMerge" tf:"allow_rebase_with_merge"`
}

// ResourceBranchPolicyMinReviewers schema and implementation for min reviewer policy resource
func ResourceBranchPolicyMergeTypes() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: mergeTypesFlattenFunc,
		ExpandFunc:  mergeTypesExpandFunc,
		PolicyType:  MergeTypes,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema

	// Dynamically create the schema based on the mergeTypePolicySettings tags
	metaField := reflect.TypeOf(mergeTypePolicySettings{})
	// Loop through the fields, adding schema for each one.
	for i := 0; i < metaField.NumField(); i++ {
		tfName := metaField.Field(i).Tag.Get("tf")
		if _, ok := settingsSchema[tfName]; ok {
			continue // skip those which are already set
		}
		if metaField.Field(i).Type == reflect.TypeOf(true) {
			settingsSchema[tfName] = &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			}
		}
	}
	return resource
}

// API to TF
func mergeTypesFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf("Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := mergeTypePolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	tipe := reflect.TypeOf(policySettings)
	for i := 0; i < tipe.NumField(); i++ {
		tfName := tipe.Field(i).Tag.Get("tf")
		ps := reflect.ValueOf(policySettings)
		if tipe.Field(i).Type == reflect.TypeOf(true) {
			settings[tfName] = ps.Field(i).Bool()
		}
	}

	d.Set(SchemaSettings, settingsList)
	return nil
}

// From TF to API
func mergeTypesExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})

	tipe := reflect.TypeOf(mergeTypePolicySettings{})
	for i := 0; i < tipe.NumField(); i++ {
		tags := tipe.Field(i).Tag
		apiName := tags.Get("json")
		tfName := tags.Get("tf")
		if _, ok := policySettings[apiName]; ok {
			continue
		}
		if tipe.Field(i).Type == reflect.TypeOf(true) {
			policySettings[apiName] = settings[tfName].(bool)
		}
	}

	return policyConfig, projectID, nil
}
