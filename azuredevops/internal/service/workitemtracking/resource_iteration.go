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

func ResourceIteration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateIteration,
		Read:   resourceReadIteration,
		Update: resourceUpdateIteration,
		Delete: resourceDeleteIteration,
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
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"finish_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
					},
				},
			},
		},
	}
}

func resourceCreateIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	postedNode := workitemtracking.WorkItemClassificationNode{
		Name: converter.String(d.Get("name").(string)),
	}

	if attributesData, ok := d.GetOk("attributes"); ok {
		attributesList := attributesData.([]interface{})
		if len(attributesList) > 0 && attributesList[0] != nil {
			attrs := attributesList[0].(map[string]interface{})
			nodeAttributes := make(map[string]interface{})

			if v, ok := attrs["start_date"].(string); ok && v != "" {
				nodeAttributes["startDate"] = v
			}
			if v, ok := attrs["finish_date"].(string); ok && v != "" {
				nodeAttributes["finishDate"] = v
			}

			if len(nodeAttributes) > 0 {
				postedNode.Attributes = &nodeAttributes
			}
		}
	}

	args := workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		PostedNode:     &postedNode,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Path:           converter.String(d.Get("path").(string)),
	}

	workItem, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("creating Iteration: %w", err)
	}

	if workItem.Id != nil {
		d.SetId(strconv.Itoa(*workItem.Id))
	} else {
		return fmt.Errorf("creating Iteration: API did not return an ID")
	}

	return resourceReadIteration(d, m)
}

func resourceReadIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)

	params := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Path:           converter.String(d.Id()),
		Depth:          converter.Int(0),
	}

	node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, params)
	if err != nil {
		if utils.ResponseWasNotFound(err) || strings.Contains(err.Error(), "VS402485") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Iteration failed: %w", err)
	}

	d.Set("name", node.Name)
	if node.Identifier != nil {
		d.Set("identifier", node.Identifier.String())
	}

	if node.Path != nil {
		d.Set("path", convertIterationNodePath(node.Path))
	}

	if node.HasChildren != nil {
		d.Set("has_children", converter.ToBool(node.HasChildren, false))
	}

	if node.Attributes != nil {
		attrs := make(map[string]interface{})
		for k, v := range *node.Attributes {
			if k == "startDate" && v != nil {
				attrs["start_date"] = v.(string)
			}
			if k == "finishDate" && v != nil {
				attrs["finish_date"] = v.(string)
			}
		}
		if len(attrs) > 0 {
			d.Set("attributes", []interface{}{attrs})
		} else {
			d.Set("attributes", nil)
		}
	} else {
		d.Set("attributes", nil)
	}

	return nil
}

func resourceUpdateIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	postedNode := workitemtracking.WorkItemClassificationNode{
		Name: converter.String(d.Get("name").(string)),
	}

	if attributesData, ok := d.GetOk("attributes"); ok {
		attributesList := attributesData.([]interface{})
		if len(attributesList) > 0 && attributesList[0] != nil {
			attrs := attributesList[0].(map[string]interface{})
			nodeAttributes := make(map[string]interface{})

			if v, ok := attrs["start_date"].(string); ok && v != "" {
				nodeAttributes["startDate"] = v
			}
			if v, ok := attrs["finish_date"].(string); ok && v != "" {
				nodeAttributes["finishDate"] = v
			}

			if len(nodeAttributes) > 0 {
				postedNode.Attributes = &nodeAttributes
			}
		}
	}

	args := workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		PostedNode:     &postedNode,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Path:           converter.String(d.Id()),
	}

	_, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("updating Iteration: %w", err)
	}

	return resourceReadIteration(d, m)
}

func resourceDeleteIteration(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	structureGroup := workitemtracking.TreeStructureGroupValues.Iterations

	// Need to get root node to find reclassify ID
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
		return fmt.Errorf("deleting Iteration: %v", err)
	}

	return nil
}

func convertIterationNodePath(path *string) string {
	itemPath := ""
	if path != nil {
		itemPathList := strings.Split(strings.ReplaceAll(*path, "\\", "/"), "/")
		if len(itemPathList) >= 3 {
			itemPath = strings.Join(itemPathList[3:], "/")
		}
	}
	return "/" + itemPath
}
