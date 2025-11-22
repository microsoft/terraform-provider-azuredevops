package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// CreateClassificationNodeSchema schema for classification node
func CreateClassificationNodeSchema(outer map[string]*schema.Schema) map[string]*schema.Schema {
	baseSchema := map[string]*schema.Schema{
		"project_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
		"path": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"fetch_children": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"has_children": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"children": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"project_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"path": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"has_children": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
	}

	for key, elem := range baseSchema {
		outer[key] = elem
	}

	return outer
}

func CreateClassificationNodeResourceSchema(structureType workitemtracking.TreeStructureGroup) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
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
		"node_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}

	if structureType == workitemtracking.TreeStructureGroupValues.Iterations {
		s["attributes"] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_date": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"finish_date": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		}
	}

	return s
}

// ReadClassificationNode reads a classification node from Azure DevOps
func ReadClassificationNode(clients *client.AggregatedClient, d *schema.ResourceData, structureType workitemtracking.TreeStructureGroup) error {
	projectID := d.Get("project_id").(string)
	depth := 0
	if d.Get("fetch_children").(bool) {
		depth = 1
	}
	params := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &structureType,
		Depth:          converter.Int(depth),
	}

	if path, ok := d.GetOk("path"); ok {
		params.Path = converter.String(strings.TrimSpace(path.(string)))
	}

	node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, params)
	if err != nil {
		// the following error will be returned in case the classification node isn't present
		// "VS402485: The Area/Iteration name is not recognized"
		d.SetId("")

		js, parseErr := json.Marshal(params)
		if parseErr != nil {
			return fmt.Errorf("Marshalling JSON. Error: %+v", parseErr)
		}
		return fmt.Errorf("getting ClassificationNode failed. %s. Error: %w", js, err)

	}

	d.SetId(node.Identifier.String())
	d.Set("name", node.Name)

	if node.Path != nil {
		d.Set("path", convertNodePath(node.Path))
	}

	if node.HasChildren != nil {
		d.Set("has_children", converter.ToBool(node.HasChildren, false))
	}

	d.Set("children", flattenClassificationChildNodes(projectID, node.Children))
	return nil
}

func CreateOrUpdateClassificationNode(clients *client.AggregatedClient, d *schema.ResourceData, structureType workitemtracking.TreeStructureGroup) error {
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
		StructureGroup: &structureType,
		Path:           converter.String(d.Get("path").(string)),
	}

	workItem, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, args)
	if err != nil {
		return err
	}

	if workItem.Identifier != nil {
		d.SetId(workItem.Identifier.String())
	} else if workItem.Id != nil {
		d.SetId(fmt.Sprintf("%d", *workItem.Id))
	}

	d.Set("has_children", converter.ToBool(workItem.HasChildren, false))
	d.Set("node_id", workItem.Id)

	return nil
}

func DeleteClassificationNode(clients *client.AggregatedClient, d *schema.ResourceData, structureType workitemtracking.TreeStructureGroup) error {
	rootNode, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		StructureGroup: &structureType,
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

	pathToDelete := d.Get("path").(string)
	if pathToDelete == "" {
		pathToDelete = d.Get("name").(string)
	}

	deleteArgs := workitemtracking.DeleteClassificationNodeArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		StructureGroup: &structureType,
		Path:           converter.String(pathToDelete),
		ReclassifyId:   &reclassifyID,
	}

	err = clients.WorkItemTrackingClient.DeleteClassificationNode(clients.Ctx, deleteArgs)

	if err != nil {
		return fmt.Errorf("wrror in getting the hierarchy: %v", err)
	}

	return nil
}

func flattenClassificationChildNodes(projectID string, nodes *[]workitemtracking.WorkItemClassificationNode) []interface{} {
	if nodes == nil {
		return nil
	}

	results := make([]interface{}, len(*nodes))
	for i, v := range *nodes {
		results[i] = flattenClassificationNode(projectID, v)
	}
	return results
}

func flattenClassificationNode(projectID string, node workitemtracking.WorkItemClassificationNode) map[string]interface{} {
	output := make(map[string]interface{})

	output["id"] = node.Identifier.String()
	if node.Name != nil {
		output["name"] = *node.Name
	}
	output["project_id"] = projectID
	output["path"] = convertNodePath(node.Path)
	output["has_children"] = converter.ToBool(node.HasChildren, false)

	return output
}

func convertNodePath(path *string) string {
	itemPath := ""
	if path != nil {
		itemPathList := strings.Split(strings.ReplaceAll(*path, "\\", "/"), "/")
		if len(itemPathList) >= 3 {
			itemPath = strings.Join(itemPathList[3:], "/")
		}
	}
	return "/" + itemPath
}
