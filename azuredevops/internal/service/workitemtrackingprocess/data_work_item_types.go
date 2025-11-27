package workitemtrackingprocess

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func DataWorkItemTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: readWorkItemTypes,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"process_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The ID of the process.",
			},
			"work_item_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A list of work item types for the process.",
				Set:         getWorkItemTypeHash,
				Elem: &schema.Resource{
					Schema: getWorkItemTypeSchema(),
				},
			},
		},
	}
}

func readWorkItemTypes(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	processId := d.Get("process_id").(string)

	expand := workitemtrackingprocess.GetWorkItemTypeExpandValues.None

	getWorkItemTypesArgs := workitemtrackingprocess.GetProcessWorkItemTypesArgs{
		ProcessId: converter.UUID(processId),
		Expand:    &expand,
	}
	retrievedWorkItemTypes, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemTypes(ctx, getWorkItemTypesArgs)
	if err != nil {
		return diag.Errorf(" Getting list of work item types for process %s: Error: %+v", processId, err)
	}

	workItemTypes := make([]any, 0)
	for _, retrievedWorkItemType := range *retrievedWorkItemTypes {
		workItemType := map[string]any{
			"reference_name": *retrievedWorkItemType.ReferenceName,
			"name":           *retrievedWorkItemType.Name,
			"description":    *retrievedWorkItemType.Description,
			"color":          fmt.Sprintf("#%s", *retrievedWorkItemType.Color),
			"icon":           *retrievedWorkItemType.Icon,
			"is_disabled":    *retrievedWorkItemType.IsDisabled,
			"inherits_from":  retrievedWorkItemType.Inherits,
			"customization":  string(*retrievedWorkItemType.Customization),
			"url":            *retrievedWorkItemType.Url,
		}
		workItemTypes = append(workItemTypes, workItemType)
	}

	d.SetId(processId)

	err = d.Set("work_item_types", workItemTypes)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getWorkItemTypeHash(v any) int {
	return tfhelper.HashString(v.(map[string]any)["reference_name"].(string))
}
