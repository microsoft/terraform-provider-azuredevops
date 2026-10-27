package workitemtracking

import (
	"context"
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

const iterationDateFormat = "2006-01-02"

func ResourceIteration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIterationCreate,
		ReadContext:   resourceIterationRead,
		UpdateContext: resourceIterationUpdate,
		DeleteContext: resourceIterationDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceIterationImport,
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
				ValidateFunc: validateClassificationNodeName,
			},
			"parent_iteration_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"finish_date"},
				ValidateFunc: validateIterationDate,
			},
			"finish_date": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"start_date"},
				ValidateFunc: validateIterationDate,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIterationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	parentID := d.Get("parent_iteration_id").(int)

	apiPath, err := resolveParentPath(clients, projectID, parentID)
	if err != nil {
		return diag.Errorf("resolving parent path: %+v", err)
	}

	postedNode := &workitemtracking.WorkItemClassificationNode{
		Name: &name,
	}
	if _, ok := d.GetOk("start_date"); ok {
		postedNode.Attributes = expandIterationDates(d)
	}

	node, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Path:           &apiPath,
		PostedNode:     postedNode,
	})
	if err != nil {
		return diag.Errorf("creating iteration path %q: %+v", name, err)
	}
	if node.Id == nil {
		return diag.Errorf("creating iteration path %q: API did not return node ID", name)
	}

	d.SetId(strconv.Itoa(*node.Id))

	return resourceIterationRead(ctx, d, m)
}

func resourceIterationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("parsing resource ID: %+v", err)
	}

	nodes, err := clients.WorkItemTrackingClient.GetClassificationNodes(clients.Ctx, workitemtracking.GetClassificationNodesArgs{
		Project: &projectID,
		Ids:     &[]int{id},
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("reading iteration path %q: %+v", d.Get("name").(string), err)
	}

	if len(*nodes) != 1 {
		return diag.Errorf("reading iteration path %q: unexpected number of nodes returned for iteration (id=%d)", d.Get("name").(string), id)
	}

	node := (*nodes)[0]
	if node.Id == nil {
		return diag.Errorf("reading iteration path %q: iteration node (id=%d) has nil ID", d.Get("name").(string), id)
	}

	if err := flattenIteration(d, &node); err != nil {
		return diag.Errorf("reading iteration path %q: %+v", d.Get("name").(string), err)
	}

	return nil
}

func resourceIterationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	oldName, newName := d.GetChange("name")
	parentID := d.Get("parent_iteration_id").(int)

	apiPath, err := resolveParentPath(clients, projectID, parentID)
	if err != nil {
		return diag.Errorf("resolving parent path: %+v", err)
	}
	fullApiPath := buildApiPath(apiPath, oldName.(string))

	postedNode := &workitemtracking.WorkItemClassificationNode{
		Name: converter.String(newName.(string)),
	}
	if d.HasChanges("start_date", "finish_date") {
		postedNode.Attributes = expandIterationDates(d)
	}

	_, err = clients.WorkItemTrackingClient.UpdateClassificationNode(clients.Ctx, workitemtracking.UpdateClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Path:           &fullApiPath,
		PostedNode:     postedNode,
	})
	if err != nil {
		return diag.Errorf("updating iteration path %q to %q: %+v", oldName, newName, err)
	}

	return resourceIterationRead(ctx, d, m)
}

func resourceIterationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	parentID := d.Get("parent_iteration_id").(int)

	apiPath, err := resolveParentPath(clients, projectID, parentID)
	if err != nil {
		return diag.Errorf("resolving parent path: %+v", err)
	}
	fullApiPath := buildApiPath(apiPath, name)

	err = clients.WorkItemTrackingClient.DeleteClassificationNode(clients.Ctx, workitemtracking.DeleteClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Path:           &fullApiPath,
	})
	if err != nil {
		return diag.Errorf("deleting iteration path %q: %+v", fullApiPath, err)
	}

	return nil
}

func resourceIterationImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import format, expected: project_id/iteration_id, got: %q", d.Id())
	}

	projectID := parts[0]
	nodeID, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid iteration_id %q, must be an integer: %+v", parts[1], err)
	}

	clients := m.(*client.AggregatedClient)

	nodes, err := clients.WorkItemTrackingClient.GetClassificationNodes(clients.Ctx, workitemtracking.GetClassificationNodesArgs{
		Project: &projectID,
		Ids:     &[]int{nodeID},
	})
	if err != nil {
		return nil, fmt.Errorf("reading iteration node (id=%d) for import: %+v", nodeID, err)
	}
	if nodes == nil || len(*nodes) == 0 {
		return nil, fmt.Errorf("iteration node (id=%d) not found", nodeID)
	}
	node := (*nodes)[0]
	d.Set("project_id", projectID)
	if err := flattenIteration(d, &node); err != nil {
		return nil, fmt.Errorf("reading iteration node (id=%d) for import: %+v", nodeID, err)
	}

	if node.Path != nil {
		parts := strings.Split(*node.Path, "\\")
		if len(parts) > 4 {
			parentApiPath := strings.Join(parts[3:len(parts)-1], "/")
			parentNode, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
				Project:        &projectID,
				StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
				Path:           &parentApiPath,
			})
			if err != nil {
				return nil, fmt.Errorf("reading parent iteration for import: %+v", err)
			}
			if parentNode.Id != nil {
				d.Set("parent_iteration_id", *parentNode.Id)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}

func flattenIteration(d *schema.ResourceData, node *workitemtracking.WorkItemClassificationNode) error {
	d.SetId(strconv.Itoa(*node.Id))
	d.Set("name", node.Name)
	d.Set("path", convertAreaNodePath(node.Path))

	startDate, finishDate := "", ""
	if node.Attributes != nil {
		var err error
		if startDate, err = flattenIterationDate((*node.Attributes)["startDate"]); err != nil {
			return fmt.Errorf("parsing startDate: %+v", err)
		}
		if finishDate, err = flattenIterationDate((*node.Attributes)["finishDate"]); err != nil {
			return fmt.Errorf("parsing finishDate: %+v", err)
		}
	}
	d.Set("start_date", startDate)
	d.Set("finish_date", finishDate)
	return nil
}

// expandIterationDates returns the node attributes for the configured dates.
// The API requires both dates or neither; unset dates are sent as null to clear them.
func expandIterationDates(d *schema.ResourceData) *map[string]interface{} {
	attributes := map[string]interface{}{
		"startDate":  nil,
		"finishDate": nil,
	}
	if v, ok := d.GetOk("start_date"); ok {
		attributes["startDate"] = v.(string) + "T00:00:00Z"
	}
	if v, ok := d.GetOk("finish_date"); ok {
		attributes["finishDate"] = v.(string) + "T00:00:00Z"
	}
	return &attributes
}

func flattenIterationDate(v interface{}) (string, error) {
	s, ok := v.(string)
	if !ok || s == "" {
		return "", nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return "", err
	}
	return t.Format(iterationDateFormat), nil
}

func validateIterationDate(v interface{}, k string) (warnings []string, errors []error) {
	if _, err := time.Parse(iterationDateFormat, v.(string)); err != nil {
		errors = append(errors, fmt.Errorf("%q must be a date in YYYY-MM-DD format, got: %q", k, v))
	}
	return warnings, errors
}
