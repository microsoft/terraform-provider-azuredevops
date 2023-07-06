package approvalsandchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

var validRepositoryTypes = []string{"azuregit", "github", "bitbucket"}

// ResourceCheckRequiredTemplate schema and implementation for required template check resources
func ResourceCheckRequiredTemplate() *schema.Resource {
	r := genBaseCheckResource(flattenCheckRequiredTemplate, expandCheckRequiredTemplate)

	r.Schema["required_template"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"repository_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "azuregit",
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
					Required:     true,
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
			var repositoryType string
			if ecMap["repositoryType"].(string) == "git" {
				repositoryType = "azuregit"
			} else {
				repositoryType = ecMap["repositoryType"].(string)
			}
			reqTempl := map[string]interface{}{
				"repository_type": repositoryType,
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
		reqTemplList := v.([]interface{})
		for _, reqTempl := range reqTemplList {
			var repositoryType string
			reqTemplMap := reqTempl.(map[string]interface{})
			if reqTemplMap["repository_type"].(string) == "azuregit" {
				repositoryType = "git"
			} else {
				repositoryType = reqTemplMap["repository_type"].(string)
			}
			extendsChecks = append(extendsChecks, map[string]interface{}{
				"repositoryType": repositoryType,
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
