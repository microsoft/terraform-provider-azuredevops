package approvalsandchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

var validRepositoryTypes = []string{"git", "github", "bitbucket"}

// ResourceCheckRequiredTemplate schema and implementation for required template check resources
func ResourceCheckRequiredTemplate() *schema.Resource {
	r := genBaseCheckResource(flattenCheckRequiredTemplate, expandCheckRequiredTemplate)

	r.Schema["required_template"] = &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"repository_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "git",
					ValidateFunc: validation.StringInSlice(validRepositoryTypes, false),
				},
				"repository_name": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"repository_ref": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"template_path": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}

	return r
}

func flattenCheckRequiredTemplate(d *schema.ResourceData, check *pipelineschecksextras.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, check, projectID)
	if err != nil {
		return err
	}

	if check.Settings == nil {
		return fmt.Errorf("settings nil")
	}

	check.Type = approvalAndCheckType.ExtendsCheck

	if check.Settings == nil {
		return fmt.Errorf("Settings nil")
	}

	var reqTemplList []map[string]interface{}
	if extendsChecksConfig, found := check.Settings.(map[string]interface{})["extendsChecks"]; found {
		extendsChecks := extendsChecksConfig.([]interface{})
		for _, ec := range extendsChecks {
			ecMap := ec.(map[string]interface{})
			reqTempl := map[string]interface{}{
				"repository_type": ecMap["repositoryType"],
				"repository_name": ecMap["repositoryName"],
				"repository_ref":  ecMap["repositoryRef"],
				"template_path":   ecMap["templatePath"],
			}
			reqTemplList = append(reqTemplList, reqTempl)
		}
	} else {
		return fmt.Errorf("extendsChecks not found")
	}
	d.Set("required_template", reqTemplList)

	return nil
}

func expandCheckRequiredTemplate(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	var extendsChecks []interface{}
	if v, ok := d.GetOk("required_template"); ok {
		reqTemplList := v.(*schema.Set).List()
		for _, reqTempl := range reqTemplList {
			reqTemplMap := reqTempl.(map[string]interface{})
			extendsChecks = append(extendsChecks, map[string]interface{}{
				"repositoryType": reqTemplMap["repository_type"],
				"repositoryName": reqTemplMap["repository_name"],
				"repositoryRef":  reqTemplMap["repository_ref"],
				"templatePath":   reqTemplMap["template_path"],
			})
		}
	}
	settings := map[string]interface{}{}
	settings["extendsChecks"] = extendsChecks
	return doBaseExpansion(d, approvalAndCheckType.ExtendsCheck, settings, nil)
}
