package approvalsandchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

var validRepositoryTypes = []string{"git", "github", "bitbucket"}

// ResourceCheckTemplate schema and implementation for required template check resources
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

	var reqTemplSet []map[string]interface{}
	if extendsCheckMap, found := check.Settings.(map[string]interface{})["extendsChecks"]; found {
		extendsChecks := extendsCheckMap.([]map[string]interface{})
		for i := range extendsChecks {
			reqTempl := map[string]interface{}{
				"repository_type": extendsChecks[i]["repositoryType"],
				"repository_name": extendsChecks[i]["repositoryName"],
				"repository_ref":  extendsChecks[i]["repositoryRef"],
				"template_path":   extendsChecks[i]["templatePath"],
			}
			reqTemplSet = append(reqTemplSet, reqTempl)
		}
	} else {
		return fmt.Errorf("extendsChecks not found")
	}
	d.Set("required_template", reqTemplSet)

	return nil
}

func expandCheckRequiredTemplate(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	var extendsChecks []map[string]interface{}
	if v, ok := d.GetOk("required_template"); ok {
		reqTemplList := v.(*schema.Set).List()
		for _, reqTempl := range reqTemplList {
			reqTemplMap := reqTempl.(map[string]interface{})
			extendsChecks = append(extendsChecks, map[string]interface{}{
				"repositoryType": reqTemplMap["repository_type"].(string),
				"repositoryName": reqTemplMap["repository_name"].(string),
				"repositoryRef":  reqTemplMap["repository_ref"].(string),
				"templatePath":   reqTemplMap["template_path"].(string),
			})
		}
	}
	settings := map[string]interface{}{}
	settings["extendsChecks"] = extendsChecks
	return doBaseExpansion(d, approvalAndCheckType.ExtendsCheck, settings, nil)
}
