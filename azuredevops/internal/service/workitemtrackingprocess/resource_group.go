package workitemtrackingprocess

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceGroup,
		ReadContext:   readResourceGroup,
		UpdateContext: updateResourceGroup,
		DeleteContext: deleteResourceGroup,
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
			"page_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID of the page to add the group to.",
			},
			"section_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID of the section to add the group to.",
			},
			"label": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Label for the group.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the group.",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Order in which the group should appear in the section.",
			},
			"visible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A value indicating if the group should be hidden or not.",
			},
		},
	}
}

func createResourceGroup(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	group := workitemtrackingprocess.Group{
		Visible: converter.Bool(d.Get("visible").(bool)),
	}

	if v, ok := d.GetOk("label"); ok {
		group.Label = converter.String(v.(string))
	}
	if v, ok := d.GetOk("id"); ok {
		group.Id = converter.String(v.(string))
	}
	if v, ok := d.GetOk("order"); ok {
		group.Order = converter.Int(v.(int))
	}

	args := workitemtrackingprocess.AddGroupArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		PageId:     converter.String(d.Get("page_id").(string)),
		SectionId:  converter.String(d.Get("section_id").(string)),
		Group:      &group,
	}

	createdGroup, err := clients.WorkItemTrackingProcessClient.AddGroup(ctx, args)
	if err != nil {
		return diag.Errorf(" Creating group. Error %+v", err)
	}

	if createdGroup.Id == nil {
		return diag.Errorf(" Created group has no ID")
	}

	d.SetId(*createdGroup.Id)
	return setWorkItemTypeGroup(d, createdGroup)
}

func readResourceGroup(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	groupId := d.Id()
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

	// Find the group by ID
	var foundGroup *workitemtrackingprocess.Group
	if workItemType.Layout != nil && workItemType.Layout.Pages != nil {
		for _, page := range *workItemType.Layout.Pages {
			if page.Sections != nil {
				for _, section := range *page.Sections {
					if section.Groups != nil {
						for _, group := range *section.Groups {
							if group.Id != nil && *group.Id == groupId {
								foundGroup = &group
								break
							}
						}
					}
					if foundGroup != nil {
						break
					}
				}
			}
			if foundGroup != nil {
				break
			}
		}
	}

	if foundGroup == nil {
		d.SetId("")
		return nil
	}

	return setWorkItemTypeGroup(d, foundGroup)
}

func updateResourceGroup(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	groupId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)
	pageId := d.Get("page_id").(string)
	sectionId := d.Get("section_id").(string)

	// Check if page_id or section_id has changed - if so, move the group first
	if d.HasChange("page_id") || d.HasChange("section_id") {
		oldPageId, _ := d.GetChange("page_id")
		oldSectionId, _ := d.GetChange("section_id")
		moveArgs := workitemtrackingprocess.MoveGroupToPageArgs{
			ProcessId:  converter.UUID(processId),
			WitRefName: converter.String(witRefName),
			PageId:     converter.String(pageId),
			SectionId:  converter.String(sectionId),
			GroupId:    &groupId,
			Group: &workitemtrackingprocess.Group{
				Id: &groupId,
			},
			RemoveFromPageId:    converter.String(oldPageId.(string)),
			RemoveFromSectionId: converter.String(oldSectionId.(string)),
		}

		_, err := clients.WorkItemTrackingProcessClient.MoveGroupToPage(ctx, moveArgs)
		if err != nil {
			return diag.Errorf(" Moving group. Error %+v", err)
		}
	}

	updateGroup := &workitemtrackingprocess.Group{
		Visible: converter.Bool(d.Get("visible").(bool)),
	}

	if v, ok := d.GetOk("label"); ok {
		updateGroup.Label = converter.String(v.(string))
	}
	if v, ok := d.GetOk("order"); ok {
		updateGroup.Order = converter.Int(v.(int))
	}

	args := workitemtrackingprocess.UpdateGroupArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: converter.String(witRefName),
		PageId:     converter.String(pageId),
		SectionId:  converter.String(sectionId),
		GroupId:    &groupId,
		Group:      updateGroup,
	}

	updatedGroup, err := clients.WorkItemTrackingProcessClient.UpdateGroup(ctx, args)
	if err != nil {
		return diag.Errorf(" Update group. Error %+v", err)
	}

	return setWorkItemTypeGroup(d, updatedGroup)
}

func deleteResourceGroup(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	groupId := d.Id()

	args := workitemtrackingprocess.RemoveGroupArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		PageId:     converter.String(d.Get("page_id").(string)),
		SectionId:  converter.String(d.Get("section_id").(string)),
		GroupId:    &groupId,
	}

	err := clients.WorkItemTrackingProcessClient.RemoveGroup(ctx, args)
	if err != nil {
		return diag.Errorf(" Delete group. Error %+v", err)
	}

	return nil
}

func setWorkItemTypeGroup(d *schema.ResourceData, group *workitemtrackingprocess.Group) diag.Diagnostics {
	if group.Id != nil {
		d.Set("id", group.Id)
	}
	if group.Label != nil {
		d.Set("label", group.Label)
	}
	if group.Order != nil {
		d.Set("order", group.Order)
	}
	if group.Visible != nil {
		d.Set("visible", group.Visible)
	}
	return nil
}
