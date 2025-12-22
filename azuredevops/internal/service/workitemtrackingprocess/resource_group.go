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
		Importer: &schema.ResourceImporter{
			StateContext: importResourceGroup,
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
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Label for the group.",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Order in which the group should appear in the section.",
			},
			"visible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A value indicating if the group should be hidden or not.",
			},
			"control": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Controls to be created with the group. Required for HTML controls which cannot be added to existing groups. This is mutally exclusive with the 'azuredevops_workitemtrackingprocess_control' resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
							Description:      "The ID of the control. This is the field reference name (e.g., System.Description).",
						},
						"label": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label for the control.",
						},
						"control_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the control (e.g., HtmlFieldControl, FieldControl).",
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
							Description: "A value indicating if the control is read only.",
						},
						"order": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Order in which the control should appear in the group.",
						},
						"metadata": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Inner text of the control.",
						},
						"watermark": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Watermark text for the textbox.",
						},
						"inherited": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "A value indicating whether this control has been inherited from a parent layout.",
						},
						"overridden": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "A value indicating whether this control has been overridden by a child layout.",
						},
						"is_contribution": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "A value indicating if the control is a contribution (extension) control.",
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
				},
			},
		},
	}
}

func importResourceGroup(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	// Import ID format: process_id/work_item_type_reference_name/page_id/section_id/group_id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid import ID format, expected: process_id/work_item_type_reference_name/page_id/section_id/group_id")
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_reference_name", parts[1])
	d.Set("page_id", parts[2])
	d.Set("section_id", parts[3])
	d.SetId(parts[4])

	return []*schema.ResourceData{d}, nil
}

func createResourceGroup(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	group := workitemtrackingprocess.Group{
		Label:   converter.String(d.Get("label").(string)),
		Visible: converter.Bool(d.Get("visible").(bool)),
	}
	//nolint:staticcheck // SA1019: d.GetOkExists is deprecated but required to distinguish between unset and zero value
	if v, ok := d.GetOkExists("order"); ok {
		group.Order = converter.Int(v.(int))
	}

	// Add controls to the group if specified
	if v, ok := d.GetOk("control"); ok {
		controlList := v.([]interface{})
		controls := make([]workitemtrackingprocess.Control, len(controlList))
		for i, c := range controlList {
			controlMap := c.(map[string]interface{})
			control := workitemtrackingprocess.Control{
				Id:       converter.String(controlMap["id"].(string)),
				Visible:  converter.Bool(controlMap["visible"].(bool)),
				ReadOnly: converter.Bool(controlMap["read_only"].(bool)),
			}
			if label, ok := controlMap["label"].(string); ok {
				control.Label = converter.String(label)
			}
			if order, ok := controlMap["order"].(int); ok {
				control.Order = converter.Int(order)
			}
			if metadata, ok := controlMap["metadata"].(string); ok {
				control.Metadata = converter.String(metadata)
			}
			if watermark, ok := controlMap["watermark"].(string); ok {
				control.Watermark = converter.String(watermark)
			}
			if isContribution, ok := controlMap["is_contribution"].(bool); ok {
				control.IsContribution = converter.Bool(isContribution)
			}
			if contribution, ok := controlMap["contribution"].([]interface{}); ok && len(contribution) > 0 {
				control.Contribution = expandContribution(contribution)
			}
			controls[i] = control
		}
		group.Controls = &controls
	}

	args := workitemtrackingprocess.AddGroupArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		PageId:     converter.String(d.Get("page_id").(string)),
		SectionId:  converter.String(d.Get("section_id").(string)),
		Group:      &group,
	}

	var createdGroup *workitemtrackingprocess.Group
	err := retryOnContributionNotFound(ctx, d.Timeout(schema.TimeoutCreate), func() error {
		var createErr error
		createdGroup, createErr = clients.WorkItemTrackingProcessClient.AddGroup(ctx, args)
		return createErr
	})
	if err != nil {
		return diag.Errorf(" Creating group. Error %+v", err)
	}

	if createdGroup.Id == nil {
		return diag.Errorf(" Created group has no ID")
	}

	d.SetId(*createdGroup.Id)
	return readResourceGroup(ctx, d, m)
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

	foundGroup := findGroupById(workItemType.Layout, groupId)
	if foundGroup == nil {
		d.SetId("")
		return nil
	}

	d.Set("label", foundGroup.Label)
	d.Set("order", foundGroup.Order)
	d.Set("visible", foundGroup.Visible)

	// Read controls if present
	if foundGroup.Controls != nil && len(*foundGroup.Controls) > 0 {
		controls := make([]map[string]interface{}, len(*foundGroup.Controls))
		for i, c := range *foundGroup.Controls {
			control := map[string]interface{}{
				"visible":    c.Visible,
				"read_only":  c.ReadOnly,
				"inherited":  c.Inherited,
				"overridden": c.Overridden,
			}
			if c.Id != nil {
				control["id"] = *c.Id
			}
			if c.Label != nil {
				control["label"] = *c.Label
			}
			if c.ControlType != nil {
				control["control_type"] = *c.ControlType
			}
			if c.Order != nil {
				control["order"] = *c.Order
			}
			if c.Metadata != nil {
				control["metadata"] = *c.Metadata
			}
			if c.Watermark != nil {
				control["watermark"] = *c.Watermark
			}
			if c.IsContribution != nil {
				control["is_contribution"] = *c.IsContribution
			}
			if c.Contribution != nil {
				control["contribution"] = flattenContribution(c.Contribution)
			}
			controls[i] = control
		}
		d.Set("control", controls)
	}
	return nil
}

