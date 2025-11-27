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

func ResourceWorkItemType() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceWorkItemType,
		ReadContext:   readResourceWorkItemType,
		UpdateContext: updateResourceWorkItemType,
		DeleteContext: deleteResourceWorkItemType,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				id := d.Id()
				parts := strings.SplitN(id, "/", 2)
				if len(parts) != 2 || strings.EqualFold(parts[0], "") || strings.EqualFold(parts[1], "") {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected process_id/reference_name", id)
				}
				d.Set("process_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
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
				Description:      "The ID of the process the work item type belongs to.",
			},
			"color": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Color hexadecimal code to represent the work item type.",
				Default:          "#009ccc",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^#[0-9a-fA-F]{6}$`), "Must be a hexadecimal color code, i.e. #009ccc")),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the work item type.",
			},
			"icon": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Icon to represent the work item type.",
				Default:     "icon_clipboard",
			},
			"inherits_from": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Parent work item type for work item type.",
			},
			"is_disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True if the work item type need to be disabled.",
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Name of work item type.",
			},
			"reference_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference name of the work item type.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Url of the work item type.",
			},
		},
	}
}

func createResourceWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	workItemTypeRequest := workitemtrackingprocess.CreateProcessWorkItemTypeRequest{
		Name:       converter.String(d.Get("name").(string)),
		IsDisabled: converter.Bool(d.Get("is_disabled").(bool)),
		Color:      convertColorToApi(d),
		Icon:       converter.String(d.Get("icon").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		workItemTypeRequest.Description = converter.String(v.(string))
	}
	if v, ok := d.GetOk("inherits_from"); ok {
		workItemTypeRequest.InheritsFrom = converter.String(v.(string))
	}

	args := workitemtrackingprocess.CreateProcessWorkItemTypeArgs{
		ProcessId:    converter.UUID(d.Get("process_id").(string)),
		WorkItemType: &workItemTypeRequest,
	}

	createdWorkItemType, err := clients.WorkItemTrackingProcessClient.CreateProcessWorkItemType(ctx, args)
	if err != nil {
		return diag.Errorf(" Creating work item type. Error %+v", err)
	}
	d.SetId(*createdWorkItemType.ReferenceName)
	return setWorkItemType(d, createdWorkItemType)
}

func readResourceWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	referenceName := d.Id()
	processId := d.Get("process_id").(string)

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &referenceName,
		Expand:     &workitemtrackingprocess.GetWorkItemTypeExpandValues.None,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting work item type with reference name: %s for process with id %s. Error: %+v", referenceName, processId, err)
	}

	return setWorkItemType(d, workItemType)
}

func updateResourceWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	referenceName := d.Id()
	processId := d.Get("process_id").(string)

	updateWorkItemType := &workitemtrackingprocess.UpdateProcessWorkItemTypeRequest{
		IsDisabled: converter.Bool(d.Get("is_disabled").(bool)),
		Color:      convertColorToApi(d),
		Icon:       converter.String(d.Get("icon").(string)),
	}
	if v, ok := d.GetOk("description"); ok {
		updateWorkItemType.Description = converter.String(v.(string))
	}

	args := workitemtrackingprocess.UpdateProcessWorkItemTypeArgs{
		ProcessId:          converter.UUID(processId),
		WitRefName:         &referenceName,
		WorkItemTypeUpdate: updateWorkItemType,
	}

	updatedWorkItemType, err := clients.WorkItemTrackingProcessClient.UpdateProcessWorkItemType(ctx, args)
	if err != nil {
		return diag.Errorf(" Update work item type. Error %+v", err)
	}

	// Note! There is a bug in the PATCH endpoint where the response has icon always set to null. POST and GET doesn't seem to have this issue.
	updatedWorkItemType.Icon = updateWorkItemType.Icon

	return setWorkItemType(d, updatedWorkItemType)
}

func deleteResourceWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	referenceName := d.Id()
	processId := d.Get("process_id").(string)

	args := workitemtrackingprocess.DeleteProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &referenceName,
	}

	err := clients.WorkItemTrackingProcessClient.DeleteProcessWorkItemType(ctx, args)
	if err != nil {
		return diag.Errorf(" Delete work item type. Error %+v", err)
	}

	return nil
}

func setWorkItemType(d *schema.ResourceData, workItemType *workitemtrackingprocess.ProcessWorkItemType) diag.Diagnostics {
	d.Set("name", workItemType.Name)
	d.Set("description", workItemType.Description)
	d.Set("is_disabled", workItemType.IsDisabled)
	d.Set("color", convertColorToResource(workItemType))
	d.Set("icon", workItemType.Icon)
	d.Set("inherits_from", workItemType.Inherits)
	d.Set("reference_name", workItemType.ReferenceName)
	d.Set("url", workItemType.Url)
	return nil
}

func convertColorToApi(d *schema.ResourceData) *string {
	return converter.String(
		strings.ReplaceAll(d.Get("color").(string), "#", ""))
}

func convertColorToResource(workItemType *workitemtrackingprocess.ProcessWorkItemType) string {
	return fmt.Sprintf("#%s", *workItemType.Color)
}
