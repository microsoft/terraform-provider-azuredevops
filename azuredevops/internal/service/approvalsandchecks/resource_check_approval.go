package approvalsandchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk/pipelineschecksextras"
)

const (
	approvalKindApproval  = "approval"
	approvalKindPreCheck  = "pre_check"
	approvalKindPostCheck = "post_check"

	approvalDefinitionRefID  = "26014962-64a0-49f4-885b-4b874119a5cc"
	preCheckDefinitionRefID  = "0f52a19b-c67e-468f-b8eb-0ae83b532c99"
	postCheckDefinitionRefID = "06441319-13fb-4756-b198-c2da116894a4"
)

// ResourceCheckApproval schema and implementation for branch check resources
func ResourceCheckApproval() *schema.Resource {
	r := genBaseCheckResource(flattenCheckApproval, expandCheckApproval)

	r.Schema["approvers"] = &schema.Schema{
		Type:     schema.TypeList,
		MinItems: 1,
		Required: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.IsUUID,
		},
	}

	r.Schema["instructions"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	r.Schema["requester_can_approve"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	r.Schema["minimum_required_approvers"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}

	r.Schema["approval_kind"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  approvalKindApproval,
		ValidateFunc: validation.StringInSlice([]string{
			approvalKindApproval,
			approvalKindPreCheck,
			approvalKindPostCheck,
		}, false),
	}

	r.Schema["timeout"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      43200,
		ValidateFunc: validation.IntBetween(1, 43200),
	}

	return r
}

func flattenCheckApproval(d *schema.ResourceData, check *pipelineschecksextras.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, check, projectID)
	if err != nil {
		return err
	}

	if check.Settings == nil {
		return fmt.Errorf("settings nil")
	}

	check.Type = approvalAndCheckType.Approval

	settings := check.Settings.(map[string]interface{})

	if instructions, found := settings["instructions"]; found {
		d.Set("instructions", instructions)
	}

	if minRequiredApprovers, found := settings["minRequiredApprovers"]; found {
		d.Set("minimum_required_approvers", minRequiredApprovers)
	}

	if definitionRefRaw, found := settings["definitionRef"]; found {
		definitionRef, ok := definitionRefRaw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("definitionRef has unexpected type %T", definitionRefRaw)
		}

		definitionRefIDRaw, found := definitionRef["id"]
		if !found {
			return fmt.Errorf("definitionRef.id not found")
		}

		definitionRefID, ok := definitionRefIDRaw.(string)
		if !ok {
			return fmt.Errorf("definitionRef.id has unexpected type %T", definitionRefIDRaw)
		}

		approvalKind, err := approvalKindFromDefinitionRefID(definitionRefID)
		if err != nil {
			return err
		}

		d.Set("approval_kind", approvalKind)
	}

	if requesterCannotBeApprover, found := settings["requesterCannotBeApprover"]; found {
		d.Set("requester_can_approve", !requesterCannotBeApprover.(bool))
	}

	if approversRaw, found := settings["approvers"]; found {
		approverIds := make([]string, 0)
		for _, approverRaw := range approversRaw.([]interface{}) {
			approver := approverRaw.(map[string]interface{})
			approverId := approver["id"].(string)

			approverIds = append(approverIds, approverId)
		}

		d.Set("approvers", approverIds)
	} else {
		return fmt.Errorf("approvers input not found")
	}

	if check.Timeout != nil {
		d.Set("timeout", *check.Timeout)
	}

	return nil
}

func expandCheckApproval(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	approvers := make([]interface{}, 0)

	if userApproversRaw, ok := d.GetOk("approvers"); ok {
		userApprovers := userApproversRaw.([]interface{})
		for _, user := range userApprovers {
			approvers = append(approvers, map[string]interface{}{
				"id": user.(string),
			})
		}
	}

	approvalKind := d.Get("approval_kind").(string)
	definitionRefID, err := definitionRefIDFromApprovalKind(approvalKind)
	if err != nil {
		return nil, "", err
	}

	settings := map[string]interface{}{
		"instructions":              d.Get("instructions").(string),
		"minRequiredApprovers":      d.Get("minimum_required_approvers").(int),
		"requesterCannotBeApprover": !d.Get("requester_can_approve").(bool),
		"approvers":                 approvers,
		"definitionRef": map[string]interface{}{
			"id": definitionRefID,
		},
	}

	return doBaseExpansion(d, approvalAndCheckType.Approval, settings, converter.ToPtr(d.Get("timeout").(int)))
}

func definitionRefIDFromApprovalKind(approvalKind string) (string, error) {
	switch approvalKind {
	case approvalKindApproval:
		return approvalDefinitionRefID, nil
	case approvalKindPreCheck:
		return preCheckDefinitionRefID, nil
	case approvalKindPostCheck:
		return postCheckDefinitionRefID, nil
	default:
		return "", fmt.Errorf("unsupported approval_kind %q", approvalKind)
	}
}

func approvalKindFromDefinitionRefID(definitionRefID string) (string, error) {
	switch definitionRefID {
	case approvalDefinitionRefID:
		return approvalKindApproval, nil
	case preCheckDefinitionRefID:
		return approvalKindPreCheck, nil
	case postCheckDefinitionRefID:
		return approvalKindPostCheck, nil
	default:
		return "", fmt.Errorf("unsupported approval definitionRef.id %q", definitionRefID)
	}
}
