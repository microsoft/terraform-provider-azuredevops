package policy

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

func ResourceRepositoryFilePathPatterns() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: filePathPatternFlattenFunc,
		ExpandFunc:  filePathPatternExpandFunc,
		PolicyType:  FilePathPattern,
	})

	resource.Schema[SchemaSettings] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"filepath_patterns": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				SchemaScope: {
					Type:     schema.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SchemaRepositoryID: {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
	return resource
}

func filePathPatternFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})
	settings["filepath_patterns"] = policySettings["filenamePatterns"].([]interface{})
	_ = d.Set(SchemaSettings, settingsList)
	return nil
}

func filePathPatternExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["filenamePatterns"] = settings["filepath_patterns"]
	return policyConfig, projectID, nil
}
