package queries

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQueryCreate,
		ReadContext:   resourceQueryRead,
		UpdateContext: resourceQueryUpdate,
		DeleteContext: resourceQueryDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"parent_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},

			"wiql": {
				Type:     schema.TypeString,
				Required: true,
				// The value of 32000 matches the restrictions in Azure DevOps.
				ValidateFunc: validation.StringLenBetween(1, 32000),
			},

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceQueryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)

	params := workitemtracking.CreateQueryArgs{
		Project: &projectID,
		Query:   converter.String(d.Get("parent_path").(string)),
		PostedQuery: &workitemtracking.QueryHierarchyItem{
			Name:     converter.String(d.Get("name").(string)),
			Wiql:     converter.String(d.Get("wiql").(string)),
			IsFolder: converter.Bool(false),
		},
	}

	resp, err := clients.WorkItemTrackingClient.CreateQuery(clients.Ctx, params)
	if err != nil {
		return diag.Errorf(" Creating query. Error: %s", err)
	}

	d.SetId(resp.Id.String())

	return resourceQueryRead(clients.Ctx, d, m)
}

func resourceQueryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	id := d.Id()

	params := workitemtracking.GetQueryArgs{
		Project: converter.String(d.Get("project_id").(string)),
		Query:   &id,
	}

	resp, err := clients.WorkItemTrackingClient.GetQuery(clients.Ctx, params)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting query with id: %s. Error: %+v", id, err)
	}

	if resp != nil {
		if resp.Path != nil {
			d.Set("path", resp.Path)
		}
		if resp.Name != nil {
			d.Set("name", resp.Name)
		}
		if resp.Wiql != nil {
			d.Set("wiql", resp.Wiql)
		}
	}
	return nil
}

func resourceQueryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	var diags diag.Diagnostics

	id := d.Id()

	params := workitemtracking.GetQueryArgs{
		Project: converter.String(d.Get("project_id").(string)),
		Query:   &id,
	}

	existing, err := clients.WorkItemTrackingClient.GetQuery(clients.Ctx, params)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting query with id: %s. Error: %+v", id, err)
	}

	updateArgs := workitemtracking.UpdateQueryArgs{
		Project:     converter.String(d.Get("project_id").(string)),
		Query:       &id,
		QueryUpdate: existing,
	}

	if d.HasChange("wiql") {
		updateArgs.QueryUpdate.Wiql = converter.String(d.Get("wiql").(string))
	}

	if d.HasChange("name") {
		updateArgs.QueryUpdate.Name = converter.String(d.Get("name").(string))
	}

	_, err = clients.WorkItemTrackingClient.UpdateQuery(clients.Ctx, updateArgs)
	if err != nil {
		return diag.Errorf(" Updating query with ID: %s. Error detail: %+v", id, err)
	}

	readDiags := resourceQueryRead(clients.Ctx, d, m)
	diags = append(diags, readDiags...)

	return diags
}

func resourceQueryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	id := d.Id()

	params := workitemtracking.DeleteQueryArgs{
		Project: converter.String(d.Get("project_id").(string)),
		Query:   &id,
	}

	err := clients.WorkItemTrackingClient.DeleteQuery(clients.Ctx, params)
	if err != nil {
		return diag.Errorf(" Deleting query with id %s: %+v", id, err)
	}
	return nil
}
