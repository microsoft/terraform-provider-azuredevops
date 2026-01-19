package branch

import (
	"maps"

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

	settingsSchema := resource.Schema["settings"].Elem.(*schema.Resource).Schema
	maps.Copy(settingsSchema, map[string]*schema.Schema{
		"reviewer_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"submitter_can_vote": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"allow_completion_with_rejects_or_waits": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"on_last_iteration_require_vote": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"settings.0.on_push_reset_approved_votes", "settings.0.on_push_reset_all_votes"},
		},

		"on_each_iteration_require_vote": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"on_push_reset_approved_votes": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"settings.0.on_last_iteration_require_vote"},
		},

		"on_push_reset_all_votes": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"settings.0.on_last_iteration_require_vote"},
		},

		"last_pusher_cannot_approve": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	})
	return resource
}

// API to TF
func minReviewersFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings["reviewer_count"] = policySettings["minimumApproverCount"]
	settings["submitter_can_vote"] = policySettings["creatorVoteCounts"]
	settings["allow_completion_with_rejects_or_waits"] = policySettings["allowDownvotes"]
	settings["on_push_reset_approved_votes"] = policySettings["resetOnSourcePush"]
	settings["on_last_iteration_require_vote"] = policySettings["requireVoteOnLastIteration"]
	settings["on_push_reset_all_votes"] = policySettings["resetRejectionsOnSourcePush"]
	settings["last_pusher_cannot_approve"] = policySettings["blockLastPusherVote"]
	settings["on_each_iteration_require_vote"] = policySettings["requireVoteOnEachIteration"]

	d.Set("settings", settingsList)
	return nil
}

// From TF to API
func minReviewersExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["minimumApproverCount"] = settings["reviewer_count"]
	policySettings["creatorVoteCounts"] = settings["submitter_can_vote"]
	policySettings["allowDownvotes"] = settings["allow_completion_with_rejects_or_waits"]
	policySettings["requireVoteOnLastIteration"] = settings["on_last_iteration_require_vote"]
	policySettings["resetOnSourcePush"] = settings["on_push_reset_approved_votes"]
	policySettings["blockLastPusherVote"] = settings["last_pusher_cannot_approve"]
	policySettings["requireVoteOnEachIteration"] = settings["on_each_iteration_require_vote"]

	resetRejectionsOnSourcePush := settings["on_push_reset_all_votes"].(bool)
	policySettings["resetRejectionsOnSourcePush"] = resetRejectionsOnSourcePush
	if resetRejectionsOnSourcePush {
		policySettings["resetOnSourcePush"] = true
	}

	return policyConfig, projectID, nil
}
