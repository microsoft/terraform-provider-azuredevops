package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

const UNIT = 1024 * 1024

func ResourceRepositoryMaxFileSize() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: fileSizeFlattenFunc,
		ExpandFunc:  fileSizeExpandFunc,
		PolicyType:  FileSize,
	})

	resource.Schema["max_file_size"] = &schema.Schema{
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validation.IntInSlice([]int{1, 2, 5, 10, 100, 200}),
	}
	return resource

}

func fileSizeFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	_ = d.Set("max_file_size", policySettings["maximumGitBlobSizeInBytes"].(float64)/UNIT)
	return nil
}

func fileSizeExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["maximumGitBlobSizeInBytes"] = d.Get("max_file_size").(int) * UNIT
	return policyConfig, projectID, nil
}
