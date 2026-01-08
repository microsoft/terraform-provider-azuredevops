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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceWorkItemType() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceWorkItemType,
		ReadContext:   readResourceWorkItemType,
		UpdateContext: updateResourceWorkItemType,
		DeleteContext: deleteResourceWorkItemType,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				processId, referenceName, err := tfhelper.ParseImportedName(d.Id(), "process_id/reference_name")
				if err != nil {
					return nil, err
				}
				d.Set("process_id", processId)
				d.SetId(referenceName)
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
			"parent_work_item_reference_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference name of the parent work item type.",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "True if the work item type is enabled.",
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
			"pages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of pages for the work item type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the page.",
						},
						"page_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the page.",
						},
						"sections": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of sections in the page.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the section.",
									},
									"groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of groups in the section.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the group.",
												},
												"controls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of controls in the group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The ID of the control.",
															},
														},
													},
												},
											},
										},
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

func createResourceWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	workItemTypeRequest := workitemtrackingprocess.CreateProcessWorkItemTypeRequest{
		Name:       converter.String(d.Get("name").(string)),
		IsDisabled: converter.Bool(!d.Get("is_enabled").(bool)),
		Color:      convertColorToApi(d),
		Icon:       converter.String(d.Get("icon").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		workItemTypeRequest.Description = converter.String(v.(string))
	}
	if v, ok := d.GetOk("parent_work_item_reference_name"); ok {
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
	if createdWorkItemType.ReferenceName == nil {
		return diag.Errorf(" Creating work item type. Reference name is nil")
	}
	d.SetId(*createdWorkItemType.ReferenceName)

	// The POST operation doesn't support layout expand, so we have to call read and risk eventual consistency problems
	return readResourceWorkItemType(ctx, d, m)
}

func readResourceWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	referenceName := d.Id()
	processId := d.Get("process_id").(string)

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &referenceName,
		Expand:     &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting work item type with reference name: %s for process with id %s. Error: %+v", referenceName, processId, err)
	}

	return flattenWorkItemType(d, workItemType)
}

func updateResourceWorkItemType(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	referenceName := d.Id()
	processId := d.Get("process_id").(string)

	updateWorkItemType := &workitemtrackingprocess.UpdateProcessWorkItemTypeRequest{
		IsDisabled: converter.Bool(!d.Get("is_enabled").(bool)),
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

	_, err := clients.WorkItemTrackingProcessClient.UpdateProcessWorkItemType(ctx, args)
	if err != nil {
		return diag.Errorf(" Update work item type. Error %+v", err)
	}

	// The PATCH operation doesn't support layout expand, so we have to call read and risk eventual consistency problems
	return readResourceWorkItemType(ctx, d, m)
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

func flattenWorkItemType(d *schema.ResourceData, workItemType *workitemtrackingprocess.ProcessWorkItemType) diag.Diagnostics {
	d.Set("name", workItemType.Name)
	d.Set("description", workItemType.Description)
	if workItemType.IsDisabled != nil {
		d.Set("is_enabled", !*workItemType.IsDisabled)
	} else {
		d.Set("is_enabled", true)
	}
	if workItemType.Color != nil {
		d.Set("color", convertColorToResource(*workItemType.Color))
	} else {
		d.Set("color", nil)
	}
	d.Set("icon", workItemType.Icon)
	d.Set("parent_work_item_reference_name", workItemType.Inherits)
	d.Set("reference_name", workItemType.ReferenceName)
	d.Set("url", workItemType.Url)

	var pages []map[string]any
	if workItemType.Layout != nil && workItemType.Layout.Pages != nil {
		for _, page := range *workItemType.Layout.Pages {
			pages = append(pages, flattenPage(page))
		}
	}
	if err := d.Set("pages", pages); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func convertColorToApi(d *schema.ResourceData) *string {
	return converter.String(
		strings.ReplaceAll(d.Get("color").(string), "#", ""))
}

func convertColorToResource(apiFormattedColor string) string {
	return fmt.Sprintf("#%s", apiFormattedColor)
}

func flattenControl(control workitemtrackingprocess.Control) map[string]any {
	controlMap := map[string]any{}
	if control.Id != nil {
		controlMap["id"] = *control.Id
	}
	return controlMap
}

func flattenGroup(group workitemtrackingprocess.Group) map[string]any {
	groupMap := map[string]any{}
	if group.Id != nil {
		groupMap["id"] = *group.Id
	}
	if group.Controls != nil {
		var controls []map[string]any
		for _, control := range *group.Controls {
			controls = append(controls, flattenControl(control))
		}
		groupMap["controls"] = controls
	}
	return groupMap
}

func flattenSection(section workitemtrackingprocess.Section) map[string]any {
	sectionMap := map[string]any{}
	if section.Id != nil {
		sectionMap["id"] = *section.Id
	}
	if section.Groups != nil {
		var groups []map[string]any
		for _, group := range *section.Groups {
			groups = append(groups, flattenGroup(group))
		}
		sectionMap["groups"] = groups
	}
	return sectionMap
}

func flattenPage(page workitemtrackingprocess.Page) map[string]any {
	pageMap := map[string]any{}
	if page.Id != nil {
		pageMap["id"] = *page.Id
	}
	if page.PageType != nil {
		pageMap["page_type"] = string(*page.PageType)
	}
	if page.Sections != nil {
		var sections []map[string]any
		for _, section := range *page.Sections {
			sections = append(sections, flattenSection(section))
		}
		pageMap["sections"] = sections
	}
	return pageMap
}
