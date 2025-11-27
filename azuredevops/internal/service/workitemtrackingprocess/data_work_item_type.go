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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataWorkItemType() *schema.Resource {
	workItemTypeSchema := getWorkItemTypeSchema()

	// Add the required input fields that are specific to the single work item type data source
	workItemTypeSchema["process_id"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
		Description:      "The ID of the process.",
	}
	workItemTypeSchema["reference_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The reference name of the work item type.",
	}

	return &schema.Resource{
		ReadContext: readDataWorkItemType,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: workItemTypeSchema,
	}
}

func getWorkItemTypeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"reference_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Reference name of the work item type.",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the work item type.",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description of the work item type.",
		},
		"color": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Color hexadecimal code to represent the work item type.",
		},
		"icon": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Icon to represent the work item type.",
		},
		"is_disabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates if the work item type is disabled.",
		},
		"inherits_from": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Parent work item type reference name.",
		},
		"customization": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Indicates the type of customization on this work item type.",
		},
		"url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "URL of the work item type.",
		},
	}
}

func workItemTypeToMap(workItemType *workitemtrackingprocess.ProcessWorkItemType) map[string]any {
	return map[string]any{
		"reference_name": workItemType.ReferenceName,
		"name":           workItemType.Name,
		"description":    workItemType.Description,
		"color":          convertColorToResource(workItemType),
		"icon":           workItemType.Icon,
		"is_disabled":    workItemType.IsDisabled,
		"inherits_from":  workItemType.Inherits,
		"customization":  string(*workItemType.Customization),
		"url":            workItemType.Url,
	}
}

func readDataWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	processId := d.Get("process_id").(string)
	referenceName := d.Get("reference_name").(string)

	expand := workitemtrackingprocess.GetWorkItemTypeExpandValues.None

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &referenceName,
		Expand:     &expand,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting work item type with reference name: %s for process with id %s. Error: %+v", referenceName, processId, err)
	}

	workItemTypeMap := workItemTypeToMap(workItemType)
	for key, value := range workItemTypeMap {
		d.Set(key, value)
	}

	d.SetId(referenceName)

	return nil
}
