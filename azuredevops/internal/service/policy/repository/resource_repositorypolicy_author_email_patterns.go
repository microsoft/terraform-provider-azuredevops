package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/policy"
)

func ResourceRepositoryPolicyAuthorEmailPatterns() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: authorEmailPatternFlattenFunc,
		ExpandFunc:  authorEmailPatternExpandFunc,
		PolicyType:  AuthorEmailPattern,
	})

	resource.Schema["author_email_patterns"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	return resource
}

func authorEmailPatternFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	_ = d.Set("author_email_patterns", policySettings["authorEmailPatterns"])
	return nil
}

func authorEmailPatternExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["authorEmailPatterns"] = d.Get("author_email_patterns").([]interface{})
	return policyConfig, projectID, nil
}
