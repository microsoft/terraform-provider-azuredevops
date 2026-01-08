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

func ResourcePage() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourcePage,
		ReadContext:   readResourcePage,
		UpdateContext: updateResourcePage,
		DeleteContext: deleteResourcePage,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: importResourcePage,
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
			"label": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The label for the page.",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Order in which the page should appear in the layout.",
			},
			"visible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A value indicating if the page should be hidden or not.",
			},
			"page_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the page (e.g., custom, history, links, attachments).",
			},
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A value indicating whether any user operations are permitted on this page.",
			},
			"inherited": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A value indicating whether this page has been inherited from a parent layout.",
			},
			"overridden": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A value indicating whether this page has been overridden by a child layout.",
			},
			"section": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The sections of the page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the section.",
						},
					},
				},
			},
		},
	}
}

func importResourcePage(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	// Import ID format: process_id/work_item_type_reference_name/page_id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format, expected: process_id/work_item_type_reference_name/page_id")
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_reference_name", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func createResourcePage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	page := workitemtrackingprocess.Page{
		Label:    converter.String(d.Get("label").(string)),
		Visible:  converter.Bool(d.Get("visible").(bool)),
		PageType: &workitemtrackingprocess.PageTypeValues.Custom,
	}
	rawConfig := d.GetRawConfig().AsValueMap()
	if order := rawConfig["order"]; !order.IsNull() {
		page.Order = converter.Int(d.Get("order").(int))
	}

	args := workitemtrackingprocess.AddPageArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		Page:       &page,
	}

	createdPage, err := clients.WorkItemTrackingProcessClient.AddPage(ctx, args)
	if err != nil {
		return diag.Errorf(" Creating page. Error %+v", err)
	}

	if createdPage.Id == nil {
		return diag.Errorf(" Created page has no ID")
	}

	d.SetId(*createdPage.Id)
	return readResourcePage(ctx, d, m)
}

func readResourcePage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	pageId := d.Id()
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

	foundPage := findPageById(workItemType.Layout, pageId)
	if foundPage == nil {
		d.SetId("")
		return nil
	}

	d.Set("label", foundPage.Label)
	d.Set("order", foundPage.Order)
	d.Set("visible", foundPage.Visible)
	d.Set("page_type", foundPage.PageType)
	d.Set("locked", foundPage.Locked)
	d.Set("inherited", foundPage.Inherited)
	d.Set("overridden", foundPage.Overridden)

	if foundPage.Sections != nil {
		sections := make([]map[string]any, len(*foundPage.Sections))
		for i, s := range *foundPage.Sections {
			section := map[string]any{}
			if s.Id != nil {
				section["id"] = *s.Id
			}
			sections[i] = section
		}
		d.Set("section", sections)
	}

	return nil
}

func updateResourcePage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	pageId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	updatePage := &workitemtrackingprocess.Page{
		Id:      &pageId,
		Label:   converter.String(d.Get("label").(string)),
		Visible: converter.Bool(d.Get("visible").(bool)),
	}
	rawConfig := d.GetRawConfig().AsValueMap()
	if order := rawConfig["order"]; !order.IsNull() {
		updatePage.Order = converter.Int(d.Get("order").(int))
	}

	args := workitemtrackingprocess.UpdatePageArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: converter.String(witRefName),
		Page:       updatePage,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdatePage(ctx, args)
	if err != nil {
		return diag.Errorf(" Update page. Error %+v", err)
	}

	return readResourcePage(ctx, d, m)
}

func deleteResourcePage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	pageId := d.Id()

	args := workitemtrackingprocess.RemovePageArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		PageId:     &pageId,
	}

	err := clients.WorkItemTrackingProcessClient.RemovePage(ctx, args)
	if err != nil {
		return diag.Errorf(" Delete page. Error %+v", err)
	}

	return nil
}

func findPageById(layout *workitemtrackingprocess.FormLayout, pageId string) *workitemtrackingprocess.Page {
	if layout == nil {
		return nil
	}
	pages := layout.Pages
	if pages == nil {
		return nil
	}
	for _, page := range *pages {
		if page.Id != nil && *page.Id == pageId {
			return &page
		}
	}
	return nil
}

