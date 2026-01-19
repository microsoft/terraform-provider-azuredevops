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

func ResourceInheritedState() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceInheritedState,
		ReadContext:   readResourceInheritedState,
		UpdateContext: updateResourceInheritedState,
		DeleteContext: deleteResourceInheritedState,
		Importer: &schema.ResourceImporter{
			StateContext: importResourceInheritedState,
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
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Name of the inherited state to manage.",
			},
			"hidden": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the state is hidden.",
			},
		},
	}
}

func importResourceInheritedState(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	// Import ID format: process_id/work_item_type_reference_name/state_name
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format, expected: process_id/work_item_type_reference_name/state_name")
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_reference_name", parts[1])
	d.Set("name", parts[2])

	// We need to look up the state by name to get its ID
	clients := m.(*client.AggregatedClient)
	state, err := findInheritedStateByName(ctx, clients, parts[0], parts[1], parts[2])
	if err != nil {
		return nil, err
	}
	if state.Id == nil {
		return nil, fmt.Errorf("state ID is nil")
	}

	d.SetId(state.Id.String())

	return []*schema.ResourceData{d}, nil
}

func createResourceInheritedState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)
	name := d.Get("name").(string)

	state, err := findInheritedStateByName(ctx, clients, processId, witRefName, name)
	if err != nil {
		return diag.FromErr(err)
	}
	if state.Id == nil {
		return diag.Errorf("state ID is nil")
	}

	d.SetId(state.Id.String())

	return updateResourceInheritedState(ctx, d, m)
}

func readResourceInheritedState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stateId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	state, err := clients.WorkItemTrackingProcessClient.GetStateDefinition(ctx, workitemtrackingprocess.GetStateDefinitionArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		StateId:    converter.UUID(stateId),
	})
	if err != nil {
		return diag.Errorf("getting state definition: %+v", err)
	}
	if state == nil {
		d.SetId("")
		return nil
	}

	d.Set("hidden", state.Hidden)

	return nil
}

func updateResourceInheritedState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stateId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	rawConfig := d.GetRawConfig().AsValueMap()
	if hiddenRaw := rawConfig["hidden"]; !hiddenRaw.IsNull() {
		hidden := hiddenRaw.True()
		if hidden {
			_, err := clients.WorkItemTrackingProcessClient.HideStateDefinition(ctx, workitemtrackingprocess.HideStateDefinitionArgs{
				ProcessId:      converter.UUID(processId),
				WitRefName:     &witRefName,
				StateId:        converter.UUID(stateId),
				HideStateModel: &workitemtrackingprocess.HideStateModel{Hidden: &hidden},
			})
			if err != nil {
				return diag.Errorf("hiding state: %+v", err)
			}
		} else {
			// Use DELETE to unhide a state
			err := clients.WorkItemTrackingProcessClient.DeleteStateDefinition(ctx, workitemtrackingprocess.DeleteStateDefinitionArgs{
				ProcessId:  converter.UUID(processId),
				WitRefName: &witRefName,
				StateId:    converter.UUID(stateId),
			})
			if err != nil {
				return diag.Errorf("unhiding state: %+v", err)
			}
		}
	}

	return readResourceInheritedState(ctx, d, m)
}

func deleteResourceInheritedState(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stateId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	err := clients.WorkItemTrackingProcessClient.DeleteStateDefinition(ctx, workitemtrackingprocess.DeleteStateDefinitionArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		StateId:    converter.UUID(stateId),
	})
	if err != nil {
		return diag.Errorf("deleting state: %+v", err)
	}

	return nil
}

func findInheritedStateByName(ctx context.Context, clients *client.AggregatedClient, processId string, witRefName string, name string) (*workitemtrackingprocess.WorkItemStateResultModel, error) {
	states, err := clients.WorkItemTrackingProcessClient.GetStateDefinitions(ctx, workitemtrackingprocess.GetStateDefinitionsArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
	})
	if err != nil {
		return nil, err
	}

	if states != nil {
		for _, state := range *states {
			if state.Name != nil && *state.Name == name {
				if state.CustomizationType == nil {
					return nil, fmt.Errorf("state %q has no customization type", name)
				}
				if *state.CustomizationType == workitemtrackingprocess.CustomizationTypeValues.Custom {
					return nil, fmt.Errorf("state %q is a custom state, use azuredevops_workitemtrackingprocess_state for custom states", name)
				}
				return &state, nil
			}
		}
	}

	return nil, fmt.Errorf("inherited state %q not found", name)
}
