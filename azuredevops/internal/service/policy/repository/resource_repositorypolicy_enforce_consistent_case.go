package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
)

func ResourceRepositoryEnforceConsistentCase() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: enforceConsistentCaseFlattenFunc,
		ExpandFunc:  enforceConsistentCaseExpandFunc,
		PolicyType:  CaseEnforcement,
	})
	resource.Schema["enforce_consistent_case"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	}
	return resource
}

func enforceConsistentCaseFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	_ = d.Set("enforce_consistent_case", policySettings["enforceConsistentCase"])
	return nil
}

func enforceConsistentCaseExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["enforceConsistentCase"] = d.Get("enforce_consistent_case").(bool)
	return policyConfig, projectID, nil
}
