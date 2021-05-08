package policy

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

type autoReviewerPolicySettings struct {
	SubmitterCanVote bool     `json:"creatorVoteCounts"`
	AutoReviewerIds  []string `json:"requiredReviewerIds"`
	PathFilters      []string `json:"filenamePatterns"`
	DisplayMessage   string   `json:"message"`
}

const (
	autoReviewerIds        = "auto_reviewer_ids"
	pathFilters            = "path_filters"
	displayMessage         = "message"
	schemaSubmitterCanVote = "submitter_can_vote"
)

// ResourceBranchPolicyAutoReviewers schema and implementation for automatic code reviewer policy resource
func ResourceBranchPolicyAutoReviewers() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: autoReviewersFlattenFunc,
		ExpandFunc:  autoReviewersExpandFunc,
		PolicyType:  AutoReviewers,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema
	settingsSchema[autoReviewerIds] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	settingsSchema[pathFilters] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	settingsSchema[displayMessage] = &schema.Schema{
		Type:     schema.TypeString,
		Required: false,
		Default:  "",
		Optional: true,
	}
	settingsSchema[schemaSubmitterCanVote] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  false,
		Optional: true,
	}

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

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings[schemaSubmitterCanVote] = policySettings.SubmitterCanVote
	settings[autoReviewerIds] = policySettings.AutoReviewerIds
	settings[pathFilters] = policySettings.PathFilters
	settings[displayMessage] = policySettings.DisplayMessage
	_ = d.Set(SchemaSettings, settingsList)
	return nil
}

func autoReviewersExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["creatorVoteCounts"] = settings[schemaSubmitterCanVote].(bool)
	policySettings["message"] = settings[displayMessage].(string)

	if value, ok := settings[autoReviewerIds]; ok {
		var reviewersID []string
		for _, item := range value.([]interface{}) {
			reviewersID = append(reviewersID, item.(string))
		}
		policySettings["requiredReviewerIds"] = reviewersID
	}

	if value, ok := settings[pathFilters]; ok {
		var pathFilters []string
		for _, item := range value.([]interface{}) {
			pathFilters = append(pathFilters, item.(string))
		}
		policySettings["filenamePatterns"] = pathFilters
	}

	return policyConfig, projectID, nil
}
