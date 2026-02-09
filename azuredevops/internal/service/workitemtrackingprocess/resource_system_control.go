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

func ResourceSystemControl() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceSystemControl,
		ReadContext:   readResourceSystemControl,
		UpdateContext: updateResourceSystemControl,
		DeleteContext: deleteResourceSystemControl,
		Importer: &schema.ResourceImporter{
			StateContext: importResourceSystemControl,
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
			"work_item_type_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID (reference name) of the work item type.",
			},
			"control_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID of the system control (e.g., System.AreaPath, System.IterationPath, System.Reason).",
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
				Default:     true,
				Description: "Whether the control should be visible.",
			},
			"control_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the control.",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the control is read-only.",
			},
		},
	}
}

func importResourceSystemControl(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	// Import ID format: process_id/work_item_type_id/control_id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format, expected: process_id/work_item_type_id/control_id")
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_id", parts[1])
	d.Set("control_id", parts[2])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func createResourceSystemControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Get("control_id").(string)

	control := workitemtrackingprocess.Control{
		Visible: converter.Bool(d.Get("visible").(bool)),
	}

	if v, ok := d.GetOk("label"); ok {
		control.Label = converter.String(v.(string))
	}

	args := workitemtrackingprocess.UpdateSystemControlArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_id").(string)),
		ControlId:  &controlId,
		Control:    &control,
	}

	updatedControl, err := clients.WorkItemTrackingProcessClient.UpdateSystemControl(ctx, args)
	if err != nil {
		return diag.Errorf("creating system control customization: %+v", err)
	}
	if updatedControl == nil {
		return diag.Errorf("updated system control is nil")
	}

	d.SetId(controlId)

	return readResourceSystemControl(ctx, d, m)
}

func readResourceSystemControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)

	args := workitemtrackingprocess.GetSystemControlsArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: converter.String(witRefName),
	}

	controls, err := clients.WorkItemTrackingProcessClient.GetSystemControls(ctx, args)
	if err != nil {
		return diag.Errorf("getting system controls: %+v", err)
	}

	// Find the specific control - GetSystemControls returns only edited controls
	var foundControl *workitemtrackingprocess.Control
	if controls != nil {
		for _, c := range *controls {
			if c.Id != nil && *c.Id == controlId {
				foundControl = &c
				break
			}
		}
	}

	if foundControl == nil {
		// Control not in the edited list means it's been reverted to default
		d.SetId("")
		return nil
	}

	if foundControl.Label != nil {
		d.Set("label", *foundControl.Label)
	}
	if foundControl.Visible != nil {
		d.Set("visible", *foundControl.Visible)
	}
	if foundControl.ControlType != nil {
		d.Set("control_type", *foundControl.ControlType)
	}
	if foundControl.ReadOnly != nil {
		d.Set("read_only", *foundControl.ReadOnly)
	}

	return nil
}

func updateResourceSystemControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()

	control := workitemtrackingprocess.Control{
		Visible: converter.Bool(d.Get("visible").(bool)),
	}

	if v, ok := d.GetOk("label"); ok {
		control.Label = converter.String(v.(string))
	}

	args := workitemtrackingprocess.UpdateSystemControlArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_id").(string)),
		ControlId:  &controlId,
		Control:    &control,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdateSystemControl(ctx, args)
	if err != nil {
		return diag.Errorf("updating system control: %+v", err)
	}

	return readResourceSystemControl(ctx, d, m)
}

func deleteResourceSystemControl(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	controlId := d.Id()

	args := workitemtrackingprocess.DeleteSystemControlArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_id").(string)),
		ControlId:  &controlId,
	}

	_, err := clients.WorkItemTrackingProcessClient.DeleteSystemControl(ctx, args)
	if err != nil {
		return diag.Errorf("deleting system control customization: %+v", err)
	}

	return nil
}
