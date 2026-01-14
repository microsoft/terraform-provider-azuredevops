package workitemtrackingprocess

import (
	"context"
	"fmt"
	"regexp"
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

func ResourceState() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceState,
		ReadContext:   readResourceState,
		UpdateContext: updateResourceState,
		DeleteContext: deleteResourceState,
		Importer: &schema.ResourceImporter{
			StateContext: importResourceState,
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
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Name of the state.",
			},
			"color": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^#[0-9a-fA-F]{6}$`), "Must be a hexadecimal color code, i.e. #b2b2b2")),
				Description:      "Color hexadecimal code to represent the state.",
			},
			"state_category": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Proposed", "In Progress", "Resolved", "Completed", "Removed"}, false)),
				Description:      "Category of the state. Valid values: Proposed, In Progress, Resolved, Completed, Removed.",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Order in which the state should appear.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the state.",
			},
		},
	}
}

func importResourceState(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	// Import ID format: process_id/work_item_type_reference_name/state_id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format, expected: process_id/work_item_type_reference_name/state_id")
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_reference_name", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func createResourceState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stateModel := workitemtrackingprocess.WorkItemStateInputModel{
		Name:          converter.String(d.Get("name").(string)),
		Color:         convertColorToApi(d),
		StateCategory: converter.String(d.Get("state_category").(string)),
	}

	if v, ok := d.GetOk("order"); ok {
		stateModel.Order = converter.Int(v.(int))
	}

	args := workitemtrackingprocess.CreateStateDefinitionArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		StateModel: &stateModel,
	}

	createdState, err := clients.WorkItemTrackingProcessClient.CreateStateDefinition(ctx, args)
	if err != nil {
		return diag.Errorf("creating state: %+v", err)
	}
	if createdState == nil {
		return diag.Errorf("created state is nil")
	}
	if createdState.Id == nil {
		return diag.Errorf("created state has no ID")
	}

	d.SetId(createdState.Id.String())

	return readResourceState(ctx, d, m)
}

func readResourceState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stateId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	args := workitemtrackingprocess.GetStateDefinitionArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: converter.String(witRefName),
		StateId:    converter.UUID(stateId),
	}

	state, err := clients.WorkItemTrackingProcessClient.GetStateDefinition(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("getting state with id %s: %+v", stateId, err)
	}
	if state == nil {
		d.SetId("")
		return nil
	}

	return flattenState(d, state)
}

func updateResourceState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stateId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

if d.HasChanges("name", "color", "order") {
		stateModel := workitemtrackingprocess.WorkItemStateInputModel{}

		if d.HasChange("name") {
			stateModel.Name = converter.String(d.Get("name").(string))
		}
		if d.HasChange("color") {
			stateModel.Color = convertColorToApi(d)
		}
		if d.HasChange("order") {
			stateModel.Order = converter.Int(d.Get("order").(int))
		}

		args := workitemtrackingprocess.UpdateStateDefinitionArgs{
			ProcessId:  converter.UUID(processId),
			WitRefName: converter.String(witRefName),
			StateId:    converter.UUID(stateId),
			StateModel: &stateModel,
		}

		_, err := clients.WorkItemTrackingProcessClient.UpdateStateDefinition(ctx, args)
		if err != nil {
			return diag.Errorf("updating state: %+v", err)
		}
	}

	return readResourceState(ctx, d, m)
}

func deleteResourceState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stateId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	args := workitemtrackingprocess.DeleteStateDefinitionArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: converter.String(witRefName),
		StateId:    converter.UUID(stateId),
	}

	err := clients.WorkItemTrackingProcessClient.DeleteStateDefinition(ctx, args)
	if err != nil {
		return diag.Errorf("deleting state: %+v", err)
	}

	return nil
}

func flattenState(d *schema.ResourceData, state *workitemtrackingprocess.WorkItemStateResultModel) diag.Diagnostics {
	if state.Name != nil {
		d.Set("name", *state.Name)
	}
	if state.Color != nil {
		d.Set("color", convertColorToResource(*state.Color))
	}
	if state.StateCategory != nil {
		d.Set("state_category", *state.StateCategory)
	}
	if state.Order != nil {
		d.Set("order", *state.Order)
	}
	if state.Url != nil {
		d.Set("url", *state.Url)
	}

	return nil
}
