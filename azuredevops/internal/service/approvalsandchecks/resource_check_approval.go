package approvalsandchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
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

	r.Schema["requestor_can_approve"] = &schema.Schema{
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

func flattenCheckApproval(d *schema.ResourceData, check *pipelineschecks.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, check, projectID)
	if err != nil {
		return err
	}

	if check.Settings == nil {
		return fmt.Errorf("settings nil")
	}

	check.Type.Id = converter.UUID("8C6F20A7-A545-4486-9777-F762FAFE0D4D")
	check.Type.Name = converter.ToPtr("Approval")

	settings := check.Settings.(map[string]interface{})

	if instructions, found := settings["instructions"]; found {
		d.Set("instructions", instructions)
	} else {
		return fmt.Errorf("instructions not found")
	}

	if minRequiredApprovers, found := settings["minRequiredApprovers"]; found {
		d.Set("minimum_required_approvers", minRequiredApprovers)
	} else {
		return fmt.Errorf("minRequiredApprovers not found")
	}

	if requesterCannotBeApprover, found := settings["requesterCannotBeApprover"]; found {
		d.Set("requestor_can_approve", !requesterCannotBeApprover.(bool))
	} else {
		d.Set("requestor_can_approve", true)
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

func expandCheckApproval(d *schema.ResourceData) (*pipelineschecks.CheckConfiguration, string, error) {
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
		"requesterCannotBeApprover": !d.Get("requestor_can_approve").(bool),
		"approvers":                 approvers,
	}

	checkType := &pipelineschecks.CheckType{
		Id:   converter.UUID("8C6F20A7-A545-4486-9777-F762FAFE0D4D"),
		Name: converter.ToPtr("Approval"),
	}

	return doBaseExpansion(d, checkType, settings, converter.ToPtr(d.Get("timeout").(int)))
}
