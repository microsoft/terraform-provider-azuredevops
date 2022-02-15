package branch

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/policy"
)

type minReviewerPolicySettings struct {
	ApprovalCount                     int  `json:"minimumApproverCount" tf:"reviewer_count"`
	SubmitterCanVote                  bool `json:"creatorVoteCounts" tf:"submitter_can_vote"`
	AllowCompletionWithRejectsOrWaits bool `json:"allowDownvotes" tf:"allow_completion_with_rejects_or_waits"`
	OnPushResetApprovedVotes          bool `json:"resetOnSourcePush" tf:"on_push_reset_approved_votes"`
	OnLastIterationRequireVote        bool `json:"requireVoteOnLastIteration" tf:"on_last_iteration_require_vote"`
	OnPushResetAllVotes               bool `json:"resetRejectionsOnSourcePush" tf:"on_push_reset_all_votes"`
	LastPusherCannotVote              bool `json:"blockLastPusherVote" tf:"last_pusher_cannot_approve"`
}

// ResourceBranchPolicyMinReviewers schema and implementation for min reviewer policy resource
func ResourceBranchPolicyMinReviewers() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: minReviewersFlattenFunc,
		ExpandFunc:  minReviewersExpandFunc,
		PolicyType:  MinReviewerCount,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema

	// Dynamically create the schema based on the minReviewerPolicySettings tags
	metaField := reflect.TypeOf(minReviewerPolicySettings{})
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
		if metaField.Field(i).Type == reflect.TypeOf(0) {
			settingsSchema[tfName] = &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			}
		}
		if conflicts, ok := metaField.Field(i).Tag.Lookup("ConflictsWith"); ok {
			if _, ok := settingsSchema[conflicts]; ok {
				settingsSchema[conflicts].ConflictsWith = []string{SchemaSettings + ".0." + tfName}
				settingsSchema[tfName].ConflictsWith = []string{SchemaSettings + ".0." + conflicts}
			}
		}
	}
	return resource
}

// API to TF
func minReviewersFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf("Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := minReviewerPolicySettings{}
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
		if tipe.Field(i).Type == reflect.TypeOf(0) {
			settings[tfName] = ps.Field(i).Int()
		}
	}

	d.Set(SchemaSettings, settingsList)
	return nil
}

// From TF to API
func minReviewersExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})

	tipe := reflect.TypeOf(minReviewerPolicySettings{})
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
		if tipe.Field(i).Type == reflect.TypeOf(0) {
			policySettings[apiName] = settings[tfName].(int)
		}
	}

	return policyConfig, projectID, nil
}
