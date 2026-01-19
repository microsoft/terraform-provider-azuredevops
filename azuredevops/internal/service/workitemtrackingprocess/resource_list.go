package workitemtrackingprocess

import (
	"context"
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

func ResourceList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListCreate,
		ReadContext:   resourceListRead,
		UpdateContext: resourceListUpdate,
		DeleteContext: resourceListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Description:      "Name of the list.",
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "string",
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"string", "integer"}, false)),
				Description:      "Data type of the list. Valid values: string, integer.",
			},
			"items": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "A list of items.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_suggested": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether items outside of the suggested list are allowed.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the list.",
			},
		},
	}
}

func resourceListCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	listType := d.Get("type").(string)
	isSuggested := d.Get("is_suggested").(bool)
	items := expandItems(d.Get("items").([]any))

	picklist := &workitemtrackingprocess.PickList{
		Name:        &name,
		Type:        &listType,
		IsSuggested: &isSuggested,
		Items:       &items,
	}

	args := workitemtrackingprocess.CreateListArgs{
		Picklist: picklist,
	}

	createdList, err := clients.WorkItemTrackingProcessClient.CreateList(ctx, args)
	if err != nil {
		return diag.Errorf(" Creating list. Error: %+v", err)
	}

	if createdList.Id == nil {
		return diag.Errorf(" Created list has no ID")
	}

	d.SetId(createdList.Id.String())
	return resourceListRead(ctx, d, m)
}

func resourceListRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	listId := d.Id()

	args := workitemtrackingprocess.GetListArgs{
		ListId: converter.UUID(listId),
	}

	list, err := clients.WorkItemTrackingProcessClient.GetList(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Reading list %s. Error: %+v", listId, err)
	}

	if list.Name != nil {
		d.Set("name", *list.Name)
	}
	if list.Type != nil {
		d.Set("type", strings.ToLower(*list.Type))
	}
	if list.IsSuggested != nil {
		d.Set("is_suggested", *list.IsSuggested)
	}
	if list.Items != nil {
		if err := d.Set("items", *list.Items); err != nil {
			return diag.Errorf(" setting items: %+v", err)
		}
	}
	if list.Url != nil {
		d.Set("url", *list.Url)
	}

	return nil
}

func resourceListUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	listId := d.Id()
	name := d.Get("name").(string)
	isSuggested := d.Get("is_suggested").(bool)
	items := expandItems(d.Get("items").([]any))

	picklist := &workitemtrackingprocess.PickList{
		Id:          converter.UUID(listId),
		Name:        &name,
		IsSuggested: &isSuggested,
		Items:       &items,
	}

	args := workitemtrackingprocess.UpdateListArgs{
		ListId:   converter.UUID(listId),
		Picklist: picklist,
	}

	_, err := clients.WorkItemTrackingProcessClient.UpdateList(ctx, args)
	if err != nil {
		return diag.Errorf(" Updating list %s. Error: %+v", listId, err)
	}

	return resourceListRead(ctx, d, m)
}

func resourceListDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	listId := d.Id()

	args := workitemtrackingprocess.DeleteListArgs{
		ListId: converter.UUID(listId),
	}

	err := clients.WorkItemTrackingProcessClient.DeleteList(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return diag.Errorf(" Deleting list %s. Error: %+v", listId, err)
	}

	return nil
}

func expandItems(input []any) []string {
	items := make([]string, len(input))
	for i, v := range input {
		items[i] = v.(string)
	}
	return items
}
