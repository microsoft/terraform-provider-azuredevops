package workitemtracking

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceArea() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateArea,
		Read:   resourceReadArea,
		Update: resourceUpdateArea,
		Delete: resourceDeleteArea,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace, validation.StringDoesNotContainAny("\\/$?*\":<>|#{}[]+=%")),
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"has_children": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreateArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	postedNode := workitemtracking.WorkItemClassificationNode{
		Name: converter.String(d.Get("name").(string)),
	}

	args := workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		PostedNode:     &postedNode,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           converter.String(d.Get("path").(string)),
	}

	workItem, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("creating Area: %w", err)
	}

	if workItem.Id != nil {
		d.SetId(strconv.Itoa(*workItem.Id))
	} else {
		return fmt.Errorf("creating Area: API did not return an ID")
	}

	return resourceReadArea(d, m)
}

func resourceReadArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)

	params := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           converter.String(d.Id()),
		Depth:          converter.Int(0),
	}

	node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, params)
	if err != nil {
		if utils.ResponseWasNotFound(err) || strings.Contains(err.Error(), "VS402485") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Area failed: %w", err)
	}

	d.Set("name", node.Name)
	if node.Identifier != nil {
		d.Set("identifier", node.Identifier.String())
	}

	if node.Path != nil {
		d.Set("path", convertAreaNodePath(node.Path))
	}

	if node.HasChildren != nil {
		d.Set("has_children", converter.ToBool(node.HasChildren, false))
	}

	return nil
}

func resourceUpdateArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	postedNode := workitemtracking.WorkItemClassificationNode{
		Name: converter.String(d.Get("name").(string)),
	}

	args := workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		PostedNode:     &postedNode,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Path:           converter.String(d.Id()),
	}

	_, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("updating Area: %w", err)
	}

	return resourceReadArea(d, m)
}

func resourceDeleteArea(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	structureGroup := workitemtracking.TreeStructureGroupValues.Areas

	rootNode, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		StructureGroup: &structureGroup,
		Path:           converter.String(""),
		Depth:          converter.Int(1),
	})
	if err != nil {
		return fmt.Errorf("error getting root node for reclassification: %v", err)
	}

	if rootNode.Id == nil {
		return fmt.Errorf("error: Could not determine the ID of the root Node")
	}
	reclassifyID := *rootNode.Id

	deleteArgs := workitemtracking.DeleteClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		StructureGroup: &structureGroup,
		Path:           converter.String(d.Id()),
		ReclassifyId:   &reclassifyID,
	}

	err = clients.WorkItemTrackingClient.DeleteClassificationNode(clients.Ctx, deleteArgs)
	if err != nil {
		return fmt.Errorf("deleting Area: %v", err)
	}

	return nil
}

func convertAreaNodePath(path *string) string {
	itemPath := ""
	if path != nil {
		itemPathList := strings.Split(strings.ReplaceAll(*path, "\\", "/"), "/")
		if len(itemPathList) >= 3 {
			itemPath = strings.Join(itemPathList[3:], "/")
		}
	}
	return "/" + itemPath
}
