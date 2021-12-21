package utils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/workitemtracking"
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

// ReadClassificationNode reads a classification node from Azure DevOps
func ReadClassificationNode(clients *client.AggregatedClient, d *schema.ResourceData, structureType workitemtracking.TreeStructureGroup) error {
	projectID := d.Get("project_id").(string)
	depth := 0
	if d.Get("fetch_children").(bool) {
		depth = 1
	}
	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &structureType,
		Depth:          converter.Int(depth),
	}

	path, pathSet := d.GetOk("path")
	if pathSet {
		args.Path = converter.String(strings.TrimSpace(path.(string)))
	}

	node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		// the following error will be returned in case the classification node isn't present
		// "VS402485: The Area/Iteration name is not recognized"
		d.SetId("")
		return fmt.Errorf("Error getting Iteration with path %q: %w", path, err)
	}

	d.SetId(node.Identifier.String())
	if args.Path == nil {
		d.Set("path", convertNodePath(node.Path))
	}
	d.Set("name", node.Name)
	d.Set("has_children", converter.ToBool(node.HasChildren, false))
	d.Set("children", flattenClassificationChildNodes(projectID, node.Children))
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
