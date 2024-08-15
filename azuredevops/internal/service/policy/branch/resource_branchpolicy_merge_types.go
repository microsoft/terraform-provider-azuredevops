package branch

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
	"maps"
)

type mergeTypePolicySettings struct {
	AllowSquash        bool `json:"allowSquash" tf:"allow_squash"`
	AllowRebase        bool `json:"allowRebase" tf:"allow_rebase_and_fast_forward"`
	AllowNoFastForward bool `json:"allowNoFastForward" tf:"allow_basic_no_fast_forward"`
	AllowRebaseMerge   bool `json:"allowRebaseMerge" tf:"allow_rebase_with_merge"`
}

func ResourceBranchPolicyMergeTypes() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: mergeTypesFlattenFunc,
		ExpandFunc:  mergeTypesExpandFunc,
		PolicyType:  MergeTypes,
	})

	settingsSchema := resource.Schema["settings"].Elem.(*schema.Resource).Schema
	maps.Copy(settingsSchema, map[string]*schema.Schema{
		"allow_squash": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
		},

		"allow_rebase_and_fast_forward": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
		},

		"allow_basic_no_fast_forward": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
		},

		"allow_rebase_with_merge": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
		},
	})

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
		return fmt.Errorf(" Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := mergeTypePolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf(" Unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings["allow_squash"] = policySettings.AllowSquash
	settings["allow_rebase_and_fast_forward"] = policySettings.AllowRebase
	settings["allow_basic_no_fast_forward"] = policySettings.AllowNoFastForward
	settings["allow_rebase_with_merge"] = policySettings.AllowRebaseMerge
	d.Set("settings", settingsList)
	return nil
}

// From TF to API
func mergeTypesExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["allowSquash"] = settings["allow_squash"]
	policySettings["allowRebase"] = settings["allow_rebase_and_fast_forward"]
	policySettings["allowNoFastForward"] = settings["allow_basic_no_fast_forward"]
	policySettings["allowRebaseMerge"] = settings["allow_rebase_with_merge"]

	return policyConfig, projectID, nil
}
