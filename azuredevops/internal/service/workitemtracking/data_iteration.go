package workitemtracking

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataIteration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIterationRead,
		Schema: map[string]*schema.Schema{
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
		},
	}
}

func dataSourceIterationRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	depth := 0
	if d.Get("fetch_children").(bool) {
		depth = 1
	}
	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Depth:          converter.Int(depth),
	}

	path, pathSet := d.GetOk("path")
	if pathSet {
		args.Path = converter.String(strings.TrimSpace(path.(string)))
	}

	iteration, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("Error getting Iteration with path %q: %w", path, err)
	}

	d.SetId(iteration.Identifier.String())
	if args.Path == nil {
		d.Set("path", convertIterationPath(iteration.Path))
	}
	d.Set("name", iteration.Name)
	d.Set("has_children", converter.ToBool(iteration.HasChildren, false))
	d.Set("children", flattenIterationNodes(projectID, iteration.Children))
	return nil
}

func flattenIterationNodes(projectID string, nodes *[]workitemtracking.WorkItemClassificationNode) []interface{} {
	if nodes == nil {
		return nil
	}

	results := make([]interface{}, len(*nodes))
	for i, v := range *nodes {
		results[i] = flattenIterationNode(projectID, v)
	}
	return results
}

func flattenIterationNode(projectID string, node workitemtracking.WorkItemClassificationNode) map[string]interface{} {
	output := make(map[string]interface{})

	output["id"] = node.Identifier.String()
	if node.Name != nil {
		output["name"] = *node.Name
	}
	output["project_id"] = projectID
	output["path"] = convertIterationPath(node.Path)
	output["has_children"] = converter.ToBool(node.HasChildren, false)

	return output
}

func convertIterationPath(path *string) string {
	itemPath := ""
	if path != nil {
		itemPathList := strings.Split(strings.ReplaceAll(*path, "\\", "/"), "/")
		if len(itemPathList) >= 3 {
			itemPath = strings.Join(itemPathList[3:], "/")
		}
	}
	return "/" + itemPath
}
