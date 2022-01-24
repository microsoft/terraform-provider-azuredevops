package branch

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/policy"
)

// ResourceBranchPolicyCommentResolution schema and implementation for min reviewer policy resource
func ResourceBranchPolicyCommentResolution() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: commentResolutionFlattenFunc,
		ExpandFunc:  commentResolutionExpandFunc,
		PolicyType:  CommentResolution,
	})
	return resource
}

func commentResolutionFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})

	d.Set(SchemaSettings, settingsList)
	return nil
}

func commentResolutionExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}
	return policyConfig, projectID, nil
}
