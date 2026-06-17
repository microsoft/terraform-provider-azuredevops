package workitemtracking

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

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
				ValidateFunc: validateClassificationNodeName,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
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
	parentID := d.Get("parent_id").(int)

	apiPath, err := resolveParentPath(clients, projectID, parentID)
	if err != nil {
		return diag.Errorf("resolving parent path: %+v", err)
	}

	node, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &apiPath,
		PostedNode: &workitemtracking.WorkItemClassificationNode{
			Name: &name,
		},
	})
	if err != nil {
		return diag.Errorf("creating area path %q: %+v", name, err)
	}
	if node.Id == nil {
		return diag.Errorf("creating area path %q: API did not return node ID", name)
	}

	d.SetId(strconv.Itoa(*node.Id))

	return resourceAreaRead(ctx, d, m)
}

func resourceAreaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("parsing resource ID: %+v", err)
	}

	node, err := clients.WorkItemTrackingClient.GetClassificationNodes(clients.Ctx, workitemtracking.GetClassificationNodesArgs{
		Project: &projectID,
		Ids:     &[]int{id},
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("reading area path %q: %+v", d.Get("name").(string), err)
	}
	if (*node)[0].Id == nil {
		return diag.Errorf("reading area path %q: API did not return node ID", d.Get("name").(string))
	}

	d.SetId(strconv.Itoa(*(*node)[0].Id))
	d.Set("name", (*node)[0].Name)
	d.Set("full_path", convertAreaNodePath((*node)[0].Path))

	return nil
}

func resourceAreaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	oldName, newName := d.GetChange("name")
	parentID := d.Get("parent_id").(int)

	apiPath, err := resolveParentPath(clients, projectID, parentID)
	if err != nil {
		return diag.Errorf("resolving parent path: %+v", err)
	}
	fullApiPath := buildApiPath(apiPath, oldName.(string))

	_, err = clients.WorkItemTrackingClient.UpdateClassificationNode(clients.Ctx, workitemtracking.UpdateClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &fullApiPath,
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
	parentID := d.Get("parent_id").(int)

	apiPath, err := resolveParentPath(clients, projectID, parentID)
	if err != nil {
		return diag.Errorf("resolving parent path: %+v", err)
	}
	fullApiPath := buildApiPath(apiPath, name)

	err = clients.WorkItemTrackingClient.DeleteClassificationNode(clients.Ctx, workitemtracking.DeleteClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           &fullApiPath,
	})
	if err != nil {
		return diag.Errorf("deleting area path %q: %+v", fullApiPath, err)
	}

	return nil
}

func resourceAreaImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import format, expected: project_id/area_id, got: %q", d.Id())
	}

	projectID := parts[0]
	nodeID, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid area_id %q, must be an integer: %+v", parts[1], err)
	}

	clients := m.(*client.AggregatedClient)

	nodes, err := clients.WorkItemTrackingClient.GetClassificationNodes(clients.Ctx, workitemtracking.GetClassificationNodesArgs{
		Project: &projectID,
		Ids:     &[]int{nodeID},
	})
	if err != nil {
		return nil, fmt.Errorf("reading area node (id=%d) for import: %+v", nodeID, err)
	}
	if nodes == nil || len(*nodes) == 0 {
		return nil, fmt.Errorf("area node (id=%d) not found", nodeID)
	}
	node := (*nodes)[0]
	d.SetId(strconv.Itoa(*node.Id))
	d.Set("project_id", projectID)
	d.Set("name", node.Name)
	d.Set("full_path", convertAreaNodePath(node.Path))

	if node.Path != nil {
		parts := strings.Split(*node.Path, "\\")
		if len(parts) > 4 {
			parentApiPath := strings.Join(parts[3:len(parts)-1], "/")
			parentNode, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
				Project:        &projectID,
				StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
				Path:           &parentApiPath,
			})
			if err != nil {
				return nil, fmt.Errorf("reading parent area for import: %+v", err)
			}
			if parentNode.Id != nil {
				d.Set("parent_id", *parentNode.Id)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}

func resolveParentPath(clients *client.AggregatedClient, projectID string, parentID int) (string, error) {
	if parentID == 0 {
		return "", nil
	}

	nodes, err := clients.WorkItemTrackingClient.GetClassificationNodes(clients.Ctx, workitemtracking.GetClassificationNodesArgs{
		Project: &projectID,
		Ids:     &[]int{parentID},
	})
	if err != nil {
		return "", fmt.Errorf("looking up parent area node (id=%d): %+v", parentID, err)
	}
	if nodes == nil || len(*nodes) == 0 {
		return "", fmt.Errorf("parent area node (id=%d) not found", parentID)
	}

	parentNode := (*nodes)[0]
	parts := strings.Split(*parentNode.Path, "\\")
	if len(parts) > 3 {
		return strings.Join(parts[3:], "/"), nil
	}
	return "", nil
}

func buildApiPath(parentPath, name string) string {
	parent := strings.Trim(parentPath, "/")
	if parent == "" {
		return name
	}
	return parent + "/" + name
}

func convertAreaNodePath(path *string) string {
	if path == nil {
		return "/"
	}
	parts := strings.Split(*path, "\\")
	if len(parts) > 3 {
		return "/" + strings.Join(parts[3:], "/")
	}
	return "/"
}

func validateClassificationNodeName(v interface{}, k string) (warnings []string, errors []error) {
	classificationNodeReservedNames := []string{
		"PRN", "CON", "NUL", "AUX", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "COM10", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	value := v.(string)

	if len(strings.TrimSpace(value)) == 0 {
		errors = append(errors, fmt.Errorf("%q must not be empty or whitespace", k))
		return warnings, errors
	}

	if len(value) > 255 {
		errors = append(errors, fmt.Errorf("%q must not exceed 255 characters, got %d", k, len(value)))
	}

	if value == "." || value == ".." {
		errors = append(errors, fmt.Errorf("%q must not be a reserved name (%q)", k, value))
	}

	upper := strings.ToUpper(value)
	if slices.Contains(classificationNodeReservedNames, upper) {
		errors = append(errors, fmt.Errorf("%q must not be a system-reserved name (%s)", k, value))
	}

	const invalidChars = `\/:"*?<>|#$&+`
	for _, ch := range value {
		if strings.ContainsRune(invalidChars, ch) {
			errors = append(errors, fmt.Errorf("%q must not contain the character %q", k, string(ch)))
			break
		}
		if unicode.IsControl(ch) {
			errors = append(errors, fmt.Errorf("%q must not contain Unicode control characters", k))
			break
		}
	}

	return warnings, errors
}
