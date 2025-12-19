package workitemtrackingprocess

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceControl() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceControl,
		ReadContext:   readResourceControl,
		UpdateContext: updateResourceControl,
		DeleteContext: deleteResourceControl,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: importResourceControl,
		},
		Schema: map[string]*schema.Schema{
			"process_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The ID of the process.",
			},
			"work_item_type_reference_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The reference name of the work item type.",
			},
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID of the group to add the control to.",
			},
			"control_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID for the control. For field controls, this is the field reference name.",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Label for the field.",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Order in which the control should appear in its group.",
			},
			"visible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A value indicating if the control should be hidden or not.",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A value indicating if the control is readonly.",
			},
			"metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Inner text of the control.",
			},
			"watermark": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Watermark text for the textbox.",
			},
			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Height of the control, for html controls.",
			},
			"control_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the control.",
			},
			"inherited": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A value indicating whether this layout node has been inherited from a parent layout.",
			},
			"overridden": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A value indicating whether this layout node has been overridden by a child layout.",
			},
			"is_contribution": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A value indicating if the layout node is contribution or not.",
			},
			"contribution": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Contribution for the control.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contribution_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
							Description:      "The id for the contribution.",
						},
						"height": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The height for the contribution.",
						},
						"inputs": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A dictionary holding key value pairs for contribution inputs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"show_on_deleted_work_item": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "A value indicating if the contribution should be shown on deleted work item.",
						},
					},
				},
			},
		},
	}
}

func importResourceControl(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	// Import ID format: process_id/work_item_type_reference_name/group_id/control_id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid import ID format, expected: process_id/work_item_type_reference_name/group_id/control_id")
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_reference_name", parts[1])
	d.Set("group_id", parts[2])
	d.Set("control_id", parts[3])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}

func createResourceControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	control := workitemtrackingprocess.Control{
		Id:       converter.String(d.Get("control_id").(string)),
		Visible:  converter.Bool(d.Get("visible").(bool)),
		ReadOnly: converter.Bool(d.Get("read_only").(bool)),
	}

	if v, ok := d.GetOk("label"); ok {
		control.Label = converter.String(v.(string))
	}
	//nolint:staticcheck // SA1019: d.GetOkExists is deprecated but required to distinguish between unset and zero value
	if v, ok := d.GetOkExists("order"); ok {
		control.Order = converter.Int(v.(int))
	}
	if v, ok := d.GetOk("metadata"); ok {
		control.Metadata = converter.String(v.(string))
	}
	if v, ok := d.GetOk("watermark"); ok {
		control.Watermark = converter.String(v.(string))
	}
	//nolint:staticcheck // SA1019: d.GetOkExists is deprecated but required to distinguish between unset and zero value
	if v, ok := d.GetOkExists("height"); ok {
		control.Height = converter.Int(v.(int))
	}
	control.IsContribution = converter.Bool(d.Get("is_contribution").(bool))
	if v, ok := d.GetOk("contribution"); ok {
		control.Contribution = expandContribution(v.([]interface{}))
	}

	args := workitemtrackingprocess.CreateControlInGroupArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		GroupId:    converter.String(d.Get("group_id").(string)),
		Control:    &control,
	}

	createdControl, err := clients.WorkItemTrackingProcessClient.CreateControlInGroup(ctx, args)
	if err != nil {
		return diag.Errorf(" Creating control. Error %+v", err)
	}

	if createdControl.Id == nil {
		return diag.Errorf(" Created control has no ID")
	}

	d.SetId(*createdControl.Id)
	return readResourceControl(ctx, d, m)
}

func readResourceControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	// Get the work item type with layout expanded
	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		Expand:     &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting work item type with reference name: %s for process with id %s. Error: %+v", witRefName, processId, err)
	}

	foundControl, foundGroupId := findControlById(workItemType.Layout, controlId)
	if foundControl == nil {
		d.SetId("")
		return nil
	}

	d.Set("group_id", foundGroupId)
	d.Set("label", foundControl.Label)
	d.Set("order", foundControl.Order)
	d.Set("visible", foundControl.Visible)
	d.Set("read_only", foundControl.ReadOnly)
	d.Set("metadata", foundControl.Metadata)
	d.Set("watermark", foundControl.Watermark)
	d.Set("height", foundControl.Height)
	d.Set("control_type", foundControl.ControlType)
	d.Set("inherited", foundControl.Inherited)
	d.Set("overridden", foundControl.Overridden)
	d.Set("is_contribution", foundControl.IsContribution)

	if foundControl.Contribution != nil {
		d.Set("contribution", flattenContribution(foundControl.Contribution))
	}

	return nil
}

func updateResourceControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)
	groupId := d.Get("group_id").(string)

	control := &workitemtrackingprocess.Control{
		Visible:  converter.Bool(d.Get("visible").(bool)),
		ReadOnly: converter.Bool(d.Get("read_only").(bool)),
	}

	if v, ok := d.GetOk("label"); ok {
		control.Label = converter.String(v.(string))
	}
	//nolint:staticcheck // SA1019: d.GetOkExists is deprecated but required to distinguish between unset and zero value
	if v, ok := d.GetOkExists("order"); ok {
		control.Order = converter.Int(v.(int))
	}
	if v, ok := d.GetOk("metadata"); ok {
		control.Metadata = converter.String(v.(string))
	}
	if v, ok := d.GetOk("watermark"); ok {
		control.Watermark = converter.String(v.(string))
	}
	//nolint:staticcheck // SA1019: d.GetOkExists is deprecated but required to distinguish between unset and zero value
	if v, ok := d.GetOkExists("height"); ok {
		control.Height = converter.Int(v.(int))
	}
	control.IsContribution = converter.Bool(d.Get("is_contribution").(bool))
	if v, ok := d.GetOk("contribution"); ok {
		control.Contribution = expandContribution(v.([]interface{}))
	}

	// Check if group_id has changed - if so, move the control
	if d.HasChange("group_id") {
		oldGroupId, _ := d.GetChange("group_id")
		moveArgs := workitemtrackingprocess.MoveControlToGroupArgs{
			ProcessId:       converter.UUID(processId),
			WitRefName:      converter.String(witRefName),
			GroupId:         converter.String(groupId),
			ControlId:       &controlId,
			Control:         control,
			RemoveFromGroupId: converter.String(oldGroupId.(string)),
		}

		_, err := clients.WorkItemTrackingProcessClient.MoveControlToGroup(ctx, moveArgs)
		if err != nil {
			return diag.Errorf(" Moving control. Error %+v", err)
		}

		return readResourceControl(ctx, d, m)
	}

	args := workitemtrackingprocess.UpdateControlArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: converter.String(witRefName),
		GroupId:    converter.String(groupId),
		ControlId:  &controlId,
		Control:    control,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdateControl(ctx, args)
	if err != nil {
		return diag.Errorf(" Update control. Error %+v", err)
	}

	return readResourceControl(ctx, d, m)
}

func deleteResourceControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()

	args := workitemtrackingprocess.RemoveControlFromGroupArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		GroupId:    converter.String(d.Get("group_id").(string)),
		ControlId:  &controlId,
	}

	err := clients.WorkItemTrackingProcessClient.RemoveControlFromGroup(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return diag.Errorf(" Delete control. Error %+v", err)
	}

	return nil
}

func findControlById(layout *workitemtrackingprocess.FormLayout, controlId string) (*workitemtrackingprocess.Control, string) {
	if layout == nil {
		return nil, ""
	}
	pages := layout.Pages
	if pages == nil {
		return nil, ""
	}
	for _, page := range *pages {
		control, groupId := findControlInPage(&page, controlId)
		if control != nil {
			return control, groupId
		}
	}
	return nil, ""
}

func findControlInPage(page *workitemtrackingprocess.Page, controlId string) (*workitemtrackingprocess.Control, string) {
	sections := page.Sections
	if sections == nil {
		return nil, ""
	}
	for _, section := range *sections {
		control, groupId := findControlInSection(&section, controlId)
		if control != nil {
			return control, groupId
		}
	}
	return nil, ""
}

func findControlInSection(section *workitemtrackingprocess.Section, controlId string) (*workitemtrackingprocess.Control, string) {
	groups := section.Groups
	if groups == nil {
		return nil, ""
	}
	for _, group := range *groups {
		control := findControlInGroup(&group, controlId)
		if control != nil {
			return control, *group.Id
		}
	}
	return nil, ""
}

func findControlInGroup(group *workitemtrackingprocess.Group, controlId string) *workitemtrackingprocess.Control {
	controls := group.Controls
	if controls == nil {
		return nil
	}
	for _, control := range *controls {
		id := control.Id
		if id == nil {
			continue
		}
		if *id == controlId {
			return &control
		}
	}
	return nil
}

func expandContribution(input []interface{}) *workitemtrackingprocess.WitContribution {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	contribution := &workitemtrackingprocess.WitContribution{}

	if v, ok := raw["contribution_id"].(string); ok && v != "" {
		contribution.ContributionId = &v
	}
	if v, ok := raw["height"].(int); ok && v != 0 {
		contribution.Height = &v
	}
	if v, ok := raw["show_on_deleted_work_item"].(bool); ok {
		contribution.ShowOnDeletedWorkItem = &v
	}
	if v, ok := raw["inputs"].(map[string]interface{}); ok && len(v) > 0 {
		inputs := make(map[string]interface{})
		for key, val := range v {
			inputs[key] = val
		}
		contribution.Inputs = &inputs
	}

	return contribution
}

func flattenContribution(contribution *workitemtrackingprocess.WitContribution) []interface{} {
	if contribution == nil {
		return nil
	}

	result := make(map[string]interface{})

	if contribution.ContributionId != nil {
		result["contribution_id"] = *contribution.ContributionId
	}
	if contribution.Height != nil {
		result["height"] = *contribution.Height
	}
	if contribution.ShowOnDeletedWorkItem != nil {
		result["show_on_deleted_work_item"] = *contribution.ShowOnDeletedWorkItem
	}
	if contribution.Inputs != nil {
		inputs := make(map[string]string)
		for key, val := range *contribution.Inputs {
			if strVal, ok := val.(string); ok {
				inputs[key] = strVal
			}
		}
		result["inputs"] = inputs
	}

	return []interface{}{result}
}
