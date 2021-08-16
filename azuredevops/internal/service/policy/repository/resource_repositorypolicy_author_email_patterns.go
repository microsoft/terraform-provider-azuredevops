package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/policy/branch"
)

func ResourceRepositoryPolicyAuthorEmailPatterns() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: authorEmailPatternFlattenFunc,
		ExpandFunc:  authorEmailPatternExpandFunc,
		PolicyType:  AuthorEmailPattern,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema
	settingsSchema["author_email_patterns"] = &schema.Schema{
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

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})
	if policySettings["authorEmailPatterns"] != nil {
		settings["author_email_patterns"] = policySettings["authorEmailPatterns"].([]interface{})
	}
	_ = d.Set(branch.SchemaSettings, settingsList)
	return nil
}

func authorEmailPatternExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["authorEmailPatterns"] = settings["author_email_patterns"]
	return policyConfig, projectID, nil
}
