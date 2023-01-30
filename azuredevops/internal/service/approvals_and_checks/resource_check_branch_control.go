package approvals_and_checks

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks"
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

	return r
}

func flattenBranchControlCheck(d *schema.ResourceData, branchControlCheck *pipelineschecks.CheckConfiguration, projectID string) error {
	doBaseFlattening(d, branchControlCheck, projectID, evaluateBranchProtectionDefId, evaluateBranchProtectionDefVersion)

	if branchControlCheck.Settings == nil {
		return fmt.Errorf("Settings nil")
	}

	var inputs map[string]interface{}

	if inputMap, found := branchControlCheck.Settings.(map[string]interface{})["inputs"]; found {
		inputs = inputMap.(map[string]interface{})
	} else {
		return fmt.Errorf("inputs not found")
	}

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

	return nil
}

func expandBranchControlCheck(d *schema.ResourceData) (*pipelineschecks.CheckConfiguration, string, error) {
	inputs := map[string]interface{}{
		"allowedBranches":          d.Get("allowed_branches").(string),
		"ensureProtectionOfBranch": strconv.FormatBool(d.Get("verify_branch_protection").(bool)),
		"allowUnknownStatusBranch": strconv.FormatBool(d.Get("ignore_unknown_protection_status").(bool)),
	}

	return doBaseExpansion(d, inputs, evaluateBranchProtectionDef)
}
