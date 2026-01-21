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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceInheritedControl() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceInheritedControl,
		ReadContext:   readResourceInheritedControl,
		UpdateContext: updateResourceInheritedControl,
		DeleteContext: deleteResourceInheritedControl,
		Importer: &schema.ResourceImporter{
			StateContext: importResourceInheritedControl,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID of the group containing the control.",
			},
			"control_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID of the inherited control to customize.",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Label for the control.",
			},
			"visible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the control should be visible.",
			},
		},
	}
}

func importResourceInheritedControl(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
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

func createResourceInheritedControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Get("control_id").(string)
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)
	groupId := d.Get("group_id").(string)

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		// Returns the layout containing the controls
		Expand: &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		return diag.Errorf("getting work item type: %+v", err)
	}
	if workItemType == nil || workItemType.Layout == nil {
		return diag.Errorf("work item type or layout is nil")
	}

	group := findGroupById(workItemType.Layout, groupId)
	if group == nil {
		return diag.Errorf("group %s not found in layout", groupId)
	}

	existingControl := findControlInGroup(group, controlId)
	if existingControl == nil {
		return diag.Errorf("control %s not found in group %s", controlId, groupId)
	}

	if existingControl.Inherited == nil || !*existingControl.Inherited {
		return diag.Errorf("control %s is not inherited, use azuredevops_workitemtrackingprocess_control to manage custom controls", controlId)
	}

	d.SetId(controlId)

	return updateResourceInheritedControl(ctx, d, m)
}

func readResourceInheritedControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)
	groupId := d.Get("group_id").(string)

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		// Returns the layout containing the controls
		Expand: &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		return diag.Errorf("getting work item type: %+v", err)
	}
	if workItemType == nil || workItemType.Layout == nil {
		return diag.Errorf("work item type or layout is nil")
	}

	group := findGroupById(workItemType.Layout, groupId)
	if group == nil {
		d.SetId("")
		return nil
	}

	control := findControlInGroup(group, controlId)
	if control == nil {
		d.SetId("")
		return nil
	}

	if control.Label != nil {
		d.Set("label", *control.Label)
	}
	if control.Visible != nil {
		d.Set("visible", *control.Visible)
	}
	return nil
}

func updateResourceInheritedControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()

	control := workitemtrackingprocess.Control{}

	rawConfig := d.GetRawConfig().AsValueMap()
	if visible := rawConfig["visible"]; !visible.IsNull() {
		control.Visible = converter.Bool(visible.True())
	}

	if v, ok := d.GetOk("label"); ok {
		control.Label = converter.String(v.(string))
	}

	args := workitemtrackingprocess.UpdateControlArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		GroupId:    converter.String(d.Get("group_id").(string)),
		ControlId:  &controlId,
		Control:    &control,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdateControl(ctx, args)
	if err != nil {
		return diag.Errorf("updating inherited control: %+v", err)
	}

	return readResourceInheritedControl(ctx, d, m)
}

func deleteResourceInheritedControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
		return diag.Errorf("reverting inherited control: %+v", err)
	}

	return nil
}
