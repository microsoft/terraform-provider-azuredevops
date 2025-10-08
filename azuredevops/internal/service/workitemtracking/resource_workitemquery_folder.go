package workitemtracking

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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceQueryFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQueryFolderCreate,
		ReadContext:   resourceQueryFolderRead,
		DeleteContext: resourceQueryFolderDelete,
		Importer: 	   tfhelper.ImportProjectQualifiedResource(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			// The ID of the parent folder.
			// It should not be specified if 'area' is specified.
			"parent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"parent_id", "area"},
			},

			// If specified, the area should be one of either 'Shared Queries' or 'My Queries'.
			// It should not be specified if 'parent_id' is specified.
			"area": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Shared Queries", "My Queries"}, false),
				ExactlyOneOf: []string{"parent_id", "area"},
			},
		},
	}
}

func resourceQueryFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)

	parent := d.Get("area").(string)
	if parent == "" {
		parent = d.Get("parent_id").(string)
	}

	params := workitemtracking.CreateQueryArgs{
		Project: &projectID,
		Query:   &parent,
		PostedQuery: &workitemtracking.QueryHierarchyItem{
			Name:     converter.String(d.Get("name").(string)),
			IsFolder: converter.Bool(true),
		},
	}

	resp, err := clients.WorkItemTrackingClient.CreateQuery(clients.Ctx, params)
	if err != nil {
		return diag.Errorf(" Creating query folder. Error: %s", err)
	}

	if resp.Id == nil {
		return diag.Errorf(" Creating query folder. Error: ID was nil")
	}

	d.SetId(resp.Id.String())
	return resourceQueryFolderRead(clients.Ctx, d, m)
}

func resourceQueryFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	if d.Id() == "" {
		return diag.Errorf(" Reading query folder: ID was not set")
	}

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
		return diag.Errorf(" Getting query folder with id: %s. Error: %+v", id, err)
	}

	if resp != nil {
		if resp.Name != nil {
			d.Set("name", resp.Name)
		}
		if resp.IsFolder != nil && !*resp.IsFolder {
			return diag.Errorf(" The query with id: %s is not a folder.", id)
		}
	}
	return nil
}

func resourceQueryFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	id := d.Id()

	params := workitemtracking.DeleteQueryArgs{
		Project: converter.String(d.Get("project_id").(string)),
		Query:   &id,
	}

	err := clients.WorkItemTrackingClient.DeleteQuery(clients.Ctx, params)
	if err != nil {
		return diag.Errorf(" Deleting query folder with id %s: %+v", id, err)
	}
	return nil
}
