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

type autoReviewerPolicySettings struct {
	SubmitterCanVote     bool     `json:"creatorVoteCounts"`
	AutoReviewerIds      []string `json:"requiredReviewerIds"`
	PathFilters          []string `json:"filenamePatterns"`
	DisplayMessage       string   `json:"message"`
	MinimumApproverCount int      `json:"minimumApproverCount"`
}

func ResourceBranchPolicyAutoReviewers() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: autoReviewersFlattenFunc,
		ExpandFunc:  autoReviewersExpandFunc,
		PolicyType:  AutoReviewers,
	})

	settingsSchema := resource.Schema["settings"].Elem.(*schema.Resource).Schema
	maps.Copy(settingsSchema, map[string]*schema.Schema{
		"auto_reviewer_ids": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"path_filters": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"message": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},

		"submitter_can_vote": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
		},

		"minimum_number_of_reviewers": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(1),
		},
	})
	return resource
}

func autoReviewersFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf("unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := autoReviewerPolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf("unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings["submitter_can_vote"] = policySettings.SubmitterCanVote
	settings["auto_reviewer_ids"] = policySettings.AutoReviewerIds
	settings["path_filters"] = policySettings.PathFilters
	settings["message"] = policySettings.DisplayMessage
	settings["minimum_number_of_reviewers"] = policySettings.MinimumApproverCount
	_ = d.Set("settings", settingsList)
	return nil
}

func autoReviewersExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["creatorVoteCounts"] = settings["submitter_can_vote"].(bool)
	policySettings["message"] = settings["message"].(string)
	policySettings["minimumApproverCount"] = settings["minimum_number_of_reviewers"].(int)

	if value, ok := settings["auto_reviewer_ids"]; ok {
		var reviewersID []string
		for _, item := range value.([]interface{}) {
			reviewersID = append(reviewersID, item.(string))
		}
		policySettings["requiredReviewerIds"] = reviewersID
	}

	if value, ok := settings["path_filters"]; ok {
		var pathFilterSettings []string
		for _, item := range value.([]interface{}) {
			pathFilterSettings = append(pathFilterSettings, item.(string))
		}
		policySettings["filenamePatterns"] = pathFilterSettings
	}

	return policyConfig, projectID, nil
}
