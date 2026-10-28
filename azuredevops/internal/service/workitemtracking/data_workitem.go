package workitemtracking

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataWorkItem() *schema.Resource {
	dataSchema := map[string]*schema.Schema{
		"id": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "The ID of the work item.",
		},
		"project_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
			Description:  "The ID of the project the work item belongs to.",
		},
	}

	for name, attribute := range workItemComputedSchema() {
		dataSchema[name] = attribute
	}

	return &schema.Resource{
		ReadContext: dataSourceWorkItemRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: dataSchema,
	}
}

func dataSourceWorkItemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	id := d.Get("id").(int)

	args := workitemtracking.GetWorkItemArgs{
		Id:      &id,
		Project: converter.String(d.Get("project_id").(string)),
		Expand:  converter.ToPtr(workitemtracking.WorkItemExpandValues.All),
	}

	workItem, err := clients.WorkItemTrackingClient.GetWorkItem(ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return diag.Errorf("Work item with ID %d was not found", id)
		}
		return diag.Errorf("Getting work item with ID %d: %+v", id, err)
	}

	workItemData := flattenWorkItemDataSource(workItem)

	d.SetId(strconv.Itoa(id))
	d.Set("title", workItemData["title"])
	d.Set("description", workItemData["description"])
	d.Set("type", workItemData["type"])
	d.Set("state", workItemData["state"])
	d.Set("area_path", workItemData["area_path"])
	d.Set("iteration_path", workItemData["iteration_path"])
	d.Set("parent_id", workItemData["parent_id"])
	d.Set("url", workItemData["url"])

	if err = d.Set("tags", workItemData["tags"]); err != nil {
		return diag.Errorf("Setting `tags`: %+v", err)
	}

	if err = d.Set("fields", workItemData["fields"]); err != nil {
		return diag.Errorf("Setting `fields`: %+v", err)
	}

	if err = d.Set("relations", workItemData["relations"]); err != nil {
		return diag.Errorf("Setting `relations`: %+v", err)
	}

	return nil
}

func workItemComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The title of the work item.",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The description of the work item.",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The type of the work item.",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The state of the work item.",
		},
		"area_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The area path of the work item.",
		},
		"iteration_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The iteration path of the work item.",
		},
		"parent_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The ID of the parent work item. `0` when the work item has no parent.",
		},
		"tags": {
			Type:        schema.TypeSet,
			Computed:    true,
			Set:         schema.HashString,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The tags of the work item.",
		},
		"url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The REST API URL of the work item.",
		},
		"fields": {
			Type:        schema.TypeMap,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "All fields of the work item, keyed by their reference name. Non string values are JSON encoded.",
		},
		"relations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"rel": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The reference name of the relation.",
					},
					"url": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The URL of the related object.",
					},
				},
			},
			Description: "The relations of the work item.",
		},
	}
}

func flattenWorkItemDataSource(workItem *workitemtracking.WorkItem) map[string]interface{} {
	flattened := map[string]interface{}{}

	if workItem == nil {
		return flattened
	}

	if workItem.Id != nil {
		flattened["id"] = *workItem.Id
	}
	if workItem.Url != nil {
		flattened["url"] = *workItem.Url
	}

	if workItem.Fields != nil {
		fields := make(map[string]string, len(*workItem.Fields))
		for key, value := range *workItem.Fields {
			fields[key] = workItemFieldToString(value)

			switch key {
			case "System.Tags":
				if tags, ok := value.(string); ok {
					flattened["tags"] = strings.Split(tags, "; ")
				}
			default:
				if name, ok := systemFieldMapping[key]; ok {
					flattened[name] = value
				}
			}
		}
		flattened["fields"] = fields
	}

	if workItem.Relations != nil {
		relations := make([]interface{}, 0, len(*workItem.Relations))
		for _, relation := range *workItem.Relations {
			relations = append(relations, map[string]interface{}{
				"rel": converter.ToString(relation.Rel, ""),
				"url": converter.ToString(relation.Url, ""),
			})
		}
		flattened["relations"] = relations
	}

	return flattened
}

func workItemFieldToString(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	default:
		encoded, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(encoded)
	}
}
