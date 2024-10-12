package branch

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
)

// ResourceBranchPolicyWorkItemLinking schema and implementation for min reviewer policy resource
func ResourceBranchPolicyWorkItemLinking() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: workItemLinkingFlattenFunc,
		ExpandFunc:  workItemLinkingExpandFunc,
		PolicyType:  WorkItemLinking,
	})
	return resource
}

func workItemLinkingFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	settingsList := d.Get("settings").([]interface{})

	d.Set("settings", settingsList)
	return nil
}

func workItemLinkingExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}
	return policyConfig, projectID, nil
}
