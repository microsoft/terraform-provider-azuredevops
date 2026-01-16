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

func ResourceInheritedPage() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceInheritedPage,
		ReadContext:   readResourceInheritedPage,
		UpdateContext: updateResourceInheritedPage,
		DeleteContext: deleteResourceInheritedPage,
		Importer: &schema.ResourceImporter{
			StateContext: importResourceInheritedPage,
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
			"page_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID of the inherited page to customize.",
			},
			"label": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Label for the page.",
			},
		},
	}
}

func importResourceInheritedPage(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	// Import ID format: process_id/work_item_type_reference_name/page_id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format, expected: process_id/work_item_type_reference_name/page_id")
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_reference_name", parts[1])
	d.Set("page_id", parts[2])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func createResourceInheritedPage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	pageId := d.Get("page_id").(string)
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		// Returns the layout containing the pages
		Expand: &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		return diag.Errorf("getting work item type: %+v", err)
	}
	if workItemType == nil || workItemType.Layout == nil {
		return diag.Errorf("work item type or layout is nil")
	}

	existingPage := findPageById(workItemType.Layout, pageId)
	if existingPage == nil {
		return diag.Errorf("page %s not found in layout", pageId)
	}

	if existingPage.Inherited == nil || !*existingPage.Inherited {
		return diag.Errorf("page %s is not inherited, use azuredevops_workitemtrackingprocess_page to manage custom pages", pageId)
	}

	d.SetId(pageId)

	return updateResourceInheritedPage(ctx, d, m)
}

func readResourceInheritedPage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	pageId := d.Id()
	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_reference_name").(string)

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		// Returns the layout containing the pages
		Expand: &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		return diag.Errorf("getting work item type: %+v", err)
	}
	if workItemType == nil || workItemType.Layout == nil {
		return diag.Errorf("work item type or layout is nil")
	}

	page := findPageById(workItemType.Layout, pageId)
	if page == nil {
		d.SetId("")
		return nil
	}

	if page.Label != nil {
		d.Set("label", *page.Label)
	}
	return nil
}

func updateResourceInheritedPage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	pageId := d.Id()

	page := workitemtrackingprocess.Page{
		Id:    &pageId,
		Label: converter.String(d.Get("label").(string)),
	}

	args := workitemtrackingprocess.UpdatePageArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		Page:       &page,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdatePage(ctx, args)
	if err != nil {
		return diag.Errorf("updating inherited page: %+v", err)
	}

	return readResourceInheritedPage(ctx, d, m)
}

func deleteResourceInheritedPage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	pageId := d.Id()

	args := workitemtrackingprocess.RemovePageArgs{
		ProcessId:  converter.UUID(d.Get("process_id").(string)),
		WitRefName: converter.String(d.Get("work_item_type_reference_name").(string)),
		PageId:     &pageId,
	}

	err := clients.WorkItemTrackingProcessClient.RemovePage(ctx, args)
	if err != nil {
		return diag.Errorf("reverting inherited page: %+v", err)
	}

	return nil
}

func findPageById(layout *workitemtrackingprocess.FormLayout, pageId string) *workitemtrackingprocess.Page {
	if layout == nil || layout.Pages == nil {
		return nil
	}
	for _, page := range *layout.Pages {
		if page.Id != nil && *page.Id == pageId {
			return &page
		}
	}
	return nil
}
