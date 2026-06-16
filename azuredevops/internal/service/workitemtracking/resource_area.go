package workitemtracking

import (
	"context"
	"fmt"
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

func ResourceArea() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAreaCreate,
		ReadContext:   resourceAreaRead,
		UpdateContext: resourceAreaUpdate,
		DeleteContext: resourceAreaDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceAreaImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
				ForceNew: true,
			},
			"full_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAreaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	parentPath := d.Get("path").(string)

	apiPath := trimLeadingSlash(parentPath)

	node, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &apiPath,
		PostedNode: &workitemtracking.WorkItemClassificationNode{
			Name: &name,
		},
	})
	if err != nil {
		return diag.Errorf("creating area path %q under %q: %+v", name, parentPath, err)
	}

	if node.Identifier == nil {
		return diag.Errorf("creating area path: identifier was nil")
	}

	d.SetId(node.Identifier.String())
	return resourceAreaRead(ctx, d, m)
}

func resourceAreaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	parentPath := d.Get("path").(string)

	apiPath := buildApiPath(parentPath, name)

	node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &apiPath,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("reading area path %q: %+v", apiPath, err)
	}

	d.SetId(node.Identifier.String())
	d.Set("name", node.Name)
	d.Set("full_path", convertAreaNodePath(node.Path))

	return nil
}

func resourceAreaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	oldName, newName := d.GetChange("name")
	parentPath := d.Get("path").(string)

	apiPath := buildApiPath(parentPath, oldName.(string))

	_, err := clients.WorkItemTrackingClient.UpdateClassificationNode(clients.Ctx, workitemtracking.UpdateClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &apiPath,
		PostedNode: &workitemtracking.WorkItemClassificationNode{
			Name: converter.String(newName.(string)),
		},
	})
	if err != nil {
		return diag.Errorf("updating area path %q to %q: %+v", oldName, newName, err)
	}

	return resourceAreaRead(ctx, d, m)
}

func resourceAreaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	parentPath := d.Get("path").(string)

	apiPath := buildApiPath(parentPath, name)

	err := clients.WorkItemTrackingClient.DeleteClassificationNode(clients.Ctx, workitemtracking.DeleteClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &apiPath,
	})
	if err != nil {
		return diag.Errorf("deleting area path %q: %+v", apiPath, err)
	}

	return nil
}

func resourceAreaImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import format, expected: project_id/path (e.g., project-uuid/ParentArea/ChildArea)")
	}

	projectID := parts[0]
	fullPath := "/" + parts[1]

	pathParts := strings.Split(strings.TrimPrefix(fullPath, "/"), "/")
	name := pathParts[len(pathParts)-1]
	parentPath := "/"
	if len(pathParts) > 1 {
		parentPath = "/" + strings.Join(pathParts[:len(pathParts)-1], "/")
	}

	d.Set("project_id", projectID)
	d.Set("name", name)
	d.Set("path", parentPath)

	clients := m.(*client.AggregatedClient)
	apiPath := parts[1]

	node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &apiPath,
		Depth:          converter.Int(0),
	})
	if err != nil {
		return nil, fmt.Errorf("reading area path for import %q: %+v", fullPath, err)
	}

	d.SetId(node.Identifier.String())
	d.Set("full_path", convertAreaNodePath(node.Path))

	return []*schema.ResourceData{d}, nil
}

func buildApiPath(parentPath, name string) string {
	parent := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(parentPath), "/"), "/")
	if parent == "" {
		return name
	}
	return parent + "/" + name
}

func trimLeadingSlash(path string) string {
	return strings.TrimPrefix(strings.TrimSpace(path), "/")
}

func convertAreaNodePath(path *string) string {
	if path == nil {
		return "/"
	}
	itemPath := strings.ReplaceAll(*path, "\\", "/")
	parts := strings.Split(itemPath, "/")
	if len(parts) > 3 {
		return "/" + strings.Join(parts[3:], "/")
	}
	return "/"
}
