package workitemtrackingprocess

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
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
			"work_item_type_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "The ID (reference name) of the work item type.",
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
	parts, err := tfhelper.ParseImportedNameParts(d.Id(), "process_id/work_item_type_id/page_id", 3)
	if err != nil {
		return nil, err
	}

	d.Set("process_id", parts[0])
	d.Set("work_item_type_id", parts[1])
	d.Set("page_id", parts[2])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func createResourceInheritedPage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	pageId := d.Get("page_id").(string)

	existingPage, diags := getInheritedPage(ctx, d, m, pageId)
	if diags != nil {
		return diags
	}
	if existingPage == nil {
		return diag.Errorf("page %s not found in layout", pageId)
	}

	d.SetId(pageId)

	return updateResourceInheritedPage(ctx, d, m)
}

func readResourceInheritedPage(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	page, diags := getInheritedPage(ctx, d, m, d.Id())
	if diags != nil {
		return diags
	}
	if page == nil {
		d.SetId("")
		return nil
	}

	if page.Label != nil {
		d.Set("label", *page.Label)
	}
	return nil
}

func getInheritedPage(ctx context.Context, d *schema.ResourceData, m any, pageId string) (*workitemtrackingprocess.Page, diag.Diagnostics) {
	clients := m.(*client.AggregatedClient)

	processId := d.Get("process_id").(string)
	witRefName := d.Get("work_item_type_id").(string)

	getWorkItemTypeArgs := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
		ProcessId:  converter.UUID(processId),
		WitRefName: &witRefName,
		Expand:     &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
	}
	workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(ctx, getWorkItemTypeArgs)
	if err != nil {
		return nil, diag.Errorf("getting work item type: %+v", err)
	}
	if workItemType == nil {
		return nil, diag.Errorf("work item type is nil")
	}
	if workItemType.Layout == nil {
		return nil, diag.Errorf("work item type layout is nil")
	}

	page := findPageById(workItemType.Layout, pageId)
	if page == nil {
		return nil, nil
	}

	if page.Inherited == nil || !*page.Inherited {
		return nil, diag.Errorf("page %s is not inherited, use azuredevops_workitemtrackingprocess_page to manage custom pages", pageId)
	}

	return page, nil
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
		WitRefName: converter.String(d.Get("work_item_type_id").(string)),
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
		WitRefName: converter.String(d.Get("work_item_type_id").(string)),
		PageId:     &pageId,
	}

	err := clients.WorkItemTrackingProcessClient.RemovePage(ctx, args)
	if err != nil {
		return diag.Errorf("reverting inherited page: %+v", err)
	}

	return nil
}
