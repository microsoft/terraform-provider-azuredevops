package approvalsandchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
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

	r.Schema["timeout"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
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

	settings := map[string]interface{}{
		"instructions":              d.Get("instructions").(string),
		"minRequiredApprovers":      d.Get("minimum_required_approvers").(int),
		"requesterCannotBeApprover": !d.Get("requester_can_approve").(bool),
		"approvers":                 approvers,
	}

	timeout := 43200 // 12 hour default
	if val, ok := d.GetOk("timeout"); ok {
		timeout = val.(int)
	}

	return doBaseExpansion(d, approvalAndCheckType.Approval, settings, converter.ToPtr(timeout))
}
