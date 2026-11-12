package workitemtracking

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const getWorkItemsBatchSize = 200

func DataWorkItems() *schema.Resource {
	workItemSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The ID of the work item.",
		},
	}

	for name, attribute := range workItemComputedSchema() {
		workItemSchema[name] = attribute
	}

	return &schema.Resource{
		ReadContext: dataSourceWorkItemsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"wiql": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ExactlyOneOf: []string{"wiql", "query_id"},
				Description:  "The Work Item Query Language (WIQL) statement used to look up the work items.",
			},
			"query_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"wiql", "query_id"},
				Description:  "The ID of a saved query used to look up the work items.",
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the project the query is scoped to.",
			},
			"work_items": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource{Schema: workItemSchema},
				Description: "The work items returned by the query.",
			},
		},
	}
}

func dataSourceWorkItemsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	wiql := d.Get("wiql").(string)
	queryID := d.Get("query_id").(string)
	projectID := d.Get("project_id").(string)

	var project *string
	if projectID != "" {
		project = &projectID
	}

	var queryResult *workitemtracking.WorkItemQueryResult
	var err error
	if queryID != "" {
		queryResult, err = clients.WorkItemTrackingClient.QueryById(ctx, workitemtracking.QueryByIdArgs{
			Id:      converter.UUID(queryID),
			Project: project,
		})
		if err != nil {
			return diag.Errorf("Running query %s: %+v", queryID, err)
		}
	} else {
		queryResult, err = clients.WorkItemTrackingClient.QueryByWiql(ctx, workitemtracking.QueryByWiqlArgs{
			Wiql:    &workitemtracking.Wiql{Query: &wiql},
			Project: project,
		})
		if err != nil {
			return diag.Errorf("Running WIQL query: %+v", err)
		}
	}

	ids := collectWorkItemIDs(queryResult)

	workItems := make([]interface{}, 0, len(ids))
	for start := 0; start < len(ids); start += getWorkItemsBatchSize {
		end := min(start+getWorkItemsBatchSize, len(ids))
		batch := ids[start:end]

		batchResult, err := clients.WorkItemTrackingClient.GetWorkItems(ctx, workitemtracking.GetWorkItemsArgs{
			Ids:         &batch,
			Project:     project,
			Expand:      converter.ToPtr(workitemtracking.WorkItemExpandValues.All),
			ErrorPolicy: converter.ToPtr(workitemtracking.WorkItemErrorPolicyValues.Omit),
		})
		if err != nil {
			return diag.Errorf("Getting work items: %+v", err)
		}
		if batchResult == nil {
			continue
		}

		for i := range *batchResult {
			workItem := (*batchResult)[i]
			if workItem.Id == nil {
				continue
			}
			workItems = append(workItems, flattenWorkItemDataSource(&workItem))
		}
	}

	d.SetId("workitems#" + strconv.Itoa(tfhelper.HashString(projectID+"|"+queryID+"|"+wiql)))
	if err = d.Set("work_items", workItems); err != nil {
		return diag.Errorf("Setting `work_items`: %+v", err)
	}

	return nil
}

func collectWorkItemIDs(queryResult *workitemtracking.WorkItemQueryResult) []int {
	ids := make([]int, 0)
	if queryResult == nil {
		return ids
	}

	seen := map[int]bool{}
	appendID := func(reference *workitemtracking.WorkItemReference) {
		if reference == nil || reference.Id == nil || seen[*reference.Id] {
			return
		}
		seen[*reference.Id] = true
		ids = append(ids, *reference.Id)
	}

	if queryResult.WorkItems != nil {
		for i := range *queryResult.WorkItems {
			appendID(&(*queryResult.WorkItems)[i])
		}
	}

	if queryResult.WorkItemRelations != nil {
		for i := range *queryResult.WorkItemRelations {
			appendID((*queryResult.WorkItemRelations)[i].Target)
		}
	}

	return ids
}