func updateResourceGroup(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	groupId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)
	pageId := d.Get("page_id").(string)
	sectionId := d.Get("section_id").(string)

	updateGroup := &workitemtrackingprocess.Group{
		Label:   converter.String(d.Get("label").(string)),
		Visible: converter.Bool(d.Get("visible").(bool)),
	}
	//nolint:staticcheck // SA1019: d.GetOkExists is deprecated but required to distinguish between unset and zero value
	if v, ok := d.GetOkExists("order"); ok {
		updateGroup.Order = converter.Int(v.(int))
	}

	// Check if page_id or section_id has changed - if so, move the group (which also updates it)
	if d.HasChange("page_id") || d.HasChange("section_id") {
		oldPageId, _ := d.GetChange("page_id")
		oldSectionId, _ := d.GetChange("section_id")
		moveArgs := workitemtrackingprocess.MoveGroupToPageArgs{
			ProcessId:           converter.UUID(processId),
			WitRefName:          converter.String(witRefName),
			PageId:              converter.String(pageId),
			SectionId:           converter.String(sectionId),
			GroupId:             &groupId,
			Group:               updateGroup,
			RemoveFromPageId:    converter.String(oldPageId.(string)),
			RemoveFromSectionId: converter.String(oldSectionId.(string)),
		}

		_, err := clients.WorkItemTrackingProcessClient.MoveGroupToPage(ctx, moveArgs)
		if err != nil {
			return diag.Errorf(" Moving group. Error %+v", err)
		}

		return readResourceGroup(ctx, d, m)
	}

	args := workitemtrackingprocess.UpdateGroupArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: converter.String(witRefName),
		PageId:     converter.String(pageId),
		SectionId:  converter.String(sectionId),
		GroupId:    &groupId,
		Group:      updateGroup,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdateGroup(ctx, args)
	if err != nil {
		return diag.Errorf(" Update group. Error %+v", err)
	}

	return readResourceGroup(ctx, d, m)
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

	err := retryOnUnexpectedException(ctx, d.Timeout(schema.TimeoutDelete), func() error {
		return clients.WorkItemTrackingProcessClient.RemoveGroup(ctx, args)
	})

	if err != nil {
		label := d.Get("label").(string)
		return diag.Errorf(" Delete group %s. Error %+v", label, err)
	}

	return nil
}

func findGroupById(layout *workitemtrackingprocess.FormLayout, groupId string) *workitemtrackingprocess.Group {
	if layout == nil {
		return nil
	}
	pages := layout.Pages
	if pages == nil {
		return nil
	}
	for _, page := range *pages {
		group := findGroupInPage(&page, groupId)
		if group != nil {
			return group
		}
	}
	return nil
}

func findGroupInPage(page *workitemtrackingprocess.Page, groupId string) *workitemtrackingprocess.Group {
	sections := page.Sections
	if sections == nil {
		return nil
	}
	for _, section := range *sections {
		group := findGroupInSection(&section, groupId)
		if group != nil {
			return group
		}
	}
	return nil
}

func findGroupInSection(section *workitemtrackingprocess.Section, groupId string) *workitemtrackingprocess.Group {
	groups := section.Groups
	if groups == nil {
		return nil
	}
	for _, group := range *groups {
		id := group.Id
		if id == nil {
			continue
		}
		if *id == groupId {
			return &group
		}
	}
	return nil
}
