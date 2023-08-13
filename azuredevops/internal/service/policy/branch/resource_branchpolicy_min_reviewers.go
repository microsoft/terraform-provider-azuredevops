package branch

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
)

// ResourceBranchPolicyMinReviewers schema and implementation for min reviewer policy resource
func ResourceBranchPolicyMinReviewers() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: minReviewersFlattenFunc,
		ExpandFunc:  minReviewersExpandFunc,
		PolicyType:  MinReviewerCount,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema

	settingsSchema["reviewer_count"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
	}

	settingsSchema["submitter_can_vote"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	settingsSchema["allow_completion_with_rejects_or_waits"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	settingsSchema["on_last_iteration_require_vote"] = &schema.Schema{
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"settings.0.on_push_reset_approved_votes", "settings.0.on_push_reset_all_votes"},
	}

	settingsSchema["on_push_reset_approved_votes"] = &schema.Schema{
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"settings.0.on_last_iteration_require_vote"},
	}

	settingsSchema["on_push_reset_all_votes"] = &schema.Schema{
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"settings.0.on_last_iteration_require_vote"},
	}

	settingsSchema["last_pusher_cannot_approve"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	return resource
}

// API to TF
func minReviewersFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings["reviewer_count"] = policySettings["minimumApproverCount"]
	settings["submitter_can_vote"] = policySettings["creatorVoteCounts"]
	settings["allow_completion_with_rejects_or_waits"] = policySettings["allowDownvotes"]
	settings["on_push_reset_approved_votes"] = policySettings["resetOnSourcePush"]
	settings["on_last_iteration_require_vote"] = policySettings["requireVoteOnLastIteration"]
	settings["on_push_reset_all_votes"] = policySettings["resetRejectionsOnSourcePush"]
	settings["last_pusher_cannot_approve"] = policySettings["blockLastPusherVote"]

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
	policySettings["minimumApproverCount"] = settings["reviewer_count"]
	policySettings["creatorVoteCounts"] = settings["submitter_can_vote"]
	policySettings["allowDownvotes"] = settings["allow_completion_with_rejects_or_waits"]
	policySettings["requireVoteOnLastIteration"] = settings["on_last_iteration_require_vote"]
	policySettings["resetOnSourcePush"] = settings["on_push_reset_approved_votes"]
	policySettings["blockLastPusherVote"] = settings["last_pusher_cannot_approve"]

	resetRejectionsOnSourcePush := settings["on_push_reset_all_votes"].(bool)
	policySettings["resetRejectionsOnSourcePush"] = resetRejectionsOnSourcePush
	if resetRejectionsOnSourcePush {
		policySettings["resetOnSourcePush"] = true
	}

	return policyConfig, projectID, nil
}
