package approvalsandchecks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

var evaluateBranchProtectionDefVersion = "0.0.1"
var evaluateBranchProtectionDefId = "86b05a0c-73e6-4f7d-b3cf-e38f3b39a75b"

var evaluateBranchProtectionDef = map[string]interface{}{
	"id":      evaluateBranchProtectionDefId,
	"name":    "evaluatebranchProtection",
	"version": evaluateBranchProtectionDefVersion,
}

// ResourceBranchControlCheck schema and implementation for branch check resources
func ResourceCheckBranchControl() *schema.Resource {
	r := genBaseCheckResource(flattenBranchControlCheck, expandBranchControlCheck)

	r.Schema["allowed_branches"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  `*`,
	}
	r.Schema["verify_branch_protection"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	r.Schema["ignore_unknown_protection_status"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	r.Schema["display_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "Managed by Terraform",
		ValidateFunc: validation.StringIsNotEmpty,
	}

	r.Schema["timeout"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1440,
		ValidateFunc: validation.IntBetween(1, 2147483647),
	}

	return r
}

func flattenBranchControlCheck(d *schema.ResourceData, branchControlCheck *pipelineschecksextras.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, branchControlCheck, projectID)
	if err != nil {
		return err
	}

	if branchControlCheck.Settings == nil {
		return fmt.Errorf("Settings nil")
	}

	branchControlCheck.Type.Id = converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7")

	if displayName, found := branchControlCheck.Settings.(map[string]interface{})["displayName"]; found {
		d.Set("display_name", displayName.(string))
	} else {
		return fmt.Errorf("displayName setting not found")
	}

	if definitionRefMap, found := branchControlCheck.Settings.(map[string]interface{})["definitionRef"]; found {
		definitionRef := definitionRefMap.(map[string]interface{})
		if id, found := definitionRef["id"]; found {
			if !strings.EqualFold(id.(string), evaluateBranchProtectionDefId) {
				return fmt.Errorf("invalid definitionRef id")
			}
		} else {
			return fmt.Errorf("definitionRef ID not found. Expect ID: %s", evaluateBranchProtectionDefId)
		}
		if version, found := definitionRef["version"]; found {
			if version != evaluateBranchProtectionDefVersion {
				return fmt.Errorf("unsupported definitionRef version. Expect version: %s", evaluateBranchProtectionDefVersion)
			}
		} else {
			return fmt.Errorf("unsupported definitionRef version")
		}
	} else {
		return fmt.Errorf("definitionRef not found")
	}

	if inputMap, found := branchControlCheck.Settings.(map[string]interface{})["inputs"]; found {
		inputs := inputMap.(map[string]interface{})
		if AllowedBranches, found := inputs["allowedBranches"]; found {
			d.Set("allowed_branches", AllowedBranches)
		} else {
			return fmt.Errorf("allowedBranches input not found")
		}
		if verifyBranchProtection, found := inputs["ensureProtectionOfBranch"]; found {
			value, err := strconv.ParseBool(verifyBranchProtection.(string))
			if err != nil {
				return err
			}
			d.Set("verify_branch_protection", value)
		} else {
			return fmt.Errorf("ensureProtectionOfBranch input not found")
		}
		if ignoreUnknownProtectionStatus, found := inputs["allowUnknownStatusBranch"]; found {
			value, err := strconv.ParseBool(ignoreUnknownProtectionStatus.(string))
			if err != nil {
				return err
			}
			d.Set("ignore_unknown_protection_status", value)
		} else {
			return fmt.Errorf("allowUnknownStatusBranch input not found")
		}
	} else {
		return fmt.Errorf("inputs not found")
	}

	if branchControlCheck.Timeout != nil {
		d.Set("timeout", *branchControlCheck.Timeout)
	}

	return nil
}

func expandBranchControlCheck(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	inputs := map[string]interface{}{
		"allowedBranches":          d.Get("allowed_branches").(string),
		"ensureProtectionOfBranch": strconv.FormatBool(d.Get("verify_branch_protection").(bool)),
		"allowUnknownStatusBranch": strconv.FormatBool(d.Get("ignore_unknown_protection_status").(bool)),
	}

	settings := map[string]interface{}{}
	settings["inputs"] = inputs
	settings["definitionRef"] = evaluateBranchProtectionDef
	settings["displayName"] = d.Get("display_name").(string)

	return doBaseExpansion(d, approvalAndCheckType.BranchProtection, settings, converter.ToPtr(d.Get("timeout").(int)))
}
