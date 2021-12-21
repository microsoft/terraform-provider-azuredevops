package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/policy"
)

func ResourceRepositoryMaxPathLength() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: pathLengthFlattenFunc,
		ExpandFunc:  pathLengthExpandFunc,
		PolicyType:  PathLength,
	})

	resource.Schema["max_path_length"] = &schema.Schema{
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validation.IntBetween(1, 10000),
	}
	return resource
}

func pathLengthFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	_ = d.Set("max_path_length", policySettings["maxPathLength"].(float64))
	return nil
}

func pathLengthExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["maxPathLength"] = d.Get("max_path_length").(int)
	return policyConfig, projectID, nil
}
