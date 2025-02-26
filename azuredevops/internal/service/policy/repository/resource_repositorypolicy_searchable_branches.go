package repository

import (
	"maps"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
)

type searchbranchesPolicySettings struct {
	SearchBranches []string `json:"searchBranches"`
}

func ResourceRepositorySearchableBranches() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: searchableBranchesFlattenFunc,
		ExpandFunc:  searchableBranchesExpandFunc,
		PolicyType:  SearchableBranches,
	})

	maps.Copy(resource.Schema, map[string]*schema.Schema{
		"searchable_branches": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"blocking": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		// API only accepts a single repository as scope.
		"repository_ids": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},
	})

	return resource
}

func searchableBranchesFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	_ = d.Set("searchable_branches", policySettings["searchBranches"].([]interface{}))
	return nil
}

func searchableBranchesExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["searchBranches"] = d.Get("searchable_branches")

	// Overriding blocking and enabled as it has no use for this policy setting
	policySettings["blocking"] = false
	policySettings["enabled"] = false

	return policyConfig, projectID, nil
}
